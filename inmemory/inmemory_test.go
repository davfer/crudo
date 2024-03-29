package inmemory

import (
	"context"
	"errors"
	"github.com/davfer/crudo/entity"
	"github.com/davfer/go-specification"
	"reflect"
	"testing"
)

type testMemoEntity struct {
	Id            string
	Attr1         string
	SomeNiceField string
	PreCreateErr  bool
	PolicyErr     bool
	SetIdErr      bool
}

func (t *testMemoEntity) GetId() entity.Id {
	return entity.Id(t.Id)
}

func (t *testMemoEntity) SetId(id entity.Id) error {
	if t.SetIdErr {
		return errors.New("error setting id")
	}
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
	if t.PreCreateErr {
		return entity.ErrEntityAlreadyExists
	}
	return nil
}

func (t *testMemoEntity) PreUpdate() error {
	return nil
}

type nilIdStrategy struct{}

func (n nilIdStrategy) Generate(k *testMemoEntity) entity.Id {
	return entity.Id("")
}

type nilPolicy struct{}

func (n nilPolicy) ApplyCreate(ctx context.Context, e *testMemoEntity, collection []*testMemoEntity) ([]*testMemoEntity, error) {
	if e.PolicyErr {
		return collection, entity.ErrEntityAlreadyExists
	}

	return append(collection, e), nil
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
			r:    NewRepository(nil, WithIdStrategy[*testMemoEntity](nilIdStrategy{})),
			ctx:  context.TODO(),
			calls: []*testMemoEntity{
				{Id: "", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "", Attr1: "attr2", SomeNiceField: "some_nice_field"},
				{Id: "", Attr1: "attr3", SomeNiceField: "some_nice_field"},
			},
			expect: []*testMemoEntity{
				{Id: "", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "", Attr1: "attr2", SomeNiceField: "some_nice_field"},
				{Id: "", Attr1: "attr3", SomeNiceField: "some_nice_field"},
			},
			wantErr: false,
		},
		{
			name: "Test Create with MRU policy",
			r: NewRepository(
				nil,
				WithIdStrategy[*testMemoEntity](nilIdStrategy{}),
				WithPolicy[*testMemoEntity](PolicyMru[*testMemoEntity]{Capacity: 2}),
			),
			ctx: context.TODO(),
			calls: []*testMemoEntity{
				{Id: "", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "", Attr1: "attr2", SomeNiceField: "some_nice_field"},
				{Id: "", Attr1: "attr3", SomeNiceField: "some_nice_field"},
			},
			expect: []*testMemoEntity{
				{Id: "", Attr1: "attr2", SomeNiceField: "some_nice_field"},
				{Id: "", Attr1: "attr3", SomeNiceField: "some_nice_field"},
			},
			wantErr: false,
		},
		{
			name: "Test Create with LRU policy",
			r: NewRepository(
				nil,
				WithIdStrategy[*testMemoEntity](nilIdStrategy{}),
				WithPolicy[*testMemoEntity](PolicyLru[*testMemoEntity]{Capacity: 2}),
			),
			ctx: context.TODO(),
			calls: []*testMemoEntity{
				{Id: "", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "", Attr1: "attr2", SomeNiceField: "some_nice_field"},
				{Id: "", Attr1: "attr3", SomeNiceField: "some_nice_field"},
			},
			expect: []*testMemoEntity{
				{Id: "", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "", Attr1: "attr2", SomeNiceField: "some_nice_field"},
			},
			wantErr: false,
		},
		{
			name: "Test Create already exists",
			r:    NewRepository([]*testMemoEntity{{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"}}),
			ctx:  context.TODO(),
			calls: []*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
			},
			expect: []*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
			},
			wantErr: true,
		},
		{
			name: "Test Create with PreCreate error",
			r:    NewRepository([]*testMemoEntity{}),
			ctx:  context.TODO(),
			calls: []*testMemoEntity{
				{Id: "", Attr1: "attr1", PreCreateErr: true},
			},
			expect:  []*testMemoEntity{},
			wantErr: true,
		},
		{
			name: "Test Create with Policy error",
			r:    NewRepository([]*testMemoEntity{}, WithIdStrategy[*testMemoEntity](nilIdStrategy{}), WithPolicy[*testMemoEntity](nilPolicy{})),
			ctx:  context.TODO(),
			calls: []*testMemoEntity{
				{Id: "", Attr1: "attr1", PolicyErr: true},
			},
			expect:  []*testMemoEntity{},
			wantErr: true,
		},
		{
			name: "Test Create with SetId error",
			r:    NewRepository([]*testMemoEntity{}, WithIdStrategy[*testMemoEntity](nilIdStrategy{})),
			ctx:  context.TODO(),
			calls: []*testMemoEntity{
				{Id: "", Attr1: "attr1", SetIdErr: true},
			},
			expect:  []*testMemoEntity{},
			wantErr: true,
		},
		{
			name: "Test Create with IdStrategy null and SetId error",
			r:    NewRepository([]*testMemoEntity{}),
			ctx:  context.TODO(),
			calls: []*testMemoEntity{
				{Id: "", Attr1: "attr1", SetIdErr: true},
			},
			expect:  []*testMemoEntity{},
			wantErr: true,
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

func TestRepository_Read(t *testing.T) {
	type testCase[K entity.Entity] struct {
		name    string
		r       Repository[K]
		ctx     context.Context
		id      entity.Id
		wantE   K
		wantErr bool
	}

	tests := []testCase[*testMemoEntity]{
		{
			name:    "Test Read",
			r:       *NewRepository([]*testMemoEntity{{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"}}, WithIdStrategy[*testMemoEntity](nilIdStrategy{})),
			ctx:     context.TODO(),
			id:      entity.Id("1"),
			wantE:   &testMemoEntity{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
			wantErr: false,
		},
		{
			name:    "Test Read not found",
			r:       *NewRepository([]*testMemoEntity{{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"}}, WithIdStrategy[*testMemoEntity](nilIdStrategy{})),
			ctx:     context.TODO(),
			id:      entity.Id("2"),
			wantE:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := tt.r.Read(tt.ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotE, tt.wantE) {
				t.Errorf("Read() gotE = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestRepository_Match(t *testing.T) {
	type testCase[K entity.Entity] struct {
		name    string
		r       Repository[K]
		ctx     context.Context
		c       specification.Criteria
		want    []K
		wantErr bool
	}
	tests := []testCase[*testMemoEntity]{
		{
			name: "Test Match",
			r: *NewRepository([]*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "2", Attr1: "attr2", SomeNiceField: "some_nice_field"},
			},
				WithIdStrategy[*testMemoEntity](nilIdStrategy{})),
			ctx: context.TODO(),
			c: specification.Attr{
				Name:       "Attr1",
				Value:      "attr2",
				Comparison: specification.ComparisonEq,
			},
			want: []*testMemoEntity{
				{Id: "2", Attr1: "attr2", SomeNiceField: "some_nice_field"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Match(tt.ctx, tt.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Match() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_MatchOne(t *testing.T) {
	type testCase[K entity.Entity] struct {
		name    string
		r       Repository[K]
		ctx     context.Context
		c       specification.Criteria
		wantK   K
		wantErr bool
	}
	tests := []testCase[*testMemoEntity]{
		{
			name: "Test Match One",
			r: *NewRepository([]*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "2", Attr1: "attr2", SomeNiceField: "some_nice_field"},
				{Id: "3", Attr1: "attr3", SomeNiceField: "some_nice_field"},
			},
				WithIdStrategy[*testMemoEntity](nilIdStrategy{})),
			ctx: context.TODO(),
			c: specification.Attr{
				Name:       "SomeNiceField",
				Value:      "some_nice_field",
				Comparison: specification.ComparisonEq,
			},
			wantK:   &testMemoEntity{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
			wantErr: false,
		},
		{
			name: "Test Match One not found",
			r: *NewRepository([]*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
			},
				WithIdStrategy[*testMemoEntity](nilIdStrategy{})),
			ctx: context.TODO(),
			c: specification.Attr{
				Name:       "SomeNiceField",
				Value:      "some_nice_field2",
				Comparison: specification.ComparisonEq,
			},
			wantK:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotK, err := tt.r.MatchOne(tt.ctx, tt.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("MatchOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotK, tt.wantK) {
				t.Errorf("MatchOne() gotK = %v, want %v", gotK, tt.wantK)
			}
		})
	}
}

func TestRepository_ReadAll(t *testing.T) {
	type testCase[K entity.Entity] struct {
		name    string
		r       Repository[K]
		ctx     context.Context
		want    []K
		wantErr bool
	}
	tests := []testCase[*testMemoEntity]{
		{
			name: "Test Read All",
			r: *NewRepository([]*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "2", Attr1: "attr2", SomeNiceField: "some_nice_field"},
			},
				WithIdStrategy[*testMemoEntity](nilIdStrategy{})),
			ctx: context.TODO(),
			want: []*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
				{Id: "2", Attr1: "attr2", SomeNiceField: "some_nice_field"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.ReadAll(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	type testCase[K entity.Entity] struct {
		name    string
		r       Repository[K]
		ctx     context.Context
		entity  K
		wantErr bool
		read    K
	}
	tests := []testCase[*testMemoEntity]{
		{
			name: "Test Update",
			r: *NewRepository([]*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
			}, WithIdStrategy[*testMemoEntity](nilIdStrategy{})),
			read:    &testMemoEntity{Id: "1", Attr1: "attr2", SomeNiceField: "some_nice_field2"},
			ctx:     context.TODO(),
			entity:  &testMemoEntity{Id: "1", Attr1: "attr2", SomeNiceField: "some_nice_field2"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Update(tt.ctx, tt.entity); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got, _ := tt.r.Read(tt.ctx, tt.entity.GetId()); !reflect.DeepEqual(got, tt.read) {
				t.Errorf("Read() got = %v, want %v", got, tt.read)
			}
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	type testCase[K entity.Entity] struct {
		name    string
		r       Repository[K]
		ctx     context.Context
		entity  K
		read    K
		wantErr bool
	}
	tests := []testCase[*testMemoEntity]{
		{
			name: "Test Delete",
			r: *NewRepository([]*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
			}, WithIdStrategy[*testMemoEntity](nilIdStrategy{})),
			ctx:     context.TODO(),
			entity:  &testMemoEntity{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
			wantErr: false,
			read:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Delete(tt.ctx, tt.entity); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got, _ := tt.r.Read(tt.ctx, tt.entity.GetId()); !reflect.DeepEqual(got, tt.read) {
				t.Errorf("Read() got = %v, want %v", got, tt.read)
			}
		})
	}
}

func TestRepository_Start(t *testing.T) {
	type testCase[K entity.Entity] struct {
		name        string
		r           Repository[K]
		ctx         context.Context
		onBootstrap func(ctx context.Context) error
		wantErr     error
	}
	tests := []testCase[*testMemoEntity]{
		{
			name: "Test Start",
			r: *NewRepository([]*testMemoEntity{
				{Id: "1", Attr1: "attr1", SomeNiceField: "some_nice_field"},
			},
				WithIdStrategy[*testMemoEntity](nilIdStrategy{})),
			ctx: context.TODO(),
			onBootstrap: func(ctx context.Context) error {
				return errors.New("error1234")
			},
			wantErr: errors.New("error1234"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Start(tt.ctx, tt.onBootstrap); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
