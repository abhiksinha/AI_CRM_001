package contracts

// ListDealsRequest defines the query parameters for listing deals.
type ListDealsRequest struct {
	Page     int    `schema:"page"`
	PageSize int    `schema:"page_size"`
	Stage    string `schema:"stage"` // Filter by deal stage
}

// CreateDealRequest defines the JSON body for creating a new deal.
type CreateDealRequest struct {
	Name              string  `json:"name"`
	Stage             string  `json:"stage"`
	Value             float64 `json:"value"`
	ExpectedCloseDate string  `json:"expected_close_date"` // e.g., "2024-12-31"
	ContactID         string  `json:"contact_id"`
	OwnerID           string  `json:"owner_id"`
}

// UpdateDealRequest defines the JSON body for updating a deal.
type UpdateDealRequest struct {
	Name              *string  `json:"name"`
	Stage             *string  `json:"stage"`
	Value             *float64 `json:"value"`
	ExpectedCloseDate *string  `json:"expected_close_date"`
}

// CreateTaskRequest defines the JSON body for creating a task.
type CreateTaskRequest struct {
	Title        string `json:"title"`
	DueDate      string `json:"due_date"`
	IsCompleted  bool   `json:"is_completed"`
	AssignedToID string `json:"assigned_to_id"`
}
