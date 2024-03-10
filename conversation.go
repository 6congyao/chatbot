package main

import (
	"errors"
	"fmt"
	"strings"
)

var templatesMap = map[string]string{
	"thank":                 "Dear %s, Your unwavering support and trust in our products/services mean the world to us. We are truly grateful for the opportunity to serve you and for the strong partnership we have built.",
	"transaction completed": "We sincerely appreciate your business and the trust you have placed in us. If you have any further questions or need assistance in the future, please don't hesitate to reach out. We value your satisfaction and look forward to serving you again.",
}

func parseEvent(event FacebookEvent) (string, error) {
	if event.Field == "messages" {
		//Make some call to ChatGPT asking if the comment posted on the
		//Page suggests user wants to give a review
		// callChatGPTToCheckIfMessageSuggestReview(event.Message) returns boolean
		if true { //Swap this out based on ChatGPT response
			message, _err := generateResponseMessage(event.Message, event.CustomerId)
			if _err == nil {
				return message, sendMessage(message, event.CustomerId, "RESPONSE")
			}
		}
	}
	return "", nil
}

func generateResponseMessage(message string, customerId string) (string, error) {
	if strings.Contains(message, "thank") {
		return fmt.Sprintf(templatesMap["thank"], customerId), nil
	}
	if strings.Contains(message, "transaction completed") {
		return templatesMap["transaction completed"], nil
	}
	return "", errors.New("no need to reply")
}
