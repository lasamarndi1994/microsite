package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

/*
* Create new SNS client
* @return *sns.Client
 */
func NewSNSClient() *sns.Client {
	fmt.Println(os.Getenv("AWS_REGION"))
	cfg, err := config.LoadDefaultConfig(context.TODO(),

		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatal("AWS config error:", err)
	}

	return sns.NewFromConfig(cfg)
}

/*
* Send SMS via SNS
* @param snsClient *sns.Client
* @param phone string
* @param message string
* @return error
 */
func SendSMS(snsClient *sns.Client, phone string, message string) error {
	_, err := snsClient.Publish(context.TODO(), &sns.PublishInput{
		Message:     &message,
		PhoneNumber: &phone, // +91xxxxxxxxxx
	})

	if err != nil {
		return err
	}

	fmt.Println("Message sent to:", phone)
	return nil
}

/*
* Send message to SNS topic
* @param snsClient *sns.Client
* @param msg string
* @return error
 */
func SendTopicMessage(snsClient *sns.Client, msg string) error {
	topicArn := os.Getenv("SNS_TOPIC_ARN")

	_, err := snsClient.Publish(context.TODO(), &sns.PublishInput{
		Message:  &msg,
		TopicArn: &topicArn,
	})
	return err
}
