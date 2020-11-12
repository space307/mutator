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
	fmt.Fprintf(out, "// MutateTo%%s return %%s filled from  pkg1.%%s\n", t.Name(), k.Name(), k.Name())
	fmt.Fprintf(out, "func (t *%%s) MutateTo%%s()  pkg1.%%s {\n", t.Name(), k.Name(), k.Name())
	fmt.Fprintf(out, "	return  pkg1.%%s{\n", k.Name())
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
	fmt.Fprintf(out, "func (t *%%s) FillFrom%%s(k pkg1.%%s) {\n", t.Name(), k.Name(), k.Name())
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

func generateFile(in Struct, to Struct) error {
	f, err := ioutil.TempFile(in.Path, "temp_mutagen")

	if err != nil {
		return err
	}

	// заполняем данные для шаблона файла данными структурами
	toPkg := fmt.Sprintf(`pkg1 "%s"`, to.Path)
	imp := fmt.Sprintf(`fmt.Fprintln(out,"	pkg1 \"%s\"")`, to.Path)
	gen := fmt.Sprintf(`generate(out,%s{},pkg1.%s{})`, in.Name, to.Name)

	fmt.Fprintln(f, fmt.Sprintf(k, in.Path, toPkg, in.Path, imp, gen))

	if err := f.Close(); err != nil {
		return err
	}

	src := f.Name()
	dest := src + "_test.go"
	err = os.Rename(src, dest)
	if err != nil {
		return err
	}

	defer func() {
		// remove temp generator
		_ = os.Remove(dest)
	}()

	// запускаем генерацию конечного кода
	execArgs := []string{"test", "-v", "-run", "TestK"}

	cmd := exec.Command("go", execArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = in.Path
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}
