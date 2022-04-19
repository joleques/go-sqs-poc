package sqs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
)

type MessageSQS struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func Send(messageSQS MessageSQS) (string, error) {
	creds := credentials.NewEnvCredentials()
	_, err := creds.Get()
	if err != nil {
		return "", err
	}
	region := os.Getenv("AWS_REGION")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region:      aws.String(region),
			Credentials: creds,
		},
	}))

	svc := sqs.New(sess)

	qURL := os.Getenv("AWS_URL_QUEUE")

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"IdMessage": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(messageSQS.Id),
			},
		},
		MessageBody: aws.String(messageSQS.Message),
		QueueUrl:    &qURL,
	})

	if err != nil {
		return "", err
	}
	messageResult := fmt.Sprintf("message send successfully, id: ", *result.MessageId)
	log.Println("Mensagem Send Queue:", *result.MessageId)
	return messageResult, nil
}
