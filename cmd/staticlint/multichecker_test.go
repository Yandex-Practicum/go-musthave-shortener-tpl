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
	foundS := false
	for _, a := range analyzers {
		if a.Name == "S1005" || a.Name == "S1006" {
			foundS = true
			break
		}
	}
	if !foundS {
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
