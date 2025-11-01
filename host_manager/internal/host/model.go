package host

// Host represents a single host entry stored in the simulated hosts file.
type Host struct {
	Domain string `json:"domain"`
	IP     string `json:"ip"`
	Type   string `json:"type"`
}

// QueryHostResponse models the OpenAPI response for querying a single host.
type QueryHostResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Host   `json:"data"`
}

// QueryHostListResponse models the OpenAPI response for listing hosts.
type QueryHostListResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    QueryHostListResult `json:"data"`
}

// QueryHostListResult wraps the host list and total for list responses.
type QueryHostListResult struct {
	Total int    `json:"total"`
	List  []Host `json:"list"`
}

// MutationResponse represents a response for create/delete operations.
type MutationResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// AddHostRequest captures the expected payload when creating a host.
type AddHostRequest struct {
	Domain string `json:"domain"` // logical host name
}

// DeleteHostRequest captures the expected payload when removing a host.
type DeleteHostRequest struct {
	Domain string `json:"domain"`
}
