package mgm

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	config *Config
	client *mongo.Client
	db     *mongo.Database
}

func NewConnection(conf *Config, dbName string, opts ...*options.ClientOptions) (*Connection, error) {
	if conf == nil {
		conf = defaultConf()
	}

	conn := &Connection{config: conf}

	var err error
	if conn.client, err = conn.NewClient(opts...); err != nil {
		return nil, err
	}

	conn.db = conn.client.Database(dbName)

	return conn, nil
}

func (c *Connection) Ctx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), c.config.CtxTimeout)
	return ctx
}

// NewClient return new mongodb client.
func (c *Connection) NewClient(opts ...*options.ClientOptions) (*mongo.Client, error) {
	client, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	if err = client.Connect(c.Ctx()); err != nil {
		return nil, err
	}

	return client, nil
}

// NewCollection return new collection with passed database
func (c *Connection) NewCollection(db *mongo.Database, name string, opts ...*options.CollectionOptions) *Collection {
	coll := db.Collection(name, opts...)

	return &Collection{Connection: c, Collection: coll}
}

func (c *Connection) Client() *mongo.Client {
	return c.client
}

func (c *Connection) Database() *mongo.Database {
	return c.db
}

// CollectionByName return new collection from default config
func (c *Connection) CollectionByName(name string, opts ...*options.CollectionOptions) *Collection {
	return c.NewCollection(c.db, name, opts...)
}

type Config struct {
	// Set to 10 second (10*time.Second) for example.
	CtxTimeout time.Duration
}

func defaultConf() *Config {
	return &Config{CtxTimeout: 10 * time.Second}
}
