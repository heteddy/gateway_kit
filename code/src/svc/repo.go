// @Author : detaohe
// @File   : repo
// @Description:
// @Date   : 2022/8/30 18:23

package svc

type Repo struct {
	// 通过redis或者数据库获取
	// todo 暂时hardcode
	repo map[string][]string
}

func NewServiceRepo() *Repo {
	repo := make(map[string][]string)
	repo["server"] = []string{
		"192.168.64.7:9192",
		"192.168.64.7:9193",
	}
	return &Repo{
		repo: repo,
	}
}

func (repo *Repo) GetServices(name string) ([]string, error) {
	return repo.repo[name], nil
}
