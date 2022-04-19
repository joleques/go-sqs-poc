package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
)

func Receive() (string, error) {

	svc, err := configSQS()
	if err != nil {
		return "", err
	}
	channel := make(chan *sqs.Message)
	go receiveMessage(channel, *svc)
	for message := range channel {
		processMessage(message)
		deleteMessage(message, *svc)
	}

	return "Processo finalisado.....", nil
}

func configSQS() (*sqs.SQS, error) {
	creds := credentials.NewEnvCredentials()
	_, err := creds.Get()
	if err != nil {
		return nil, err
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
	return svc, nil
}

func receiveMessage(channel chan<- *sqs.Message, svc sqs.SQS) {
	qURL := os.Getenv("AWS_URL_QUEUE")
	for {
		result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            &qURL,
			MaxNumberOfMessages: aws.Int64(5),
			VisibilityTimeout:   aws.Int64(60), // 60 seconds
			WaitTimeSeconds:     aws.Int64(0),
		})

		if err != nil {
			log.Println("failed to fetch sqs message %v", err)
		}

		for _, message := range result.Messages {
			log.Println("Mensagem Recebida..:", *message.MessageId, "Body:", *message.Body)
			channel <- message
		}
	}
}

func processMessage(message *sqs.Message) {
	log.Println("Mensagem Processada:", *message.MessageId, "Body:", *message.Body)
}

func deleteMessage(message *sqs.Message, svc sqs.SQS) {
	qURL := os.Getenv("AWS_URL_QUEUE")
	messageHandle := message.ReceiptHandle
	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &qURL,
		ReceiptHandle: messageHandle,
	})
	if err != nil {
		log.Fatalf("failed delete sqs message %v", err)
	}
	log.Println("Mensagem Deletada..:", *message.MessageId, "Body:", *message.Body)
}
