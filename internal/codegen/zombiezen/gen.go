package zombiezen

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"text/template"

	"github.com/delaneyj/toolbelt"
	"github.com/samber/lo"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

//go:embed templates/*.go.tpl
var templates embed.FS

func Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {

	tmpls, err := template.New("queries").Funcs(template.FuncMap{
		"goType": func(id, typ string) string {
			switch typ {
			case "text":
				return "string"
			case "integer":
				return "int64"
			default:
				panic(fmt.Sprintf("unhandled type %s for column %s", typ, id))
			}
		},
		"bindType": func(typ string) string {
			switch typ {
			case "text":
				return "BindText"
			case "integer":
				return "BindInt64"
			default:
				panic(fmt.Sprintf("unhandled type %s", typ))
			}
		},
		"Pascal": toolbelt.Pascal,
		"Camel":  toolbelt.Camel,
	}).ParseFS(templates, "templates/*.go.tpl")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	queryCtx := GenerateQueriesContext{
		PackageName: toolbelt.ToCasedString(req.Settings.Codegen.Out),
		Queries: lo.Map(req.Queries, func(q *plugin.Query, qi int) GenerateQueryContext {
			queryCtx := GenerateQueryContext{
				Name: toolbelt.ToCasedString(q.Name),
				Params: lo.Map(q.Params, func(p *plugin.Parameter, pi int) GenerateField {
					param := GenerateField{
						Column:  int(p.Number),
						Name:    toolbelt.ToCasedString(p.Column.Name),
						SQLType: toolbelt.ToCasedString(toSQLType(p.Column.Type.Name)),
						GoType:  toolbelt.ToCasedString(toGoType(p.Column.Type.Name)),
					}
					return param
				}),
			}

			if len(q.Columns) > 0 {
				queryCtx.HasResponse = true
				queryCtx.ResponseFields = lo.Map(q.Columns, func(c *plugin.Column, ci int) GenerateField {
					col := GenerateField{
						Column:  ci + 1,
						Name:    toolbelt.ToCasedString(c.Name),
						SQLType: toolbelt.ToCasedString(toSQLType(c.Type.Name)),
						GoType:  toolbelt.ToCasedString(toGoType(c.Type.Name)),
					}
					return col
				})
				queryCtx.ResponseType = toolbelt.ToCasedString(q.Name + "Res")
				queryCtx.ResponseHasMultiple = q.Cmd == ":many"
				queryCtx.SQL = q.Text
			}
			return queryCtx
		}),
	}

	buf := &bytes.Buffer{}
	if err := tmpls.ExecuteTemplate(buf, "queries.go.tpl", queryCtx); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return &plugin.GenerateResponse{
		Files: []*plugin.File{
			{
				Name:     "queries.go",
				Contents: buf.Bytes(),
			},
		},
	}, nil
}

func toSQLType(typ string) string {
	switch typ {
	case "text":
		return "text"
	case "integer":
		return "int64"
	default:
		panic(fmt.Sprintf("unhandled type %s", typ))
	}
}

func toGoType(typ string) string {
	switch typ {
	case "text":
		return "string"
	case "integer":
		return "int64"
	default:
		panic(fmt.Sprintf("unhandled type %s", typ))
	}
}

type GenerateField struct {
	Column  int
	Name    toolbelt.CasedString
	SQLType toolbelt.CasedString
	GoType  toolbelt.CasedString
}

type GenerateQueryContext struct {
	Name                toolbelt.CasedString
	Params              []GenerateField
	SQL                 string
	HasResponse         bool
	ResponseType        toolbelt.CasedString
	ResponseFields      []GenerateField
	ResponseHasMultiple bool
}

type GenerateQueriesContext struct {
	PackageName toolbelt.CasedString
	Queries     []GenerateQueryContext
}
