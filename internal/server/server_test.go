package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"fastdfs-migration-system/internal/config"
)

func TestServer_HealthCheck(t *testing.T) {
	// 创建测试配置
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
			Host: "localhost",
		},
		Logging: config.LoggingConfig{
			Level: "error", // 减少测试时的日志输出
		},
	}

	// 创建服务器实例
	server := New(cfg)

	// 创建测试请求
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 创建响应记录器
	rr := httptest.NewRecorder()

	// 执行请求
	server.router.ServeHTTP(rr, req)

	// 检查状态码
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 检查响应内容
	expected := `{"status":"ok"`
	if !contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want to contain %v",
			rr.Body.String(), expected)
	}
}

func TestServer_PingAPI(t *testing.T) {
	cfg := &config.Config{
		Logging: config.LoggingConfig{
			Level: "error",
		},
	}

	server := New(cfg)

	req, err := http.NewRequest("GET", "/api/v1/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	server.router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"message":"pong"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
		(len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
