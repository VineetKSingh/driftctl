package aws_test

import (
	"testing"
	"time"

	"github.com/cloudskiff/driftctl/test"
	"github.com/cloudskiff/driftctl/test/acceptance"
)

func TestAcc_Aws_EipAssociation(t *testing.T) {
	acceptance.Run(t, acceptance.AccTestCase{
		TerraformVersion: "0.15.5",
		Paths:            []string{"./testdata/acc/aws_eip_association"},
		Args:             []string{"scan", "--filter", "Type=='aws_eip' || Type=='aws_eip_association'", "--tf-provider-version", "3.44.0", "--deep"},
		RetryDestroy: acceptance.RetryConfig{
			Attempts: 3,
			Delay:    5 * time.Second,
		},
		Checks: []acceptance.AccCheck{
			{
				Env: map[string]string{
					"AWS_REGION": "us-east-1",
				},
				Check: func(result *test.ScanResult, stdout string, err error) {
					if err != nil {
						t.Fatal(err)
					}
					result.AssertInfrastructureIsInSync()
					result.AssertManagedCount(2)
				},
			},
		},
	})
}
