package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/carlmjohnson/requests"
	"github.com/ybbus/httpretry"
)

const (
	FB_API_ENDPOINT      = "https://graph.facebook.com"
	FB_PAGE_ACCESS_TOKEN = "EAAato4HK0VABO7tZAbDv8lC5EM5AQzplnOr2tWslZCNwVg6j7ZBz6wPaEE2GsJQZAkDxHVKkGutWLL8v1ZCdqYUOZBNh19IMm1Nc8njPIDBmrW8JCOZCg14HlCIhBtrnbZAJQt6P8VrYFTdCsTAt0iSVCYZAvS6CTfLuZBZACA1dLRWsq4LfKZBMExOn3HxJ7sQVgV8VaZCktP04SmBEZD"
	FB_PAGE_ID           = "122103960950002350"
	LATEST_API_VERSION   = "17.0"
	API_SEPARATOR        = "/"
	MESSAGES_API_NAME    = "messages"
	baseEndpoint         = FB_API_ENDPOINT + API_SEPARATOR + LATEST_API_VERSION + API_SEPARATOR + FB_PAGE_ID + API_SEPARATOR + MESSAGES_API_NAME
)

func sendMessage(message string, customerId string, messageType string) error {
	customerIdMap := make(map[string]string)
	customerIdMap["id"] = customerId
	var customerIdJSON, _customerIdErr = json.Marshal(&customerIdMap)

	messageMap := make(map[string]string)
	messageMap["text"] = message
	var messageMapJSON, _messageErr = json.Marshal(&messageMap)

	if _customerIdErr == nil || _messageErr == nil {
		return errors.New("JSON serialization error")
	}

	cl := httpretry.NewDefaultClient() //Used for retries
	var err = requests.
		URL(baseEndpoint).
		Param("access_token", FB_PAGE_ACCESS_TOKEN).
		Param("recipient", string(customerIdJSON)).
		Param("message", string(messageMapJSON)).
		Param("message_type", messageType).
		Client(cl).
		Fetch(context.Background())

	return err
}
