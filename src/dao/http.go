// @Author : detaohe
// @File   : http.go
// @Description:
// @Date   : 2022/9/6 21:34

package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type HttpRule struct {
	ID             *primitive.ObjectID `json:"id" `
	Name           string
	Description    string
	Category       int       `json:"category"  description:"匹配类型 domain=域名, url_prefix=url前缀"`
	Rule           string    `json:"rule"  description:"type=domain表示域名，type=url_prefix时表示url前缀"`
	IsHttps        int       `json:"need_https"  description:"type=支持https 1=支持"`
	IsWebsocket    int       `json:"need_websocket"  description:"启用websocket 1=启用"`
	IsStripUri     int       `json:"need_strip_uri"  description:"启用strip_uri 1=启用"`
	UrlRewrite     string    `json:"url_rewrite" description:"url重写功能，每行一个	"`
	HeaderTransfer string    `json:"header_transfer"  description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue	"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      int64     `json:"deleted_at"` // 删除时间
}
