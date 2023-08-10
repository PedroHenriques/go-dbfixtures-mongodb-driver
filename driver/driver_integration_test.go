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

type driverIntegrationTestSuite struct {
	suite.Suite

	mongoCliente *mongo.Client
	mongoDbName  string
}

func (suite *driverIntegrationTestSuite) SetupSuite() {
	ConnUrl := "mongodb://testmongo:27017"
	opts := options.Client().ApplyURI(ConnUrl)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	suite.mongoCliente = client

	suite.mongoDbName = "testCol"
}

func (suite *driverIntegrationTestSuite) TearDownSuite() {
	suite.mongoCliente.Disconnect(context.TODO())
}

func (suite *driverIntegrationTestSuite) SetupTest() {
	database := suite.mongoCliente.Database(suite.mongoDbName)
	err := database.Drop(context.TODO())
	if err != nil {
		panic(err)
	}
}

func (suite *driverIntegrationTestSuite) TestNewItShouldReturnAnInstanceOfTheInterfaceIDriver() {
	sut := driver.New(
		suite.mongoCliente, "some name", &options.DatabaseOptions{},
	)

	require.Implements(suite.T(), (*dbfixtures.IDriver)(nil), sut)
}

func (suite *driverIntegrationTestSuite) TestTruncateItShouldDeleteAllDocumentsFromTheProvidedCollectionsAndReturnNil() {
	database := suite.mongoCliente.Database(suite.mongoDbName)
	col1 := database.Collection("col1")
	col2 := database.Collection("col2")
	col3 := database.Collection("col3")

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
	insRes3, _ := col3.InsertMany(
		context.TODO(),
		[]interface{}{
			testDocument{Name: "doc 31", Age: 873425},
			testDocument{Name: "doc 32", Age: 78326},
			testDocument{Name: "doc 33", Age: 53},
		},
	)

	require.Equal(suite.T(), 2, len(insRes1.InsertedIDs))
	require.Equal(suite.T(), 1, len(insRes2.InsertedIDs))
	require.Equal(suite.T(), 3, len(insRes3.InsertedIDs))

	sut := driver.New(
		suite.mongoCliente, suite.mongoDbName, &options.DatabaseOptions{},
	)
	res := sut.Truncate([]string{"col1", "col3"})

	require.Nil(suite.T(), res)

	col1Count, _ := col1.CountDocuments(context.TODO(), bson.D{})
	col2Count, _ := col2.CountDocuments(context.TODO(), bson.D{})
	col3Count, _ := col3.CountDocuments(context.TODO(), bson.D{})

	require.EqualValues(suite.T(), 0, col1Count)
	require.EqualValues(suite.T(), 1, col2Count)
	require.EqualValues(suite.T(), 0, col3Count)
}

func (suite *driverIntegrationTestSuite) TestInsertFixturesItShouldInsertTheProvidedDocumentsInTheSpecifiedCollectionAndReturnNil() {
	database := suite.mongoCliente.Database(suite.mongoDbName)
	col1 := database.Collection("col1")
	col2 := database.Collection("col2")
	col3 := database.Collection("col3")

	col1CountBefore, _ := col1.CountDocuments(context.TODO(), bson.D{})
	col2CountBefore, _ := col2.CountDocuments(context.TODO(), bson.D{})
	col3CountBefore, _ := col3.CountDocuments(context.TODO(), bson.D{})

	require.EqualValues(suite.T(), 0, col1CountBefore)
	require.EqualValues(suite.T(), 0, col2CountBefore)
	require.EqualValues(suite.T(), 0, col3CountBefore)

	expectedDocuments := &[]interface{}{
		testDocument{Name: "doc 31", Age: 24},
		testDocument{Name: "doc 16", Age: 987},
	}

	sut := driver.New(
		suite.mongoCliente, suite.mongoDbName, &options.DatabaseOptions{},
	)
	res := sut.InsertFixtures("col3", *expectedDocuments)

	require.Nil(suite.T(), res)

	col1Count, _ := col1.CountDocuments(context.TODO(), bson.D{})
	require.EqualValues(suite.T(), 0, col1Count)
	col2Count, _ := col2.CountDocuments(context.TODO(), bson.D{})
	require.EqualValues(suite.T(), 0, col2Count)

	colDocs, _ := col3.Find(context.TODO(), bson.D{})
	actualDocuments := &[]testDocument{}
	colDocs.All(context.TODO(), actualDocuments)

	require.EqualValues(suite.T(), len(*expectedDocuments), len(*actualDocuments))
	for i, actualDocument := range *actualDocuments {
		require.EqualExportedValues(suite.T(), (*expectedDocuments)[i], actualDocument)
	}
}

func (suite *driverIntegrationTestSuite) TestCloseItShouldReturnNil() {
	sut := driver.New(
		suite.mongoCliente, suite.mongoDbName, &options.DatabaseOptions{},
	)
	res := sut.Close()

	require.Nil(suite.T(), res)
}
