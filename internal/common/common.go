package common

import (
	"scsPro/internal/model"
	"time"
)

// Module 定义导航或模块项的结构
type Module struct {
	Name        string
	Description string
	Link        string
}

// NavItems 全局导航项
var NavItems = []Module{
	{Name: "首页", Link: "/"},
	{Name: "博客", Link: "/blog"},
	{Name: "技能", Link: "/skills"},
	{Name: "项目", Link: "/projects"},
	{Name: "爱好", Link: "/hobbies"},
	{Name: "时间线", Link: "/timeline"},
	{Name: "资源", Link: "/resources"},
	{Name: "关于", Link: "/about"},
	{Name: "联系我", Link: "/contact"},
}

// GetCommonData 获取页面公共数据
func GetCommonData() (map[string]interface{}, error) {
	popularArticles, err := model.GetPopularArticles(10)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"NavItems":        NavItems,
		"PopularArticles": popularArticles,
		"UserStatus":      "未登录",
		"CurrentTime":     time.Now(),
	}, nil
}