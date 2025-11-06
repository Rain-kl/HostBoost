package tool

// DomainDetail represents detailed information about a domain.
type DomainDetail struct {
	Organization string `json:"organization"`
	IP           string `json:"ip"`
	ISP          string `json:"isp"`
}

// DetailResponse models the OpenAPI response for domain details.
type DetailResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Data    *DomainDetail `json:"data,omitempty"`
}
