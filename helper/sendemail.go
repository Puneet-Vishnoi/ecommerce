package helper

import (
	"ecommerce-project/constant"
	"ecommerce-project/types"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmailSendGrid(req types.Verification)(types.Verification, error){
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == ""{
		return req, errors.New("SENDGRID_API_KEY environment variable is not set")
	}

	// Create a SendGrid client
	client := sendgrid.NewSendClient(apiKey)

	// Set up the email message
	from := mail.NewEmail("Sender Name", constant.Sender)
	to := mail.NewEmail("Recipient Name", req.Email)
	subject := "OTP verification mail"

	otp := Randomnum()
	req.Otp = int64(otp)
	htmlContent := "<p>This is a test otp forverification <strong>" + strconv.Itoa(otp) + "</strong> </p>"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	// Send the email message
	_, err := client.Send(message)
	if err != nil {
		return req, err
	}

	return req, nil;
}

func Randomnum() int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
  return rng.Intn(1000)+1000 // OTP length of 4 digits
}