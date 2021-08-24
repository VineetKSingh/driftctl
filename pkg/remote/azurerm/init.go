package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/cloudskiff/driftctl/pkg/alerter"
	"github.com/cloudskiff/driftctl/pkg/output"
	"github.com/cloudskiff/driftctl/pkg/remote/azurerm/repository"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	"github.com/cloudskiff/driftctl/pkg/remote/common"
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/resource/azurerm"
	"github.com/cloudskiff/driftctl/pkg/terraform"
)

func Init(
	// Version required by the user
	version string,
	// Util to send alert
	alerter *alerter.Alerter,
	// Library that contains all providers
	providerLibrary *terraform.ProviderLibrary,
	// Library that contains the enumerators and details fetcher for each supported resources
	remoteLibrary *common.RemoteLibrary,
	// progress display
	progress output.Progress,
	// Repository for all resource schema
	resourceSchemaRepository *resource.SchemaRepository,
	// Factory used to create driftctl resource
	factory resource.ResourceFactory,
	// Drifctl config directory (in which terraform provider is downloaded)
	configDir string) error {

	// Define the default version of terraform provider to be used. When the user does not require a specific one
	if version == "" {
		version = "2.71.0"
	}

	// This is this actual terraform provider creation
	provider, err := NewAzureTerraformProvider(version, progress, configDir)
	if err != nil {
		return err
	}
	// And then initialisation
	err = provider.Init()
	if err != nil {
		return err
	}

	providerConfig := provider.GetConfig()
	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{})
	if err != nil {
		return err
	}
	con := armcore.NewDefaultConnection(cred, nil)

	// You'll need to create a new cache that will be use to cache fetched resources lists
	c := cache.New(100)

	// Deserializer is used to convert cty value return by terraform provider to driftctl Resource
	_ = resource.NewDeserializer(factory)

	storageAccountRepo := repository.NewStorageRepository(con, providerConfig, c)

	// Adding the provider to the library
	providerLibrary.AddProvider(terraform.AZURE, provider)

	remoteLibrary.AddEnumerator(NewAzurermStorageAccountEnumerator(storageAccountRepo, factory))
	remoteLibrary.AddEnumerator(NewAzurermStorageContainerEnumerator(storageAccountRepo, factory))

	err = resourceSchemaRepository.Init(terraform.AZURE, version, provider.Schema())
	if err != nil {
		return err
	}
	azurerm.InitResourcesMetadata(resourceSchemaRepository)

	return nil
}
