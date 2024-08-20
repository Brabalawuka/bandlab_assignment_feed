package mongodb

import (
    "context"
    "fmt"
    "sync"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/cloudwego/hertz/pkg/common/hlog"
)

// Config holds the configuration for MongoDB
type Config struct {
    URL      string
    Database string
}

// MongoService defines the interface for MongoDB operations
type MongoService interface {
    GetClient() *mongo.Client
    GetDatabase() *mongo.Database
    Collection(name string) *mongo.Collection
}

// MongoClient is a wrapper for the MongoDB client
type MongoClient struct {
    client   *mongo.Client
    database *mongo.Database
}

var (
    once        sync.Once
    mongoClient MongoService
)

// Initialize sets up the MongoDB client
func Initialize(cfg *Config) error {
    var initError error
    once.Do(func() {
        if err := validateConfig(cfg); err != nil {
            initError = err
            hlog.Error("[mongoservice] invalid configuration: %w", err)
            return
        }

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URL))
        if err != nil {
            initError = fmt.Errorf("[mongoservice] failed to connect to MongoDB: %w", err)
            hlog.Error(initError)
            return
        }

        // Ping the database
        if err = client.Ping(ctx, nil); err != nil {
            initError = fmt.Errorf("[mongoservice] failed to ping MongoDB: %w", err)
            hlog.Error(initError)
            return
        }

        mongoClient = &MongoClient{
            client:   client,
            database: client.Database(cfg.Database),
        }
    })
    return initError
}

func GetMongoService() MongoService {
    return mongoClient
}

// validateConfig checks if the provided configuration is valid
func validateConfig(cfg *Config) error {
    if cfg.URL == "" || cfg.Database == "" {
        return fmt.Errorf("invalid configuration: all fields must be non-empty")
    }
    return nil
}

// GetClient returns the global MongoDB client
func (c *MongoClient) GetClient() *mongo.Client {
    if c == nil {
        hlog.Error("[mongoservice] MongoDB client is not initialized")
        return nil
    }
    return c.client
}

// GetDatabase returns the global MongoDB database
func (c *MongoClient) GetDatabase() *mongo.Database {
    if c == nil {
        hlog.Error("[mongoservice] MongoDB client is not initialized")
        return nil
    }
    return c.database
}

// Collection returns a handle to a MongoDB collection
func (c *MongoClient) Collection(name string) *mongo.Collection {
    if c == nil {
        hlog.Error("[mongoservice] MongoDB client is not initialized")
        return nil
    }
    return c.database.Collection(name)
}