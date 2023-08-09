package driver_test

import (
	"context"
	"testing"

	"github.com/PedroHenriques/go-dbfixtures-mongodb-driver/driver"
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

	require.Implements(suite.T(), (*driver.IDriver)(nil), sut)
}

func (suite *driverIntegrationTestSuite) TestTruncateItShouldDeleteAllDocumentsFromTheProvidedCollectionsAndReturnNil() {
	database := suite.mongoCliente.Database(suite.mongoDbName)
	col1 := database.Collection("col1")
	col2 := database.Collection("col2")
	col3 := database.Collection("col3")

	insRes1, _ := col1.InsertMany(
		context.TODO(),
		[]interface{}{
			testDocument{name: "doc 11", age: 3},
			testDocument{name: "doc 12", age: 33},
		},
	)
	insRes2, _ := col2.InsertMany(
		context.TODO(),
		[]interface{}{
			testDocument{name: "doc 21", age: 1},
		},
	)
	insRes3, _ := col3.InsertMany(
		context.TODO(),
		[]interface{}{
			testDocument{name: "doc 31", age: 873425},
			testDocument{name: "doc 32", age: 78326},
			testDocument{name: "doc 33", age: 53},
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

	sut := driver.New(
		suite.mongoCliente, suite.mongoDbName, &options.DatabaseOptions{},
	)
	res := sut.InsertFixtures(
		"col3",
		[]interface{}{
			testDocument{name: "doc 31", age: 24},
			testDocument{name: "doc 16", age: 987},
		},
	)

	require.Nil(suite.T(), res)

	col1Count, _ := col1.CountDocuments(context.TODO(), bson.D{})
	col2Count, _ := col2.CountDocuments(context.TODO(), bson.D{})
	col3Count, _ := col3.CountDocuments(context.TODO(), bson.D{})

	require.EqualValues(suite.T(), 0, col1Count)
	require.EqualValues(suite.T(), 0, col2Count)
	require.EqualValues(suite.T(), 2, col3Count)
}

func (suite *driverIntegrationTestSuite) TestCloseItShouldReturnNil() {
	sut := driver.New(
		suite.mongoCliente, suite.mongoDbName, &options.DatabaseOptions{},
	)
	res := sut.Close()

	require.Nil(suite.T(), res)
}

func TestDriverSuite(t *testing.T) {
	suite.Run(t, new(driverIntegrationTestSuite))
}

type testDocument struct {
	name string
	age  int
}
