// Package revive предоставляет анализатор revive.
package revive

import "os/exec"

// RunRevive запускает revive и возвращает его вывод
func RunRevive() (string, error) {
	cmd := exec.Command("revive", "-config", "revive/config.toml", "../../...")
	output, err := cmd.CombinedOutput()
	return string(output), err
}
