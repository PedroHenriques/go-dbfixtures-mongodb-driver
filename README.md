[![Coverage Status](https://coveralls.io/repos/github/PedroHenriques/go-dbfixtures-mongodb-driver/badge.svg?branch=main)](https://coveralls.io/github/PedroHenriques/go-dbfixtures-mongodb-driver?branch=main)
![ci workflow](https://github.com/PedroHenriques/go-dbfixtures-mongodb-driver/actions/workflows/ci.yml/badge.svg?branch=main)
![cd workflow](https://github.com/PedroHenriques/go-dbfixtures-mongodb-driver/actions/workflows/cd.yml/badge.svg)

# Fixtures Manager MongoDB Driver

An abstraction layer for the [mongodb package](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo) to facilitate handling database fixtures for testing purposes, in a MongoDB database.  
This package is ment to be used in conjunction with the [dbfixtures package](https://pkg.go.dev/github.com/PedroHenriques/go-dbfixtures/dbfixtures), but can also be used by itself.

## Installation

```sh
go get github.com/PedroHenriques/go-dbfixtures-mongodb-driver
```

## Usage

This package exposes the `func New(mongoClient *mongo.Client, dbName string, dbOpts *options.DatabaseOptions) dbfixtures.IDriver` function that returns an instance of the driver.  

**Note:** For detailed information about the `options.DatabaseOptions` argument, please consult the [MongoDB Golang driver's documentation]https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.12.1/mongo/options#DatabaseOptions).

An instance of the driver exposes the following interface

```go
type IDriver interface {
	// clears the specified "tables" of any content
	Truncate(tableNames []string) error

	// inserts the supplied "rows" into the specified "table"
	InsertFixtures(tableName string, fixtures []interface{}) error

	// cleanup and terminate the connection to the database
	Close() error
}
```

### Example

```go
package driver_test

import (
	"context"

	"github.com/PedroHenriques/go-dbfixtures-mongodb-driver/driver"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type driverE2eTestSuite struct {
	suite.Suite

	mongoCliente *mongo.Client
	mongoDbName  string
}

func (suite *driverE2eTestSuite) SetupSuite() {
	ConnUrl := "mongodb://testmongo:27017"
	opts := options.Client().ApplyURI(ConnUrl)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	suite.mongoCliente = client

	suite.mongoDbName = "testCol"
}

func (suite *driverE2eTestSuite) TearDownSuite() {
	suite.mongoCliente.Disconnect(context.TODO())
}

func (suite *driverE2eTestSuite) TestItShouldWork() {
	database := suite.mongoCliente.Database(suite.mongoDbName)
	col1 := database.Collection("col1")

	insRes1, _ := col1.InsertMany(
		context.TODO(),
		[]interface{}{
			testDocument{Name: "doc 11", Age: 3},
			testDocument{Name: "doc 12", Age: 33},
		},
	)

	require.Equal(suite.T(), 2, len(insRes1.InsertedIDs))

	expectedDocuments := &[]interface{}{
		testDocument{Name: "doc 26", Age: 86},
	}

	driver := driver.New(
		suite.mongoCliente, suite.mongoDbName, &options.DatabaseOptions{},
	)

	err := driver.InsertFixtures("col1", *expectedDocuments)

	require.Nil(suite.T(), err)

	colDocs, _ := col1.Find(context.TODO(), bson.D{})
	actualDocuments := &[]testDocument{}
	colDocs.All(context.TODO(), actualDocuments)

	require.EqualValues(suite.T(), len(*expectedDocuments), len(*actualDocuments))
	for i, actualDocument := range *actualDocuments {
		require.EqualExportedValues(suite.T(), (*expectedDocuments)[i], actualDocument)
	}
}
```

## Testing This Package

* `cd` into the package's root directory
* run `sh cli/test.sh -b -gv 1.20`