package actions

import (
	"log"

	"github.com/zxc7563598/github-webhook-listener/internal/config"
)

func ExecuteActions(actions []config.ActionConfig, payload []byte) {
	for _, a := range actions {
		switch a.Type {
		case "shell":
			log.Printf("[action] executing shell: %s", a.Command)
			runShell(a.Command)
		default:
			log.Printf("[action] unknown type: %s", a.Type)
		}
	}
}
