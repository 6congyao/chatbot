package main

func parseEvent(event FacebookEvent) error {
	if event.Field == "feed" {
		//Make some call to ChatGPT asking if the comment posted on the
		//Page suggests user wants to give a review
		// callChatGPTToCheckIfMessageSuggestReview(event.Message) returns boolean
		if true { //Swap this out based on ChatGPT response
			return sendMessage("Please feel free to provide a review about our Page!",
				event.CustomerId, "UPDATE")
		}
	} else {
		// review := event.Message
		// DB Config Object- sqlCfg = sql.config(...)
		// db, err = sql.Open("mysql", cfg.FormatDSN())
		return nil
	}
	return nil
}
