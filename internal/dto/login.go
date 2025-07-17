package dto

type RegisterInput struct {
	VerifyKey     string `json:"verifyKey" validate:"required"`
	VerifyType    int    `json:"verifyType" validate:"required"`
	VerifyPurpose string `json:"verifyPurpose" validate:"required"`
}

type VerifyInput struct {
	VerifyKey  string `json:"verifyKey" validate:"required"`
	VerifyCode string `json:"verifyCode" validate:"required"`
}

type VerifyOutput struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
