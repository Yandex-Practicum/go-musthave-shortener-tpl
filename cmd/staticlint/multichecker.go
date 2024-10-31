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
	"honnef.co/go/tools/staticcheck"
	"log"
	"strings"
)

// main - точка входа.
func main() {
	var mychecks []*analysis.Analyzer

	// определяем мультипроверку
	mychecks = append(mychecks,
		noexit.NoExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	)

	// Добавляем все SA-анализаторы из пакета staticcheck
	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name[:2] == "SA" { // фильтр SA-анализаторов
			mychecks = append(mychecks, a.Analyzer)
		}
	}

	// Подключаем анализатор ST1000 (класс ST)
	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name == "ST1000" {
			mychecks = append(mychecks, a.Analyzer)
			break
		}
	}

	// Подключаем публичные анализаторы
	// Запускаем revive
	reviveOutput, err := revive.RunRevive()
	if err != nil {
		log.Fatalf("ошибка запуска revive: %v\n%s", err, reviveOutput)
	}

	// Выводим результаты revive
	fmt.Println("Revive Output:")
	fmt.Println(strings.TrimSpace(reviveOutput))

	// Запускаем multichecker с добавленными анализаторами
	multichecker.Main(
		mychecks...,
	)
}
