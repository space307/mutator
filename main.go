package main

import (
	"fmt"
	"go/ast"
	"os"
)

func main() {
	// директория проекта, в котором ищем структуры для генерации
	dir := "."

	// парсим ВСЕ структуры из проекта, по специальному комменту находим соответствия между структурами
	// пары структур пишем в p.Connections
	p := Parser{All: make(map[string]*ast.StructType), Connections: make(map[Struct]Struct)}
	if err := p.Parse(dir); err != nil {
		fmt.Printf("Error parsing: %s", err)
	}

	//Проходимся по всем парам и генерим файлы с кодом преобразования
	for from, to := range p.Connections {
		err := generateFile(from, to)
		if err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(1)
		}
	}
}
