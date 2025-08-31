package mongo_test

import (
	"context"
	"testing"

	"github.com/davfer/crudo/entity"
	"github.com/davfer/crudo/mongo/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type testEntity struct {
	ID            bson.ObjectID `bson:"_id"`
	Attr1         string        `bson:"attr_1"`
	SomeNiceField string        `bson:"some_nice_field"`
}

func (t testEntity) GetID() entity.ID {
	return mongo.NewIDFromObjectID(t.ID)
}

func (t testEntity) SetID(id entity.ID) error {
	t.ID = mongo.ToMustObjectID(id)
	return nil
}

func (t testEntity) GetResourceID() (string, error) {
	return t.Attr1, nil
}

func (t testEntity) SetResourceID(s string) error {
	t.Attr1 = s
	return nil
}

func TestNewMongoRepository(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:latest")
	defer func() {
		if err := testcontainers.TerminateContainer(mongodbContainer); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}
	t.Logf("mongodb uri: %s", uri)

	client, err := mongo2.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("failed to connect to mongo: %s", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			t.Fatalf("failed to disconnect from mongo: %s", err)
		}
	}()
	db := client.Database("test")
	collection := db.Collection("test")

	repo := mongo.NewMongoRepository[*testEntity](collection)
	err = repo.Start(ctx, func(ctx context.Context) error {
		t.Logf("onBootstrap called")
		return nil
	})
	if err != nil {
		t.Fatalf("failed to start repository: %s", err)
	}

	res, err := repo.Create(ctx, &testEntity{
		ID:    bson.NewObjectID(),
		Attr1: "test",
	})
	if err != nil {
		t.Fatalf("failed to create entity: %s", err)
	}
	if res.GetID().String() == "" {
		t.Fatalf("entity id is empty")
	}
	t.Logf("created entity with id: %s", res.GetID().String())

	res2, err := repo.Read(ctx, res.GetID())
	if err != nil {
		t.Fatalf("failed to read entity: %s", err)
	}
	if !res2.GetID().Equals(res.GetID()) {
		t.Fatalf("entity id is not test")
	}
	if res2.Attr1 != res.Attr1 {
		t.Fatalf("entity attr1 is not test")
	}
	t.Logf("read entity with id: %s", res2.GetID().String())

	ress, err := repo.ReadAll(ctx)
	if len(ress) != 1 {
		t.Fatalf("failed to read all entities: %s", err)
	}
	t.Logf("read all entities: %d", len(ress))

	err = repo.Delete(ctx, res)
	if err != nil {
		t.Fatalf("failed to delete entity: %s", err)
	}
	t.Logf("deleted entity with id: %s", res.GetID().String())
}
