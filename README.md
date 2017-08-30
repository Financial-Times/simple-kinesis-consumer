# Simple Kinesis Consumer (simple-kinesis-consumer)

__Basic consumer for AWS Kinesis stream to be solely used as a TESTING tool in current state__

# Installation

For the first time:

`go get github.com/Financial-Times/simple-kinesis-consumer`

or update:

`go get -u github.com/Financial-Times/simple-kinesis-consumer`

to build:

`go build simple-kinesis-consumer`

to run:

`./simple-kinesis-consumer --kinesisStreamName={{stream_name}}`


The app assumes that you have correctly set up your AWS credentials by either using the `~/.aws/credentials` file:

```
[default]
aws_access_key_id = AKID1234567890
aws_ secret_access_key = MY-SECRET-KEY
```

or the default AWS environment variables otherwise requests will return 401 Unauthorised

```
AWS_ACCESS_KEY_ID=AKID1234567890
AWS_SECRET_ACCESS_KEY=MY-SECRET-KEY
```