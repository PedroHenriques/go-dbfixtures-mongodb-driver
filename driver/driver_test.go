package driver_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type testDocument struct {
	Name string
	Age  int
}

func TestDriverSuite(t *testing.T) {
	suite.Run(t, new(driverIntegrationTestSuite))
	suite.Run(t, new(driverE2eTestSuite))
}
