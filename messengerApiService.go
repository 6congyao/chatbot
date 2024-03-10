package main

import (
	"context"

	"github.com/carlmjohnson/requests"
	"github.com/ybbus/httpretry"
)

const (
	FB_API_ENDPOINT      = "https://graph.facebook.com"
	FB_PAGE_ACCESS_TOKEN = "EAAzg2pLO8LoBOzZAXxXCEQLWGDxFlM1D2vmv3wyEcAQtTOxOLGRodkFdoggJiFritaOUKErKeVbcPDvq6AcNihwaXnI7NHAFcUi0rYTxCc1jO9AKAtSC4a76ywP7c6yB13VyDPpjOiWYw4Wrb4sHHSrptxHH2wUf9U7CXKWphifd0GjZAiSIKreMZCQS2pAWtwHZAumE"
	FB_PAGE_ID           = "257662950760575"
	LATEST_API_VERSION   = "v19.0"
	API_SEPARATOR        = "/"
	MESSAGES_API_NAME    = "messages"
	baseEndpoint         = FB_API_ENDPOINT + API_SEPARATOR + LATEST_API_VERSION + API_SEPARATOR + FB_PAGE_ID + API_SEPARATOR + MESSAGES_API_NAME
)

type FacebookRecipient struct {
	Id string `json:"id,omitempty"`
}

type FacebookMessage struct {
	Text string `json:"text,omitempty"`
}

type FacebookRequest struct {
	Recipient     FacebookRecipient `json:"recipient,omitempty"`
	Message       FacebookMessage   `json:"message,omitempty"`
	MessagingType string            `json:"messaging_type,omitempty"`
}

func sendMessage(message string, customerId string, messageType string) error {
	rcp := FacebookRecipient{
		Id: customerId,
	}
	msg := FacebookMessage{
		Text: message,
	}
	req := FacebookRequest{
		Recipient:     rcp,
		Message:       msg,
		MessagingType: messageType,
	}

	cl := httpretry.NewDefaultClient() //Used for retries
	var err = requests.
		URL(baseEndpoint).
		Param("access_token", FB_PAGE_ACCESS_TOKEN).
		BodyJSON(&req).
		Client(cl).
		Fetch(context.Background())

	return err
}
