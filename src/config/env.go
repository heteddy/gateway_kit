// @Author : detaohe
// @File   : env.go
// @Description:
// @Date   : 2022/9/4 20:52

package config

func initEnvironments() {
	Environments = make(map[string]string, 10)
	Environments["mongo.host"] = "MONGO_HOST"
	Environments["mongo.port"] = "MONGO_PORT"
	Environments["mongo.user"] = "MONGO_USER"
	Environments["mongo.dbName"] = "MONGO_DBNAME"
	Environments["mongo.replicaSet"] = "MONGO_REPLICASET"
	//Environments["mongo.mongo"] = "MONGO_HOST"
}
