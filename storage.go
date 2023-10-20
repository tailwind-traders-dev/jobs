//go:build mage
// +build mage

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/magefile/mage/mg"
)

type Azure mg.Namespace

func (Azure) ServiceBus() error {
	queueName := os.Getenv("AZURE_SERVICEBUS_QUEUE_NAME")
	if queueName == "" {
		return errors.New("AZURE_SERVICEBUS_QUEUE_NAME environment variable not found")
	}

	t := Test{}
	client, err := t.GetServiceBusClientDefault()
	if err != nil {
		return err
	}

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	if err != nil {
		return err
	}

	ctx := context.Background()
	messages, err := receiver.PeekMessages(ctx, 5, nil)
	if err != nil {
		return err
	}
	for _, x := range messages {
		fmt.Printf("%s\n", x.Body)
	}

	return nil
}

type Test struct{}

func (Test) GetServiceBusClientDefault() (*azservicebus.Client, error) {
	namespace, ok := os.LookupEnv("AZURE_SERVICEBUS_HOSTNAME") //ex: myservicebus.servicebus.windows.net
	if !ok {
		return nil, errors.New("AZURE_SERVICEBUS_HOSTNAME environment variable not found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		return nil, err
	}

	client, err := azservicebus.NewClient(namespace, cred, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
