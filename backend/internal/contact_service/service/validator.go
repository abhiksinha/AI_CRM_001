package service

import (
	"CRM/internal/contact_service/contracts"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
)

// ValidateCreateContactRequest performs validation on the CreateContactRequest struct.
func ValidateCreateContactRequest(req contracts.CreateContactRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.FirstName, validation.Required.Error("first name is required")),
		validation.Field(&req.Email, validation.Required.Error("email is required"), is.EmailFormat),
		// Phone is now optional. If provided, it must match the E.164 format.
		validation.Field(&req.Phone,
			validation.When(req.Phone != "", validation.Match(regexp.MustCompile(`^\+\d{1,15}$`)).Error("phone number must be in E.164 format (e.g., +12125552368)")),
		),
	)
}

// ValidateUpdateContactRequest performs validation on the UpdateContactRequest struct.
func ValidateUpdateContactRequest(req contracts.UpdateContactRequest) error {
	return validation.ValidateStruct(&req,
		// If email is provided, it must be a valid format.
		validation.Field(&req.Email, validation.When(req.Email != nil, is.EmailFormat)),
		// If phone is provided, it must be a valid format.
		validation.Field(&req.Phone,
			validation.When(req.Phone != nil, validation.Match(regexp.MustCompile(`^\+\d{1,15}$`)).Error("phone number must be in E.164 format")),
		),
	)
}

// ValidateAddNoteRequest performs validation on the AddNoteRequest struct.
func ValidateAddNoteRequest(req contracts.AddNoteRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Content, validation.Required.Error("note content cannot be empty"), validation.Length(1, 4096)),
	)
}
