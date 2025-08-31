package mongo_test

import (
	"reflect"
	"testing"

	"github.com/davfer/crudo/entity"
	"github.com/davfer/crudo/mongo/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestId_ToMustObjectId(t *testing.T) {
	objId, _ := bson.ObjectIDFromHex("5f3e3e3e3e3e3e3e3e3e3e3e")
	tests := []struct {
		name  string
		i     entity.ID
		want  bson.ObjectID
		panic string
	}{
		{
			name: "Test to must object id",
			i:    "5f3e3e3e3e3e3e3e3e3e3e3e",
			want: objId,
		},
		{
			name:  "Test to must object id",
			i:     "ppppp",
			panic: "could not convert ppppp to ObjectId",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.panic == "" {
				if got := mongo.ToMustObjectID(tt.i); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ToMustObjectID() = %v, want %v", got, tt.want)
				}
			} else {
				defer func() {
					if r := recover(); r != nil {
						if r != tt.panic {
							t.Errorf("ToMustObjectID() = %v, want %v", r, tt.panic)
						}
					}
				}()
				mongo.ToMustObjectID(tt.i)
			}
		})
	}
}

func TestId_TryObjectId(t *testing.T) {
	objId, _ := bson.ObjectIDFromHex("5f3e3e3e3e3e3e3e3e3e3e3e")
	tests := []struct {
		name string
		i    entity.ID
		want *bson.ObjectID
	}{
		{
			name: "Test try object id",
			i:    "ppppp",
			want: nil,
		},
		{
			name: "Test try object id",
			i:    "5f3e3e3e3e3e3e3e3e3e3e3e",
			want: &objId,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := mongo.TryObjectID(tt.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TryObjectID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIdFromObjectId(t *testing.T) {
	objId, _ := bson.ObjectIDFromHex("5f3e3e3e3e3e3e3e3e3e3e3e")
	tests := []struct {
		name string
		id   bson.ObjectID
		want entity.ID
	}{
		{
			name: "Test new id from object id",
			id:   objId,
			want: "5f3e3e3e3e3e3e3e3e3e3e3e",
		},
		{
			name: "Test new id from object id",
			id:   bson.ObjectID{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := mongo.NewIDFromObjectID(tt.id); got != tt.want {
				t.Errorf("NewIDFromObjectID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIdFromObjectIds(t *testing.T) {
	tests := []struct {
		name string
		ids  map[string]entity.ID
		want entity.ID
	}{
		{
			name: "Test new id from object ids",
			ids: map[string]entity.ID{
				"i": "5f3e3e3e3e3e3e3e3e3e3e3e",
				"j": "5f3e3e3e3e3e3e3e3e3e3e3e",
			},
			want: `{"i":"5f3e3e3e3e3e3e3e3e3e3e3e","j":"5f3e3e3e3e3e3e3e3e3e3e3e"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := mongo.NewIDFromObjectIDs(tt.ids); got != tt.want {
				t.Errorf("NewIDFromObjectIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
