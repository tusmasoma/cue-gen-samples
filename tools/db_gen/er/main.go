package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/Masterminds/sprig"
	"github.com/guregu/null/zero"
)

func main() {
	user_gen_exec()
	master_gen_exec()
}

func user_gen_exec() {
	// CUE のコンテキスト作成
	ctx := cuecontext.New()

	// CUE のスキーマをロード
	instances := load.Instances(
		[]string{
			"schema/db/user/schema.cue",
			"schema/db/user/user.cue",
			"schema/db/user/user_profile.cue",
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

	// `relations` フィールドを取得
	relations := value.LookupPath(cue.ParsePath("relations"))
	if !relations.Exists() {
		fmt.Println("Warning: `relations` field not found in CUE schema")
	} else {
		// JSON 経由でデコード
		jsonBytes, err := relations.MarshalJSON()
		if err != nil {
			fmt.Println("Error marshaling CUE relations to JSON:", err)
			return
		}

		var relData Relations
		err = json.Unmarshal(jsonBytes, &relData)
		if err != nil {
			fmt.Println("Error unmarshaling JSON to Relations:", err)
			return
		}

		// 各テーブルに `relations` をマッピング
		for _, rel := range relData {
			if table, exists := tables[rel.Target.TableName]; exists {
				table.Relations = append(table.Relations, rel)
			}
		}
	}

	// SQL テンプレートの読み込み
	templateFile, err := os.ReadFile("templates/db_gen/db/er/er.puml.tmpl")
	if err != nil {
		fmt.Println("Error reading template file:", err)
		return
	}

	// テンプレートをパース
	tmpl, err := template.New("puml").Funcs(getTmplFuncMap()).Parse(string(templateFile))
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
	err = os.WriteFile("db/er/er_user_db_gen.puml", output.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing SQL file:", err)
		return
	}

	fmt.Println("SQL schema generated successfully: generated_user.sql")
}

func master_gen_exec() {
	// CUE のコンテキスト作成
	ctx := cuecontext.New()

	// CUE のスキーマをロード
	instances := load.Instances(
		[]string{
			"schema/db/master/schema.cue",
			"schema/db/master/region.cue",
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

	// SQL テンプレートの読み込み
	templateFile, err := os.ReadFile("templates/db_gen/db/er/er.puml.tmpl")
	if err != nil {
		fmt.Println("Error reading template file:", err)
		return
	}

	// テンプレートをパース
	tmpl, err := template.New("puml").Funcs(getTmplFuncMap()).Parse(string(templateFile))
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
	err = os.WriteFile("db/er/er_master_db_gen.puml", output.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing SQL file:", err)
		return
	}

	fmt.Println("SQL schema generated successfully: generated_master.sql")
}

func getTmplFuncMap() template.FuncMap {
	funcMap := sprig.TxtFuncMap()
	myFuncMap := template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		// 追加
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
