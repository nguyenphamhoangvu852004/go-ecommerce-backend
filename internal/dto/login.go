package dto

type RegisterInput struct {
	VerifyKey     string `json:"verifyKey" validate:"required"`
	VerifyType    int    `json:"verifyType" validate:"required"`
	VerifyPurpose string `json:"verifyPurpose" validate:"required"`
}
