// @Author : detaohe
// @File   : resp.go
// @Description:
// @Date   : 2022/4/26 8:27 PM

package util

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	ContentType = "application/json;charset=utf-8"
)

func EncodeResp(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", ContentType)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		return err
	}
	return nil
}
