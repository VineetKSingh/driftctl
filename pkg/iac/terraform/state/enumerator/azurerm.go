package enumerator

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/pkg/errors"
	"github.com/snyk/driftctl/pkg/helpers/azure"
	"github.com/snyk/driftctl/pkg/iac/config"
)

type AzureRMEnumerator struct {
	containerName, objectPath string
	containerClient           azblob.ContainerClient
	origin                    string
}

func NewAzureRMEnumerator(config config.SupplierConfig) (*AzureRMEnumerator, error) {
	splitPath := strings.Split(config.Path, "/")
	if len(splitPath) < 2 || splitPath[1] == "" {
		return nil, errors.Errorf("Unable to parse azurerm backend storage splitPath: %s. Must be CONTAINER/PATH/TO/OBJECT", config.Path)
	}
	containerName := splitPath[0]
	objectPath := strings.Join(splitPath[1:], "/")

	credential, err := azure.GetBlobSharedKey()
	if err != nil {
		return nil, err
	}
	container, err := azblob.NewContainerClientWithSharedKey(
		fmt.Sprintf(
			"https://%s.blob.core.windows.net/%s",
			credential.AccountName(),
			containerName,
		),
		credential,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &AzureRMEnumerator{
		containerName:   containerName,
		objectPath:      objectPath,
		containerClient: container,
		origin:          config.String(),
	}, nil
}

func (s *AzureRMEnumerator) Origin() string {
	return s.origin
}

func (s *AzureRMEnumerator) Enumerate() ([]string, error) {

	// prefix should contain everything that does not have a glob pattern should be the glob matcher string
	prefix, pattern := GlobS3(s.objectPath)
	fullPattern := strings.Join([]string{prefix, pattern}, "/")
	fullPattern = strings.Trim(fullPattern, "/")

	pager := s.containerClient.ListBlobsFlat(&azblob.ContainerListBlobFlatSegmentOptions{
		Prefix: &prefix,
	})

	files := make([]string, 0)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		for _, v := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			if *v.Properties.ContentLength == 0 {
				continue
			}
			if match, _ := doublestar.Match(fullPattern, *v.Name); match {
				files = append(files, strings.Join([]string{s.containerName, *v.Name}, "/"))
			}
		}
	}

	if err := pager.Err(); err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.Errorf("no Terraform state was found for %s, exiting", s.origin)
	}

	return files, nil
}
