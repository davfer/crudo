package inmemory

import (
	"testing"

	"github.com/davfer/crudo/entity"
)

func TestUuidIdStrategy_Generate(t *testing.T) {
	t.Run("Test uuid", func(t *testing.T) {
		strategy := UuidIdStrategy[entity.Entity]{}
		if got := strategy.Generate(&testMemoEntity{}); got.IsEmpty() {
			t.Errorf("Generate() = %v", got)
		}
	})
}
