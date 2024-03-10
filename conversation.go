package main

import (
	"errors"
	"fmt"
	"strings"
)

var template = "Dear %s, Your unwavering support and trust in our products/services mean the world to us. We are truly grateful for the opportunity to serve you and for the strong partnership we have built."

func parseEvent(event FacebookEvent) error {
	if event.Field == "messages" {
		//Make some call to ChatGPT asking if the comment posted on the
		//Page suggests user wants to give a review
		// callChatGPTToCheckIfMessageSuggestReview(event.Message) returns boolean
		if true { //Swap this out based on ChatGPT response
			message, _err := generateResponseMessage(event.Message, event.CustomerId)
			if _err == nil {
				return sendMessage(message, event.CustomerId, "RESPONSE")
			}
		}
	} else {
		// review := event.Message
		// DB Config Object- sqlCfg = sql.config(...)
		// db, err = sql.Open("mysql", cfg.FormatDSN())
		return nil
	}
	return nil
}

func generateResponseMessage(message string, customerId string) (string, error) {
	if strings.Contains(message, "thank") {
		return fmt.Sprintf(template, customerId), nil
	}
	return "", errors.New("no need to reply")
}
