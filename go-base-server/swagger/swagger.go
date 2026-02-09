package swagger

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"go-base-server/router/registry"

	"github.com/gin-gonic/gin"
)

// SwaggerInfo API 基本信息
type SwaggerInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	BasePath    string `json:"basePath"`
	Host        string `json:"host"`
}

// SwaggerDoc Swagger 文档结构
type SwaggerDoc struct {
	Swagger     string                          `json:"swagger"`
	Info        SwaggerInfo                     `json:"info"`
	Host        string                          `json:"host"`
	BasePath    string                          `json:"basePath"`
	Paths       map[string]map[string]Operation `json:"paths"`
	Definitions map[string]Definition           `json:"definitions"`
	Tags        []Tag                           `json:"tags"`
	SecurityDef map[string]SecurityDefinition   `json:"securityDefinitions,omitempty"`
}

// Operation API 操作
type Operation struct {
	Tags        []string              `json:"tags"`
	Summary     string                `json:"summary"`
	Description string                `json:"description,omitempty"`
	OperationID string                `json:"operationId"`
	Consumes    []string              `json:"consumes"`
	Produces    []string              `json:"produces"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []map[string][]string `json:"security,omitempty"`
}

// Parameter 参数
type Parameter struct {
	Name        string     `json:"name"`
	In          string     `json:"in"`
	Description string     `json:"description,omitempty"`
	Required    bool       `json:"required"`
	Type        string     `json:"type,omitempty"`
	Schema      *SchemaRef `json:"schema,omitempty"`
}

// SchemaRef Schema 引用
type SchemaRef struct {
	Ref string `json:"$ref,omitempty"`
}

// Response 响应
type Response struct {
	Description string     `json:"description"`
	Schema      *SchemaRef `json:"schema,omitempty"`
}

// Definition 定义
type Definition struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

// Property 属性
type Property struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Format      string `json:"format,omitempty"`
	Example     string `json:"example,omitempty"`
}

// Tag 标签
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// SecurityDefinition 安全定义
type SecurityDefinition struct {
	Type string `json:"type"`
	Name string `json:"name"`
	In   string `json:"in"`
}

// Config Swagger 配置
type Config struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
}

// Generator Swagger 生成器
type Generator struct {
	config      Config
	definitions map[string]Definition
}

// NewGenerator 创建生成器
func NewGenerator(config Config) *Generator {
	return &Generator{
		config:      config,
		definitions: make(map[string]Definition),
	}
}

// Generate 生成 Swagger 文档
func (g *Generator) Generate() *SwaggerDoc {
	routes := registry.GetAllRoutes()

	doc := &SwaggerDoc{
		Swagger: "2.0",
		Info: SwaggerInfo{
			Title:       g.config.Title,
			Description: g.config.Description,
			Version:     g.config.Version,
		},
		Host:        g.config.Host,
		BasePath:    g.config.BasePath,
		Paths:       make(map[string]map[string]Operation),
		Definitions: make(map[string]Definition),
		Tags:        make([]Tag, 0),
		SecurityDef: map[string]SecurityDefinition{
			"Bearer": {
				Type: "apiKey",
				Name: "Authorization",
				In:   "header",
			},
		},
	}

	// 收集所有 tags
	tagMap := make(map[string]bool)

	for _, route := range routes {
		// 移除 BasePath 前缀
		path := route.Path
		if strings.HasPrefix(path, g.config.BasePath) {
			path = strings.TrimPrefix(path, g.config.BasePath)
		}
		if path == "" {
			path = "/"
		}

		method := strings.ToLower(route.Method)

		if doc.Paths[path] == nil {
			doc.Paths[path] = make(map[string]Operation)
		}

		op := Operation{
			Tags:        []string{route.Group},
			Summary:     route.Summary,
			Description: route.Description,
			OperationID: g.generateOperationID(route),
			Consumes:    []string{"application/json"},
			Produces:    []string{"application/json"},
			Responses: map[string]Response{
				"200": {
					Description: "成功",
					Schema:      &SchemaRef{Ref: "#/definitions/Response"},
				},
			},
		}

		// 添加请求参数
		if route.Request != nil {
			if method == "post" || method == "put" || method == "patch" {
				// POST/PUT/PATCH 使用 body 参数
				typeName := g.addDefinition(route.Request)
				op.Parameters = append(op.Parameters, Parameter{
					Name:        "body",
					In:          "body",
					Description: "请求参数",
					Required:    true,
					Schema:      &SchemaRef{Ref: "#/definitions/" + typeName},
				})
			} else if method == "get" || method == "delete" {
				// GET/DELETE 使用 query 参数
				queryParams := g.parseQueryParams(route.Request)
				op.Parameters = append(op.Parameters, queryParams...)
			}
		}

		// 添加认证
		if route.NeedAuth {
			op.Security = []map[string][]string{
				{"Bearer": {}},
			}
		}

		doc.Paths[path][method] = op
		tagMap[route.Group] = true
	}

	// 添加 tags
	for tag := range tagMap {
		doc.Tags = append(doc.Tags, Tag{Name: tag})
	}

	// 添加通用响应定义
	doc.Definitions["Response"] = Definition{
		Type: "object",
		Properties: map[string]Property{
			"code":    {Type: "integer", Description: "状态码"},
			"message": {Type: "string", Description: "消息"},
			"data":    {Type: "object", Description: "数据"},
		},
	}

	// 合并自定义定义
	for k, v := range g.definitions {
		doc.Definitions[k] = v
	}

	return doc
}

// addDefinition 添加结构体定义
func (g *Generator) addDefinition(obj interface{}) string {
	if obj == nil {
		return ""
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	name := t.Name()
	if _, exists := g.definitions[name]; exists {
		return name
	}

	def := Definition{
		Type:       "object",
		Properties: make(map[string]Property),
		Required:   make([]string, 0),
	}

	g.parseStructFields(t, &def)

	g.definitions[name] = def
	return name
}

// parseStructFields 递归解析结构体字段
func (g *Generator) parseStructFields(t reflect.Type, def *Definition) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 处理匿名嵌入结构体（如 PageRequest）
		if field.Anonymous {
			embeddedType := field.Type
			if embeddedType.Kind() == reflect.Ptr {
				embeddedType = embeddedType.Elem()
			}
			if embeddedType.Kind() == reflect.Struct {
				g.parseStructFields(embeddedType, def)
			}
			continue
		}

		// 获取 json tag
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			// 尝试使用 form tag
			jsonTag = field.Tag.Get("form")
			if jsonTag == "" || jsonTag == "-" {
				continue
			}
		}
		jsonName := strings.Split(jsonTag, ",")[0]

		// 获取字段描述: 优先从 comment tag 获取
		description := field.Tag.Get("comment")
		if description == "" {
			// 从 label tag 获取
			description = field.Tag.Get("label")
		}

		prop := Property{
			Type:        g.goTypeToSwagger(field.Type),
			Description: description,
		}

		// 解析 binding tag 获取是否必填
		bindingTag := field.Tag.Get("binding")
		if strings.Contains(bindingTag, "required") {
			def.Required = append(def.Required, jsonName)
		}

		def.Properties[jsonName] = prop
	}
}

// parseQueryParams 解析 query 参数
func (g *Generator) parseQueryParams(obj interface{}) []Parameter {
	if obj == nil {
		return nil
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	var params []Parameter
	g.parseQueryParamsRecursive(t, &params)
	return params
}

// parseQueryParamsRecursive 递归解析 query 参数
func (g *Generator) parseQueryParamsRecursive(t reflect.Type, params *[]Parameter) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 处理匿名嵌入结构体
		if field.Anonymous {
			embeddedType := field.Type
			if embeddedType.Kind() == reflect.Ptr {
				embeddedType = embeddedType.Elem()
			}
			if embeddedType.Kind() == reflect.Struct {
				g.parseQueryParamsRecursive(embeddedType, params)
			}
			continue
		}

		// 优先使用 form tag, 然后 json tag
		paramName := field.Tag.Get("form")
		if paramName == "" || paramName == "-" {
			paramName = field.Tag.Get("json")
			if paramName == "" || paramName == "-" {
				continue
			}
		}
		paramName = strings.Split(paramName, ",")[0]

		// 获取描述
		description := field.Tag.Get("comment")
		if description == "" {
			description = field.Tag.Get("label")
		}

		// 检查是否必填
		bindingTag := field.Tag.Get("binding")
		required := strings.Contains(bindingTag, "required")

		*params = append(*params, Parameter{
			Name:        paramName,
			In:          "query",
			Description: description,
			Required:    required,
			Type:        g.goTypeToSwagger(field.Type),
		})
	}
}

// goTypeToSwagger Go 类型转 Swagger 类型
func (g *Generator) goTypeToSwagger(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Ptr:
		return g.goTypeToSwagger(t.Elem())
	default:
		return "object"
	}
}

// generateOperationID 生成操作 ID
func (g *Generator) generateOperationID(route registry.RouteInfo) string {
	// 从路径生成操作ID
	path := strings.ReplaceAll(route.Path, "/", "_")
	path = strings.ReplaceAll(path, "-", "_")
	path = strings.TrimPrefix(path, "_")
	return strings.ToLower(route.Method) + "_" + path
}

// RegisterRoutes 注册 Swagger 路由
func RegisterRoutes(r *gin.Engine, config Config) {
	gen := NewGenerator(config)
	doc := gen.Generate()

	// Swagger UI 和 JSON
	r.GET("/swagger/*any", func(c *gin.Context) {
		path := c.Param("any")
		switch path {
		case "/doc.json":
			c.JSON(http.StatusOK, doc)
		case "/", "/index.html", "":
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, swaggerUIHTML(config.Host))
		default:
			c.Status(http.StatusNotFound)
		}
	})
}

// ToJSON 转为 JSON
func (doc *SwaggerDoc) ToJSON() string {
	b, _ := json.MarshalIndent(doc, "", "  ")
	return string(b)
}

// swaggerUIHTML Swagger UI 页面
func swaggerUIHTML(host string) string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: '#swagger-ui',
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.SwaggerUIStandalonePreset
                ],
                layout: "BaseLayout"
            });
        }
    </script>
</body>
</html>`
}
