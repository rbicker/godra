package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database describe all database interaction.
type Database interface {
	Connect() error
	Disconnect() error
	FindUserByName(string) (*User, error)
	FindUserByID(string) (*User, error)
}

// MGO implements the database interface, representing a mongodb connection.
type MGO struct {
	url     string
	dbname  string
	colname string
}

var _ Database = MGO{}

var client *mongo.Client
var col *mongo.Collection

// NewMongoConnection creates a new mongo database connection.
// It takes functional parameters to change default options
// such as the mongo url
// It returns the newly created server or an error if
// something went wrong.
func NewMongoConnection(opts ...func(*MGO) error) (Database, error) {
	// create server with default options
	var m = MGO{
		url:     "mongodb://localhost:27017",
		dbname:  "db",
		colname: "users",
	}

	// run functional options
	for _, op := range opts {
		err := op(&m)
		if err != nil {
			return nil, fmt.Errorf("setting option failed: %w", err)
		}
	}
	return m, nil
}

// SetURL changes the url to which the connection should be established.
func SetURL(url string) func(*MGO) error {
	return func(m *MGO) error {
		m.url = url
		return nil
	}
}

// SetDBName changes the name of the mongodb database.
func SetDBName(dbname string) func(*MGO) error {
	return func(m *MGO) error {
		m.dbname = dbname
		return nil
	}
}

// SetCollectionName changes the name of the mongodb collection.
func SetCollectionName(colname string) func(*MGO) error {
	return func(m *MGO) error {
		m.colname = colname
		return nil
	}
}

// Connect establishes a connection to a mongodb server.
func (m MGO) Connect() error {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(m.url))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	col = client.Database(m.dbname).Collection(m.colname)
	if err != nil {
		return err
	}
	return nil
}

// Disconnect closes the connection to the mongodb server.
func (MGO) Disconnect() error {
	return client.Disconnect(context.TODO())
}
