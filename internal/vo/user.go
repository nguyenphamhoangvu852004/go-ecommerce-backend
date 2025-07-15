package vo

type UserRegistrationRequest struct {
	Email   string `json:"email" validate:"required"`
	Purpose string `json:"purpose" validate:"required"`
}
