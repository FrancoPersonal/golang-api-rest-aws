package clients

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDynamoDBClient(t *testing.T) {
	c := NewDynamoDBClient()
	require.NotNil(t, c)
}

func TestNewSQLServerClient(t *testing.T) {
	c := NewSQLServerClient()
	require.NotNil(t, c)
}
