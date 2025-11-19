package actions

import (
	"context"
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	// DefaultShellTimeout 默认shell命令执行超时时间
	DefaultShellTimeout = 5 * time.Minute
	// DefaultWorkDir 默认工作目录（如果未指定）
	DefaultWorkDir = "/tmp"
)

func runShell(cmdStr string) {
	runShellWithTimeout(cmdStr, DefaultShellTimeout, "")
}

func runShellWithTimeout(cmdStr string, timeout time.Duration, workDir string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)

	// 设置工作目录
	if workDir != "" {
		cmd.Dir = workDir
	} else {
		// 如果没有指定工作目录，使用当前用户的主目录或/tmp
		if homeDir := os.Getenv("HOME"); homeDir != "" {
			cmd.Dir = homeDir
		} else {
			cmd.Dir = DefaultWorkDir
		}
	}

	// 设置环境变量（继承当前环境，但可以添加额外的）
	cmd.Env = os.Environ()

	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("[shell] 命令执行超时 (超过 %v): %s", timeout, cmdStr)
		} else {
			log.Printf("[shell] 执行错误: %v, 命令: %s", err, cmdStr)
		}
		if len(out) > 0 {
			log.Printf("[shell] 错误输出: %s", string(out))
		}
		return
	}

	if len(out) > 0 {
		log.Printf("[shell] 输出: %s", string(out))
	} else {
		log.Printf("[shell] 命令执行成功: %s", cmdStr)
	}
}
