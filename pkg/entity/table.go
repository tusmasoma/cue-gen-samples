package entity

import (
	"sort"

	"github.com/guregu/null/zero"
)

type Tables []*Table
type Table struct {
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	ColumnMap          map[string]*Column `json:"columns"`
	InterleaveInParent string             `json:"interleave_in_parent"`
	RowDeletionPolicy  *RowDeletionPolicy `json:"row_deletion_policy"`
	Indexes            Indexes            `json:"indexes"`

	IsUser   bool `json:"is_user"`
	IsMaster bool `json:"is_master"`

	Todo    string `json:"todo"`
	Comment string `json:"comment"`

	Relations Relations `json:"relations"`
}

func (t *Table) GetName() string {
	return t.Name
}

type Columns []*Column
type Column struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Pk          *int64 `json:"pk"`
	Description string `json:"description"`

	Size      *int    `json:"size"`
	IsMaxSize bool    `json:"is_max_size"`
	ArrayType *string `json:"array_type"`
	IsNull    *bool   `json:"is_null"`
}

type Indexes []*Index
type Index struct {
	Keys     Keys `json:"keys"`
	IsUnique bool `json:"unique"`
}

type Keys []*Key
type Key struct {
	Column string `json:"column"`
	Desc   bool   `json:"desc"`
}

type RowDeletionPolicy struct {
	Column  string `json:"column"`
	TtlDays int    `json:"ttl_days"`
}

func (c *Column) HasSize() bool {
	return c.Size != nil || c.IsMaxSize
}

func (c *Column) IsPrimaryKey() bool {
	return c.Pk != nil
}

func (c Column) IsSoftDeleteColumn() bool {
	return c.Name == "deleted_at"
}

func (c Column) IsNullable() bool {
	return zero.BoolFromPtr(c.IsNull).ValueOrZero() || c.IsSoftDeleteColumn()
}

func (c Column) IsCreatedAtColumn() bool {
	return c.Name == "created_at"
}

func (c Column) IsUpdatedAtColumn() bool {
	return c.Name == "updated_at"
}

func (c Column) SQLType() string {
	switch c.Type {
	case "array":
		return "array<" + *c.ArrayType + ">"
	case "enum":
		return "int64"
	default:
		return c.Type
	}
}

func (t *Table) Columns() Columns {
	res := make(Columns, 0, len(t.ColumnMap))
	for i := range t.ColumnMap {
		res = append(res, t.ColumnMap[i])
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

func (t *Table) ColumnsWithoutPrimaryKeys() Columns {
	res := make(Columns, 0, len(t.ColumnMap))
	columns := t.Columns()
	for i := range columns {
		if !columns[i].IsPrimaryKey() {
			res = append(res, columns[i])
		}
	}
	return res
}

func (t *Table) PrimaryKeys() Columns {
	res := make(Columns, 0, 2)
	columns := t.Columns()
	for i := range columns {
		if columns[i].IsPrimaryKey() {
			res = append(res, columns[i])
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return *res[i].Pk < *res[j].Pk
	})
	return res
}

func (c *Column) GoType() string {
	// TODO nullable
	switch c.Type {
	case "array":
		// TODO: array対応(mepでは現状使用していない)
		panic("unexpected type")
	case "bool":
		if c.IsNullable() {
			panic("bool is not nullable")
		}
		return "bool"
	case "bytes":
		if c.IsNullable() {
			panic("bytes is not nullable")
		}
		return "[]byte"
	case "date":
		if c.IsNullable() {
			panic("date is not nullable")
		}
		return "time.Time"
	case "float64":
		if c.IsNullable() {
			panic("float64 is not nullable")
		}
		return "float64"
	case "int64":
		if c.IsNullable() {
			panic("int64 is not nullable")
		}
		return "int64"
	case "json":
		// TODO: json対応(mepでは現状使用していない)
		panic("unexpected type")
	case "numeric":
		if c.IsNullable() {
			panic("numeric is not nullable")
		}
		return "big.Rat"
	case "string":
		if c.IsNullable() {
			return "zeronull.String"
		}
		return "string"
	case "timestamp":
		if c.IsSoftDeleteColumn() {
			return "spanner.NullTime"
		} else {
			return "time.Time"
		}
	default:
		panic("unexpected type")
	}
}
