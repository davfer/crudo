package entity_test

import (
	"reflect"
	"testing"

	"github.com/davfer/crudo/entity"
)

func TestId_Equals(t *testing.T) {
	tests := []struct {
		name string
		i1   entity.ID
		i2   entity.ID
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
		i    entity.ID
		want map[string]entity.ID
	}{
		{
			name: "Test compound ids",
			i:    entity.ID(`{"i":"5f3e3e3e3e3e3e3e3e3e3e3e", "j":"5f3e3e3e3e3e3e3e3e3e3e3e"}`),
			want: map[string]entity.ID{
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
			if got := tt.i.GetCompoundIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCompoundIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestId_IsCompound(t *testing.T) {
	tests := []struct {
		name string
		i    entity.ID
		want bool
	}{
		{
			name: "Test compound",
			i:    entity.ID(`{"i":"5f3e3e3e3e3e3e3e3e3e3e3e", "j":"5f3e3e3e3e3e3e3e3e3e3e3e"}`),
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
		i    entity.ID
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
		i    entity.ID
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

func TestNewIdFromString(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want entity.ID
	}{
		{
			name: "Test new id from string",
			id:   "5f3e3e3e3e3e3e3e3e3e3e3e",
			want: "5f3e3e3e3e3e3e3e3e3e3e3e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := entity.NewIDFromString(tt.id); got != tt.want {
				t.Errorf("NewIDFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
