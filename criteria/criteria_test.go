package criteria

import (
	"math"
	"testing"
)

type testEntity struct {
	id            string
	Attr1         string
	someOtherAttr bool
	IntNice       int
	Floating      float64
}

func TestAnd_IsSatisfiedBy(t *testing.T) {
	tests := []struct {
		name     string
		operands []Criteria
		value    any
		want     bool
	}{
		{
			name:     "Test empty operands",
			operands: []Criteria{},
			value:    testEntity{},
			want:     true,
		},
		{
			name: "Test single operand true",
			operands: []Criteria{
				Attr{
					Name:       "Attr1",
					Value:      "test",
					Comparison: ComparisonEq,
				},
			},
			value: testEntity{
				Attr1: "test",
			},
			want: true,
		},
		{
			name: "Test single operand false",
			operands: []Criteria{
				Attr{
					Name:       "Attr1",
					Value:      "testa",
					Comparison: ComparisonEq,
				},
			},
			value: testEntity{
				Attr1: "test",
			},
			want: false,
		},
		{
			name: "Test single operand true, true",
			operands: []Criteria{
				Attr{
					Name:       "Attr1",
					Value:      "test",
					Comparison: ComparisonEq,
				},
				Attr{
					Name:       "IntNice",
					Value:      1,
					Comparison: ComparisonGt,
				},
			},
			value: testEntity{
				Attr1:   "test",
				IntNice: 2,
			},
			want: true,
		},
		{
			name: "Test single operand true, false",
			operands: []Criteria{
				Attr{
					Name:       "Attr1",
					Value:      "test",
					Comparison: ComparisonEq,
				},
				Attr{
					Name:       "IntNice",
					Value:      1,
					Comparison: ComparisonGt,
				},
			},
			value: testEntity{
				Attr1:   "test",
				IntNice: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := And{
				Operands: tt.operands,
			}
			if got := a.IsSatisfiedBy(tt.value); got != tt.want {
				t.Errorf("IsSatisfiedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttr_IsSatisfiedBy(t *testing.T) {
	type fields struct {
		Name       string
		Value      any
		Comparison Comparator
	}
	tests := []struct {
		name   string
		fields fields
		value  any
		want   bool
	}{
		{
			name: "Test eq true",
			fields: fields{
				Name:       "Floating",
				Value:      math.Pi,
				Comparison: ComparisonEq,
			},
			value: testEntity{
				Floating: math.Pi,
			},
			want: true,
		},
		{
			name: "Test eq false",
			fields: fields{
				Name:       "IntNice",
				Value:      1,
				Comparison: ComparisonEq,
			},
			value: testEntity{
				IntNice: 2,
			},
			want: false,
		},
		{
			name: "Test ne true",
			fields: fields{
				Name:       "IntNice",
				Value:      1,
				Comparison: ComparisonNe,
			},
			value: testEntity{
				IntNice: 2,
			},
			want: true,
		},
		{
			name: "Test ne false",
			fields: fields{
				Name:       "IntNice",
				Value:      1,
				Comparison: ComparisonNe,
			},
			value: testEntity{
				IntNice: 1,
			},
			want: false,
		},
		{
			name: "Test gt true",
			fields: fields{
				Name:       "Floating",
				Value:      1.0,
				Comparison: ComparisonGt,
			},
			value: testEntity{
				Floating: 1.1,
			},
			want: true,
		},
		{
			name: "Test gt string true",
			fields: fields{
				Name:       "Attr1",
				Value:      "a",
				Comparison: ComparisonGt,
			},
			value: testEntity{
				Attr1: "a",
			},
			want: false,
		},
		{
			name: "Test gt false",
			fields: fields{
				Name:       "Floating",
				Value:      1.0,
				Comparison: ComparisonGt,
			},
			value: testEntity{
				Floating: 1.0,
			},
		},
		{
			name: "Test gte true",
			fields: fields{
				Name:       "Floating",
				Value:      2.34,
				Comparison: ComparisonGte,
			},
			value: testEntity{
				Floating: 2.34,
			},
			want: true,
		},
		{
			name: "Test gte false",
			fields: fields{
				Name:       "IntNice",
				Value:      12,
				Comparison: ComparisonGte,
			},
			value: testEntity{
				IntNice: 1,
			},
			want: false,
		},

		{
			name: "Test gte false",
			fields: fields{
				Name:       "Attr1",
				Value:      "a",
				Comparison: ComparisonGte,
			},
			value: testEntity{
				Attr1: "a",
			},
			want: true,
		},
		{
			name: "Test lt true",
			fields: fields{
				Name:       "IntNice",
				Value:      12,
				Comparison: ComparisonLt,
			},
			value: testEntity{
				IntNice: 1,
			},
			want: true,
		},
		{
			name: "Test lt false",
			fields: fields{
				Name:       "IntNice",
				Value:      12,
				Comparison: ComparisonLt,
			},
			value: testEntity{
				IntNice: 12,
			},
		},
		{
			name: "Test lte true",
			fields: fields{
				Name:       "IntNice",
				Value:      12,
				Comparison: ComparisonLte,
			},
			value: testEntity{
				IntNice: 12,
			},
			want: true,
		},
		{
			name: "Test lte false",
			fields: fields{
				Name:       "IntNice",
				Value:      12,
				Comparison: ComparisonLte,
			},
			value: testEntity{
				IntNice: 13,
			},
			want: false,
		},
		{
			name: "Test non struct",
			fields: fields{
				Name:       "IntNice",
				Value:      12,
				Comparison: ComparisonLte,
			},
			value: 12,
			want:  false,
		},
		{
			name: "Test non existing field",
			fields: fields{
				Name:       "NonExisting",
				Value:      12,
				Comparison: ComparisonLte,
			},
			value: testEntity{},
			want:  false,
		},
		{
			name: "Test non existing comparison",
			fields: fields{
				Name:       "IntNice",
				Value:      12,
				Comparison: "NonExisting",
			},
			value: testEntity{},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Attr{
				Name:       tt.fields.Name,
				Value:      tt.fields.Value,
				Comparison: tt.fields.Comparison,
			}
			if got := a.IsSatisfiedBy(tt.value); got != tt.want {
				t.Errorf("IsSatisfiedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOr_IsSatisfiedBy(t *testing.T) {
	tests := []struct {
		name     string
		operands []Criteria
		value    any
		want     bool
	}{
		{
			name:     "Test empty operands",
			operands: []Criteria{},
			value:    testEntity{},
			want:     true,
		},
		{
			name: "Test single operand true",
			operands: []Criteria{
				Attr{
					Name:       "IntNice",
					Value:      42,
					Comparison: ComparisonGt,
				},
			},
			value: testEntity{
				IntNice: 50,
			},
			want: true,
		},
		{
			name: "Test single operand false",
			operands: []Criteria{
				Attr{
					Name:       "Attr1",
					Value:      "testa",
					Comparison: ComparisonEq,
				},
			},
			value: testEntity{
				Attr1: "test",
			},
			want: false,
		},
		{
			name: "Test double operand true, true",
			operands: []Criteria{
				Attr{
					Name:       "Attr1",
					Value:      "test",
					Comparison: ComparisonEq,
				},
				Attr{
					Name:       "IntNice",
					Value:      1,
					Comparison: ComparisonGt,
				},
			},
			value: testEntity{
				Attr1:   "test",
				IntNice: 2,
			},
			want: true,
		},
		{
			name: "Test double operand true, false",
			operands: []Criteria{
				Attr{
					Name:       "Attr1",
					Value:      "test",
					Comparison: ComparisonEq,
				},
				Attr{
					Name:       "IntNice",
					Value:      1,
					Comparison: ComparisonGt,
				},
			},
			value: testEntity{
				Attr1:   "test",
				IntNice: 1,
			},
			want: true,
		},
		{
			name: "Test double operand false, false",
			operands: []Criteria{
				Attr{
					Name:       "Attr1",
					Value:      "test",
					Comparison: ComparisonEq,
				},
				Attr{
					Name:       "IntNice",
					Value:      1,
					Comparison: ComparisonGt,
				},
			},
			value: testEntity{
				Attr1:   "testa",
				IntNice: 0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Or{
				Operands: tt.operands,
			}
			if got := o.IsSatisfiedBy(tt.value); got != tt.want {
				t.Errorf("IsSatisfiedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNot_IsSatisfiedBy(t *testing.T) {
	type fields struct {
		Operand Criteria
	}
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		operand Criteria
		value   any
		want    bool
	}{
		{
			name: "Test true",
			operand: Attr{
				Name:       "Attr1",
				Value:      "test",
				Comparison: ComparisonEq,
			},
			value: testEntity{
				Attr1: "testa",
			},
			want: true,
		},
		{
			name: "Test false",
			operand: Attr{
				Name:       "Attr1",
				Value:      "test",
				Comparison: ComparisonEq,
			},
			value: testEntity{
				Attr1: "test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Not{
				Operand: tt.operand,
			}
			if got := n.IsSatisfiedBy(tt.value); got != tt.want {
				t.Errorf("IsSatisfiedBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
