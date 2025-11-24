package host

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

var (
	// ErrHostExists indicates attempting to create a duplicate host entry.
	ErrHostExists = errors.New("host already exists")
	// ErrHostNotFound indicates the requested host entry is absent.
	ErrHostNotFound = errors.New("host not found")
)

// FileRepository manages hosts persisted in a JSON file to simulate /etc/hosts.
type FileRepository struct {
	Path  string // exported for use by Service to create syncer with same path
	hosts map[string]Host
	mu    sync.RWMutex
}

// NewFileRepository constructs a file-backed repository and loads existing state.
func NewFileRepository(path string) (*FileRepository, error) {
	repo := &FileRepository{
		Path:  path,
		hosts: make(map[string]Host),
	}

	if err := repo.ensureFile(); err != nil {
		return nil, fmt.Errorf("ensure file: %w", err)
	}

	if err := repo.load(); err != nil {
		return nil, fmt.Errorf("load data: %w", err)
	}

	return repo, nil
}

// List returns the complete host collection.
func (r *FileRepository) List() []Host {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Host, 0, len(r.hosts))
	for _, h := range r.hosts {
		result = append(result, h)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Domain < result[j].Domain
	})

	return result
}

// Get fetches a host by domain.
func (r *FileRepository) Get(domain string) (Host, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	host, ok := r.hosts[domain]
	if !ok {
		return Host{}, ErrHostNotFound
	}

	return host, nil
}

// Create persists a new host entry.
func (r *FileRepository) Create(host Host) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.hosts[host.Domain]; exists {
		return ErrHostExists
	}

	r.hosts[host.Domain] = host

	return r.persistLocked()
}

// Delete removes a host by domain.
func (r *FileRepository) Delete(domain string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.hosts[domain]; !exists {
		return ErrHostNotFound
	}

	delete(r.hosts, domain)

	return r.persistLocked()
}

// ListByType returns all hosts with the specified type.
func (r *FileRepository) ListByType(hostType string) []Host {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Host, 0)
	for _, h := range r.hosts {
		if h.Type == hostType {
			result = append(result, h)
		}
	}

	return result
}

// UpdateIP updates the IP address of a host by domain.
func (r *FileRepository) UpdateIP(domain, newIP string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	host, exists := r.hosts[domain]
	if !exists {
		return ErrHostNotFound
	}

	host.IP = newIP
	r.hosts[domain] = host

	return r.persistLocked()
}

func (r *FileRepository) ensureFile() error {
	dir := filepath.Dir(r.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		if err := os.WriteFile(r.Path, []byte("[]"), 0o644); err != nil {
			return err
		}
	}

	return nil
}

func (r *FileRepository) load() error {
	raw, err := os.ReadFile(r.Path)
	if err != nil {
		return err
	}

	var entries []Host
	if len(raw) == 0 {
		entries = []Host{}
	} else if err := json.Unmarshal(raw, &entries); err != nil {
		return err
	}

	for _, h := range entries {
		r.hosts[h.Domain] = h
	}

	return nil
}

func (r *FileRepository) persistLocked() error {
	entries := make([]Host, 0, len(r.hosts))
	for _, h := range r.hosts {
		entries = append(entries, h)
	}

	payload, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(r.Path), "hosts-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(payload); err != nil {
		tmp.Close()
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), r.Path)
}
