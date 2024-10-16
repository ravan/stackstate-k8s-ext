package sync

import (
	"github.com/ravan/stackstate-client/stackstate/receiver"
	"github.com/ravan/stackstate-k8s-ext/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSync(t *testing.T) {
	conf := getConfig(t)
	factory, err := Sync(&conf.Kubernetes)
	require.NoError(t, err)
	assert.Equal(t, 1, factory.GetComponentCount())

	sts := receiver.NewClient(&conf.StackState, &conf.Instance)
	err = sts.Send(factory)
	require.NoError(t, err)
}

func getConfig(t *testing.T) *config.Configuration {
	require.NoError(t, os.Setenv("CONFIG_FILE", "./conf.yaml"))
	c, err := config.GetConfig()
	require.NoError(t, err)
	return c
}
