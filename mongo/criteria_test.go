package mongo

import (
	"github.com/davfer/crudo/entity"
	"github.com/davfer/go-specification"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
)

type testUnsupportedCriteria struct {
}

type testMongerEntity struct {
	Id            string `bson:"_id,omitempty"`
	Attr1         string `bson:"attr_1,omitempty"`
	SomeNiceField string `bson:"another_col_name"`
	NotBsoned     string
}

func (t *testMongerEntity) GetId() entity.Id {
	return entity.Id("")
}

func (t *testMongerEntity) SetId(id entity.Id) error {
	t.Id = string(id)
	return nil
}

func (t *testMongerEntity) GetResourceId() (string, error) {
	return t.Attr1, nil
}

func (t *testMongerEntity) SetResourceId(s string) error {
	t.Attr1 = s
	return nil
}

func (t *testMongerEntity) PreCreate() error {
	return nil
}

func (t *testMongerEntity) PreUpdate() error {
	return nil
}

func (t *testUnsupportedCriteria) IsSatisfiedBy(c any) bool {
	return true
}

func TestConvertToMongoCriteria(t *testing.T) {
	tests := []struct {
		name    string
		crit    specification.Criteria
		want    Criteria
		wantErr bool
	}{
		{
			name: "Test And criteria",
			crit: specification.And{
				Operands: []specification.Criteria{
					specification.Attr{
						Name:       "Attr1",
						Value:      12,
						Comparison: specification.ComparisonEq,
					},
				},
			},
			want: MongoAnd{
				Operands: []Criteria{
					Attr{
						Name:       "attr_1",
						Value:      12,
						Comparison: specification.ComparisonEq,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test Or criteria",
			crit: specification.Or{
				Operands: []specification.Criteria{
					specification.Attr{
						Name:       "SomeNiceField",
						Value:      12,
						Comparison: specification.ComparisonEq,
					},
				},
			},
			want: MongoOr{
				Operands: []Criteria{
					Attr{
						Name:       "another_col_name",
						Value:      12,
						Comparison: specification.ComparisonEq,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Test Not criteria",
			crit: specification.Not{
				Operand: specification.Attr{
					Name:       "SomeNiceField",
					Value:      12,
					Comparison: specification.ComparisonEq,
				},
			},
			want: MongoNot{
				Operand: Attr{
					Name:       "another_col_name",
					Value:      12,
					Comparison: specification.ComparisonEq,
				},
			},
			wantErr: false,
		},
		{
			name: "Test attr criteria",
			crit: specification.Attr{
				Name:       "Attr1",
				Value:      12,
				Comparison: specification.ComparisonEq,
			},
			want: Attr{
				Name:       "attr_1",
				Value:      12,
				Comparison: specification.ComparisonEq,
			},
			wantErr: false,
		},
		{
			name:    "Test unsupported criteria",
			crit:    &testUnsupportedCriteria{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test attr criteria field not found",
			crit:    specification.Attr{Name: "Attr2", Value: 12, Comparison: specification.ComparisonEq},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test attr criteria field not bson",
			crit:    specification.Attr{Name: "NotBsoned", Value: 12, Comparison: specification.ComparisonEq},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test and error propagation",
			crit:    specification.And{Operands: []specification.Criteria{&testUnsupportedCriteria{}}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test or error propagation",
			crit:    specification.Or{Operands: []specification.Criteria{&testUnsupportedCriteria{}}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test not error propagation",
			crit:    specification.Not{Operand: &testUnsupportedCriteria{}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToMongoCriteria(tt.crit, &testMongerEntity{Attr1: "asd"})
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToMongoCriteria() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToMongoCriteria() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoAnd_GetExpression(t *testing.T) {
	tests := []struct {
		name     string
		operands []Criteria
		want     bson.M
	}{
		{
			name: "Test one operand",
			operands: []Criteria{
				Attr{
					Name:       "some_column",
					Value:      12,
					Comparison: specification.ComparisonEq,
				},
			},
			want: bson.M{"$and": []bson.M{{"some_column": bson.M{"$eq": 12}}}},
		},
		{
			name: "Test two operands",
			operands: []Criteria{
				Attr{
					Name:       "some_column",
					Value:      12,
					Comparison: specification.ComparisonEq,
				},
				Attr{
					Name:       "some_column",
					Value:      12,
					Comparison: specification.ComparisonEq,
				},
			},
			want: bson.M{"$and": []bson.M{{"some_column": bson.M{"$eq": 12}}, {"some_column": bson.M{"$eq": 12}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := MongoAnd{
				Operands: tt.operands,
			}
			if got := a.GetExpression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoAttr_GetExpression(t *testing.T) {
	tests := []struct {
		name       string
		attr       string
		value      any
		comparison specification.Comparator
		want       bson.M
	}{
		{
			name:       "Test eq",
			attr:       "some_column",
			value:      12,
			comparison: specification.ComparisonEq,
			want:       bson.M{"some_column": bson.M{"$eq": 12}},
		},
		{
			name:       "Test gt",
			attr:       "some_column",
			value:      12,
			comparison: specification.ComparisonGt,
			want:       bson.M{"some_column": bson.M{"$gt": 12}},
		},
		{
			name:       "Test gte",
			attr:       "some_column",
			value:      12,
			comparison: specification.ComparisonGte,
			want:       bson.M{"some_column": bson.M{"$gte": 12}},
		},
		{
			name:       "Test lt",
			attr:       "some_column",
			value:      12,
			comparison: specification.ComparisonLt,
			want:       bson.M{"some_column": bson.M{"$lt": 12}},
		},
		{
			name:       "Test lte",
			attr:       "some_column",
			value:      12,
			comparison: specification.ComparisonLte,
			want:       bson.M{"some_column": bson.M{"$lte": 12}},
		},
		{
			name:       "Test ne",
			attr:       "some_column",
			value:      12,
			comparison: specification.ComparisonNe,
			want:       bson.M{"some_column": bson.M{"$ne": 12}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Attr{
				Name:       tt.attr,
				Value:      tt.value,
				Comparison: tt.comparison,
			}
			if got := a.GetExpression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoNot_GetExpression(t *testing.T) {
	tests := []struct {
		name    string
		operand Criteria
		want    bson.M
	}{
		{
			name: "Test operand",
			operand: Attr{
				Name:       "some_column",
				Value:      12,
				Comparison: specification.ComparisonEq,
			},
			want: bson.M{"$not": bson.M{"some_column": bson.M{"$eq": 12}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := MongoNot{
				Operand: tt.operand,
			}
			if got := n.GetExpression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoOr_GetExpression(t *testing.T) {
	tests := []struct {
		name     string
		operands []Criteria
		want     bson.M
	}{
		{
			name: "Test one operand",
			operands: []Criteria{
				Attr{
					Name:       "some_column",
					Value:      12,
					Comparison: specification.ComparisonEq,
				},
			},
			want: bson.M{"$or": []bson.M{{"some_column": bson.M{"$eq": 12}}}},
		},
		{
			name: "Test two operands",
			operands: []Criteria{
				Attr{
					Name:       "some_column",
					Value:      12,
					Comparison: specification.ComparisonEq,
				},
				Attr{
					Name:       "another_column",
					Value:      "howard",
					Comparison: specification.ComparisonNe,
				},
			},
			want: bson.M{"$or": []bson.M{{"some_column": bson.M{"$eq": 12}}, {"another_column": bson.M{"$ne": "howard"}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := MongoOr{
				Operands: tt.operands,
			}
			if got := o.GetExpression(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
