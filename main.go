package main

import (
	"github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
	"github.com/Financial-Times/simple-kinesis-consumer/kinesis"
	"os"
	"time"
)

func main() {
	app := cli.App("simple-kinesis-consumer", "Basic consumer for AWS Kinesis stream")

	kinesisStreamName := app.String(cli.StringOpt{
		Name:   "kinesisStreamName",
		Desc:   "Kinesis stream to subscribe to",
		EnvVar: "KINESIS_STREAM_NAME",
	})
	kinesisRegion := app.String(cli.StringOpt{
		Name:   "kinesisRegion",
		Value:  "eu-west-1",
		Desc:   "AWS region the kinesis stream is located",
		EnvVar: "KINESIS_REGION",
	})
	logLevel := app.String(cli.StringOpt{
		Name:   "logLevel",
		Value:  "info",
		Desc:   "App log level",
		EnvVar: "LOG_LEVEL",
	})

	app.Action = func() {

		log.SetFormatter(&log.JSONFormatter{})
		lvl, err := log.ParseLevel(*logLevel)
		if err != nil {
			log.WithField("LOG_LEVEL", *logLevel).Warn("Cannot parse log level, setting it to INFO.")
			lvl = log.InfoLevel
		}
		log.SetLevel(lvl)

		log.WithFields(log.Fields{
			"LOG_LEVEL":             *logLevel,
			"KINESIS_STREAM_NAME":   *kinesisStreamName,
			"KINESIS_REGION":	 *kinesisRegion,
		}).Info("Starting app with arguments")

		kinesisClient, err := kinesis.NewClient(*kinesisStreamName, *kinesisRegion)
		if err != nil {
			log.WithError(err).Fatal("Error creating Kinesis client")
		}

		var shardIterator string = ""
		for {
			output, err := kinesisClient.GetRecordsFromStream(shardIterator)
			if err != nil {
				log.WithError(err).Error("Error consuming from stream!")
				break
			}
			if len(output.Records) < 1 {
				log.Info("No records consumed\n")
			}
			for _, record := range output.Records {
				log.Infof("Consumed %s from kinesis stream", record.Data)
			}

			shardIterator = *output.NextShardIterator
			time.Sleep(15000 * time.Millisecond)
		}
	}
	app.Run(os.Args)
}
