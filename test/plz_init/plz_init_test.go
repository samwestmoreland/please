package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thought-machine/please/src/plzinit"
)

func TestPlzInit(t *testing.T) {
	options := map[string]string{"foo.bar": "bar"}
	plzinit.InitConfigFile("./test_config", options)
	_, err := os.Stat("test_config")
	require.NoError(t, err)
}
