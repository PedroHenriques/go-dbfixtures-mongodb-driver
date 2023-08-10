package driver_test

import (
	"context"

	"github.com/PedroHenriques/go-dbfixtures-mongodb-driver/driver"
	"github.com/PedroHenriques/go-dbfixtures/dbfixtures"
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

func (suite *driverE2eTestSuite) SetupTest() {
	database := suite.mongoCliente.Database(suite.mongoDbName)
	err := database.Drop(context.TODO())
	if err != nil {
		panic(err)
	}
}

func (suite *driverE2eTestSuite) TestTheDriverWorksWithTheCoreLibrary() {
	database := suite.mongoCliente.Database(suite.mongoDbName)
	col1 := database.Collection("col1")
	col2 := database.Collection("col2")

	insRes1, _ := col1.InsertMany(
		context.TODO(),
		[]interface{}{
			testDocument{Name: "doc 11", Age: 3},
			testDocument{Name: "doc 12", Age: 33},
		},
	)
	insRes2, _ := col2.InsertMany(
		context.TODO(),
		[]interface{}{
			testDocument{Name: "doc 21", Age: 1},
		},
	)

	require.Equal(suite.T(), 2, len(insRes1.InsertedIDs))
	require.Equal(suite.T(), 1, len(insRes2.InsertedIDs))

	expectedDocuments := &[]interface{}{
		testDocument{Name: "doc 26", Age: 86},
		testDocument{Name: "doc 27", Age: 87},
	}

	driver := driver.New(
		suite.mongoCliente, suite.mongoDbName, &options.DatabaseOptions{},
	)

	fixtureHandler := dbfixtures.New(driver)

	err := fixtureHandler.InsertFixtures(
		[]string{"col2"},
		map[string][]interface{}{
			"col2": *expectedDocuments,
		},
	)

	require.Nil(suite.T(), err)

	colDocs, _ := col2.Find(context.TODO(), bson.D{})
	actualDocuments := &[]testDocument{}
	colDocs.All(context.TODO(), actualDocuments)

	require.EqualValues(suite.T(), len(*expectedDocuments), len(*actualDocuments))
	for i, actualDocument := range *actualDocuments {
		require.EqualExportedValues(suite.T(), (*expectedDocuments)[i], actualDocument)
	}
}
