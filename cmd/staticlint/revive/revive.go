package revive

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// RunRevive запускает revive и возвращает его вывод.
func RunRevive() (string, error) {
	// Определяем путь к исполняемому файлу
	executablePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("не удалось определить путь к исполняемому файлу: %v", err)
	}

	// Формируем путь к конфигурационному файлу относительно местоположения исполняемого файла
	configPath := filepath.Join(filepath.Dir(executablePath), "cmd/staticlint/revive/config.toml")

	// Устанавливаем таймаут и запускаем команду
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "revive", "-config", configPath, "./...") // Используем ./..., чтобы проверить только текущий каталог
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("превышено время ожидания revive")
	}
	if err != nil {
		fmt.Printf("Ошибка выполнения revive: %v\nВывод: %s\n", err, output)
	}

	return string(output), err
}
