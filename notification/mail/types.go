package mail

type TemplateDataWelcomeEmail struct {
	Email string
	FirstName string
}

type TemplateDataVerificationEmail struct {
	Email string
	FirstName string
	Otp string
}
