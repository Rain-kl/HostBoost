package host

import (
	"errors"
	"fmt"
	"strings"
)

// Service coordinates host operations and validation.
type Service struct {
	repo *FileRepository
}

// NewService instantiates a host service.
func NewService(repo *FileRepository) *Service {
	return &Service{repo: repo}
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
	req.Domain = normalizeDomain(req.Domain)
	if req.Domain == "" {
		return errors.New("domain is required")
	}

	host := Host{
		Domain: req.Domain,
		IP:     "6.6.6.6",
		Type:   "cloudflare",
	}

	if err := s.repo.Create(host); err != nil {
		if errors.Is(err, ErrHostExists) {
			return fmt.Errorf("host %s already exists", req.Domain)
		}
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

	return nil
}

func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(strings.ToLower(domain))
	return domain
}
