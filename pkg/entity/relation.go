package entity

import "strings"

type Relations []*Relation
type Relation struct {
	Source *TableRelation `json:"source"`
	Target *TableRelation `json:"target"`
}

type TableRelation struct {
	TableName string `json:"table_name"`
	Column    string `json:"column"`
	Zero      bool   `json:"zero"`
	Many      bool   `json:"many"`
}

func (r *Relation) RelString() string {
	builder := strings.Builder{}
	if r.Source.Many {
		builder.WriteString("}")
	} else {
		builder.WriteString("|")
	}
	if r.Source.Zero {
		builder.WriteString("o")
	} else {
		builder.WriteString("|")
	}
	builder.WriteString("--")

	if r.Target.Zero {
		builder.WriteString("o")
	} else {
		builder.WriteString("|")
	}
	if r.Target.Many {
		builder.WriteString("{")
	} else {
		builder.WriteString("|")
	}
	return builder.String()
}
