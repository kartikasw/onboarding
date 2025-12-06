package request

type ResetPasswordRequest struct {
	Email string `form:"email" binding:"required,validEmail"`
}

type ResetPassword struct {
	Email          string `form:"email" binding:"required,validEmail"`
	OTP            string `form:"otp" binding:"required"`
	NewPassword    string `form:"new_password" binding:"required,validPassword"`
	VerifyPassword string `form:"verify_password" binding:"required"`
}
