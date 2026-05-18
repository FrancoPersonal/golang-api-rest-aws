package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	l := New()
	require.NotNil(t, l)
}
