// @Author : detaohe
// @File   : service.go
// @Description:
// @Date   : 2022/4/17 8:15 PM

package dao

type Client struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}
