/*
Package driver provides functionality to manage seeds in a MongoDb database.
*/
package driver

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: replace with an import of this interface from the core package.
type IDriver interface {
	// clears the specified "tables" of any content
	Truncate(tableNames []string) error

	// inserts the supplied "rows" into the specified "table"
	InsertFixtures(tableName string, fixtures []interface{}) error

	// cleanup and terminate the connection to the database
	Close() error
}

/*
New creates and returns an instance of the MongoDb fixtures driver.
*/
func New(
	mongoClient *mongo.Client, dbName string, dbOpts *options.DatabaseOptions,
) IDriver {
	database := mongoClient.Database(dbName, dbOpts)

	return &driver{
		database: database,
	}
}

type driver struct {
	database *mongo.Database
}

/*
Truncate clears the specified "tables" of any content.
*/
func (driver *driver) Truncate(tableNames []string) error {
	var col *mongo.Collection

	for _, tableName := range tableNames {
		col = driver.database.Collection(tableName)
		
		err := col.Drop(context.TODO())
		if err != nil {
			return err
		}
	}

	return nil
}

/*
InsertFixtures inserts the supplied "rows" into the specified "table".
*/
func (driver *driver) InsertFixtures(tableName string, fixtures []interface{}) error {
	collection := driver.database.Collection(tableName)
	_, err := collection.InsertMany(context.TODO(), fixtures)
	if err != nil {
		return err
	}

	return nil
}

/*
Close cleanup and terminate the connection to the database.
*/
func (driver *driver) Close() error {
	return nil
}