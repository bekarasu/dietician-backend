package response

import "time"

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	IsVerified bool      `json:"isVerified"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type RegisterResponse struct {
	OTPToken string `json:"otpToken"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
