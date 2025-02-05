package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sqweek/dialog"
)

func main() {
	// Открываем диалог выбора исходного файла
	inputFilePath, err := dialog.File().Title("Выберите исходный CSV файл").Load()
	if err != nil {
		log.Fatalf("Ошибка выбора файла: %v", err)
	}

	// Открываем исходный CSV файл
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v", err)
	}
	defer inputFile.Close()

	// Создаем CSV reader с разделителем ";"
	reader := csv.NewReader(inputFile)
	reader.Comma = ';' // Устанавливаем разделитель для CSV

	// Пропускаем первую строку (заголовки)
	_, err = reader.Read()
	if err != nil {
		log.Fatalf("Ошибка чтения первой строки: %v", err)
	}

	// Открываем диалог выбора папки для сохранения итогового файла
	outputFolderPath, err := dialog.Directory().Title("Выберите папку для сохранения итогового файла").Browse()
	if err != nil {
		log.Fatalf("Ошибка выбора папки: %v", err)
	}

	// Создаем путь для итогового файла
	outputFilePath := outputFolderPath + "\\output.txt"

	// Создаем новый файл для записи результатов
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("Ошибка создания файла: %v", err)
	}
	defer outputFile.Close()

	// Создаем writer для записи в файл с использованием буферизации
	writer := bufio.NewWriter(outputFile)

	// Открываем файл как строку, оборачиваем в to_clob('
	writer.WriteString("to_clob(\n'")

	// Читаем строки и обрабатываем их
	var lineCount int
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() != "EOF" {
				log.Fatalf("Ошибка чтения CSV: %v", err)
			}
			break // Конец файла
		}

		// Заменяем одинарные кавычки на две
		for i := range record {
			record[i] = strings.ReplaceAll(record[i], "'", "''")
		}

		// Записываем текущую строку в итоговый файл, добавляем разделитель ;
		writer.WriteString(strings.Join(record, ";"))

		// Если это 15-я строка, то нужно завершить её одинарной кавычкой
		lineCount++
		if lineCount%15 == 0 {
			// Добавляем `')||to_clob(` на новую строку
			writer.WriteString("\n')||to_clob(\n'")

			// Следующую строку начинаем с `'`
			// Добавляем символ `'` в начало текущей строки
			for i := range record {
				record[i] = "'" + record[i]
			}
		} else {
			// Если строка не 15-я, просто добавляем новую строку
			writer.WriteString("\n")
		}
	}

	// Закрываем CLOB обертку
	writer.WriteString("')")

	// Сохраняем изменения в файле
	writer.Flush()

	fmt.Println("Файл успешно обработан!")
}
