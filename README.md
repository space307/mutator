# generator

[![Release](https://img.shields.io/github/release/space307/mutator.svg)](https://github.com/space307/mutator/releases/latest)
[![License](https://img.shields.io/github/license/space307/mutator.svg)](https://raw.githubusercontent.com/space307/mutator/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/space307/mutator)](https://goreportcard.com/report/github.com/space307/mutator)


Генератор сканит все структуры в проекте и ищет структуры со специальным комментарием вида
` mutagento <struct package> <struct name>`

При генерации создаётся в директории с такой структурой временный файл теста, который генерирует код, в котором содержаться методы перевода в эту структуру.

Пример

есть в коде в пакете "inner" структура

```
type Test1 struct {
    Field int
}
```
и
```
// mutagento inner Test1
type MyStruct struct {
    Field int
    AnotherField int
}
```

В результате генерации рядом с файлом структуры будет сгенерирован файл с методами

```
func (t *MyStruct) MutateToTest1()  inner.Test1 {
	return  inner.Test1{
		Field: t.Field,
	}
}

func (t *MyStruct) FillFromTest1(k inner.Test1) {
	t.Field = k.Field
}
```