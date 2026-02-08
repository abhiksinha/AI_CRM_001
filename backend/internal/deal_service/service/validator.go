package service

import (
	"CRM/internal/deal_service/contracts"

	"github.com/go-ozzo/ozzo-validation/v4"
)

// ValidateCreateDealRequest performs validation on the CreateDealRequest struct.
func ValidateCreateDealRequest(req contracts.CreateDealRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required.Error("deal name is required")),
		validation.Field(&req.Stage, validation.Required.Error("deal stage is required")),
		validation.Field(&req.Value, validation.Min(0.0).Error("value must be non-negative")),
		validation.Field(&req.ContactID, validation.Required.Error("contact id is required")),
		validation.Field(&req.OwnerID, validation.Required.Error("owner id is required")),
		validation.Field(&req.ExpectedCloseDate, validation.Date("2006-01-02").Error("expected close date must be in YYYY-MM-DD format")),
	)
}

// ValidateUpdateDealRequest performs validation on the UpdateDealRequest struct.
func ValidateUpdateDealRequest(req contracts.UpdateDealRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Value, validation.When(req.Value != nil, validation.Min(0.0).Error("value must be non-negative"))),
		validation.Field(&req.ExpectedCloseDate, validation.When(req.ExpectedCloseDate != nil, validation.Date("2006-01-02").Error("expected close date must be in YYYY-MM-DD format"))),
	)
}

// ValidateCreateTaskRequest performs validation on the CreateTaskRequest struct.
func ValidateCreateTaskRequest(req contracts.CreateTaskRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Title, validation.Required.Error("task title is required"), validation.Length(1, 255)),
		validation.Field(&req.AssignedToID, validation.Required.Error("assigned to id is required")),
	)
}
