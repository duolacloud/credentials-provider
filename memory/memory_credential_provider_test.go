package gorm

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMemoryCredentialProvider(t *testing.T) {
	provider := NewMemoryCredentialProvider()

	app := "douyin"
	clientId := uuid.NewString()
	credentialType := "password"
	t.Logf("clientId: %s", clientId)

	options, err := provider.Get(context.TODO(), app, clientId, credentialType)
	assert.Nil(t, options)
	assert.Nil(t, err)

	err = provider.Set(context.TODO(), app, clientId, credentialType, map[string]any{"username": "root", "password": "root"})
	assert.Nil(t, err)

	options, err = provider.Get(context.TODO(), app, clientId, credentialType)
	assert.Nil(t, err)
	assert.Equal(t, "root", options["password"])
	assert.Equal(t, "root", options["username"])

	err = provider.Set(context.TODO(), app, clientId, credentialType, map[string]any{"username": "root", "password": "secret"})
	assert.Nil(t, err)

	options, err = provider.Get(context.TODO(), app, clientId, credentialType)
	assert.Nil(t, err)
	assert.Equal(t, "root", options["username"])
	assert.Equal(t, "secret", options["password"])
}
