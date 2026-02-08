package contracts

// DealResponse is the standard response for a single deal.
type DealResponse struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Stage             string  `json:"stage"`
	Value             float64 `json:"value"`
	ExpectedCloseDate string  `json:"expected_close_date"`
	ContactID         string  `json:"contact_id"`
	OwnerID           string  `json:"owner_id"`
	CreatedAt         int64   `json:"created_at"`
	UpdatedAt         int64   `json:"updated_at"`
}

// ListDealsResponse is the response for listing multiple deals.
type ListDealsResponse struct {
	Data       []DealResponse `json:"data"`
	TotalCount int64          `json:"total_count"`
}

// TaskResponse is the standard response for a single task.
type TaskResponse struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	DueDate      int64  `json:"due_date"`
	IsCompleted  bool   `json:"is_completed"`
	AssignedToID string `json:"assigned_to_id"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

// ListTasksResponse is the response for listing multiple tasks.
type ListTasksResponse struct {
	Data       []TaskResponse `json:"data"`
	TotalCount int64          `json:"total_count"`
}
