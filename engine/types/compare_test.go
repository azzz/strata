package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValueCompare_NullValuesAreEqual(t *testing.T) {
	t.Parallel()

	result, err := NewNullValue().Compare(NewNullValue())
	require.NoError(t, err)
	assert.Equal(t, Equal, result)
}
