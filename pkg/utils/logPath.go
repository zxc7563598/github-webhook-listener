package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GenerateRequestLogPath 生成一个用于存放单次请求结果的随机文件路径
//
// 路径结构示例：
//
//	<exe_dir>/logs/shell/20260127/req_1706341234567_a3f9c2.log
//
// 设计说明：
//  1. 以「可执行文件所在目录」为基准，而不是运行目录
//  2. 按日期分目录，方便定位与清理
//  3. 文件名包含时间戳 + 随机串，避免并发冲突
func GenerateRequestLogPath() (string, error) {
	baseDir, err := GetExecutableDir()
	if err != nil {
		return "", fmt.Errorf("GetExecutableDir 失败: %w", err)
	}
	// 日期目录：YYYYMMDD
	dateDir := time.Now().Format("20060102")
	// logs/shell/YYYYMMDD
	dirPath := filepath.Join(baseDir, "logs", "shell", dateDir)
	// 确保目录存在
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("创建日志目录失败: %w", err)
	}
	// 生成随机文件名
	fileName, err := randomFileName()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirPath, fileName), nil
}

func GetExecutableDir() (string, error) {
	// 获取可执行文件路径（可能是软链接）
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("获取可执行路径失败: %w", err)
	}
	// 解析软链接，拿到真实路径（非常重要）
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", fmt.Errorf("绝对路径获取失败: %w", err)
	}
	// 可执行文件所在目录
	return filepath.Dir(exePath), nil
}

func randomFileName() (string, error) {
	ts := time.Now().UnixMilli()
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random bytes failed: %w", err)
	}
	return fmt.Sprintf(
		"req_%d_%s.log",
		ts,
		hex.EncodeToString(b),
	), nil
}
