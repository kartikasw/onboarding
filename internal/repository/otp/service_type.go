package otp

type ServiceType struct {
	Code string
	Name string
}

var (
	ServiceForgotPassword = ServiceType{
		Code: "forgot",
		Name: "Forgot Password",
	}
)
