package mongo

import (
	"encoding/json"
	"fmt"

	"github.com/davfer/crudo/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewIDFromObjectID(id primitive.ObjectID) entity.ID {
	if id.IsZero() {
		return ""
	}
	return entity.ID(id.Hex())
}

func NewIDFromObjectIDs(ids map[string]entity.ID) entity.ID {
	jsoned, _ := json.Marshal(ids)

	return entity.ID(jsoned)
}

func ToMustObjectID(i entity.ID) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(i.String())
	if err != nil {
		panic(fmt.Sprintf("could not convert %s to ObjectId", i))
	}

	return id
}

func TryObjectID(i entity.ID) *primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(i.String())
	if err != nil {
		return nil
	}

	return &id
}
