package inmemory_test

import (
	"testing"

	"github.com/davfer/crudo/entity"
	"github.com/davfer/crudo/inmemory"
)

func TestUuidIdStrategy_Generate(t *testing.T) {
	t.Run("Test uuid", func(t *testing.T) {
		strategy := inmemory.UuidIdStrategy[entity.Entity]{}
		if got := strategy.Generate(&testMemoEntity{}); got.IsEmpty() {
			t.Errorf("Generate() = %v", got)
		}
	})
}
