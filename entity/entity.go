package entity

import (
	"encoding/json"
	"fmt"
)

var ErrEntityNotFound = fmt.Errorf("entity not found")
var ErrEntityAlreadyExists = fmt.Errorf("entity already exists")
var ErrIdNotEmpty = fmt.Errorf("id is not empty")
var ErrResourceIdNotEmpty = fmt.Errorf("resource id is not empty")
var ErrResourceIdNotSupported = fmt.Errorf("resource id is not supported")

type Entity interface {
	GetID() ID                      // GetId should return internal identifier ID of the entity
	SetID(ID) error                 // SetId is called when the system is assigning an ID to the entity if applies
	GetResourceID() (string, error) // GetResourceId is an opinionated way to identify publicly an entity (slug, name, id, etc)
	SetResourceID(string) error     // SetResourceId @deprecated is called when the system is assigning a resource id to the entity if applies
}

type EventfulEntity interface {
	PreCreate() error // PreCreate is called before the entity is created
	PreUpdate() error // PreUpdate is called before the entity is updated
}

type ID string

func (i ID) String() string {
	return string(i)
}

func (i ID) IsEmpty() bool {
	return i == ""
}

func (i ID) Equals(i2 ID) bool {
	return i.String() == i2.String()
}

func (i ID) IsCompound() bool {
	// {"i":""}
	if len(i.String()) > 8 && i.String()[0] == '{' && i.String()[len(i.String())-1] == '}' {
		return true
	}

	return false
}

func (i ID) GetCompoundIDs() map[string]ID {
	var strIds map[string]string
	err := json.Unmarshal([]byte(i.String()), &strIds)
	if err != nil {
		return nil
	}

	ids := make(map[string]ID)
	for k, v := range strIds {
		ids[k] = NewIDFromString(v)
	}

	return ids
}

func NewIDFromString(id string) ID {
	return ID(id)
}
