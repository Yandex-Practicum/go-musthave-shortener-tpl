package revive

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func RunRevive() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "revive", "-config", "cmd/staticlint/revive/config.toml", "./...")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Ошибка выполнения revive: %v\nВывод: %s\n", err, output)
		return string(output), err
	}
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("превышено время ожидания revive")
	}

	return string(output), nil
}
