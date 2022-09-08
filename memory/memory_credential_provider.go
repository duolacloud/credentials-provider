package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/duolacloud/credentials-provider"
)

// 内存凭证管理配置项
type MemoryCredentialProviderOptions struct {
}

type MemoryCredentialProviderOption func(*MemoryCredentialProviderOptions)

// 创建内存凭证管理
func NewMemoryCredentialProvider(opts ...MemoryCredentialProviderOption) credentials.CredentialProvider {
	provider := &MemoryCredentialProvider{
		options:        &MemoryCredentialProviderOptions{},
		credentialsMap: make(map[string]map[string]any),
	}
	for _, opt := range opts {
		opt(provider.options)
	}
	return provider
}

// 基于 gorm 的凭证管理
type MemoryCredentialProvider struct {
	options        *MemoryCredentialProviderOptions // 配置
	credentialsMap map[string]map[string]any        // 凭证内容
	lock           sync.RWMutex                     // 并发锁
}

// 设置凭证
func (p *MemoryCredentialProvider) Set(ctx context.Context, app, key, credentialType string, options map[string]any) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	cacheKey := strings.Join([]string{app, key, credentialType}, "|")
	p.credentialsMap[cacheKey] = options
	return nil
}

// 查询凭证
func (p *MemoryCredentialProvider) Get(ctx context.Context, app, key, credentialType string) (options map[string]any, err error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	cacheKey := strings.Join([]string{app, key, credentialType}, "|")
	if options, ok := p.credentialsMap[cacheKey]; ok {
		return options, nil
	}
	return nil, nil
}
