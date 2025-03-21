package handler

import (
	"scsPro/internal/model"
	"time"
)

type Module struct {
	Name        string
	Description string
	Link        string
}

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