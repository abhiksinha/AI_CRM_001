package service

import (
	"CRM/internal/user_service/contracts"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// ... (existing validators)

func ValidateCreateApiKeyRequest(req contracts.CreateApiKeyRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required.Error("user_id is required")),
	)
}

func ValidateMatchApiKeyRequest(req contracts.MatchApiKeyRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ApiKey, validation.Required.Error("api_key is required")),
	)
}

func ValidateExpireApiKeyRequest(req contracts.ExpireApiKeyRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ApiKey, validation.Required.Error("api_key is required")),
	)
}
func ValidateCreateUserRequest(req contracts.CreateUserRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.FirstName, validation.Required.Error("first name is required")),
		validation.Field(&req.LastName, validation.Required.Error("last name is required")),
		validation.Field(&req.Email, validation.Required.Error("email is required"), is.EmailFormat),
		validation.Field(&req.Password, validation.Required.Error("password is required"), validation.By(func(value interface{}) error {
			p, _ := value.(contracts.Password)
			return validation.Validate(string(p), validation.Length(8, 100))
		})),
		validation.Field(&req.Role, validation.Required.Error("role is required"), validation.In("admin", "user").Error("role must be 'admin' or 'user'")),
	)
}
func ValidateUpdateUserRequest(req contracts.UpdateUserRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Role, validation.When(req.Role != nil, validation.In("admin", "user").Error("role must be 'admin' or 'user'"))),
	)
}
func ValidateListUsersRequest(req contracts.ListUsersRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Page, validation.Min(1).Error("page must be greater than 0")),
		validation.Field(&req.PageSize, validation.Min(1).Error("page_size must be greater than 0")),
		validation.Field(&req.Role, validation.When(req.Role != "", validation.In("admin", "user").Error("role must be 'admin' or 'user'"))),
	)
}
func ValidateVerifyPasswordRequest(req contracts.VerifyPasswordRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ID, validation.Required.Error("user id is required")),
		validation.Field(&req.Password, validation.Required.Error("password is required")),
	)
}
