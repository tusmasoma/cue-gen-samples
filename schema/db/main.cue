package db

import (
	"list"
	"github.com/tusmasoma/cue-gen-samples/schema/db/user"
	"github.com/tusmasoma/cue-gen-samples/schema/db/master"
)

user_data:   user.data_with_default_column
master_data: master.data_with_default_column
relations: list.Concat([master.relations, user.relations])