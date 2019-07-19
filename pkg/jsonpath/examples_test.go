package jsonpath_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/VirtusLab/go-extended/pkg/json"
	"github.com/VirtusLab/go-extended/pkg/jsonpath"
)

func ExampleJSONPath_EvalResults_simple() {
	js := `{
	"welcome":{
		"message":["Good Morning", "Hello World!"]
	}
}`
	expression := "{$.welcome.message[1]}"

	data, err := json.ToInterface(strings.NewReader(js))
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}

	results, err := jsonpath.New(expression).ExecuteToInterface(data)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
	fmt.Println(results)

	// Output:
	// Hello World!
}

func ExampleJSONPath_EvalResults_array_filter() {
	bs := []byte(`[
	{"key":"a","value" : "I"},
	{"key":"b","value" : "II"},
	{"key":"c","value" : "III"}
]`)

	expression := `{$[?(@.key=="b")].value}`

	data, err := json.ToInterface(bytes.NewReader(bs))
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}

	results, err := jsonpath.New(expression).ExecuteToInterface(data)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}
	fmt.Println(results)

	// Output:
	// II
}

func ExampleJSONPath_EvalResults_goessner() {
	js := `{ "store": {
    "book": [ 
      { "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      { "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      { "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      { "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}`

	expressions := []string{
		`{$.store.book[*].author}`,
		`{$..author}`,
		`{$.store..price}`,
		`{$..book[2].title}`,
		`{$..book[-1:].title}`,
		`{$..book[0,1].title}`,
		`{$..book[:2].title}`,
		`{$..book[?(@.isbn)].title}`,
		`{$..book[?(@.price < 10.0)].title}`,
		`{$.store.bicycle.color}`,
		`{$..book[?(@.author=='J. R. R. Tolkien')].title}`,
	}

	data, err := json.ToInterface(strings.NewReader(js))
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}

	for _, expression := range expressions {
		writer := new(bytes.Buffer)
		err := jsonpath.New(expression).AllowMissingKeys(true).Execute(writer, data)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
		}
		fmt.Println(writer.String())
	}

	// Output:
	// Nigel Rees Evelyn Waugh Herman Melville J. R. R. Tolkien
	// Nigel Rees Evelyn Waugh Herman Melville J. R. R. Tolkien
	// 8.95 12.99 8.99 22.99 19.95
	// Moby Dick
	// The Lord of the Rings
	// Sayings of the Century Sword of Honour
	// Sayings of the Century Sword of Honour
	// Moby Dick The Lord of the Rings
	// Sayings of the Century Moby Dick
	// red
	// The Lord of the Rings
}

func ExampleJSONPath_EvalResults_kubectl_docs() {
	js := `{
  "kind": "List",
  "items":[
    {
      "kind":"None",
      "metadata":{"name":"127.0.0.1"},
      "status":{
        "capacity":{"cpu":"4"},
        "addresses":[{"type": "LegacyHostIP", "address":"127.0.0.1"}]
      }
    },
    {
      "kind":"None",
      "metadata":{"name":"127.0.0.2"},
      "status":{
        "capacity":{"cpu":"8"},
        "addresses":[
          {"type": "LegacyHostIP", "address":"127.0.0.2"},
          {"type": "another", "address":"127.0.0.3"}
        ]
      }
    }
  ],
  "users":[
    {
      "name": "myself",
      "user": {}
    },
    {
      "name": "e2e",
      "user": {"username": "admin", "password": "secret"}
    }
  ]
}`

	expressions := []string{
		`{.kind}`,
		`{['kind']}`,
		`{..name}`,
		`{.items[*].metadata.name}`,
		`{.users[0].name}`,
		`{.items[*]['metadata.name','status.capacity']}`,
		`{.users[?(@.name=="e2e")].user.password}`,
		`{range .items[*]}[{.metadata.name},{.status.capacity}]{end}`,
	}

	data, err := json.ToInterface(strings.NewReader(js))
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
	}

	for _, expression := range expressions {
		results, err := jsonpath.New(expression).ExecuteToInterface(data)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
		}
		fmt.Println(results)
	}

	// Output:
	// List
	// List
	// [127.0.0.1 127.0.0.2 myself e2e]
	// [127.0.0.1 127.0.0.2]
	// myself
	// [127.0.0.1 127.0.0.2 map[cpu:4] map[cpu:8]]
	// secret
	// [[ 127.0.0.1 , map[cpu:4] ] [ 127.0.0.2 , map[cpu:8] ]]
}
