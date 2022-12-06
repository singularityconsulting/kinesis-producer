package producer

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

func ExampleSimple() {
	logger := &StdLogger{log.New(os.Stdout, "", log.LstdFlags)}
	client := kinesis.New(session.New(aws.NewConfig()))
	pr := New(&Config{
		StreamName:   "test",
		BacklogCount: 2000,
		Client:       client,
		Logger:       logger,
	})

	pr.Start()

	// Handle failures
	go func() {
		for r := range pr.NotifyFailures() {
			// r contains `Data`, `PartitionKey` and `Error()`
			logger.Error("detected put failure", r.Error)
		}
	}()

	go func() {
		for i := 0; i < 5000; i++ {
			err := pr.Put([]byte("foo"), "bar")
			if err != nil {
				logger.Error("error producing", err)
			}
		}
	}()

	time.Sleep(3 * time.Second)
	pr.Stop()
}
