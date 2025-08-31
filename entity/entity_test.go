package entity

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestId_Equals(t *testing.T) {
	tests := []struct {
		name string
		i1   Id
		i2   Id
		want bool
	}{
		{
			name: "Test equals",
			i1:   "5f3e3e3e3e3e3e3e3e3e3e3e",
			i2:   "5f3e3e3e3e3e3e3e3e3e3e3e",
			want: true,
		},
		{
			name: "Test not equals",
			i1:   "5f3e3e3e3e3e3e3e3e3e3e3e",
			i2:   "5f3e3e3e3e3e3e3e3e3e3e3f",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i1.Equals(tt.i2); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestId_GetCompoundIds(t *testing.T) {
	tests := []struct {
		name string
		i    Id
		want map[string]Id
	}{
		{
			name: "Test compound ids",
			i:    Id(`{"i":"5f3e3e3e3e3e3e3e3e3e3e3e", "j":"5f3e3e3e3e3e3e3e3e3e3e3e"}`),
			want: map[string]Id{
				"i": "5f3e3e3e3e3e3e3e3e3e3e3e",
				"j": "5f3e3e3e3e3e3e3e3e3e3e3e",
			},
		},
		{
			name: "Test not compound ids",
			i:    "5f3e3e3e3e3e3e3e3e3e3e3e",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.GetCompoundIds(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCompoundIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestId_IsCompound(t *testing.T) {
	tests := []struct {
		name string
		i    Id
		want bool
	}{
		{
			name: "Test compound",
			i:    Id(`{"i":"5f3e3e3e3e3e3e3e3e3e3e3e", "j":"5f3e3e3e3e3e3e3e3e3e3e3e"}`),
			want: true,
		},
		{
			name: "Test not compound",
			i:    "5f3e3e3e3e3e3e3e3e3e3e3e",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsCompound(); got != tt.want {
				t.Errorf("IsCompound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestId_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		i    Id
		want bool
	}{
		{
			name: "Test empty",
			i:    "",
			want: true,
		},
		{
			name: "Test not empty",
			i:    "5f3e3e3e3e3e3e3e3e3e3e3e",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestId_String(t *testing.T) {
	tests := []struct {
		name string
		i    Id
		want string
	}{
		{
			name: "Test string",
			i:    "5f3e3e3e3e3e3e3e3e3e3e3e",
			want: "5f3e3e3e3e3e3e3e3e3e3e3e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestId_ToMustObjectId(t *testing.T) {
	objId, _ := primitive.ObjectIDFromHex("5f3e3e3e3e3e3e3e3e3e3e3e")
	tests := []struct {
		name  string
		i     Id
		want  primitive.ObjectID
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
			if tt.panic == "" {
				if got := tt.i.ToMustObjectId(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ToMustObjectId() = %v, want %v", got, tt.want)
				}
			} else {
				defer func() {
					if r := recover(); r != nil {
						if r != tt.panic {
							t.Errorf("ToMustObjectId() = %v, want %v", r, tt.panic)
						}
					}
				}()
				tt.i.ToMustObjectId()
			}
		})
	}
}

func TestId_TryObjectId(t *testing.T) {
	objId, _ := primitive.ObjectIDFromHex("5f3e3e3e3e3e3e3e3e3e3e3e")
	tests := []struct {
		name string
		i    Id
		want *primitive.ObjectID
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
			if got := tt.i.TryObjectId(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TryObjectId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIdFromObjectId(t *testing.T) {
	objId, _ := primitive.ObjectIDFromHex("5f3e3e3e3e3e3e3e3e3e3e3e")
	tests := []struct {
		name string
		id   primitive.ObjectID
		want Id
	}{
		{
			name: "Test new id from object id",
			id:   objId,
			want: "5f3e3e3e3e3e3e3e3e3e3e3e",
		},
		{
			name: "Test new id from object id",
			id:   primitive.ObjectID{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIdFromObjectId(tt.id); got != tt.want {
				t.Errorf("NewIdFromObjectId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIdFromObjectIds(t *testing.T) {
	tests := []struct {
		name string
		ids  map[string]Id
		want Id
	}{
		{
			name: "Test new id from object ids",
			ids: map[string]Id{
				"i": "5f3e3e3e3e3e3e3e3e3e3e3e",
				"j": "5f3e3e3e3e3e3e3e3e3e3e3e",
			},
			want: `{"i":"5f3e3e3e3e3e3e3e3e3e3e3e","j":"5f3e3e3e3e3e3e3e3e3e3e3e"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIdFromObjectIds(tt.ids); got != tt.want {
				t.Errorf("NewIdFromObjectIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIdFromString(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want Id
	}{
		{
			name: "Test new id from string",
			id:   "5f3e3e3e3e3e3e3e3e3e3e3e",
			want: "5f3e3e3e3e3e3e3e3e3e3e3e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIdFromString(tt.id); got != tt.want {
				t.Errorf("NewIdFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
