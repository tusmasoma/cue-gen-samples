package user

data: i_user_profile: {
	description: "ユーザープロフィール情報"
	columns: {
		profile_id: {pk: 1, type: "string", size: 36, description: "プロフィールID"}
		user_id: {type: "string", size: 36, description: "ユーザーID"}
		bio: {type: "string", is_max_size: true, description: "自己紹介"}
		website: {type: "string", size: 255, description: "ウェブサイトURL"}
	}
}

i_relations: user_profile_relations: [
	{
		source: {table_name: data.i_user.name, column: data.i_user.columns.user_id.name}
		target: {table_name: data.i_user_profile.name, column: data.i_user_profile.columns.user_id.name, zero: false}
	},
]
