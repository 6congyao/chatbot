package main

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/ybbus/httpretry"
)

type Recipient struct {
	Id string `json:"id,omitempty"`
}

type Message struct {
	Text string `json:"text,omitempty"`
}

type Request struct {
	Recipient     Recipient `json:"recipient,omitempty"`
	Message       Message   `json:"message,omitempty"`
	MessagingType string    `json:"messaging_type,omitempty"`
}

// func main() {
// 	SendMessage("msg for test2", "7436026376434696", "RESPONSE")
// }

func SendMessage(message string, customerId string, messageType string) error {
	rcp := Recipient{
		Id: customerId,
	}
	msg := Message{
		Text: message,
	}
	req := Request{
		Recipient:     rcp,
		Message:       msg,
		MessagingType: messageType,
	}
	var res string
	// json, _ := json.MarshalIndent(req, "", " ")
	// fmt.Println(string(json))

	cl := httpretry.NewDefaultClient() //Used for retries
	var err = requests.
		URL("https://graph.facebook.com/v19.0/257662950760575/messages").
		Param("access_token", "EAAzg2pLO8LoBOzZAXxXCEQLWGDxFlM1D2vmv3wyEcAQtTOxOLGRodkFdoggJiFritaOUKErKeVbcPDvq6AcNihwaXnI7NHAFcUi0rYTxCc1jO9AKAtSC4a76ywP7c6yB13VyDPpjOiWYw4Wrb4sHHSrptxHH2wUf9U7CXKWphifd0GjZAiSIKreMZCQS2pAWtwHZAumE").
		BodyJSON(&req).ToString(&res).
		Client(cl).
		Fetch(context.Background())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	return err
}