package credentials

import (
	"time"

	"gorm.io/datatypes"
)

type Credential struct {
	ID        string            `gorm:"primaryKey"` // 凭证ID，用户自定义设置
	Options   datatypes.JSONMap // 凭证内容
	CreatedAt *time.Time        // 创建时间
	UpdatedAt time.Time         // 更新时间
}
