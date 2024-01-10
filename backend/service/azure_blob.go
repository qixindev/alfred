package service

import (
	"alfred/backend/pkg/global"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"io"
	"time"
)

const (
	containerName = "picture"
)

func getCredential() (*azblob.SharedKeyCredential, *azblob.Client, error) {
	accountName := global.CONFIG.AzureBlob.AccountName
	accountKey := global.CONFIG.AzureBlob.AccountKey
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, nil, err
	}

	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     time.Now().UTC(),
		ExpiryTime:    time.Now().UTC().Add(time.Hour * 24 * 30 * 12 * 200), // 200年
		Permissions:   to.Ptr(sas.BlobPermissions{Read: true, Create: true, Write: true, Tag: true}).String(),
		ContainerName: containerName,
	}.SignWithSharedKey(credential)
	if err != nil {
		return nil, nil, err
	}
	sasURL := fmt.Sprintf("https://%s.blob.core.chinacloudapi.cn/?%s", accountName, sasQueryParams.Encode())
	upClient, err := azblob.NewClientWithNoCredential(sasURL, nil)
	if err != nil {
		return nil, nil, err
	}

	return credential, upClient, nil
}

func UploadFileToAzureBlob(fileName string, src io.Reader) (string, error) {
	accountName := global.CONFIG.AzureBlob.AccountName
	credential, upClient, err := getCredential()
	if err != nil {
		return "", err
	}

	if _, err = upClient.UploadStream(context.Background(), containerName, fileName, src, nil); err != nil {
		return "", err
	}

	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     time.Now().UTC().Add(time.Second * -10),
		ExpiryTime:    time.Now().UTC().Add(time.Hour * 24 * 30 * 12 * 200), // 200年
		Permissions:   to.Ptr(sas.ContainerPermissions{Read: true, List: true}).String(),
		ContainerName: containerName,
		BlobName:      fileName,
		ContentType:   "image/png",
	}.SignWithSharedKey(credential)
	if err != nil {
		return "", err
	}
	sasURL := fmt.Sprintf("https://%s.blob.core.chinacloudapi.cn/%s/%s?%s", accountName, containerName, fileName, sasQueryParams.Encode())

	return sasURL, nil
}
