package app

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type sqsExporter struct {
	einterfaces.GenericSocketExporter
	streamName  *string
	region      *string
	queueUrl    *string
	whitelist   map[string]bool
	sqsInstance *sqs.SQS
	session     *session.Session
}

func (k *sqsExporter) InitExporter(c *model.Config) {
	k.region = c.SocketExporterSettings.Region
	k.queueUrl = c.SocketExporterSettings.QueueUrl
	k.session, _ = session.NewSession(&aws.Config{Region: aws.String(*k.region)})
	k.sqsInstance = sqs.New(k.session)
	k.whitelist = *c.SocketExporterSettings.WhitelistedEvents
}

func (k *sqsExporter) Export(c *model.Config, message *model.WebSocketEvent) bool {
	if k.sqsInstance == nil || k.session == nil {
		k.InitExporter(c)
	}
	if k.whitelist[message.EventType()] == true {
		message := message.ToJson()
		_, err := k.sqsInstance.SendMessage(&sqs.SendMessageInput{
			MessageBody: &message,
			QueueUrl:    k.queueUrl,
		})
		if err != nil {
			fmt.Println(err)
			return false
		}
		fmt.Println("exported")
		return true
	}
	return false
}
