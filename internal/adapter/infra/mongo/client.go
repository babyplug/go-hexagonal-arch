package mongo

import (
	"clean-arch/internal/adapter/config"
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	_client *client
	_cOnce  = sync.Once{}
)

type client struct {
	*mongo.Client
}

type Client interface {
	Connect(ctx context.Context) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
	Disconnect(ctx context.Context) error
	ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error)
	ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (mongo.ListDatabasesResult, error)
	NumberSessionsInProgress() int
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	StartSession(opts ...*options.SessionOptions) (mongo.Session, error)
	Timeout() *time.Duration
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
	UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(mongo.SessionContext) error) error
	Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
}

type Collection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error)
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	Drop(ctx context.Context) error
	Indexes() mongo.IndexView
}

func New(ctx context.Context, cfg *config.Config) (Client, error) {
	var err error
	_cOnce.Do(func() {
		var c *mongo.Client
		c, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
		if err != nil {
			return
		}

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err = c.Ping(ctx, readpref.Primary()); err != nil {
			return
		}

		_client = &client{c}
	})

	return _client, nil
}

func Reset() {
	_cOnce = sync.Once{}
}
