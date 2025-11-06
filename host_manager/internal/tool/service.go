package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// IPGeoInfo represents the geographical and ISP information of an IP address
type IPGeoInfo struct {
	Organization    string  `json:"organization"`
	Longitude       float64 `json:"longitude"`
	City            string  `json:"city"`
	Timezone        string  `json:"timezone"`
	ISP             string  `json:"isp"`
	Offset          int     `json:"offset"`
	Region          string  `json:"region"`
	ASN             int     `json:"asn"`
	ASNOrganization string  `json:"asn_organization"`
	Country         string  `json:"country"`
	IP              string  `json:"ip"`
	Latitude        float64 `json:"latitude"`
	ContinentCode   string  `json:"continent_code"`
	CountryCode     string  `json:"country_code"`
	RegionCode      string  `json:"region_code"`
}

// DNSResolver defines the interface for domain name resolution
type DNSResolver interface {
	// ResolveDomain resolves a domain name to IP addresses
	ResolveDomain(domain string) ([]string, error)
}

// IPGeoService defines the interface for IP geolocation services
type IPGeoService interface {
	// GetIPInfo retrieves geographical information for the given IP
	// If ip is empty, it returns information for the caller's IP
	GetIPInfo(ip string) (*IPGeoInfo, error)
}

// DefaultDNSResolver implements DNSResolver using Go's standard net package
type DefaultDNSResolver struct {
	timeout time.Duration
}

// NewDefaultDNSResolver creates a new DNS resolver with the specified timeout
func NewDefaultDNSResolver(timeout time.Duration) *DefaultDNSResolver {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &DefaultDNSResolver{
		timeout: timeout,
	}
}

// ResolveDomain resolves a domain name to IP addresses
func (r *DefaultDNSResolver) ResolveDomain(domain string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	resolver := &net.Resolver{}
	ips, err := resolver.LookupHost(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve domain %s: %w", domain, err)
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no IP addresses found for domain %s", domain)
	}

	return ips, nil
}

// IPSBGeoService implements IPGeoService using api.ip.sb
type IPSBGeoService struct {
	baseURL    string
	httpClient *http.Client
}

// NewIPSBGeoService creates a new IP geolocation service using api.ip.sb
func NewIPSBGeoService() *IPSBGeoService {
	return &IPSBGeoService{
		baseURL: "https://api.ip.sb/geoip",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetIPInfo retrieves geographical information for the given IP
func (s *IPSBGeoService) GetIPInfo(ip string) (*IPGeoInfo, error) {
	url := s.baseURL
	if ip != "" {
		url = fmt.Sprintf("%s/%s", s.baseURL, ip)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Host", "api.ip.sb")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query IP info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var info IPGeoInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &info, nil
}

// ToolService aggregates DNS resolver and IP geo service
type ToolService struct {
	dnsResolver  DNSResolver
	ipGeoService IPGeoService
}

// NewToolService creates a new tool service with default implementations
func NewToolService() *ToolService {
	return &ToolService{
		dnsResolver:  NewDefaultDNSResolver(5 * time.Second),
		ipGeoService: NewIPSBGeoService(),
	}
}

// NewToolServiceWithProviders creates a new tool service with custom providers
func NewToolServiceWithProviders(resolver DNSResolver, geoService IPGeoService) *ToolService {
	return &ToolService{
		dnsResolver:  resolver,
		ipGeoService: geoService,
	}
}

// ResolveDomain resolves a domain name to IP addresses
func (s *ToolService) ResolveDomain(domain string) ([]string, error) {
	return s.dnsResolver.ResolveDomain(domain)
}

// GetIPInfo retrieves geographical information for the given IP
func (s *ToolService) GetIPInfo(ip string) (*IPGeoInfo, error) {
	return s.ipGeoService.GetIPInfo(ip)
}

// SetDNSResolver sets a custom DNS resolver
func (s *ToolService) SetDNSResolver(resolver DNSResolver) {
	s.dnsResolver = resolver
}

// SetIPGeoService sets a custom IP geolocation service
func (s *ToolService) SetIPGeoService(geoService IPGeoService) {
	s.ipGeoService = geoService
}
