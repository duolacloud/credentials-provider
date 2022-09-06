package credentials

import (
	"context"
	"strings"

	gorm_repo "github.com/duolacloud/crud-core-gorm/repositories"
	core_cache "github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/repositories"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CredentialProvider interface {
	// 设置凭证
	Set(ctx context.Context, app, key, credentialType string, options map[string]any) error

	// 获取凭证
	Get(ctx context.Context, app, key, credentialType string) (options map[string]any, err error)
}

type CacheCredentialProvider struct {
	repo repositories.CrudRepository[Credential, Credential, map[string]any]
}

func (p *CacheCredentialProvider) Set(ctx context.Context, app, key, credentialType string, options map[string]any) error {
	id := strings.Join([]string{app, key, credentialType}, "|")
	credential, err := p.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if credential == nil {
		_, err = p.repo.Create(ctx, &Credential{
			ID:      id,
			Options: options,
		})
	} else {
		_, err = p.repo.Update(ctx, id, &map[string]any{
			"options": datatypes.JSONMap(options),
		})
	}
	return err
}

func (p *CacheCredentialProvider) Get(ctx context.Context, app, key, credentialType string) (options map[string]any, err error) {
	id := strings.Join([]string{app, key, credentialType}, "|")
	credential, err := p.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if credential == nil {
		return nil, nil
	}
	return credential.Options, nil
}

func NewCacheCredentialProvider(db *gorm.DB, cache core_cache.Cache) CredentialProvider {
	return &CacheCredentialProvider{
		repo: repositories.NewCacheRepository[Credential, Credential, map[string]any](
			gorm_repo.NewGormCrudRepository[Credential, Credential, map[string]any](db),
			cache,
		),
	}
}
