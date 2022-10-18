// @Author : detaohe
// @File   : flow_test.go
// @Description:
// @Date   : 2022/10/18 17:28

package dao

import (
	"context"
	"gateway_kit/config"
	"gateway_kit/util/mongodb"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestServiceHourDao_GetServicesDetail(t *testing.T) {
	//hosts: 192.168.64.5:27017,192.168.64.6:27017,192.168.64.7:27017
	//user:
	//pass:
	//database: gateway
	//replica: rs0
	config.All.Tables.RequestHour = "tb_request_hour"
	config.All.Tables.ServiceHour = "tb_service_hour"
	//config.All.Tables.RequestHour = "tb_request_hour"
	c := mongodb.Config{
		Hosts:    []string{"192.168.64.5:27017", "192.168.64.6:27017", "192.168.64.7:27017"},
		Database: "gateway",
		Replica:  "rs0",
	}
	config.InitMongo(c, "debug")
	config.InitLogger("./", "gateway_kit")
	d := NewServiceHourDao()
	convey.Convey("GetServicesDetail", t, func() {
		d.GetServicesDetail(context.Background(), "后台测试服务64_5", time.Now().AddDate(0, 0, -1), time.Now())
	})
	//convey.Convey("get sum", t, func() {
	//	d.GetSum(context.Background(), time.Now().AddDate(0, 0, -1), time.Now())
	//})
}
