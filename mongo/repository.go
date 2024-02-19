package mongo

import (
	"context"
	"fmt"
	"github.com/davfer/crudo/entity"
	"github.com/davfer/go-specification"
	"github.com/davfer/go-specification/mongo/repository"
	mongoSpec "github.com/davfer/go-specification/mongo/resolver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
)

type Repository[K entity.Entity] struct {
	criteriaRepo repository.CriteriaRepository[K]
	Collection   *mongo.Collection
	logger       *logrus.Entry
}

func NewMongoRepository[K entity.Entity](collection *mongo.Collection) *Repository[K] {
	return NewMongoRepositoryWithOpts[K](collection, nil)
}

func NewMongoRepositoryWithOpts[K entity.Entity](collection *mongo.Collection, logger *logrus.Entry) *Repository[K] {
	if logger == nil {
		l := logrus.New()
		l.SetOutput(io.Discard)
		logger = l.WithField("repository", collection.Name())
	}

	return &Repository[K]{
		criteriaRepo: repository.CriteriaRepository[K]{
			Collection: collection,
			Converter:  mongoSpec.NewMongoResolver(),
		},
		Collection: collection,
		logger:     logger.WithField("repository", collection.Name()),
	}
}

func (r *Repository[K]) Start(ctx context.Context, onBootstrap func(ctx context.Context) error) error {
	r.logger.Info("bootstrapping mongo repository")

	if exists, err := collectionExists(ctx, r.Collection); err == nil && !exists && onBootstrap != nil {
		r.logger.Debug("sending onBootstrap event")
		return onBootstrap(ctx)
	}

	return nil
}
func (r *Repository[K]) Create(ctx context.Context, e K) (K, error) {
	r.logger.WithField("entity", e).Debug("creating entity")

	err := e.PreCreate()
	if err != nil {
		r.logger.WithError(err).Error("error pre creating entity")
		return e, errors.Wrap(err, "error pre creating entity")
	}

	insertResult, err := r.Collection.InsertOne(ctx, e)
	if err != nil {
		r.logger.WithError(err).Error("error inserting entity")
		return e, err
	}

	id := entity.NewIdFromObjectId(insertResult.InsertedID.(primitive.ObjectID))
	err = e.SetId(id)
	if err != nil {
		r.logger.WithError(err).Error("error setting id")
		return e, err
	}

	r.logger.WithField("id", insertResult.InsertedID).Debug("entity created")
	return e, nil
}

func (r *Repository[K]) Read(ctx context.Context, id entity.Id) (e K, err error) {
	r.logger.WithField("id", id).Debug("reading entity")

	err = r.Collection.FindOne(ctx, r.getMongoSearchIdentifier(id)).Decode(&e)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return e, entity.ErrEntityNotFound
		}

		r.logger.WithError(err).Error("error reading entity")
		return e, errors.Wrap(err, fmt.Sprintf("error reading entity %s", id))
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
	r.logger.Debug("reading all entities")

	var entities []K
	cursor, err := r.Collection.Find(ctx, bson.D{{}})
	if err != nil {
		r.logger.WithError(err).Error("error finding all")
		return nil, errors.Wrap(err, "error finding all")
	}
	if err = cursor.All(ctx, &entities); err != nil {
		r.logger.WithError(err).Error("error reading all")
		return nil, errors.Wrap(err, "error reading all")
	}

	if len(entities) == 0 {
		return []K{}, nil
	}

	return entities, nil
}

func (r *Repository[K]) Update(ctx context.Context, entity K) error {
	r.logger.WithField("id", entity.GetId()).Debug("updating entity")

	err := entity.PreUpdate()
	if err != nil {
		r.logger.WithError(err).Error("error pre updating entity")
		return errors.Wrap(err, "error pre updating entity")
	}

	_, err = r.Collection.UpdateOne(ctx, r.getMongoSearchIdentifier(entity.GetId()), bson.M{"$set": entity})
	if err != nil {
		r.logger.WithError(err).Error("error updating entity")
		return errors.Wrap(err, fmt.Sprintf("error updating entity %s", entity.GetId()))
	}

	return nil
}

func (r *Repository[K]) Delete(ctx context.Context, entity K) error {
	r.logger.WithField("id", entity.GetId()).Debug("deleting entity")
	_, err := r.Collection.DeleteOne(ctx, r.getMongoSearchIdentifier(entity.GetId()))
	if err != nil {
		r.logger.WithError(err).Error("error deleting entity")
		return errors.Wrap(err, fmt.Sprintf("error deleting entity %s", entity.GetId()))
	}

	return nil
}

func (r *Repository[K]) getMongoSearchIdentifier(id entity.Id) bson.M {
	if id.IsCompound() {
		var m bson.M

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
		return false, errors.Wrap(err, "could not list collections")
	}

	for _, c := range collections {
		if c == col.Name() {
			return true, nil
		}
	}

	return false, nil
}
