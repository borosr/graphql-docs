package docs

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
)

const typUndefined = "Undefined"

func Init(s graphql.Schema, c Config) {
	sb := strings.Builder{}
	d := docs{c}
	d.run(s.QueryType(), &sb)
	d.run(s.MutationType(), &sb)
	d.run(s.SubscriptionType(), &sb)
	fmt.Println(sb.String())
}

func (d docs) run(o *graphql.Object, out *strings.Builder) {
	if o == nil {
		return
	}
	format := d.tabs(0)+"%s"
	if d.c.pretty {
		format += "{"
	}
	out.WriteString(fmt.Sprintf(format+"\n", strings.ToLower(o.PrivateName)))
	d.walkFields(o.Fields(), 0, out)
}

func (d docs) walkFields(fdm graphql.FieldDefinitionMap, i int, out *strings.Builder) {
	i++
	for k, f := range fdm {
		if f == nil {
			continue
		}
		d.display(i, k, f.Args, out)
		if d.c.pretty {
			if _, prim := getTypeOrValue(f.Type, true); !prim {
				out.WriteString("{")
			}
		}
		out.WriteString("\n")
		d.castField(f.Type, i, out)
	}
	if d.c.pretty {
		out.WriteString(d.tabs(i-1)+"}\n")
	}
}

func (d docs) display(i int, k string, args []*graphql.Argument, out *strings.Builder) {
	format := d.tabs(i) + "%s"
	if a := d.walkArgs(args); len(a) == 0 {
		out.WriteString(fmt.Sprintf(format, k))
	} else {
		format += "(%s)"
		out.WriteString(fmt.Sprintf(format, k, a))
	}
}


func (d docs) castField(f graphql.Output, i int, out *strings.Builder) {
	if o, ok := f.(*graphql.Object); ok {
		d.walkFields(o.Fields(), i, out)
	} else if o, ok := f.(*graphql.List); ok {
		d.castField(o.OfType, i, out)
	}
}

func (d docs) walkArgs(args []*graphql.Argument) string {
	var result = make([]string, 0)
	for _, arg := range args {
		if arg == nil {
			continue
		}
		result = d.castInputObj(arg.Type, arg.PrivateName, result)
	}
	if len(result) == 0 {
		return ""
	}
	return strings.Join(result, ",")
}

func (d docs) castInputObj(typ graphql.Input, parent string, data []string) []string {
	if o, ok := typ.(*graphql.InputObject); ok {
		fields := o.Fields()
		var r = make([]string, 0)
		for _, f := range fields {
			r = d.castInputObj(f.Type, f.PrivateName, r)
		}
		return append(data, parent+":{"+strings.Join(r, ",")+"}")
	} else {

		typeOrValue, _ := getTypeOrValue(typ, d.c.pretty)
		return append(data, parent+":"+typeOrValue)
	}
}

func getTypeOrValue(typ graphql.Input, pretty bool) (string, bool) {
	t := getType(typ)
	if pretty {
		switch t {
		case graphql.String.Name():
			return "\"\"", true
		case graphql.Int.Name():
			return "0", true
		case graphql.Float.Name():
			return "0.0", true
		case graphql.Boolean.Name():
			return "false", true
		default:
			return typUndefined, false
		}
	}
	return t, t == graphql.String.Name() || t == graphql.Int.Name() || t == graphql.Float.Name() || t == graphql.Boolean.Name()
}

func getType(typ graphql.Input) string {
	if typ == nil {
		return typUndefined
	} else {
		return typ.Name()
	}
}

func (d docs) tabs(n int) string {
	if n == 0 {
		return ""
	}
	b := strings.Builder{}
	for i := 0; i < n; i++ {
		if !d.c.pretty {
			b.WriteString("|")
		}
		b.WriteString("\t")
	}
	return b.String()
}

type docs struct {
	c Config
}

type Config struct {
	pretty bool
}
