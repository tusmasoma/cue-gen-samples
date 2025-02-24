package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"text/template"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/tusmasoma/cue-gen-samples/pkg/entity"
	"github.com/tusmasoma/cue-gen-samples/pkg/util"
	"golang.org/x/tools/imports"
)

func main() {
	infra_user_model_gen_exec()
}

func infra_user_model_gen_exec() {
	// CUE のコンテキスト作成
	ctx := cuecontext.New()

	// CUE のスキーマをロード
	instances := load.Instances(
		[]string{
			"schema/db/main.cue",
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
	data := value.LookupPath(cue.ParsePath("user_data"))
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
	if relations.Exists() {
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

	// 各エンティティごとに `model.go` & `models.go` を生成
	for _, table := range tables {
		generateModelFiles(table)
	}
}

// `model.go` & `models.go` の生成
func generateModelFiles(table *entity.Table) {
	outputDir := filepath.Join("pkg/infra/spanner/model/user", table.Name)
	os.MkdirAll(outputDir, os.ModePerm)

	// `model.go` の生成
	generateFileFromTemplate(
		"templates/db_gen/infra/spanner/model/user/model.go.tmpl",
		filepath.Join(outputDir, "model_db_gen.go"),
		table,
	)

	// `models.go` の生成
	generateFileFromTemplate(
		"templates/db_gen/infra/spanner/model/user/models.go.tmpl",
		filepath.Join(outputDir, "models_db_gen.go"),
		table,
	)
}

// テンプレートを適用して Go コードを生成
func generateFileFromTemplate(templatePath, outputFile string, table *entity.Table) {
	// テンプレートを読み込む
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error reading template file for %s: %v\n", table.Name, err)
		return
	}

	// テンプレートをパース
	tmpl, err := template.New(filepath.Base(templatePath)).Funcs(util.GetTmplFuncMap()).Parse(string(tmplContent))
	if err != nil {
		fmt.Printf("Error parsing template for %s: %v\n", table.Name, err)
		return
	}

	// テンプレートを適用
	var output bytes.Buffer
	err = tmpl.Execute(&output, table)
	if err != nil {
		fmt.Printf("Error executing template for %s: %v\n", table.Name, err)
		return
	}

	// Goコードをフォーマット
	formattedOutput, err := format.Source(output.Bytes())
	if err != nil {
		fmt.Printf("Error formatting Go code for %s: %v\n", table.Name, err)
		formattedOutput = output.Bytes()
	}

	// goimports で import の整理
	formattedOutput, err = imports.Process(outputFile, formattedOutput, nil)
	if err != nil {
		fmt.Printf("Error running goimports for %s: %v\n", table.Name, err)
	}

	// Go ファイルに保存
	err = os.WriteFile(outputFile, formattedOutput, 0644)
	if err != nil {
		fmt.Printf("Error writing file for %s: %v\n", table.Name, err)
		return
	}

	fmt.Printf("Generated: %s\n", outputFile)
}
