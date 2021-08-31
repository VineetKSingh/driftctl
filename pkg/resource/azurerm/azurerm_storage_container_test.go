package azurerm_test

import (
	"testing"

	"github.com/cloudskiff/driftctl/test"
	"github.com/cloudskiff/driftctl/test/acceptance"
)

func TestAcc_Azure_StorageContainer(t *testing.T) {
	acceptance.Run(t, acceptance.AccTestCase{
		TerraformVersion: "0.14.4",
		Paths:            []string{"./testdata/acc/azurerm_storage_container"},
		Args: []string{
			"scan",
			"--to", "azure+tf",
			"--filter", "Type=='azurerm_storage_container' || Type=='azurerm_storage_account'",
		},
		Checks: []acceptance.AccCheck{
			{
				Check: func(result *test.ScanResult, stdout string, err error) {
					if err != nil {
						t.Fatal(err)
					}
					result.AssertInfrastructureIsInSync()
					// We should have 1 storage account resource and 3 storage container
					result.AssertManagedCount(4)
				},
			},
		},
	})
}
