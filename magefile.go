//go:build mage
// +build mage

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Hello is our default mage target which we also call
// by default within our Docker container
func Hello() error {
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("hello %s\n", now)
	return nil
}

// Goodbye is an alternative mage target we can call
func Goodbye() error {
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("goodbye %s\n", now)
	return nil
}

type ServiceBus mg.Namespace

// Send sends a single message to the service bus queue
func (ServiceBus) Send(message string) error {
	queueOrTopic := os.Getenv("AZURE_SERVICEBUS_QUEUE_NAME")
	if queueOrTopic == "" {
		return errors.New("AZURE_SERVICEBUS_QUEUE_NAME environment variable not found")
	}

	ctx := context.Background()
	client, err := getServiceBusClient()
	if err != nil {
		return err
	}
	defer client.Close(ctx)

	sender, err := client.NewSender(queueOrTopic, nil)
	if err != nil {
		return nil
	}
	defer sender.Close(ctx)

	sbMessage := &azservicebus.Message{
		Body: []byte(message),
	}
	err = sender.SendMessage(ctx, sbMessage, nil)
	if err != nil {
		return err
	}

	return nil
}

// SendMessageBatch sends a batch of 10 messages to the
// service bus queue with a 1 second delay between each
// message
func (ServiceBus) SendMessageBatch() error {
	queueOrTopic := os.Getenv("AZURE_SERVICEBUS_QUEUE_NAME")
	if queueOrTopic == "" {
		return errors.New("AZURE_SERVICEBUS_QUEUE_NAME environment variable not found")
	}

	ctx := context.Background()
	client, err := getServiceBusClient()
	if err != nil {
		return err
	}
	defer client.Close(ctx)

	sender, err := client.NewSender(queueOrTopic, nil)
	if err != nil {
		return err
	}
	defer sender.Close(ctx)

	batch, err := sender.NewMessageBatch(ctx, nil)
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		now := time.Now().UTC().Format(time.RFC3339)
		message := azservicebus.Message{
			Body: []byte(fmt.Sprintf("hello %s", now)),
		}
		err := batch.AddMessage(&message, nil)
		if err != nil {
			return err
		}
		fmt.Printf("Added message: %s\n", message.Body)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("Sending messages\n")
	return sender.SendMessageBatch(ctx, batch, nil)
}

// Receive receives and completes a batch of 5 messages
// from the service bus queue with a 1 second delay between
// each message. This will receive messages indefinitely.
func (ServiceBus) Receive() error {
	queueOrTopic := os.Getenv("AZURE_SERVICEBUS_QUEUE_NAME")
	if queueOrTopic == "" {
		return errors.New("AZURE_SERVICEBUS_QUEUE_NAME environment variable not found")
	}

	ctx := context.Background()
	client, err := getServiceBusClient()
	if err != nil {
		return err
	}
	defer client.Close(ctx)

	receiver, err := client.NewReceiverForQueue(queueOrTopic, nil)
	if err != nil {
		return err
	}
	defer receiver.Close(ctx)

	count := 5
	for {
		messages, err := receiver.ReceiveMessages(ctx, count, nil)
		if err != nil {
			return err
		}

		for _, message := range messages {
			fmt.Printf("Received message: %s\n", message.Body)
			err = receiver.CompleteMessage(ctx, message, nil)
			if err != nil {
				return err
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// ReceiveAll receives and completes all messages from the
// Service Bus queue, in batches of 5, and exits when complete
func (ServiceBus) ReceiveAll() error {
	queueOrTopic := os.Getenv("AZURE_SERVICEBUS_QUEUE_NAME")
	if queueOrTopic == "" {
		return errors.New("AZURE_SERVICEBUS_QUEUE_NAME environment variable not found")
	}

	ctx := context.Background()
	client, err := getServiceBusClient()
	if err != nil {
		return err
	}
	defer client.Close(ctx)

	receiver, err := client.NewReceiverForQueue(queueOrTopic, nil)
	if err != nil {
		return err
	}
	defer receiver.Close(ctx)

	count := 5
	for {
		timeout, _ := context.WithTimeout(ctx, 5*time.Second)
		messages, err := receiver.ReceiveMessages(timeout, count, nil)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Printf("No messages in 5s, exiting\n")
				return nil
			}
			return err
		}

		if len(messages) == 0 {
			return nil
		}

		for _, message := range messages {
			fmt.Printf("Received message: %s\n", message.Body)
			err = receiver.CompleteMessage(ctx, message, nil)
			if err != nil {
				return err
			}
		}
	}
}

func getServiceBusClient() (*azservicebus.Client, error) {
	// todo: merge with getServiceBusClientDefault and use
	// AZURE_SERVICEBUS_CONNECTION_STRING if exists, otherwise
	// use azidentity.NewDefaultAzureCredential by default

	connectionString, ok := os.LookupEnv("AZURE_SERVICEBUS_CONNECTION_STRING") //ex: myservicebus.servicebus.windows.net
	if !ok {
		return nil, errors.New("AZURE_SERVICEBUS_CONNECTION_STRING environment variable not found")
	}

	client, err := azservicebus.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getServiceBusClientDefault() (*azservicebus.Client, error) {
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

type Docker mg.Namespace

// BuildDev builds the container image, "jobs", with --no-cache
// and Dockerfile which builds a static binary and
// multi-stage builds to utilize a distroless image
func (Docker) Build() error {
	cmd1 := []string{
		"docker",
		"build",
		"--no-cache",
		"-t",
		"jobs",
		".",
	}
	return sh.RunV(cmd1[0], cmd1[1:]...)
}

// BuildDev builds the container image, "jobs", with --no-cache
// and dev.Dockerfile which uses the golang:latest image,
// installs mage and vim, for more interactive development
func (Docker) BuildDev() error {
	cmd1 := []string{
		"docker",
		"build",
		"--no-cache",
		"-f",
		"dev.Dockerfile",
		"-t",
		"jobs",
		".",
	}
	return sh.RunV(cmd1[0], cmd1[1:]...)
}

// Run runs the jobs container with the mage target
func (Docker) Run(target string) error {
	cmd1 := []string{
		"docker",
		"run",
		"-it",
		"jobs",
		"mage",
		target,
	}
	return sh.RunV(cmd1[0], cmd1[1:]...)
}

type Deploy mg.Namespace

// ContainerApps deploys the Container App(s) via containerapp.bicep
// into the provided <resource group>
// Requires: AZURE_SERVICEBUS_CONNECTION_STRING
func (Deploy) ContainerApps(resourceGroup string) error {
	serviceBusConnection := os.Getenv("AZURE_SERVICEBUS_CONNECTION_STRING")
	if serviceBusConnection == "" {
		return errors.New("AZURE_SERVICEBUS_CONNECTION_STRING environment variable not found")
	}
	cmd1 := []string{
		"az",
		"deployment",
		"group",
		"create",
		"--resource-group",
		resourceGroup,
		"--template-file",
		"deploy/azure-container-apps/containerapp.bicep",
		"--parameters",
		"service_bus_connection=" + serviceBusConnection,
	}
	return sh.RunV(cmd1[0], cmd1[1:]...)
}

// Empty empties the <resource group> via empty.bicep
func (Deploy) Empty(resourceGroup string) error {
	cmd1 := []string{
		"az",
		"deployment",
		"group",
		"create",
		"--resource-group",
		resourceGroup,
		"--mode",
		"Complete",
		"--template-file",
		"deploy/azure-container-apps/empty.bicep",
	}
	return sh.RunV(cmd1[0], cmd1[1:]...)
}

// Group creates the <resource group> in <location>
func (Deploy) Group(name, location string) error {
	cmd1 := []string{
		"az",
		"group",
		"create",
		"--name",
		name,
		"--location",
		location,
	}
	return sh.RunV(cmd1[0], cmd1[1:]...)
}
