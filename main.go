package main

import (
	"flag"
	"log"
	"q/cmd/smoothmq"
	"q/cmd/tester"
	"q/models"
	"q/queue/sqlite"

  "github.com/spf13/viper"
)

type DefaultTenantManager struct{}

func (tm *DefaultTenantManager) GetTenant() int64 {
	return 1
}

func (tm *DefaultTenantManager) GetAWSSecretKey(accessKey string, region string) (int64, string, error) {
	return int64(1), "YOUR_SECRET_ACCESS_KEY", nil
}

func NewDefaultTenantManager() models.TenantManager {
	return &DefaultTenantManager{}
}

func Run(tm models.TenantManager, queue models.Queue) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var runTester bool
	var numSenders, numReceivers, numMessagesPerGoroutine int
	var endpoint string

	flag.BoolVar(&runTester, "tester", false, "Run in test mode")

	flag.IntVar(&numSenders, "senders", 0, "Number of send goroutines")
	flag.IntVar(&numMessagesPerGoroutine, "messages", 1, "Number of messages to send per goroutine")
	flag.IntVar(&numReceivers, "receivers", 0, "Number of receive goroutines")
	flag.StringVar(&endpoint, "endpoint", "http://localhost:4001", "SQS endpoint for testing")

	flag.Parse()

	if runTester {
		tester.Run(numSenders, numReceivers, numMessagesPerGoroutine, endpoint)
	} else {
		smoothmq.Run(tm, queue)
	}
}

func init() {
  viper.AutomaticEnv()

  viper.BindEnv("ui-port", "UI_PORT")
  viper.BindEnv("sqs-port", "SQS_PORT")

  viper.SetDefault("ui-port", "4000")
  viper.SetDefault("sqs-port", "4001")
}

func main() {
	tenantManager := NewDefaultTenantManager()
	queue := sqlite.NewSQLiteQueue()

	Run(tenantManager, queue)
}
