package utils

var OTP = OTPGenerator{}

type OTPGeneratorInterface interface {
	GenerateOTP() string
}

type OTPGenerator struct {}

func (o OTPGenerator) GenerateOTP() string {
	charset = numbers
	return RandomString(6)
}
