package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"text/template"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/Masterminds/sprig"
	"github.com/guregu/null/zero"
)

func main() {
	// CUE のコンテキスト作成
	ctx := cuecontext.New()

	// CUE のスキーマをロード
	instances := load.Instances(
		[]string{
			"schema/db/user/schema.cue",
			"schema/db/user/user.cue",
		},
		nil,
	)
	if len(instances) == 0 {
		fmt.Println("No CUE files found")
		return
	}

	// インスタンスを解析
	value := ctx.BuildInstance(instances[0])
	if value.Err() != nil {
		fmt.Println("Error building CUE instance:", value.Err())
		return
	}

	// `data` フィールドを取得
	data := value.LookupPath(cue.ParsePath("data"))
	if !data.Exists() {
		fmt.Println("Error: `data` field not found in CUE schema")
		return
	}

	// Go の構造体に変換
	var tables map[string]*Table
	err := data.Decode(&tables)
	if err != nil {
		fmt.Println("Error decoding CUE data:", err)
		return
	}

	log.Print(tables)

	// SQL テンプレートの読み込み
	templateFile, err := os.ReadFile("templates/db_gen/db/ddl/user.sql.tmpl")
	if err != nil {
		fmt.Println("Error reading template file:", err)
		return
	}

	// テンプレートをパース
	tmpl, err := template.New("sql").Funcs(getTmplFuncMap()).Parse(string(templateFile))
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	// SQL の出力
	var output bytes.Buffer
	err = tmpl.Execute(&output, tables)
	if err != nil {
		fmt.Println("Error generating SQL:", err)
		return
	}

	// SQL をファイルに保存
	err = os.WriteFile("db/ddl/user_db_gen.sql", output.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing SQL file:", err)
		return
	}

	fmt.Println("SQL schema generated successfully: generated_user.sql")
}

func getTmplFuncMap() template.FuncMap {
	funcMap := sprig.TxtFuncMap()
	myFuncMap := template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		// 追加s
	}
	for i := range myFuncMap {
		funcMap[i] = myFuncMap[i]
	}
	return funcMap
}

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

func (c *Column) GetName() string {
	return c.Name
}

func (c *Column) HasSize() bool {
	return c.Size != nil || c.IsMaxSize
}

func (c *Column) IsPrimaryKey() bool {
	return c.Pk != nil
}

func (c Column) IsBool() bool {
	return c.Type == "bool"
}

func (c Column) IsFloat64() bool {
	return c.Type == "float64"
}

func (c Column) IsInt() bool {
	return false
}

func (c Column) IsInt64() bool {
	return c.Type == "int64"
}

func (c Column) IsString() bool {
	return c.Type == "string"
}

func (c Column) IsNullString() bool {
	return c.Type == "string" && c.IsNullable()
}

func (c Column) IsEnum() bool {
	return c.Type == "enum"
}

func (c Column) IsTime() bool {
	return c.Type == "timestamp" || c.Type == "date"
}

func (c Column) IsCreatedAtColumn() bool {
	return c.Name == "created_at"
}

func (c Column) IsUpdatedAtColumn() bool {
	return c.Name == "updated_at"
}

func (c Column) IsSoftDeleteColumn() bool {
	return c.Name == "deleted_at"
}

func (c Column) IsNumeric() bool {
	return c.Type == "numeric"
}

func (c Column) IsNullable() bool {
	return zero.BoolFromPtr(c.IsNull).ValueOrZero() || c.IsSoftDeleteColumn()
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
