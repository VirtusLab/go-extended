package yaml_test

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/VirtusLab/go-extended/pkg/yaml"
)

func ExampleToInterface_empty() {
	y := ``

	data, err := yaml.ToInterface(strings.NewReader(y))
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
	fmt.Println(data)
	fmt.Println(reflect.TypeOf(data))
	fmt.Println(reflect.ValueOf(data).Kind())

	// Output:
	// map[]
	// map[string]interface {}
	// map
}

func ExampleToInterface_simple() {
	y := `---
welcome:
  message:
  - "Good Morning"
  - "Hello World!"
`

	data, err := yaml.ToInterface(strings.NewReader(y))
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
	fmt.Println(data)
	fmt.Println(reflect.TypeOf(data))
	fmt.Println(reflect.ValueOf(data).Kind())

	// Output:
	// map[welcome:map[message:[Good Morning Hello World!]]]
	// map[string]interface {}
	// map
}

func ExampleToInterface_multi() {
	y := `---
data:
  en: "Hello World!"
---
data:
  pl: "Witaj Świecie!"
`

	data, err := yaml.ToInterface(strings.NewReader(y))
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
	fmt.Println(data)
	fmt.Println(reflect.TypeOf(data))
	fmt.Println(reflect.ValueOf(data).Kind())

	// Output:
	// [map[data:map[en:Hello World!]] map[data:map[pl:Witaj Świecie!]]]
	// []interface {}
	// slice
}
