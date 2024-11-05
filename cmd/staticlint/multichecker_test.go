// main_test.go
package main

import (
	"strings"
	"testing"
)

// Тестируем configureAnalyzers, чтобы проверить правильность настройки анализаторов.
func TestConfigureAnalyzers(t *testing.T) {
	analyzers := configureAnalyzers()

	// Пример проверки, что добавлен анализатор ST1000
	foundST1000 := false
	for _, a := range analyzers {
		if a.Name == "S1005" || a.Name == "ST1000" {
			foundST1000 = true
			break
		}
	}
	if !foundST1000 {
		t.Error("анализатор ST1000 не был добавлен")
	}

	// Пример проверки, что добавлен хотя бы один SA-анализатор
	foundSA := false
	for _, a := range analyzers {
		if strings.HasPrefix(a.Name, "SA") {
			foundSA = true
			break
		}
	}
	if !foundSA {
		t.Error("не было добавлено SA-анализаторов")
	}
}

// Тестируем runRevive, проверяя успешный запуск revive.
func TestRunRevive(t *testing.T) {
	output, err := runRevive()
	if err != nil {
		t.Fatalf("runRevive завершился с ошибкой: %v", err)
	}
	if len(output) == 0 {
		t.Error("runRevive вернул пустой вывод")
	}
}
