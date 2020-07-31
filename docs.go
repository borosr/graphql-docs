package docs

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
)

func Init(s graphql.Schema) {
	run(s.QueryType())
	run(s.MutationType())
	run(s.SubscriptionType())
}

func run(o *graphql.Object) {
	if o == nil {
		return
	}
	fmt.Printf(tabs(0)+"%s\n", o.PrivateName)
	walkFields(o.Fields(), 0)
}

func walkFields(fdm graphql.FieldDefinitionMap, i int) {
	i++
	for k, f := range fdm {
		if f == nil {
			continue
		}
		display(i, k, f.Args)
		castField(f.Type, i)
	}
}

func display(i int, k string, args []*graphql.Argument) {
	format := tabs(i) + "%s"
	if a := walkArgs(args); len(a) == 0 {
		fmt.Printf(format+"\n", k)
	} else {
		fmt.Printf(format+"(%s)\n", k, a)
	}
}

func walkArgs(args []*graphql.Argument) string {
	var result = make([]string, 0)
	for _, arg := range args {
		if arg == nil {
			continue
		}
		result = castInputObj(arg.Type, arg.PrivateName, result)
	}
	if len(result) == 0 {
		return ""
	}
	return strings.Join(result, ",")
}

func castInputObj(typ graphql.Input, parent string, d []string) []string {
	if o, ok := typ.(*graphql.InputObject); ok {
		fields := o.Fields()
		var r = make([]string, 0)
		for _, f := range fields {
			r = castInputObj(f.Type, f.PrivateName, r)
		}
		return append(d, parent+":{"+strings.Join(r, ",")+"}")
	} else {
		var t string
		if typ == nil {
			t = "Undefined"
		} else {
			t = typ.Name()
		}
		return append(d, parent+":"+t)
	}
}

func castField(f graphql.Output, i int) {
	if o, ok := f.(*graphql.Object); ok {
		walkFields(o.Fields(), i)
	} else if o, ok := f.(*graphql.List); ok {
		castField(o.OfType, i)
	}
}

func tabs(n int) string {
	if n == 0 {
		return ""
	}
	b := strings.Builder{}
	for i := 0; i < n; i++ {
		b.WriteString("|\t")
	}
	return b.String()
}
