package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// шаблон файла, который генерируем и запускаем для генерации конечного кода
var k = `package %s

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	%s
)

func TestK(t *testing.T) {
	out := bytes.NewBuffer([]byte{})

	fmt.Fprintln(out,"package %s")
	fmt.Fprintln(out,"import (")
	%s
	fmt.Fprintln(out,")")
	
	//generate(out, %%s{}, %%s{})
	%s

	f, err := os.Create("mutated.go")
	if err != nil {
		t.Fail()
	}
	_, err = f.Write(out.Bytes())
	if err != nil {
		t.Fail()
	}
}

func generate(out *bytes.Buffer, ex interface{}, in interface{}) {
	t := reflect.TypeOf(ex)
	k := reflect.TypeOf(in)

	fmt.Fprintln(out,"")
	fmt.Fprintf(out, "// MutateTo%%s return %%s filled from  %%s\n", t.Name(), k.Name(), k.String())
	fmt.Fprintf(out, "func (t *%%s) MutateTo%%s()  %%s {\n", t.Name(), k.Name(), k.String())
	fmt.Fprintf(out, "	return  %%s{\n", k.String())
	for i :=0;i< t.NumField();i++ {
		for j :=0;j< k.NumField();j++ {
			if k.Field(j).Name == t.Field(i).Name && k.Field(j).Type == t.Field(i).Type {
				fmt.Fprintf(out, "		%%s: t.%%s,\n", t.Field(i).Name, t.Field(i).Name)
			}

			//collect wrong types
		}
	}
	fmt.Fprintln(out, "	}")
	fmt.Fprintln(out, "}")

	fmt.Fprintln(out,"")
	fmt.Fprintf(out, "// FillFrom%%s fill %%s from %%s values\n", t.Name(), k.Name(), k.Name())
	fmt.Fprintf(out, "func (t *%%s) FillFrom%%s(k %%s) {\n", t.Name(), k.Name(), k.String())
	for i :=0;i< t.NumField();i++ {
		for j :=0;j< k.NumField();j++ {
			if k.Field(j).Name == t.Field(i).Name && k.Field(j).Type == t.Field(i).Type {
				fmt.Fprintf(out, "	t.%%s = k.%%s\n", t.Field(i).Name, t.Field(i).Name)
			}

			//collect wrong types
		}
	}
	fmt.Fprintln(out, "}")
}
`

func generateFile(dir string, st []StructPair) error {
	f, err := ioutil.TempFile("./"+dir, "temp_mutagen")

	if err != nil {
		return fmt.Errorf("generate temp file: %w", err)
	}

	var (
		toPkg, imp, gen string
		pkgNum          int
	)

	pack := make(map[string]string)
	for _, item := range st {
		if _, ok := pack[item.To.Path]; !ok {
			pack[item.To.Path] = fmt.Sprintf("pkg%d", pkgNum)
			pkgNum++
		}
	}

	for p, n := range pack {
		toPkg += fmt.Sprintf("%s \"%s\"\n", n, p)
		imp += fmt.Sprintf("fmt.Fprintln(out,\"	\\\"%s\\\"\")\n", p)
	}

	// заполняем данные для шаблона файла данными структурами
	for _, item := range st {
		pkgName := pack[item.To.Path]
		gen += fmt.Sprintf("generate(out,%s{},%s.%s{})\n", item.From.Name, pkgName, item.To.Name)
	}

	if _, err = fmt.Fprintln(f, fmt.Sprintf(k, dir, toPkg, dir, imp, gen)); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	src := f.Name()
	dest := src + "_test.go"
	err = os.Rename(src, dest)
	if err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}

	defer func() {
		// remove temp generator
		if er := os.Remove(dest); er != nil {
			fmt.Printf("Error removing temporary file: %s", er)
		}
	}()

	// запускаем генерацию конечного кода
	cmd := exec.Command("go", "test", "-v", "-run", "TestK")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("run temp file: %w", err)
	}

	cmd = exec.Command("go", "fmt", "mutated.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("run fmt for mutated.go: %w", err)
	}

	return nil
}
