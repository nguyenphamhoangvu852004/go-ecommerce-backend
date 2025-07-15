package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type MailRequest struct {
	ToMail      string `json:"toEmail"`
	MessageBody string `json:"messageBody"`
	Subject     string `json:"subject"`
	Attachment  string `json:"attachment"`
}

func SendEmailToGoByAPI(OTP int32, email string, purpose string) error {
	postUrl := "http://localhost:8081/api/v1/email"

	mailRequest := &MailRequest{
		ToMail:      email,
		MessageBody: fmt.Sprintf("Verify OTP is %d", OTP),
		Subject:     "Verify OTP " + purpose,
		Attachment:  "path/to/email",
	}

	// convert struct to json
	rqBody, err := json.Marshal(mailRequest)
	if err != nil {
		return err
	}

	//create request
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(rqBody))
	if err != nil {
		return err
	}

	//exec
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Print("Response Status: ", resp)
	return nil
}
