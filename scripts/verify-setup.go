package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	fmt.Println("🚀 FastDFS迁移系统 - 环境验证脚本")
	fmt.Println("=====================================")

	checks := []func() bool{
		checkGoVersion,
		checkDependencies,
		checkBuild,
		checkTests,
		checkServer,
	}

	passed := 0
	total := len(checks)

	for _, check := range checks {
		if check() {
			passed++
		}
		fmt.Println()
	}

	fmt.Printf("验证结果: %d/%d 通过\n", passed, total)
	if passed == total {
		fmt.Println("✅ 环境验证成功！项目已准备就绪。")
		fmt.Println("\n下一步:")
		fmt.Println("1. 阅读 HANDOVER.md 了解项目状态")
		fmt.Println("2. 阅读 DEVELOPMENT.md 了解开发流程")
		fmt.Println("3. 查看 .kiro/specs/ 目录下的规格文档")
		fmt.Println("4. 开始下一个任务的开发")
	} else {
		fmt.Println("❌ 环境验证失败，请检查上述问题。")
		os.Exit(1)
	}
}

func checkGoVersion() bool {
	fmt.Print("🔍 检查Go版本... ")

	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Go未安装: %v\n", err)
		return false
	}

	version := string(output)
	// 检查是否包含go1.21或更高版本
	if strings.Contains(version, "go1.21") || strings.Contains(version, "go1.22") ||
		strings.Contains(version, "go1.23") || strings.Contains(version, "go1.24") ||
		strings.Contains(version, "go1.25") {
		fmt.Printf("✅ %s\n", strings.TrimSpace(version))
		return true
	}

	fmt.Printf("⚠️  Go版本可能过低: %s\n", strings.TrimSpace(version))
	fmt.Println("   建议使用Go 1.21+")
	return false
}

func checkDependencies() bool {
	fmt.Print("📦 检查项目依赖... ")

	cmd := exec.Command("go", "mod", "download")
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ 依赖下载失败: %v\n", err)
		return false
	}

	cmd = exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ 依赖整理失败: %v\n", err)
		return false
	}

	fmt.Println("✅ 依赖已就绪")
	return true
}

func checkBuild() bool {
	fmt.Print("🔨 检查项目构建... ")

	// 清理之前的构建
	os.RemoveAll("bin")

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("go", "build", "-o", "bin/migration-system.exe", "./cmd/server")
	} else {
		cmd = exec.Command("go", "build", "-o", "bin/migration-system", "./cmd/server")
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ 构建失败: %v\n", err)
		return false
	}

	fmt.Println("✅ 构建成功")
	return true
}

func checkTests() bool {
	fmt.Print("🧪 运行测试... ")

	cmd := exec.Command("go", "test", "./...")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ 测试失败: %v\n", err)
		fmt.Printf("输出: %s\n", string(output))
		return false
	}

	fmt.Println("✅ 所有测试通过")
	return true
}

func checkServer() bool {
	fmt.Print("🌐 检查服务器启动... ")

	// 启动服务器
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("bin/migration-system.exe")
	} else {
		cmd = exec.Command("bin/migration-system")
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ 服务器启动失败: %v\n", err)
		return false
	}

	// 等待服务器启动
	time.Sleep(2 * time.Second)

	// 测试健康检查端点
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		fmt.Printf("❌ 无法连接到服务器: %v\n", err)
		cmd.Process.Kill()
		return false
	}
	resp.Body.Close()

	// 测试API端点
	resp, err = http.Get("http://localhost:8080/api/v1/ping")
	if err != nil {
		fmt.Printf("❌ API端点无响应: %v\n", err)
		cmd.Process.Kill()
		return false
	}
	resp.Body.Close()

	// 关闭服务器
	cmd.Process.Kill()

	fmt.Println("✅ 服务器正常运行")
	return true
}
