// Package main предоставляет multichecker для проверки кода на стандартные ошибки,
package main

import (
	"fmt"
	"github.com/kamencov/go-musthave-shortener-tpl/cmd/staticlint/noexit"
	"github.com/kamencov/go-musthave-shortener-tpl/cmd/staticlint/revive"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"log"
	"strings"
)

// configureAnalyzers настраивает и возвращает список анализаторов.
func configureAnalyzers() []*analysis.Analyzer {
	var analyzers []*analysis.Analyzer
	analyzers = append(analyzers,
		noexit.NoExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	)

	// Добавляем все SA-анализаторы из пакета staticcheck
	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name[:2] == "SA" {
			analyzers = append(analyzers, a.Analyzer)
		}
	}

	// Подключаем анализатор S1005 (класс S)
	for _, a := range simple.Analyzers {
		if a.Analyzer.Name == "S1005" || a.Analyzer.Name == "S1006" {
			analyzers = append(analyzers, a.Analyzer)
			break
		}
	}

	return analyzers
}

// runRevive запускает revive и возвращает вывод и ошибку.
func runRevive() (string, error) {
	return revive.RunRevive()
}

func main() {
	// Настройка анализаторов
	mychecks := configureAnalyzers()

	// Запуск revive
	reviveOutput, err := runRevive()
	if err != nil {
		log.Fatalf("ошибка запуска revive: %v\n%s", err, reviveOutput)
	}

	// Вывод результатов revive
	fmt.Println("Revive Output:")
	fmt.Println(strings.TrimSpace(reviveOutput))

	// Запуск multichecker
	multichecker.Main(mychecks...)
}
