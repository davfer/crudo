package entity

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrEntityNotFound = errors.New("entity not found")
var ErrEntityAlreadyExists = errors.New("entity already exists")
var ErrIdNotEmpty = fmt.Errorf("id is not empty")
var ErrResourceIdNotEmpty = fmt.Errorf("resource id is not empty")
var ErrResourceIdNotSupported = fmt.Errorf("resource id is not supported")

type Entity interface {
	GetId() Id
	SetId(Id) error
	GetResourceId() (string, error)
	SetResourceId(string) error
	PreCreate() error
	PreUpdate() error
}

type Id string

func (i Id) String() string {
	return string(i)
}

func (i Id) IsEmpty() bool {
	return i == ""
}

func (i Id) Equals(i2 Id) bool {
	return i.String() == i2.String()
}

func (i Id) ToMustObjectId() primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(i.String())
	if err != nil {
		panic(fmt.Sprintf("could not convert %s to ObjectId", i))
	}

	return id
}

func (i Id) TryObjectId() *primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(i.String())
	if err != nil {
		return nil
	}

	return &id
}

func (i Id) IsCompound() bool {
	// {"i":""}
	if len(i.String()) > 8 && i.String()[0] == '{' && i.String()[len(i.String())-1] == '}' {
		return true
	}

	return false
}

func (i Id) GetCompoundIds() map[string]Id {
	var strIds map[string]string
	err := json.Unmarshal([]byte(i.String()), &strIds)
	if err != nil {
		return nil
	}

	ids := make(map[string]Id)
	for k, v := range strIds {
		ids[k] = NewIdFromString(v)
	}

	return ids
}

func NewIdFromString(id string) Id {
	return Id(id)
}
func NewIdFromObjectId(id primitive.ObjectID) Id {
	if id.IsZero() {
		return Id("")
	}
	return Id(id.Hex())
}

func NewIdFromObjectIds(ids map[string]Id) Id {
	jsoned, err := json.Marshal(ids)
	if err != nil {
		return Id("")
	}

	return Id(jsoned)
}
