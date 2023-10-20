//go:build mage
// +build mage

package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/magefile/mage/mg"
)

type Azure mg.Namespace

func (Azure) Storage() error {
	storageAccountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	if storageAccountName == "" {
		return errors.New("AZURE_STORAGE_ACCOUNT_NAME environment variable not found")
	}
	storageContainerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")
	if storageContainerName == "" {
		return errors.New("AZURE_STORAGE_CONTAINER_NAME environment variable not found")
	}

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	url1 := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccountName)
	client, err := azblob.NewClient(url1, credential, nil)
	if err != nil {
		return err
	}

	pager := client.NewListBlobsFlatPager(storageContainerName,
		&azblob.ListBlobsFlatOptions{
			Include: azblob.ListBlobsInclude{
				Snapshots: true,
				Versions:  true,
			},
		})

	ctx := context.Background()
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, blob := range resp.Segment.BlobItems {
			url2 := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", storageAccountName, storageContainerName, *blob.Name)
			fmt.Println(url2)
		}
	}
	return nil
}

func (Azure) KeyVault() error {
	vaultURI := os.Getenv("AZURE_KEY_VAULT_URI")
	if vaultURI == "" {
		return errors.New("AZURE_KEY_VAULT_URI environment variable not found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	client, err := azsecrets.NewClient(vaultURI, cred, nil)
	if err != nil {
		return err
	}

	ctx := context.Background()
	pager := client.NewListSecretsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, secret := range page.Value {
			fmt.Printf("%s\n", *secret.ID)
		}
	}

	return nil
}

func (Azure) ServiceBus() error {
	queueName := os.Getenv("AZURE_SERVICEBUS_QUEUE_NAME")
	if queueName == "" {
		return errors.New("AZURE_SERVICEBUS_QUEUE_NAME environment variable not found")
	}

	namespace := os.Getenv("AZURE_SERVICEBUS_HOSTNAME")
	if namespace == "" {
		return errors.New("AZURE_SERVICEBUS_HOSTNAME environment variable not found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	client, err := azservicebus.NewClient(namespace, cred, nil)
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

func (Azure) Postgres() error {
	return errors.New("Not yet implemented")
}
