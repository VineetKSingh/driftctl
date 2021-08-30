package aws_test

import (
	"testing"

	"github.com/cloudskiff/driftctl/test"
	"github.com/cloudskiff/driftctl/test/acceptance"
)

func TestAcc_AwsRoute53Record_WithFQDNAsId(t *testing.T) {
	acceptance.Run(t, acceptance.AccTestCase{
		TerraformVersion: "0.15.5",
		Paths:            []string{"./testdata/acc/aws_route53_record"},
		Args:             []string{"scan", "--filter", "Type=='aws_route53_record'", "--deep"},
		Checks: []acceptance.AccCheck{
			{
				Env: map[string]string{
					"AWS_REGION": "us-east-1",
				},
				Check: func(result *test.ScanResult, stdout string, err error) {
					if err != nil {
						t.Fatal(err)
					}
					result.AssertDriftCountTotal(0)
					result.Equal(0, result.Summary().TotalDeleted)
					result.Equal(8, result.Summary().TotalManaged)
				},
			},
		},
	})
}

func TestAcc_AwsRoute53Record_WithAlias(t *testing.T) {
	acceptance.Run(t, acceptance.AccTestCase{
		TerraformVersion: "0.15.5",
		Paths:            []string{"./testdata/acc/aws_route53_record_with_alias"},
		Args:             []string{"scan", "--filter", "Type=='aws_route53_record'", "--deep"},
		Checks: []acceptance.AccCheck{
			{
				Check: func(result *test.ScanResult, stdout string, err error) {
					if err != nil {
						t.Fatal(err)
					}
					result.AssertDriftCountTotal(0)
					result.Equal(2, result.Summary().TotalManaged)
				},
			},
		},
	})
}
