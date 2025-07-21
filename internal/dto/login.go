package dto

type (
	RegisterInput struct {
		VerifyKey     string `json:"verifyKey" validate:"required"`
		VerifyType    int    `json:"verifyType" validate:"required"`
		VerifyPurpose string `json:"verifyPurpose" validate:"required"`
	}
	RegisterOutput struct {
	}
)

type (
	VerifyInput struct {
		VerifyKey  string `json:"verifyKey" validate:"required"`
		VerifyCode string `json:"verifyCode" validate:"required"`
	}

	VerifyOutput struct {
		Token   string `json:"token"`
		Message string `json:"message"`
	}
)

type (
	UpdateUserPasswordInput struct {
		UserToken    string `json:"userToken"`
		UserPassword string `json:"userPassword"`
	}
	UpdateUserPasswordOutput struct {
	}
)

type (
	LoginUserInput struct {
		UserAccount  string `json:"userAccount"`
		UserPassword string `json:"userPassword"`
	}

	LoginUserOutput struct {
		Token   string `json:"token"`
		Message string `json:"message"`
	}
)

type (
	SetupTwoFactorAuthInput struct {
		UserId            uint32 `json:"userId"`
		TwoFactorAuthType string `json:"twoFactorAuthType"`
		TwoFactorEmail    string `json:"twoFactorEmail"`
	}
	SetupTwoFactorAuthOutput struct {
	}
)

type (
	TwoFactorVerifyInput struct {
		UserId        uint32 `json:"userId"`
		TwoFactorCode string `json:"twoFactorCode"`
	}

	TwoFactorVerifyOutput struct {
	}
)
