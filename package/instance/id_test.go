package instance

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)
import "math/rand"

func randID() ID {
	return ID{
		Name:     strconv.Itoa(rand.Int()),
		Instance: strconv.Itoa(rand.Int()),
	}
}

func TestID_Equals(t *testing.T) {
	for i := 0; i < 10; i++ {
		id1 := randID()
		id2 := ID{
			Name:     id1.Name,
			Instance: id1.Instance,
		}
		id3 := randID()
		assert.True(t, id1.equals(id2))
		assert.True(t, id2.equals(id1))
		assert.False(t, id1.equals(id3))
		assert.False(t, id3.equals(id1))
		assert.False(t, id2.equals(id3))
		assert.False(t, id3.equals(id2))
	}
}
