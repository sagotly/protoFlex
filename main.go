package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root := "." // Корневой каталог проекта
	outputFile := "test_functions_report.txt"

	// Создаём или очищаем файл для записи отчёта
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Ошибка при создании файла отчёта: %v", err)
	}
	defer file.Close()

	// Создаём новое файловое множество
	fset := token.NewFileSet()

	// Счётчик тестовых функций
	testFuncCount := 0

	// Функция для обработки каждого Go-файла
	processFile := func(path string) {
		// Парсим файл
		node, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			log.Printf("Ошибка при разборе файла %s: %v", path, err)
			return
		}

		// Проходимся по всем объявлениям в файле
		for _, decl := range node.Decls {
			// Проверяем, является ли объявление функцией
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				// Проверяем, начинается ли имя функции с "Test"
				if strings.HasPrefix(funcDecl.Name.Name, "Test") {
					testFuncCount++
					line := fmt.Sprintf("Найдена тестовая функция: %s в файле %s\n", funcDecl.Name.Name, path)
					file.WriteString(line)
				}
			}
		}
	}

	// Обходим все файлы в корневом каталоге и его подкаталогах
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Обрабатываем только Go-файлы
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			processFile(path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Ошибка при обходе директории: %v", err)
	}

	// Добавляем статистику в конец файла
	summary := fmt.Sprintf("\nОбщее количество тестовых функций: %d\n", testFuncCount)
	file.WriteString(summary)

	fmt.Printf("Анализ завершён. Отчёт сохранён в файле %s\n", outputFile)
}
