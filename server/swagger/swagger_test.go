package swagger

import (
	"testing"

	"server/router/registry"

	"github.com/gin-gonic/gin"
)

type swaggerEmbeddedPage struct {
	Page     int `form:"page" comment:"页码"`
	PageSize int `form:"page_size" comment:"每页数量" binding:"required"`
}

type swaggerQueryRequest struct {
	swaggerEmbeddedPage
	Keyword string `form:"keyword" comment:"关键字"`
}

type swaggerBodyRequest struct {
	Name  string `json:"name" comment:"名称" binding:"required"`
	Score int    `json:"score" comment:"分数"`
}

func TestGeneratorUsesSharedFieldRulesForQueryAndBody(t *testing.T) {
	registry.Register("GET", "/phase3/swagger-query-test", "测试分组", "查询测试", func(*gin.Context) {},
		registry.WithRequest(swaggerQueryRequest{}))
	registry.Register("POST", "/phase3/swagger-body-test", "测试分组", "提交测试", func(*gin.Context) {},
		registry.WithRequest(swaggerBodyRequest{}))

	doc := NewGenerator(Config{
		Title:    "test",
		BasePath: "/api/v1",
	}).Generate()

	queryRoute := doc.Paths["/phase3/swagger-query-test"]["get"]
	if len(queryRoute.Parameters) != 3 {
		t.Fatalf("expected 3 query parameters, got %d", len(queryRoute.Parameters))
	}
	if queryRoute.Parameters[1].Name != "page_size" || !queryRoute.Parameters[1].Required {
		t.Fatalf("unexpected page_size parameter: %#v", queryRoute.Parameters[1])
	}

	bodyRoute := doc.Paths["/phase3/swagger-body-test"]["post"]
	if len(bodyRoute.Parameters) != 1 || bodyRoute.Parameters[0].Schema == nil {
		t.Fatalf("expected body schema parameter, got %#v", bodyRoute.Parameters)
	}

	def := doc.Definitions["swaggerBodyRequest"]
	if def.Properties["name"].Type != "string" {
		t.Fatalf("unexpected body field type: %#v", def.Properties["name"])
	}
	if len(def.Required) != 1 || def.Required[0] != "name" {
		t.Fatalf("unexpected required fields: %#v", def.Required)
	}
}
