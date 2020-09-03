package app

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

type kinesisExporter struct {
	einterfaces.GenericSocketExporter
	streamName      *string
	region          *string
	partitionKey    *string
	whitelist       map[string]bool
	kinesisInstance *kinesis.Kinesis
}

func (k *kinesisExporter) InitExporter(c *model.Config) {
	k.region = c.SocketExporterSettings.Region
	k.partitionKey = c.SocketExporterSettings.PartitionKey
	sess := session.New(&aws.Config{Region: aws.String(*k.region)})
	k.kinesisInstance = kinesis.New(sess)
	k.streamName = c.SocketExporterSettings.StreamName
	k.whitelist = *c.SocketExporterSettings.WhitelistedEvents
}

func (k *kinesisExporter) Export(c *model.Config, message *model.WebSocketEvent) bool {
	if k.kinesisInstance == nil {
		k.InitExporter(c)
	}
	if k.whitelist[message.EventType()] == true {
		_, err := k.kinesisInstance.PutRecord(&kinesis.PutRecordInput{
			Data:         []byte(message.ToJson()),
			StreamName:   k.streamName,
			PartitionKey: k.partitionKey,
		})
		if err != nil {
			fmt.Println(err)
			return false
		}
		return true
	}
	return false
}
