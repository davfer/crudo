package inmemory

import (
	"context"
	"github.com/davfer/crudo/entity"
	"reflect"
	"testing"
)

type testMemoEntity struct {
	Id            string
	Attr1         string
	SomeNiceField string
}

func (t *testMemoEntity) GetId() entity.Id {
	return entity.Id(t.Id)
}

func (t *testMemoEntity) SetId(id entity.Id) error {
	t.Id = string(id)
	return nil
}

func (t *testMemoEntity) GetResourceId() (string, error) {
	return t.Attr1, nil
}

func (t *testMemoEntity) SetResourceId(s string) error {
	t.Attr1 = s
	return nil
}

func (t *testMemoEntity) PreCreate() error {
	return nil
}

func (t *testMemoEntity) PreUpdate() error {
	return nil
}

func TestRepository_Create(t *testing.T) {
	type testCase[K entity.Entity] struct {
		name    string
		r       *Repository[K]
		ctx     context.Context
		calls   []K
		expect  []K
		wantErr bool
	}
	tests := []testCase[*testMemoEntity]{
		{
			name: "Test Create no policy",
			r:    NewInMemoryRepository[*testMemoEntity](nil),
			ctx:  context.TODO(),
			calls: []*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "2", Attr1: "attr2", SomeNiceField: "some_nice_field"},
				{Id: "3", Attr1: "attr3", SomeNiceField: "some_nice_field"},
			},
			expect: []*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "2", Attr1: "attr2", SomeNiceField: "some_nice_field"},
				{Id: "3", Attr1: "attr3", SomeNiceField: "some_nice_field"},
			},
		},
		{
			name: "Test Create with MRU policy",
			r:    NewInMemoryRepositoryWithPolicy[*testMemoEntity](nil, PolicyMru[*testMemoEntity]{Capacity: 2}),
			ctx:  context.TODO(),
			calls: []*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "2", Attr1: "attr2", SomeNiceField: "some_nice_field"},
				{Id: "3", Attr1: "attr3", SomeNiceField: "some_nice_field"},
			},
			expect: []*testMemoEntity{
				{Id: "2", Attr1: "attr2", SomeNiceField: "some_nice_field"},
				{Id: "3", Attr1: "attr3", SomeNiceField: "some_nice_field"},
			},
		},
		{
			name: "Test Create with LRU policy",
			r:    NewInMemoryRepositoryWithPolicy[*testMemoEntity](nil, PolicyLru[*testMemoEntity]{Capacity: 2}),
			ctx:  context.TODO(),
			calls: []*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "2", Attr1: "attr2", SomeNiceField: "some_nice_field"},
				{Id: "3", Attr1: "attr3", SomeNiceField: "some_nice_field"},
			},
			expect: []*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "2", Attr1: "attr2", SomeNiceField: "some_nice_field"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, c := range tt.calls {
				_, err := tt.r.Create(tt.ctx, c)
				if (err != nil) != tt.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(tt.r.Collection, tt.expect) {
				t.Errorf("Create() got = %v, want %v", tt.r.Collection, tt.expect)
			}
		})
	}
}
