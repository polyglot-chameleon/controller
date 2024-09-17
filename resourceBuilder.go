package crud

import (
	"fmt"
	"os"
	"reflect"
)

type resourceBuilder struct {
	data string
}

func (rb *resourceBuilder) init() {
	rb.data = "package crud\n\ntype Resource struct {\n"
}

func (rb *resourceBuilder) add(tableCol tableInfo) {
	rb.data += fmt.Sprintf("\t%s\t%s\n", tableCol.name, rb.translate(tableCol.ctype))
}

func (rb *resourceBuilder) close() {
	rb.data += "}"
}

func (rb *resourceBuilder) build() {
	os.WriteFile("resource.go", []byte(rb.data), 0666)
}

func (rb *resourceBuilder) translate(dbtype string) string {
	var out string
	switch dbtype {
	case "INTEGER":
		out = reflect.Int.String()
	case "VARCHAR":
		out = reflect.String.String()
	}
	return out
}

type tableInfo struct {
	cid      int
	name     string
	ctype    string
	notnull  bool
	dflt_val string
	pk       bool
}
