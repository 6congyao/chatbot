package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Ddb struct {
}

func (d *Ddb) store(event FacebookEvent, replyMsg string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
		// log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"id":        &types.AttributeValueMemberS{Value: event.Id},
			"message":   &types.AttributeValueMemberS{Value: event.Message},
			"reply":     &types.AttributeValueMemberS{Value: replyMsg},
			"recipient": &types.AttributeValueMemberS{Value: event.CustomerId},
		},
		TableName: aws.String("chat-history"),
	}
	_, err = svc.PutItem(context.TODO(), input)
	if err != nil {
		return err
	}

	log.Printf("Successfully stored data with ID: %s", event.Id)

	return nil
}
