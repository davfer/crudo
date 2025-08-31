package mongo

import (
	"context"
	"errors"
	"fmt"

	"github.com/davfer/archit/patterns/opts"
	"github.com/davfer/go-specification"
	"github.com/davfer/go-specification/mongo/v2/repository"
	mongoSpec "github.com/davfer/go-specification/mongo/v2/resolver"
	"github.com/go-logr/logr"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/davfer/crudo/entity"
)

type Repository[K entity.Entity] struct {
	criteriaRepo repository.CriteriaRepository[K]
	Collection   *mongo.Collection
	logger       logr.Logger
}

func WithLogger[K entity.Entity](logger logr.Logger) opts.Opt[Repository[K]] {
	return func(c Repository[K]) Repository[K] {
		c.logger = logger
		return c
	}
}

func NewMongoRepository[K entity.Entity](collection *mongo.Collection, o ...opts.Opt[Repository[K]]) *Repository[K] {
	r := opts.New[Repository[K]](o...)

	r.Collection = collection
	r.criteriaRepo = repository.CriteriaRepository[K]{
		Collection: collection,
		Converter:  mongoSpec.NewMongoConverter(),
	}

	return &r
}

func (r *Repository[K]) Start(ctx context.Context, onBootstrap func(ctx context.Context) error) error {
	r.logger.Info("bootstrapping mongo repository")

	if exists, err := collectionExists(ctx, r.Collection); err == nil && !exists && onBootstrap != nil {
		r.logger.Info("sending onBootstrap event")
		return onBootstrap(ctx)
	}

	return nil
}
func (r *Repository[K]) Create(ctx context.Context, e K) (K, error) {
	r.logger.V(5).Info("creating entity", "entity", e)

	if ee, ok := entity.Entity(e).(entity.EventfulEntity); ok {
		err := ee.PreCreate()
		if err != nil {
			r.logger.Error(err, "error pre creating entity")
			return e, fmt.Errorf("error pre creating entity: %w", err)
		}
	}

	insertResult, err := r.Collection.InsertOne(ctx, e)
	if err != nil {
		r.logger.Error(err, "error inserting entity")
		return e, err
	}

	id := entity.NewIdFromObjectId(insertResult.InsertedID.(primitive.ObjectID))
	err = e.SetId(id)
	if err != nil {
		r.logger.Error(err, "error setting id")
		return e, err
	}

	r.logger.V(2).Info("entity created", "id", insertResult.InsertedID)
	return e, nil
}

func (r *Repository[K]) Read(ctx context.Context, id entity.Id) (e K, err error) {
	r.logger.V(5).Info("reading entity", "id", id)

	err = r.Collection.FindOne(ctx, r.getMongoSearchIdentifier(id)).Decode(&e)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return e, entity.ErrEntityNotFound
		}

		r.logger.Error(err, "error reading entity")
		return e, fmt.Errorf("error reading : entity %s: %w", id, err)
	}

	return e, nil
}

func (r *Repository[K]) Match(ctx context.Context, c specification.Criteria) ([]K, error) {
	return r.criteriaRepo.Match(ctx, c)
}

func (r *Repository[K]) MatchOne(ctx context.Context, c specification.Criteria) (K, error) {
	return r.criteriaRepo.MatchOne(ctx, c)
}

func (r *Repository[K]) ReadAll(ctx context.Context) ([]K, error) {
	r.logger.V(5).Info("reading all entities")

	var entities []K
	cursor, err := r.Collection.Find(ctx, bson.D{{}})
	if err != nil {
		r.logger.Error(err, "error finding all")
		return nil, fmt.Errorf("error finding all: %w", err)
	}
	if err = cursor.All(ctx, &entities); err != nil {
		r.logger.Error(err, "error reading all")
		return nil, fmt.Errorf("error reading all: %w", err)
	}

	if len(entities) == 0 {
		return []K{}, nil
	}

	return entities, nil
}

func (r *Repository[K]) Update(ctx context.Context, e K) error {
	r.logger.V(5).Info("updating entity", "id", e.GetId())

	if ee, ok := entity.Entity(e).(entity.EventfulEntity); ok {
		err := ee.PreUpdate()
		if err != nil {
			r.logger.Error(err, "error pre updating entity")
			return fmt.Errorf("error pre updating entity: %w", err)
		}
	}

	_, err := r.Collection.UpdateOne(ctx, r.getMongoSearchIdentifier(e.GetId()), bson.M{"$set": e})
	if err != nil {
		r.logger.Error(err, "error updating entity")
		return fmt.Errorf("error updating entity %s: %w", e.GetId(), err)
	}

	return nil
}

func (r *Repository[K]) Delete(ctx context.Context, entity K) error {
	r.logger.V(5).Info("deleting entity", "id", entity.GetId())
	_, err := r.Collection.DeleteOne(ctx, r.getMongoSearchIdentifier(entity.GetId()))
	if err != nil {
		r.logger.Error(err, "error deleting entity")
		return fmt.Errorf("error deleting entity %s: %w", entity.GetId(), err)
	}

	return nil
}

func (r *Repository[K]) getMongoSearchIdentifier(id entity.Id) bson.M {
	if id.IsCompound() {
		m := bson.M{}
		for k, v := range id.GetCompoundIds() {
			m[k] = v.TryObjectId()
			if m[k] == nil {
				m[k] = v.String()
			}
		}

		return m
	}

	return bson.M{"_id": id.ToMustObjectId()}
}

func collectionExists(ctx context.Context, col *mongo.Collection) (bool, error) {
	collections, err := col.Database().ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return false, fmt.Errorf("could not list collections: %w", err)
	}

	for _, c := range collections {
		if c == col.Name() {
			return true, nil
		}
	}

	return false, nil
}
