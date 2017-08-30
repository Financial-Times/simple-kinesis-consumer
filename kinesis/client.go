package kinesis

import (
	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	log "github.com/sirupsen/logrus"
)

type Client interface {
	GetRecordsFromStream(startReadingFrom string) (kinesis.GetRecordsOutput, error)
	AddRecordToStream(record string, conceptType string) error
	Healthcheck() fthealth.Check
}

type KinesisClient struct {
	streamName string
	svc 	   *kinesis.Kinesis
}

func NewClient(streamName string, region string) (Client, error) {
	sess := session.Must(session.NewSession())
	svc := kinesis.New(sess, &aws.Config{
		Region: aws.String(region),
	})

	return &KinesisClient{
		streamName: streamName,
		svc: svc,
	}, nil
}

func (c *KinesisClient) GetRecordsFromStream(shardIterator string) (kinesis.GetRecordsOutput, error) {
	var si string
	if shardIterator == ""  {
		describeStreamInput := &kinesis.DescribeStreamInput{
			StreamName: aws.String(c.streamName),
		}
		describeStreamOutput, err := c.svc.DescribeStream(describeStreamInput)
		if err != nil {
			log.WithError(err).Error("There was an error returning details of stream")
			return kinesis.GetRecordsOutput{}, err
		}

		input := &kinesis.GetShardIteratorInput{
			ShardId: describeStreamOutput.StreamDescription.Shards[0].ShardId,
			ShardIteratorType: aws.String("LATEST"),
			StreamName: aws.String(c.streamName),
		}

		shardI, _ := c.svc.GetShardIterator(input)
		si = *shardI.ShardIterator
	} else {
		si = shardIterator
	}

	getRecordInput := &kinesis.GetRecordsInput{
		Limit: aws.Int64(10),
		ShardIterator: aws.String(si),
	}
	output, err := c.svc.GetRecords(getRecordInput)
	if err != nil {
		return kinesis.GetRecordsOutput{}, err
	}
	return *output, nil
}
