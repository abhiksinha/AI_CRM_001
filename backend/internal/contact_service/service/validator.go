package service

import (
	"CRM/internal/contact_service/contracts"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// ValidateCreateContactRequest performs validation on the CreateContactRequest struct.
func ValidateCreateContactRequest(req contracts.CreateContactRequest) error {
	return validation.ValidateStruct(&req,
		// FirstName is required.
		validation.Field(&req.FirstName, validation.Required.Error("first name is required")),

		// Email is required and must be in a valid format.
		validation.Field(&req.Email, validation.Required.Error("email is required"), is.EmailFormat),

		// Phone is required and must match the specified regex for country code and numeric format.
		validation.Field(&req.Phone,
			validation.Required.Error("phone number is required"),
			validation.Match(regexp.MustCompile(`^\+\d{1,15}$`)).Error("phone number must be in E.164 format (e.g., +12125552368)"),
		),

		// LastName and OwnerID are optional, so no validation rules are applied here.
	)
}
