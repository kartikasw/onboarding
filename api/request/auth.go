package request

type RegisterRequest struct {
	Email    string `form:"email" binding:"required,validEmail"`
	Password string `form:"password" binding:"required,validPassword"`
}

type LoginRequest struct {
	Email    string `form:"email" binding:"required,validEmail"`
	Password string `form:"password" binding:"required"`
}
