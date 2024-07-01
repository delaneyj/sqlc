package {{.PackageName.Lower}}

import (
    "fmt"
    "zombiezen.com/go/sqlite"
)

{{- define "fillResponse"}}
row := {{.ResponseType.Pascal}}{
{{- range .ResponseFields}}
    {{.Name.Pascal}} : stmt.Column{{.SQLType.Pascal}}({{.Column}}),
{{- end}}
}
{{- end}}

{{- range .Queries}}
    {{if .HasResponse}}
type {{.ResponseType.Pascal}} struct {
        {{- range .ResponseFields}}
    {{.Name.Pascal}} {{.GoType.Lower}} `json:"{{.Name.Lower}}"`
        {{- end}}
}
    {{- end}}

func {{.Name.Pascal}}(
    tx *sqlite.Conn,
    {{- range .Params}}
    {{.Name.Lower}} {{.GoType.Lower}},
    {{- end}}
) (
    {{- if .HasResponse}}
    res {{if .ResponseHasMultiple}}[]{{else}}*{{end}}{{.ResponseType.Pascal}},
    {{- end}}
    err error,
) {
    // Prepare statement into cache
    stmt := tx.Prep(`{{.SQL}}`)
    defer stmt.Reset()

    {{ if len .Params -}}
    // Bind parameters
    {{- range .Params}}
    stmt.Bind{{.SQLType.Pascal}}({{.Column}}, {{.Name.Lower}})
    {{- end}}
    {{- end}}

    // Execute query
    {{- if .HasResponse}}
        {{- if .ResponseHasMultiple}}
    for {
        if hasRow, err := stmt.Step(); err != nil {
            return res, fmt.Errorf("failed to execute {{.Name.Lower}} SQL: %w", err)
        } else if !hasRow {
            break
        }
            {{template "fillResponse" .}}

        res = append(res, row)
    }
        {{- else}}
    if hasRow, err := stmt.Step(); err != nil {
        return res, err
    } else if hasRow {
            {{template "fillResponse" .}}
        res = &row
    }
        {{- end}}
    {{- else}}
    if _, err := stmt.Step(); err != nil {
        return fmt.Errorf("failed to execute {{.Name.Lower}} SQL: %w", err)
    }
    {{- end}}

    {{ if .HasResponse -}}
    return res, nil
    {{ else -}}
    return nil
    {{ end -}}
}
{{- end}}
