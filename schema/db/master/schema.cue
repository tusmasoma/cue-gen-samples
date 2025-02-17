package master

import (
	"list"
	"github.com/tusmasoma/cue-gen-samples/schema/db/def/spanner"
)

data: [Name=_]: spanner.#master_table & {
	name: Name
}

m_relations: {}
relations: [...spanner.#relation]
relations: list.FlattenN([for v in m_relations {v}], 1)

data_with_default_column: {for d in data {"\(d.name)": d & {
	columns: {
		created_at: {type: "timestamp", description: "生成日時"}
		updated_at: {type: "timestamp", description: "更新日時"}
	}
}}}
