// @Author : detaohe
// @File   : repo
// @Description:
// @Date   : 2022/8/30 18:23

package svc

/*
提供两种方式注册服务，
1. 适用于k8s的直接调用gateway的接口，写入client信息，gateway写入数据库并同步到redis中
2. 提供一个sdk写入到etcd, gateway通过etcd获取client的信息
*/

type Repo struct {
	// 通过redis或者数据库获取
	// todo 暂时hardcode
	repo    map[string][]string
	addrSet map[string]struct{}
}

func NewServiceRepo() *Repo {
	repo := make(map[string][]string)
	repo["server"] = []string{
		"192.168.64.7:9192",
		"192.168.64.7:9193",
	}
	addrSet := make(map[string]struct{})
	addrSet["192.168.64.7:9192"] = struct{}{}
	addrSet["192.168.64.7:9193"] = struct{}{}
	return &Repo{
		repo:    repo,
		addrSet: addrSet,
	}
}

func (repo *Repo) GetServices(name string) ([]string, error) {
	return repo.repo[name], nil
}
func (repo *Repo) Add(name string, addr string) {

}
