// @Author : detaohe
// @File   : resp.go
// @Description:
// @Date   : 2022/4/26 8:27 PM

package util

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhT "github.com/go-playground/validator/v10/translations/zh"
	"net/http"
)

const (
	ContentType = "application/json;charset=utf-8"
)

var Translator ut.Translator
var DefaultValidator *validator.Validate

func TranslateValidator() {
	DefaultValidator = validator.New()
	uni := ut.New(en.New(), zh.New())
	Translator, _ = uni.GetTranslator("zh")
	_ = zhT.RegisterDefaultTranslations(DefaultValidator, Translator)
}

//func EncodeHttpResp(_ context.Context, w http.ResponseWriter, response interface{}) error {
//	encoder := json.NewEncoder(w)
//	if err := encoder.Encode(response); err != nil {
//		return err
//	}
//	return nil
//}

type GinResponse struct {
	c *gin.Context
}

func NewGinResponse(c *gin.Context) *GinResponse {
	return &GinResponse{
		c: c,
	}
}

func (r *GinResponse) ToResp(data interface{}) {
	if data == nil {
		r.c.JSON(http.StatusOK, nil)
	} else {
		response := gin.H{"status": "success", "data": data}
		r.c.JSON(http.StatusOK, response)
	}
}
func (r *GinResponse) ToError(err error, msg ...interface{}) {
	//response := gin.H{"status": "error", "msg": msg}
	switch err.(type) {
	case validator.ValidationErrors:
		_validate, _ := err.(validator.ValidationErrors)
		errMap := _validate.Translate(Translator)

		var _msg string
		for k, v := range errMap {
			_msg += k + ":" + v + "\n"
		}
		response := gin.H{"status": "error", "data": gin.H{"msg": _msg}}
		r.c.JSON(http.StatusBadRequest, response)
	default:
		response := gin.H{"status": "error", "data": gin.H{"msg": msg}}
		r.c.JSON(http.StatusInternalServerError, response)

	}
}
