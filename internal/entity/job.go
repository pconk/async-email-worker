package entity

type EmailJob struct {
	ID      string `json:"id"`
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Message string `json:"message" validate:"required"`
	Retry   int    `json:"retry"`
}
