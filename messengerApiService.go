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
	id string
}

type FacebookMessage struct {
	text string
}

type FacebookRequest struct {
	recipient      FacebookRecipient
	message        FacebookMessage
	messaging_type string
}

func sendMessage(message string, customerId string, messageType string) error {
	// customerIdMap := make(map[string]string)
	// customerIdMap["id"] = customerId
	// var customerIdJSON, _customerIdErr = json.Marshal(&customerIdMap)

	// messageMap := make(map[string]string)
	// messageMap["text"] = message
	// var messageMapJSON, _messageErr = json.Marshal(&messageMap)

	// if _customerIdErr != nil || _messageErr != nil {
	// 	return errors.New("JSON serialization error")
	// }
	rcp := FacebookRecipient{
		id: customerId,
	}
	msg := FacebookMessage{
		text: message,
	}
	req := FacebookRequest{
		recipient:      rcp,
		message:        msg,
		messaging_type: messageType,
	}

	cl := httpretry.NewDefaultClient() //Used for retries
	var err = requests.
		URL(baseEndpoint).
		Param("access_token", FB_PAGE_ACCESS_TOKEN).
		ContentType("application/json").
		BodyJSON(&req).
		Client(cl).
		Fetch(context.Background())

	return err
}
