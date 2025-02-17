package user

data: i_user: {
	description: "ユーザー"
	columns: {
		user_id: {pk: 1, type: "string", size: 36, description: "ユーザーID"}
		name:    {type: "string", size: 20, description: "名前"}
        email:   {type: "string", size: 255, description: "メールアドレス"}
	}
}
