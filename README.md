# generator

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