package enumerator

import (
	"github.com/sirupsen/logrus"
	"github.com/snyk/driftctl/pkg/iac/config"
	"github.com/snyk/driftctl/pkg/iac/terraform/state/backend"
)

type StateEnumerator interface {
	Origin() string
	Enumerate() ([]string, error)
}

func GetEnumerator(config config.SupplierConfig) (StateEnumerator, error) {

	switch config.Backend {
	case backend.BackendKeyFile:
		return NewFileEnumerator(config), nil
	case backend.BackendKeyS3:
		return NewS3Enumerator(config), nil
	case backend.BackendKeyAzureRM:
		return NewAzureRMEnumerator(config)
	}

	logrus.WithFields(logrus.Fields{
		"backend": config.Backend,
	}).Debug("No enumerator for backend")

	return nil, nil
}
