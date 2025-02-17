package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/tusmasoma/cue-gen-samples/pkg/entity"
	"github.com/tusmasoma/cue-gen-samples/pkg/util"
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
	var tables map[string]*entity.Table
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

		var relData entity.Relations
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
	tmpl, err := template.New("puml").Funcs(util.GetTmplFuncMap()).Parse(string(templateFile))
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
	var tables map[string]*entity.Table
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
	tmpl, err := template.New("puml").Funcs(util.GetTmplFuncMap()).Parse(string(templateFile))
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
