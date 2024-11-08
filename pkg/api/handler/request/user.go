package request

type RegisterUserRequest struct {
	Username        string `json:"username" binding:"required,alphanum,min=3,max=20"`    // Ensures username is alphanumeric and within length limits
	Email           string `json:"email" binding:"required,email"`                       // Validates email format
	Password        string `json:"password" binding:"required,min=8"`                    // Ensures password is at least 8 characters
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"` // Ensures passwords match
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
