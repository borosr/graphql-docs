# GraphQL Docs
A helper project for go-graphql, to build html documentation from schema

## Usage
`GO111MODULE=on go get github.com/borosr/graphql-docs`

```go
package main

import (
 "github.com/borosr/graphql-docs"
 "github.com/graphql-go/graphql"
)


func main() {
    s, _ := graphql.NewSchema(graphql.SchemaConfig{
        Query: Query,
    })
    docs.Init(s, Config{html: true}) // creates returns the documentation's content
}
```
### Config values:
```go
Config{
    sysout:   true,
    pretty:   true,
    html:     true,
    md:       true,
    json:     true,
    filename: "graphql-docs",
}
```

### Output on console (with `sysout: true, pretty: true`)
```json
query{
	graph(repo:"",filter:{value:"",filterIn:""}){
		nodes{
			attributes{
				functions
				methods{
					name
					receiver{
						name
						fields{
							name
							type
						}
					}
				}
				structs{
					name
					fields{
						name
						type
					}
				}
			}
			title
			weight
			color
			font{
				color
			}
			id
			label
		}
		edges{
			from
			to
			arrows
			attributes{
				functions
				methods{
					name
					receiver{
						name
						fields{
							name
							type
						}
					}
				}
				structs{
					name
					fields{
						name
						type
					}
				}
			}
			title
		}
	}
}
```
### Output in `html` or `md`
You can check them, after pulling this repo and a `go test .`
# Licence
[Licence](LICENCE)
