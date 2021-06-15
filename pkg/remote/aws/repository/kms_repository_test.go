package repository

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	awstest "github.com/cloudskiff/driftctl/test/aws"
	"github.com/stretchr/testify/mock"

	"github.com/r3labs/diff/v2"
	"github.com/stretchr/testify/assert"
)

func Test_KMSRepository_ListAllKeys(t *testing.T) {
	tests := []struct {
		name    string
		mocks   func(client *awstest.MockFakeKMS)
		want    []*kms.KeyListEntry
		wantErr error
	}{
		{
			name: "List only customer keys",
			mocks: func(client *awstest.MockFakeKMS) {
				client.On("ListKeysPages",
					&kms.ListKeysInput{},
					mock.MatchedBy(func(callback func(res *kms.ListKeysOutput, lastPage bool) bool) bool {
						callback(&kms.ListKeysOutput{
							Keys: []*kms.KeyListEntry{
								{KeyId: aws.String("1")},
								{KeyId: aws.String("2")},
								{KeyId: aws.String("3")},
							},
						}, true)
						return true
					})).Return(nil).Once()
				client.On("DescribeKey",
					&kms.DescribeKeyInput{
						KeyId: aws.String("1"),
					}).Return(&kms.DescribeKeyOutput{
					KeyMetadata: &kms.KeyMetadata{
						KeyId:      aws.String("1"),
						KeyManager: aws.String("CUSTOMER"),
					},
				}, nil).Once()
				client.On("DescribeKey",
					&kms.DescribeKeyInput{
						KeyId: aws.String("2"),
					}).Return(&kms.DescribeKeyOutput{
					KeyMetadata: &kms.KeyMetadata{
						KeyId:      aws.String("2"),
						KeyManager: aws.String("AWS"),
					},
				}, nil).Once()
				client.On("DescribeKey",
					&kms.DescribeKeyInput{
						KeyId: aws.String("3"),
					}).Return(&kms.DescribeKeyOutput{
					KeyMetadata: &kms.KeyMetadata{
						KeyId:      aws.String("3"),
						KeyManager: aws.String("AWS"),
					},
				}, nil).Once()
			},
			want: []*kms.KeyListEntry{
				{KeyId: aws.String("1")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(1)
			client := awstest.MockFakeKMS{}
			tt.mocks(&client)
			r := &kmsRepository{
				client: &client,
				cache:  store,
			}
			got, err := r.ListAllKeys()
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllKeys()
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				assert.IsType(t, []*kms.KeyListEntry{}, store.Get("kmsListAllKeys"))
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}

func Test_KMSRepository_ListAllAliases(t *testing.T) {
	tests := []struct {
		name    string
		mocks   func(client *awstest.MockFakeKMS)
		want    []*kms.AliasListEntry
		wantErr error
	}{
		{
			name: "List only customer aliases",
			mocks: func(client *awstest.MockFakeKMS) {
				client.On("ListAliasesPages",
					&kms.ListAliasesInput{},
					mock.MatchedBy(func(callback func(res *kms.ListAliasesOutput, lastPage bool) bool) bool {
						callback(&kms.ListAliasesOutput{
							Aliases: []*kms.AliasListEntry{
								{AliasName: aws.String("alias/1")},
								{AliasName: aws.String("alias/foo/2")},
								{AliasName: aws.String("alias/aw/3")},
								{AliasName: aws.String("alias/aws/4")},
								{AliasName: aws.String("alias/aws/5")},
								{AliasName: aws.String("alias/awss/6")},
								{AliasName: aws.String("alias/aws7")},
							},
						}, true)
						return true
					})).Return(nil).Once()
			},
			want: []*kms.AliasListEntry{
				{AliasName: aws.String("alias/1")},
				{AliasName: aws.String("alias/foo/2")},
				{AliasName: aws.String("alias/aw/3")},
				{AliasName: aws.String("alias/awss/6")},
				{AliasName: aws.String("alias/aws7")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(1)
			client := awstest.MockFakeKMS{}
			tt.mocks(&client)
			r := &kmsRepository{
				client: &client,
				cache:  store,
			}
			got, err := r.ListAllAliases()
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllAliases()
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				assert.IsType(t, []*kms.AliasListEntry{}, store.Get("kmsListAllAliases"))
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %v -> %v", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}
