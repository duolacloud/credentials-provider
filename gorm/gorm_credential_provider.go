package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/duolacloud/credentials-provider"
	gorm_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/repositories"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// gorm 凭证管理配置项
type GormCredentialProviderOptions struct {
	cache         cache.Cache                          // 缓存策略
	cacheRepoOpts []repositories.CacheRepositoryOption // 缓存选项
	gormRepoOpts  []gorm_repo.GormCrudRepositoryOption // repository 的配置项
}

type GormCredentialProviderOption func(*GormCredentialProviderOptions)

// 为凭证管理添加缓存
func WithCache(c cache.Cache) GormCredentialProviderOption {
	return func(gcpo *GormCredentialProviderOptions) {
		gcpo.cache = c
	}
}

// gorm repository 的配置项
func WithGormRepositoryOptions(
	opts ...gorm_repo.GormCrudRepositoryOption,
) GormCredentialProviderOption {
	return func(gcpo *GormCredentialProviderOptions) {
		gcpo.gormRepoOpts = opts
	}
}

// cache repository 的配置项
func WithCacheRepositoryOptions(
	opts ...repositories.CacheRepositoryOption,
) GormCredentialProviderOption {
	return func(gcpo *GormCredentialProviderOptions) {
		gcpo.cacheRepoOpts = opts
	}
}

// 创建 gorm 凭证管理
func NewGormCredentialProvider(
	db *gorm.DB,
	opts ...GormCredentialProviderOption,
) credentials.CredentialProvider {

	provider := &GormCredentialProvider{
		options: &GormCredentialProviderOptions{
			gormRepoOpts: []gorm_repo.GormCrudRepositoryOption{},
			cacheRepoOpts: []repositories.CacheRepositoryOption{
				// 默认缓存时间，超时后清除，再次查询重建缓存
				repositories.WithExpiration(12 * time.Hour),
			},
		},
	}

	for _, opt := range opts {
		opt(provider.options)
	}

	gormRepo := gorm_repo.NewGormCrudRepository[Credential, Credential, map[string]any](
		db,
		provider.options.gormRepoOpts...,
	)

	if provider.options.cache != nil {
		cacheRepo := repositories.NewCacheRepository[Credential, Credential, map[string]any](
			gormRepo,
			provider.options.cache,
			provider.options.cacheRepoOpts...,
		)
		provider.repo = cacheRepo
	} else {
		provider.repo = gormRepo
	}

	return provider
}

// 基于 gorm 的凭证管理
type GormCredentialProvider struct {
	repo    repositories.CrudRepository[Credential, Credential, map[string]any] // 凭证的 repository
	options *GormCredentialProviderOptions                                      // 配置
}

// 设置凭证
func (p *GormCredentialProvider) Set(ctx context.Context, app, key, credentialType string, options map[string]any) error {
	// 数据库是 id + type 的联合主键
	id := fmt.Sprintf("%s|%s", app, key)
	primaryKeys := map[string]any{
		"id":   id,
		"type": credentialType,
	}

	credential, err := p.repo.Get(ctx, primaryKeys)
	if err != nil {
		return err
	}

	if credential == nil {
		// TODO 并发时可能会重复创建，这里先直接返回错误简单处理
		_, err = p.repo.Create(ctx, &Credential{
			ID:      id,
			Type:    credentialType,
			Options: options,
		})
	} else {
		_, err = p.repo.Update(ctx, primaryKeys, &map[string]any{
			"options": datatypes.JSONMap(options),
		})
	}
	return err
}

// 查询凭证
func (p *GormCredentialProvider) Get(ctx context.Context, app, key, credentialType string) (options map[string]any, err error) {
	// 数据库是 id + type 的联合主键
	id := fmt.Sprintf("%s|%s", app, key)
	primaryKeys := map[string]any{
		"id":   id,
		"type": credentialType,
	}

	credential, err := p.repo.Get(ctx, primaryKeys)
	if err != nil {
		return nil, err
	}
	if credential == nil {
		return nil, nil
	}
	return credential.Options, nil
}
