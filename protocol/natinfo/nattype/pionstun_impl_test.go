package nattype

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPionStunImpl_GetNATType(t *testing.T) {
	t.Skip("Skip test for PionStunImpl_GetNATType")
	impl := &PionStunImpl{}

	result, err := impl.GetNATType()
	assert.NoError(t, err)

	assert.Equal(t, Symmetric, result)
}
