package aws

import (
	"context"
	"testing"

	"github.com/cloudskiff/driftctl/pkg/remote/aws/repository"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	remoteerror "github.com/cloudskiff/driftctl/pkg/remote/error"
	testresource "github.com/cloudskiff/driftctl/test/resource"

	resourceaws "github.com/cloudskiff/driftctl/pkg/resource/aws"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/cloudskiff/driftctl/pkg/parallel"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cloudskiff/driftctl/test/mocks"

	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/terraform"
	"github.com/cloudskiff/driftctl/test"
	"github.com/cloudskiff/driftctl/test/goldenfile"
)

func TestIamPolicySupplier_Resources(t *testing.T) {

	cases := []struct {
		test    string
		dirName string
		mocks   func(repo *repository.MockIAMRepository)
		err     error
	}{
		{
			test:    "no iam custom policies",
			dirName: "iam_policy_empty",
			mocks: func(repo *repository.MockIAMRepository) {
				repo.On("ListAllPolicies").Once().Return([]*iam.Policy{}, nil)
			},
			err: nil,
		},
		{
			test:    "iam multiples custom policies",
			dirName: "iam_policy_multiple",
			mocks: func(repo *repository.MockIAMRepository) {
				repo.On("ListAllPolicies").Once().Return([]*iam.Policy{
					{
						Arn: aws.String("arn:aws:iam::929327065333:policy/policy-0"),
					},
					{
						Arn: aws.String("arn:aws:iam::929327065333:policy/policy-1"),
					},
					{
						Arn: aws.String("arn:aws:iam::929327065333:policy/policy-2"),
					},
				}, nil)
			},
			err: nil,
		},
		{
			test:    "cannot list iam custom policies",
			dirName: "iam_policy_empty",
			mocks: func(repo *repository.MockIAMRepository) {
				repo.On("ListAllPolicies").Once().Return(nil, awserr.NewRequestFailure(nil, 403, ""))
			},
			err: remoteerror.NewResourceEnumerationError(awserr.NewRequestFailure(nil, 403, ""), resourceaws.AwsIamPolicyResourceType),
		},
	}
	for _, c := range cases {
		shouldUpdate := c.dirName == *goldenfile.Update

		providerLibrary := terraform.NewProviderLibrary()
		supplierLibrary := resource.NewSupplierLibrary()

		repo := testresource.InitFakeSchemaRepository("aws", "3.19.0")
		resourceaws.InitResourcesMetadata(repo)
		factory := terraform.NewTerraformResourceFactory(repo)

		deserializer := resource.NewDeserializer(factory)
		if shouldUpdate {
			provider, err := InitTestAwsProvider(providerLibrary)
			if err != nil {
				t.Fatal(err)
			}
			supplierLibrary.AddSupplier(NewIamPolicySupplier(provider, deserializer, repository.NewIAMRepository(provider.session, cache.New(0))))
		}

		t.Run(c.test, func(tt *testing.T) {
			fakeIam := repository.MockIAMRepository{}
			c.mocks(&fakeIam)

			provider := mocks.NewMockedGoldenTFProvider(c.dirName, providerLibrary.Provider(terraform.AWS), shouldUpdate)
			s := &IamPolicySupplier{
				provider,
				deserializer,
				&fakeIam,
				terraform.NewParallelResourceReader(parallel.NewParallelRunner(context.TODO(), 10)),
			}
			got, err := s.Resources()
			assert.Equal(tt, c.err, err)

			mock.AssertExpectationsForObjects(tt)
			test.CtyTestDiff(got, c.dirName, provider, deserializer, shouldUpdate, t)
		})
	}
}
