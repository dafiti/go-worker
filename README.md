# Golang Worker

[![Build Status](https://img.shields.io/travis/dafiti/go-worker/master.svg?style=flat-square)](https://travis-ci.org/dafiti/go-worker)
[![Coverage Status](https://img.shields.io/coveralls/dafiti/go-worker/master.svg?style=flat-square)](https://coveralls.io/github/dafiti/go-worker?branch=master)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/dafiti/go-worker)

Go-Worker

## Installation
```sh
go get github.com/dafiti/go-worker
```

## Documentation

Read the full documentation at [https://godoc.org/github.com/dafiti/go-worker](https://godoc.org/github.com/dafiti/go-worker).

## Example

See an example of usage with the SQS as queue:

```go
package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type (
	SqsMessagesReceiver struct {
		SqsClient           sqsiface.SQSAPI
		QueueUrl            string
		MaxNumberOfMessages int64
		VisibilityTimeout   int64
		WaitTimeSeconds     int64
	}

	LogHandler struct{}
)

func (r *SqsMessagesReceiver) Receive() []Message {
	params := &sqs.ReceiveMessageInput{
		QueueUrl:            &r.QueueUrl,
		MaxNumberOfMessages: aws.Int64(r.MaxNumberOfMessages),
		VisibilityTimeout:   aws.Int64(r.VisibilityTimeout),
		WaitTimeSeconds:     aws.Int64(r.WaitTimeSeconds),
	}
	result, err := r.SqsClient.ReceiveMessage(params)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var messages []Message
	for _, m := range result.Messages {
		messages = append(messages, Message{Body: m.Body})
	}

	return messages
}

func (h *LogHandler) Handle(messages *[]Message) (bool, error) {
	for _, m := range *messages {
		log.Println(m.Body)
	}

	return true, nil
}

func main() {
	cred := credentials.AnonymousCredentials
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: cred,
	})
	if err != nil {
		panic(err.Error())
	}
	sqsClient := sqs.New(cred)

	receiver := &SqsMessagesReceiver{
		SqsClient:           sqsClient,
		QueueUrl:            "sqs-queue-url",
		MaxNumberOfMessages: 10,
		VisibilityTimeout:   20,
		WaitTimeSeconds:     10,
	}

	handler := &LogHandler{}
	worker := &Worker{MaxWorkers: 5}

	worker.Run(receiver, handler)
}
```

## License

This project is released under the MIT licence. See [LICENCE](https://github.com/dafiti/go-worker/blob/master/LICENSE) for more details.
