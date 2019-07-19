package json_test

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/VirtusLab/go-extended/pkg/json"
)

func ExampleToInterface_empty() {
	js := ``

	data, err := json.ToInterface(strings.NewReader(js))
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
	js := `{
	"welcome":{
		"message":["Good Morning", "Hello World!"]
	}
}`

	data, err := json.ToInterface(strings.NewReader(js))
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
