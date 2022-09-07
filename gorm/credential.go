package gorm

import (
	"time"

	"gorm.io/datatypes"
)

// 凭证的 gorm 实体定义
type Credential struct {
	ID        string            `gorm:"primaryKey"` // 凭证ID，由应用和 key 组成，例如： douyin|1234567
	Type      string            `gorm:"primaryKey"` // 凭证类型，例如：access_token、password
	Options   datatypes.JSONMap // 凭证内容
	CreatedAt time.Time         // 创建时间
	UpdatedAt time.Time         // 更新时间
}
