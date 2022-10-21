// @Author : detaohe
// @File   : jwt
// @Description:
// @Date   : 2022/10/20 20:28

package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"gateway_kit/config"
	"gateway_kit/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

func ExtractJWTBody(raw string) (string, error) {
	ss := strings.Split(raw, ".")
	if len(ss) == 3 {
		return ss[1], nil
	}
	return "", errors.New("jwt format not supported")
}

func DecodeJWTBody(body string) (string, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(body)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}
		if body, err := ExtractJWTBody(token); err != nil {
			util.NewGinResponse(c).ToError(err, "认证失败")
			return
		} else {
			if decoded, err := DecodeJWTBody(body); err != nil {
				util.NewGinResponse(c).ToError(err, "认证失败")
				return
			} else {
				authInfo := make(map[string]interface{})
				if err := json.Unmarshal([]byte(decoded), &authInfo); err != nil {
					util.NewGinResponse(c).ToError(err, "认证失败")
					return
				}
				for k, v := range authInfo {
					c.Set(k, v)
					// 检查是否为string类型
					if k == util.JwtKeyUserID {
						if userid, existed := c.Request.Header[util.JwtKeyUserID]; !existed {
							c.Request.Header.Set(k, v.(string))
						} else {
							config.Logger.Info("userid has existed", zap.Strings(k, userid))
							c.Request.Header.Add(k, v.(string))
						}

					}

				}
				c.Next()
			}
		}
	}
}
