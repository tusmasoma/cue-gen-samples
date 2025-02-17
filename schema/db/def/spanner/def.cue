package spanner

import (
	"list"
)

#user_table: #table & {
	is_user: true
}

#master_table: #table & {
	is_master: true
}

#table: {
	name:        string & =~"^[a-z0-9_]*$"
	description: string
	columns: [Name=_]: #column & {
		name: Name
	}
	interleave_in_parent?: string
	indexes?: [...#index]
	row_deletion_policy?: #row_deletion_policy
	isUseInOpenAPI:       bool | *false
	is_user:              bool | *false
	is_master:            bool | *false
	comment?:             string
	todo?:                string
}

#index: {
	keys: [...#key]
	unique: bool | *false
}

#key: {
	column: string
	desc:   bool | *false
}

#row_deletion_policy: {
	column:   string
	ttl_days: int
}

#column: {
	#columnCommon
	#columnForOpenAPI
}

#columnCommon: {
	name: =~"^[a-z0-9_]*$"
	pk?:  int
	type: "array" | "bool" | "bytes" | "date" | "float64" | "int64" | "json" | "numeric" | "string" | "timestamp"
	if list.Contains(["string", "bytes"], type) {
		is_max_size: bool | *false
		if is_max_size == false {
			size: int & >=1
		}
	}
	if type == "array" {
		array_type: string
	}
	description: string
	is_null?:    bool
}

#columnForOpenAPI: {
	isUseInOpenAPI: bool | *true
	typeOpenAPI:    string | *"string"
}

#table_relation: {
	table_name: =~"^[a-z0-9_]*$"
	column:     string
	zero:       bool | *false
	many:       bool | *false
}

#relation: {
	source: #table_relation
	target: #table_relation
}
