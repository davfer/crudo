package inmemory

import (
	"github.com/davfer/crudo/entity"
	"testing"
)

func TestUuidIdStrategy_Generate(t *testing.T) {
	t.Run("Test uuid", func(t *testing.T) {
		strategy := UuidIdStrategy[entity.Entity]{}
		if got := strategy.Generate(&testMemoEntity{}); got.IsEmpty() {
			t.Errorf("Generate() = %v", got)
		}
	})
}
