package request

type RegisterRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=150"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	PhoneNumber string `json:"phone_number"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}
