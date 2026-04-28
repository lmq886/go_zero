package handler

import (
	"fmt"
	"net/http"
	"strings"
)

// SwaggerUIHTML Swagger UI 页面 HTML 模板
// 使用 CDN 加载 Swagger UI 资源
const SwaggerUIHTML = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>企业级后台管理系统 API 文档</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin: 0;
            background: #fafafa;
        }
        .topbar {
            background: linear-gradient(90deg, #2b6cb0 0%, #3182ce 100%) !important;
        }
        .topbar-wrapper .link a {
            color: white !important;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            // 获取当前页面的基础 URL
            var baseUrl = window.location.origin;
            
            // 构建 Swagger 文档 URL
            var specUrl = baseUrl + '/swagger.json';
            
            // 初始化 Swagger UI
            const ui = SwaggerUIBundle({
                url: specUrl,
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                // 配置认证
                supportedSubmitMethods: ['get', 'post', 'put', 'delete', 'patch'],
                // 中文显示
                validatorUrl: null,
                oauth2RedirectUrl: window.location.origin + '/swagger/oauth2-redirect.html',
                // 授权配置
                persistAuthorization: true,
                // 显示已认证的授权信息
                showCommonExtensions: true,
                // 尝试清除认证后显示
                filter: true,
                // 请求拦截器，添加 Authorization 头
                requestInterceptor: function(request) {
                    // 获取保存的 token
                    var tokenKey = 'authorized';
                    var authData = localStorage.getItem(tokenKey);
                    if (authData) {
                        try {
                            var auth = JSON.parse(authData);
                            if (auth && auth.BearerAuth && auth.BearerAuth.value) {
                                request.headers['Authorization'] = 'Bearer ' + auth.BearerAuth.value;
                            }
                        } catch (e) {
                            // 忽略解析错误
                        }
                    }
                    return request;
                }
            });
            
            // 配置认证
            ui.initOAuth({
                clientId: "your-client-id",
                clientSecret: "your-client-secret-if-required",
                realm: "your-realms",
                appName: "企业级后台管理系统",
                scopeSeparator: " ",
                scopes: {},
                additionalQueryStringParams: {},
                useBasicAuthenticationWithAccessCodeGrant: false,
                usePkceWithAuthorizationCodeGrant: false
            });
            
            window.ui = ui;
        };
    </script>
</body>
</html>
`

// SwaggerJSON 简化的 Swagger JSON 文档
var SwaggerJSON = []byte(`{
	"openapi": "3.0.3",
	"info": {
		"title": "企业级后台管理系统 API",
		"description": "包含用户管理、角色权限、菜单管理、日志管理等完整后台功能的 RESTful API 服务。\n\n## 认证方式\n\n除公开接口外，所有接口都需要通过 Bearer Token 进行认证。\n\n**获取 Token**：\n1. 调用 POST /api/v1/auth/login 登录获取 accessToken\n2. 点击右上角 Authorize 按钮，输入 accessToken（不需要 Bearer 前缀）\n3. 点击 Authorize 保存认证信息",
		"version": "v1.0.0",
		"contact": {
			"name": "GoZero Enterprise Team",
			"email": "support@gozero.com"
		},
		"license": {
			"name": "MIT",
			"url": "https://opensource.org/licenses/MIT"
		}
	},
	"servers": [
		{
			"url": "http://localhost:8080",
			"description": "本地开发环境"
		}
	],
	"tags": [
		{
			"name": "认证模块",
			"description": "用户认证相关接口（登录、注册、登出、刷新令牌）"
		},
		{
			"name": "健康检查",
			"description": "服务健康状态检查"
		}
	],
	"security": [
		{
			"BearerAuth": []
		}
	],
	"paths": {
		"/health": {
			"get": {
				"tags": ["健康检查"],
				"summary": "健康检查",
				"description": "检查服务运行状态，无需认证",
				"operationId": "healthCheck",
				"security": [],
				"responses": {
					"200": {
						"description": "服务正常运行",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"status": {"type": "string", "description": "服务状态", "example": "ok"},
										"version": {"type": "string", "description": "服务版本", "example": "v1.0.0"},
										"uptime": {"type": "string", "description": "运行时间", "example": "service is running"}
									}
								}
							}
						}
					},
					"500": {
						"description": "服务异常",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"code": {"type": "integer", "description": "状态码", "example": 500},
										"message": {"type": "string", "description": "错误消息", "example": "服务内部错误"}
									}
								}
							}
						}
					}
				}
			}
		},
		"/api/v1/auth/login": {
			"post": {
				"tags": ["认证模块"],
				"summary": "用户登录",
				"description": "使用用户名和密码登录，返回访问令牌和刷新令牌",
				"operationId": "login",
				"security": [],
				"requestBody": {
					"required": true,
					"content": {
						"application/json": {
							"schema": {
								"type": "object",
								"required": ["username", "password"],
								"properties": {
									"username": {"type": "string", "description": "用户名", "example": "admin"},
									"password": {"type": "string", "description": "密码", "example": "admin123"}
								}
							},
							"example": {
								"username": "admin",
								"password": "admin123"
							}
						}
					}
				},
				"responses": {
					"200": {
						"description": "登录成功",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"code": {"type": "integer", "description": "状态码", "example": 0},
										"message": {"type": "string", "description": "消息", "example": "success"},
										"accessToken": {"type": "string", "description": "访问令牌（用于后续请求认证）", "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
										"refreshToken": {"type": "string", "description": "刷新令牌（用于获取新的访问令牌）", "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
										"expiresAt": {"type": "integer", "description": "过期时间戳（秒）", "example": 1714070400},
										"userInfo": {
											"type": "object",
											"description": "用户信息",
											"properties": {
												"id": {"type": "integer", "description": "用户ID", "example": 1},
												"username": {"type": "string", "description": "用户名", "example": "admin"},
												"email": {"type": "string", "description": "邮箱", "example": "admin@example.com"},
												"phone": {"type": "string", "description": "手机号", "example": "13800138000"},
												"nickname": {"type": "string", "description": "昵称", "example": "管理员"},
												"avatar": {"type": "string", "description": "头像", "example": "/avatars/admin.png"},
												"status": {"type": "integer", "description": "状态（1: 启用, 0: 禁用）", "example": 1},
												"roles": {"type": "array", "description": "角色列表", "items": {"type": "string"}, "example": ["super_admin"]},
												"permissions": {"type": "array", "description": "权限列表", "items": {"type": "string"}, "example": ["system:user:list", "system:user:create"]},
												"createdAt": {"type": "integer", "description": "创建时间戳", "example": 1713811200}
											}
										}
									}
								},
								"example": {
									"code": 0,
									"message": "success",
									"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsInVzZXJuYW1lIjoiYWRtaW4iLCJleHAiOjE3MTQwNzA0MDAsImlhdCI6MTcxNDA2MzIwMH0.K1vpbF9Qx7K5P8u0w4e3r2t1y0u9i8o7p6a5s4d3f2g1h0j9k8l7m6n5b4v3c2x1z",
									"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsInVzZXJuYW1lIjoiYWRtaW4iLCJleHAiOjE3MTQ2NjgwMDAsImlhdCI6MTcxNDA2MzIwMH0.X1y2z3a4b5c6d7e8f9g0h1i2j3k4l5m6n7o8p9q0r1s2t3u4v5w6x7y8z9a0b1",
									"expiresAt": 1714070400,
									"userInfo": {
										"id": 1,
										"username": "admin",
										"email": "admin@example.com",
										"phone": "13800138000",
										"nickname": "管理员",
										"avatar": "/avatars/admin.png",
										"status": 1,
										"roles": ["super_admin"],
										"permissions": ["system:user:list", "system:user:create", "system:user:update", "system:user:delete"],
										"createdAt": 1713811200
									}
								}
							}
						}
					},
					"400": {
						"description": "请求参数错误",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"code": {"type": "integer", "example": 400},
										"message": {"type": "string", "example": "请求参数错误"}
									}
								}
							}
						}
					},
					"401": {
						"description": "用户名或密码错误",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"code": {"type": "integer", "example": 401},
										"message": {"type": "string", "example": "用户名或密码错误"}
									}
								}
							}
						}
					}
				}
			}
		},
		"/api/v1/auth/register": {
			"post": {
				"tags": ["认证模块"],
				"summary": "用户注册",
				"description": "创建新用户账号",
				"operationId": "register",
				"security": [],
				"requestBody": {
					"required": true,
					"content": {
						"application/json": {
							"schema": {
								"type": "object",
								"required": ["username", "password"],
								"properties": {
									"username": {"type": "string", "description": "用户名", "example": "newuser"},
									"password": {"type": "string", "description": "密码", "example": "password123"},
									"email": {"type": "string", "description": "邮箱", "example": "newuser@example.com"},
									"phone": {"type": "string", "description": "手机号", "example": "13800138001"},
									"nickname": {"type": "string", "description": "昵称", "example": "新用户"}
								}
							},
							"example": {
								"username": "newuser",
								"password": "password123",
								"email": "newuser@example.com",
								"phone": "13800138001",
								"nickname": "新用户"
							}
						}
					}
				},
				"responses": {
					"200": {
						"description": "注册成功",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"code": {"type": "integer", "example": 0},
										"message": {"type": "string", "example": "success"},
										"userId": {"type": "integer", "description": "用户ID", "example": 2}
									}
								}
							}
						}
					},
					"400": {
						"description": "请求参数错误",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"code": {"type": "integer", "example": 400},
										"message": {"type": "string", "example": "请求参数错误"}
									}
								}
							}
						}
					},
					"409": {
						"description": "用户名已存在",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"code": {"type": "integer", "example": 409},
										"message": {"type": "string", "example": "用户名已存在"}
									}
								}
							}
						}
					}
				}
			}
		},
		"/api/v1/auth/refresh": {
			"post": {
				"tags": ["认证模块"],
				"summary": "刷新令牌",
				"description": "使用刷新令牌获取新的访问令牌",
				"operationId": "refreshToken",
				"security": [],
				"requestBody": {
					"required": true,
					"content": {
						"application/json": {
							"schema": {
								"type": "object",
								"required": ["refreshToken"],
								"properties": {
									"refreshToken": {"type": "string", "description": "刷新令牌", "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
								}
							},
							"example": {
								"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsInVzZXJuYW1lIjoiYWRtaW4iLCJleHAiOjE3MTQ2NjgwMDAsImlhdCI6MTcxNDA2MzIwMH0.X1y2z3a4b5c6d7e8f9g0h1i2j3k4l5m6n7o8p9q0r1s2t3u4v5w6x7y8z9a0b1"
							}
						}
					}
				},
				"responses": {
					"200": {
						"description": "刷新成功",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"code": {"type": "integer", "example": 0},
										"message": {"type": "string", "example": "success"},
										"accessToken": {"type": "string", "description": "新的访问令牌", "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
										"expiresAt": {"type": "integer", "description": "过期时间戳", "example": 1714070400}
									}
								}
							}
						}
					},
					"400": {
						"description": "刷新令牌无效或已过期",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"code": {"type": "integer", "example": 400},
										"message": {"type": "string", "example": "刷新令牌无效"}
									}
								}
							}
						}
					}
				}
			}
		},
		"/api/v1/auth/logout": {
			"post": {
				"tags": ["认证模块"],
				"summary": "用户登出",
				"description": "注销当前登录用户（需要认证）",
				"operationId": "logout",
				"responses": {
					"200": {
						"description": "登出成功",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"code": {"type": "integer", "example": 0},
										"message": {"type": "string", "example": "success"}
									}
								}
							}
						}
					},
					"401": {
						"description": "未授权",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"code": {"type": "integer", "example": 401},
										"message": {"type": "string", "example": "未授权"}
									}
								}
							}
						}
					}
				}
			}
		}
	},
	"components": {
		"securitySchemes": {
			"BearerAuth": {
				"type": "http",
				"scheme": "bearer",
				"bearerFormat": "JWT",
				"description": "JWT 认证方式。请先调用登录接口获取 accessToken，然后点击右上角 Authorize 按钮，输入 accessToken（不需要 Bearer 前缀）。\n\n格式：直接输入 accessToken 值\n示例：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
			}
		}
	}
}`)

// SwaggerUIHandler 提供 Swagger UI 页面
func SwaggerUIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(SwaggerUIHTML))
	}
}

// SwaggerJSONHandler 提供 Swagger JSON 文档
func SwaggerJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// 获取请求的协议和主机
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		host := r.Host

		// 替换服务器 URL 为当前请求的主机
		swaggerJSON := strings.Replace(
			string(SwaggerJSON),
			`"url": "http://localhost:8080"`,
			fmt.Sprintf(`"url": "%s://%s"`, scheme, host),
			1,
		)

		w.Write([]byte(swaggerJSON))
	}
}

// SwaggerYAMLHandler 提供 Swagger YAML 文档
func SwaggerYAMLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
		w.Write([]byte(`openapi: 3.0.3
info:
  title: 企业级后台管理系统 API
  description: 包含用户管理、角色权限、菜单管理、日志管理等完整后台功能的 RESTful API 服务
  version: v1.0.0
servers:
  - url: http://localhost:8080
    description: 本地开发环境
`))
	}
}

// SwaggerRedirectHandler 重定向到 Swagger UI
func SwaggerRedirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	}
}
