package gorm

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	redis_cache "github.com/duolacloud/crud-cache-redis"
	"github.com/duolacloud/crud-core/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestCredentialProvider(t *testing.T) {
	cache, _ := redis_cache.NewRedisCache()
	provider := NewGormCredentialProvider(
		setupDB(),
		WithCache(cache),
		WithCacheRepositoryOptions(
			repositories.WithExpiration(5*time.Second),
		),
	)

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

	time.Sleep(6 * time.Second)
	cachedOptions := make(map[string]any)
	cache.Get(context.TODO(), "", &cachedOptions)
	assert.Empty(t, cachedOptions)

	options, err = provider.Get(context.TODO(), app, clientId, credentialType)
	assert.Nil(t, err)
	assert.Equal(t, "root", options["username"])
	assert.Equal(t, "secret", options["password"])
}

func setupDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	dsn := "root:root@(localhost)/credentials_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, dberr := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if dberr != nil {
		panic(dberr)
	}

	dberr = db.AutoMigrate(&Credential{})
	if dberr != nil {
		panic(dberr)
	}

	return db
}
