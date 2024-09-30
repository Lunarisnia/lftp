package filesystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OpenFile(t *testing.T) {
	t.Run("Expect to keep reading until EOF", func(t *testing.T) {
		r, err := OpenFile("../../tests/6bytes", 3)
		assert.Nil(t, err)

		p := make([]byte, 3)
		r.Read(p)
		assert.Equal(t, "123", string(p))

		r.Read(p)
		assert.Equal(t, "456", string(p))
	})
	t.Run("Expect to ran out properly", func(t *testing.T) {
		r, err := OpenFile("../../tests/7bytes", 3)
		assert.Nil(t, err)

		p := make([]byte, 3)
		r.Read(p)
		assert.Equal(t, "123", string(p))

		r.Read(p)
		assert.Equal(t, "456", string(p))

		r.Read(p)
		assert.Equal(t, "7\n6", string(p))
	})
	t.Run("Expect to return 16 as the size because it is smaller than 16", func(t *testing.T) {
		r, err := OpenFile("../../tests/7bytes", 3)
		assert.Nil(t, err)
		assert.Equal(t, 16, r.Size())
	})
}
