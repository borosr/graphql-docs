package docs

import (
	"fmt"
	"log"
	"strings"

	"github.com/graphql-go/graphql"
)

const (
	typUndefined    = "Undefined"
	defaultFilename = "documentation"
)

func Init(s graphql.Schema, c Config) string {
	d := &docs{c: c, sOut: &strings.Builder{}, patterns: make([]pattern, 0)}
	d.run(s.QueryType())
	d.run(s.MutationType())
	d.run(s.SubscriptionType())
	out := d.sOut.String()
	if c.sysout {
		fmt.Println(out)
	}
	if c.html {
		if err := d.createHtml(); err != nil {
			log.Fatal(err)
		}
		return d.html
	}
	if c.md {
		if err := d.createMd(); err != nil {
			log.Fatal(err)
		}
		return d.md
	}
	return out
}

func (d *docs) run(o *graphql.Object) {
	if o == nil {
		return
	}
	format := d.tabs(0) + "%s"
	if d.c.pretty {
		format += "{"
	}
	d.sOut.WriteString(fmt.Sprintf(format+"\n", strings.ToLower(o.PrivateName)))
	d.walkFields(o.Fields(), 0)
}

func (d *docs) walkFields(fdm graphql.FieldDefinitionMap, i int) {
	i++
	for k, f := range fdm {
		if f == nil {
			continue
		}
		d.addPattern(pattern{
			Name:        f.Name,
			Description: f.Description,
			Typ:         getType(f.Type),
			Kind:        kindField,
		})
		d.display(i, k, f)
		if d.c.pretty {
			if _, prim := getTypeOrValue(f.Type, true); !prim {
				d.sOut.WriteString("{")
			}
		}
		d.sOut.WriteString("\n")
		d.castField(f.Type, i)
	}
	if d.c.pretty {
		d.sOut.WriteString(d.tabs(i-1) + "}\n")
	}
}

func (d *docs) display(i int, k string, f *graphql.FieldDefinition) {
	format := d.tabs(i) + "%s"
	if a := d.walkArgs(f.Args, f.Name); len(a) == 0 {
		d.sOut.WriteString(fmt.Sprintf(format, k))
	} else {
		format += "(%s)"
		d.sOut.WriteString(fmt.Sprintf(format, k, a))
	}
}

func (d *docs) castField(f graphql.Output, i int) {
	if o, ok := f.(*graphql.Object); ok {
		d.walkFields(o.Fields(), i)
	} else if o, ok := f.(*graphql.List); ok {
		d.castField(o.OfType, i)
	}
}

func (d *docs) walkArgs(args []*graphql.Argument, fieldName string) string {
	var result = make([]string, 0)
	for _, arg := range args {
		if arg == nil {
			continue
		}
		result = d.castInputObj(arg.Type, arg.PrivateName, arg.PrivateDescription, fieldName, result)
	}
	if len(result) == 0 {
		return ""
	}
	return strings.Join(result, ",")
}

func (d *docs) castInputObj(typ graphql.Input, parentName, parentDescription, fieldName string, data []string) []string {
	d.addPattern(pattern{
		Name:        parentName,
		Description: parentDescription,
		Typ:         getType(typ),
		Kind:        kindArg,
		Parent:      fieldName,
	})
	if o, ok := typ.(*graphql.InputObject); ok {
		fields := o.Fields()
		var r = make([]string, 0)
		for _, f := range fields {
			r = d.castInputObj(f.Type, f.PrivateName, f.PrivateDescription, fieldName, r)
		}
		return append(data, parentName+":{"+strings.Join(r, ",")+"}")
	} else {
		typeOrValue, _ := getTypeOrValue(typ, d.c.pretty)
		return append(data, parentName+":"+typeOrValue)
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
	c        Config
	sOut     *strings.Builder
	md       string
	html     string
	patterns []pattern
}

type Config struct {
	sysout   bool
	pretty   bool
	html     bool
	md       bool
	filename string
}

func (d *docs) addPattern(h pattern) {
	for _, t := range d.patterns {
		if t.Name == h.Name && t.Description == h.Description {
			return
		}
	}
	d.patterns = append(d.patterns, h)
}

func (d *docs) createHtml() error {
	var err error
	d.html, err = buildHtml(d.patterns, d.c.filename)
	return err
}

func (d *docs) createMd() error {
	var err error
	d.md, err = buildMd(d.patterns, d.c.filename)
	return err
}
