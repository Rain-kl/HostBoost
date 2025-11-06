package host

import (
	"errors"
	"fmt"
	"hostMgr/hostsync"
	"hostMgr/internal/extSvc"
	"hostMgr/internal/opt"
	"log"
	"strings"
)

// Service coordinates host operations and validation.
type Service struct {
	repo   *FileRepository
	syncer *hostsync.Syncer
}

// NewService instantiates a host service.
func NewService(repo *FileRepository) *Service {
	return &Service{
		repo:   repo,
		syncer: hostsync.NewSyncer("hosts.json"),
	}
}

// ListHosts retrieves all registered hosts.
func (s *Service) ListHosts() []Host {
	return s.repo.List()
}

// GetHost returns the host by domain or a wrapped error if missing.
func (s *Service) GetHost(domain string) (Host, error) {
	domain = normalizeDomain(domain)
	if domain == "" {
		return Host{}, errors.New("domain is required")
	}

	host, err := s.repo.Get(domain)
	if err != nil {
		return Host{}, err
	}

	return host, nil
}

// CreateHost validates and registers a new host entry.
func (s *Service) CreateHost(req AddHostRequest) error {
	cdnType := "cloudflare"
	req.Domain = normalizeDomain(req.Domain)
	if req.Domain == "" {
		return errors.New("domain is required")
	}

	optSvc, ok := extSvc.OptService.(*opt.Service)
	if !ok || optSvc == nil {
		return errors.New("opt service not initialized")
	}

	_, optIp, err := optSvc.GetCurrentOpt(cdnType)
	if err != nil {
		return err
	}
	host := Host{
		Domain: req.Domain,
		IP:     optIp.IP,
		Type:   cdnType,
	}

	if err := s.repo.Create(host); err != nil {
		if errors.Is(err, ErrHostExists) {
			return fmt.Errorf("host %s already exists", req.Domain)
		}
		return err
	}

	// 同步到系统 hosts 文件
	if err := s.syncer.Sync(); err != nil {
		log.Printf("Warning: failed to sync hosts to system: %v", err)
		return err
	}

	return nil
}

// DeleteHost removes a host by domain.
func (s *Service) DeleteHost(domain string) error {
	domain = normalizeDomain(domain)
	if domain == "" {
		return errors.New("domain is required")
	}

	if err := s.repo.Delete(domain); err != nil {
		if errors.Is(err, ErrHostNotFound) {
			return fmt.Errorf("host %s not found", domain)
		}
		return err
	}

	// 同步到系统 hosts 文件
	if err := s.syncer.Sync(); err != nil {
		log.Printf("Warning: failed to sync hosts to system: %v", err)
	}

	return nil
}

// UpdateHostsByType updates the IP address for all hosts of the specified type.
// It fetches the current optimal IP for the given type from opt service,
// then updates all hosts with that type to use this IP.
func (s *Service) UpdateHostsByType(hostType string) (int, error) {
	if hostType == "" {
		return 0, errors.New("host type is required")
	}

	// Get current optimal IP for the specified type
	optSvc, ok := extSvc.OptService.(*opt.Service)
	if !ok || optSvc == nil {
		return 0, errors.New("opt service not initialized")
	}

	_, optInfo, err := optSvc.GetCurrentOpt(hostType)
	if err != nil {
		return 0, fmt.Errorf("failed to get current opt for type %s: %w", hostType, err)
	}

	newIP := optInfo.IP
	if newIP == "" {
		return 0, fmt.Errorf("no IP found for type %s", hostType)
	}

	// Get all hosts with the specified type
	hosts := s.repo.ListByType(hostType)
	if len(hosts) == 0 {
		return 0, fmt.Errorf("no hosts found with type %s", hostType)
	}

	// Update each host's IP
	updatedCount := 0
	var updateErrors []string

	for _, host := range hosts {
		if err := s.repo.UpdateIP(host.Domain, newIP); err != nil {
			updateErrors = append(updateErrors, fmt.Sprintf("failed to update %s: %v", host.Domain, err))
		} else {
			updatedCount++
		}
	}

	if len(updateErrors) > 0 {
		return updatedCount, fmt.Errorf("updated %d hosts with IP %s, but encountered errors: %s", updatedCount, newIP, strings.Join(updateErrors, "; "))
	}

	return updatedCount, nil
}

func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(strings.ToLower(domain))
	return domain
}
