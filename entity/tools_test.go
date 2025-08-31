package entity_test

import (
	"testing"

	"github.com/davfer/crudo/entity"
)

type testToolEntity struct {
	Id            string
	Attr1         string
	SomeNiceField string
}

func (t *testToolEntity) GetID() entity.ID {
	return entity.ID(t.Id)
}

func (t *testToolEntity) SetID(id entity.ID) error {
	t.Id = string(id)
	return nil
}

func (t *testToolEntity) GetResourceID() (string, error) {
	return t.Attr1, nil
}

func (t *testToolEntity) SetResourceID(s string) error {
	t.Attr1 = s
	return nil
}

func (t *testToolEntity) PreCreate() error {
	return nil
}

func (t *testToolEntity) PreUpdate() error {
	return nil
}

func TestContainsId(t *testing.T) {
	type testCase[K entity.Entity] struct {
		name     string
		entities []K
		e        K
		want     bool
	}
	tests := []testCase[*testToolEntity]{
		{
			name: "Test Contains",
			entities: []*testToolEntity{
				{Id: "1"},
				{Id: "2"},
			},
			e:    &testToolEntity{Id: "1"},
			want: true,
		},
		{
			name: "Test Not Contains",
			entities: []*testToolEntity{
				{Id: "1"},
			},
			e:    &testToolEntity{Id: "2"},
			want: false,
		},
		{
			name: "Test Empty",
			entities: []*testToolEntity{
				{Id: "1"},
			},
			e:    &testToolEntity{Attr1: "new attr"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entity.Contains(tt.entities, tt.e); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
