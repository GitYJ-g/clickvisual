package permission

import (
	"io/ioutil"

	"github.com/gotomicro/ego/core/elog"
	"gopkg.in/yaml.v3"

	"github.com/clickvisual/clickvisual/api/internal/invoker"
)

type Service struct {
	c        *Config
	resource Resource
}

type Config struct {
	ResFilePath string
}

func New(c *Config) *Service {
	return &Service{
		c:        c,
		resource: loadMenuTree(c.ResFilePath),
	}
}

func loadMenuTree(filePath string) Resource {
	menu := Resource{}

	resourceContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		invoker.Logger.Panic("Read Resource File Failed", elog.String("err", err.Error()))
	}

	err = yaml.Unmarshal(resourceContent, &menu)
	if err != nil {
		invoker.Logger.Panic("Unmarshall %s failed: %s", elog.String("filePath", filePath), elog.String("err", err.Error()))
	}

	return menu
}

func (s *Service) AdminMenuList() []MenuTreeItem {
	return s.resource.Permission
}

func (s *Service) UserMenuList() []MenuTreeItem {
	res := make([]MenuTreeItem, 0)
	for _, p := range s.resource.Permission {
		if p.Name == "log" || p.Name == "alarm" {
			res = append(res, p)
		}
	}
	return res
}
