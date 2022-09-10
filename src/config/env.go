// @Author : detaohe
// @File   : env.go
// @Description:
// @Date   : 2022/9/4 20:52

package config

func initEnvironments() {
	Environments = make(map[string]string, 10)
	Environments["mongodb.host"] = "MONGO_HOST"
	Environments["mongodb.port"] = "MONGO_PORT"
	Environments["mongodb.user"] = "MONGO_USER"
	Environments["mongodb.dbName"] = "MONGO_DBNAME"
	Environments["mongodb.replicaSet"] = "MONGO_REPLICASET"
	//Environments["mongodb.mongodb"] = "MONGO_HOST"
}
