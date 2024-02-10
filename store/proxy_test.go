package store

import (
	"context"
	"github.com/davfer/crudo"
	"github.com/davfer/crudo/criteria"
	"github.com/davfer/crudo/entity"
	"github.com/davfer/crudo/notifier"
	"reflect"
	"testing"
)

type testProxyEntity struct {
	Id            string
	Attr1         string
	SomeNiceField string
}

func (t *testProxyEntity) GetId() entity.Id {
	return entity.Id(t.Id)
}

func (t *testProxyEntity) SetId(id entity.Id) error {
	t.Id = string(id)
	return nil
}

func (t *testProxyEntity) GetResourceId() (string, error) {
	return t.Attr1, nil
}

func (t *testProxyEntity) SetResourceId(s string) error {
	t.Attr1 = s
	return nil
}

func (t *testProxyEntity) PreCreate() error {
	return nil
}

func (t *testProxyEntity) PreUpdate() error {
	return nil
}

type spyRepository[K entity.Entity] struct {
	entities []K
	calls    []string
}

func (s *spyRepository[K]) Start(ctx context.Context, onBootstrap func(context.Context) error) error {
	s.calls = append(s.calls, "Start")
	s.entities = []K{}
	return nil
}

func (s *spyRepository[K]) Create(ctx context.Context, e K) (K, error) {
	e.SetId(entity.Id("attr1"))
	s.calls = append(s.calls, "Create")
	s.entities = append(s.entities, e)
	return e, nil
}

func (s *spyRepository[K]) Read(ctx context.Context, id entity.Id) (K, error) {
	s.calls = append(s.calls, "Read")
	return s.entities[0], nil
}

func (s *spyRepository[K]) Update(ctx context.Context, e K) error {
	s.calls = append(s.calls, "Update")
	return nil
}

func (s *spyRepository[K]) Delete(ctx context.Context, e K) error {
	s.calls = append(s.calls, "Delete")
	return nil
}

func (s *spyRepository[K]) Match(ctx context.Context, c criteria.Criteria) ([]K, error) {
	s.calls = append(s.calls, "Match")
	return s.entities, nil
}

func (s *spyRepository[K]) MatchOne(ctx context.Context, c criteria.Criteria) (K, error) {
	s.calls = append(s.calls, "MatchOne")
	return s.entities[0], nil
}

func (s *spyRepository[K]) ReadAll(ctx context.Context) ([]K, error) {
	s.calls = append(s.calls, "ReadAll")
	return s.entities, nil
}

func TestProxyStore_Create(t *testing.T) {
	type testCase[K entity.Entity] struct {
		name    string
		store   *ProxyStore[K]
		ctx     context.Context
		item    K
		want    K
		wantErr bool
	}
	tests := []testCase[*testProxyEntity]{
		{
			name:    "Test Create",
			store:   NewProxyStore[*testProxyEntity](),
			ctx:     context.TODO(),
			item:    &testProxyEntity{Attr1: "attr1", SomeNiceField: "someNiceField"},
			want:    &testProxyEntity{Id: "attr1", Attr1: "attr1", SomeNiceField: "someNiceField"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spyRepo := &spyRepository[*testProxyEntity]{}

			// load repo
			tt.store.Load(tt.ctx, spyRepo)

			// create entity
			got, err := tt.store.Create(tt.ctx, tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// check entity
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}

			// check remote repository
			if spyRepo.calls[1] != "Create" {
				t.Errorf("Create() got = %v, want %v", spyRepo.calls[0], "Create")
			}
		})
	}
}

func TestProxyStore_Delete(t *testing.T) {
	type testCase[K entity.Entity] struct {
		name    string
		store   *ProxyStore[K]
		ctx     context.Context
		item    K
		want    K
		wantErr bool
	}
	tests := []testCase[*testProxyEntity]{
		{
			name:    "Test Delete",
			store:   NewProxyStore[*testProxyEntity](),
			ctx:     context.TODO(),
			item:    &testProxyEntity{Attr1: "attr1", SomeNiceField: "someNiceField"},
			want:    &testProxyEntity{Id: "attr1", Attr1: "attr1", SomeNiceField: "someNiceField"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spyRepo := &spyRepository[*testProxyEntity]{}

			// load repo
			tt.store.Load(tt.ctx, spyRepo)

			// create entity
			err := tt.store.Delete(tt.ctx, tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// check remote repository
			if spyRepo.calls[1] != "Delete" {
				t.Errorf("Delete() got = %v, want %v", spyRepo.calls[0], "Delete")
			}
		})
	}
}

func TestProxyStore_Load(t *testing.T) {
	type args[K entity.Entity] struct {
		ctx  context.Context
		repo crudo.Repository[K]
	}
	type testCase[K entity.Entity] struct {
		name    string
		r       ProxyStore[K]
		args    args[K]
		wantErr bool
	}
	tests := []testCase[*testProxyEntity]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Load(tt.args.ctx, tt.args.repo); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProxyStore_Match(t *testing.T) {
	type args struct {
		ctx context.Context
		c   criteria.Criteria
	}
	type testCase[K entity.Entity] struct {
		name    string
		r       ProxyStore[K]
		args    args
		want    []K
		wantErr bool
	}
	tests := []testCase[*testProxyEntity]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Match(tt.args.ctx, tt.args.c)
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

func TestProxyStore_MatchOne(t *testing.T) {
	type args struct {
		ctx context.Context
		c   criteria.Criteria
	}
	type testCase[K entity.Entity] struct {
		name    string
		r       ProxyStore[K]
		args    args
		wantK   K
		wantErr bool
	}
	tests := []testCase[*testProxyEntity]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotK, err := tt.r.MatchOne(tt.args.ctx, tt.args.c)
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

func TestProxyStore_On(t *testing.T) {
	type args[K entity.Entity] struct {
		event    string
		observer notifier.ObserverCallback[K]
	}
	type testCase[K entity.Entity] struct {
		name    string
		r       ProxyStore[K]
		args    args[K]
		wantErr bool
	}
	tests := []testCase[*testProxyEntity]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.On(tt.args.event, tt.args.observer); (err != nil) != tt.wantErr {
				t.Errorf("On() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProxyStore_OnHydrate(t *testing.T) {
	type args[K entity.Entity] struct {
		onHydrate HydrateFunc[K]
	}
	type testCase[K entity.Entity] struct {
		name string
		r    ProxyStore[K]
		args args[K]
	}
	tests := []testCase[*testProxyEntity]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.OnHydrate(tt.args.onHydrate)
		})
	}
}

func TestProxyStore_Read(t *testing.T) {
	type args struct {
		ctx context.Context
		id  entity.Id
	}
	type testCase[K entity.Entity] struct {
		name    string
		r       ProxyStore[K]
		args    args
		wantE   K
		wantErr bool
	}
	tests := []testCase[*testProxyEntity]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotE, err := tt.r.Read(tt.args.ctx, tt.args.id)
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

func TestProxyStore_ReadAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type testCase[K entity.Entity] struct {
		name    string
		r       ProxyStore[K]
		args    args
		want    []K
		wantErr bool
	}
	tests := []testCase[*testProxyEntity]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.ReadAll(tt.args.ctx)
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

func TestProxyStore_Refresh(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type testCase[K entity.Entity] struct {
		name    string
		r       ProxyStore[K]
		args    args
		wantErr bool
	}
	tests := []testCase[*testProxyEntity]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Refresh(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Refresh() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProxyStore_Start(t *testing.T) {
	type args struct {
		ctx         context.Context
		onBootstrap func(ctx context.Context) error
	}
	type testCase[K entity.Entity] struct {
		name    string
		r       ProxyStore[K]
		args    args
		wantErr bool
	}
	tests := []testCase[*testProxyEntity]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Start(tt.args.ctx, tt.args.onBootstrap); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProxyStore_Update(t *testing.T) {
	type args[K entity.Entity] struct {
		ctx    context.Context
		entity K
	}
	type testCase[K entity.Entity] struct {
		name    string
		r       ProxyStore[K]
		args    args[K]
		wantErr bool
	}
	tests := []testCase[*testProxyEntity]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Update(tt.args.ctx, tt.args.entity); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
