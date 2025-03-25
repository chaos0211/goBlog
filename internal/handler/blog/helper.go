package blog

import (
	"html/template"
	"github.com/gin-gonic/gin"
)

func RegisterTemplateFuncs(r *gin.Engine) {
	r.SetFuncMap(template.FuncMap{
		"sub": sub,
		"add": add,
	})
}

func sub(a, b int) int {
	return a - b
}

func add(a, b int) int {
	return a + b
}