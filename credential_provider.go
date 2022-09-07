package credentials

import (
	"context"
)

// 定义凭证管理的接口
type CredentialProvider interface {

	// 设置凭证
	// 参数 app 是凭证的应用，如：douyin、weibo 等
	// 参数 key 是凭证的归属，如：用户 id、应用 id 等
	// 参数 credentialType 表明了凭证的类型，如：access_token、password 等
	// 参数 options 是凭证的内容，由用户自定义，并根据凭证的类型解析
	Set(ctx context.Context, app, key, credentialType string, options map[string]any) error

	// 获取凭证内容
	Get(ctx context.Context, app, key, credentialType string) (options map[string]any, err error)
}
