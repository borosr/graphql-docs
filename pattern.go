package docs

const (
	kindField = "field"
	kindArg   = "argument"
)

type kind string

type pattern struct {
	Name        string
	Description string
	Typ         string
	Kind        kind
	Parent      string
}
