package master

data: m_regions: {
    description: "リージョン"
    columns: {
        region_id:    {pk: 1, type: "string", size: 36, description: "リージョンID"}
        region_name:  {type: "string", size: 100, description: "リージョン名"}
        country_code: {type: "string", size: 10, description: "国コード（例: 'JP', 'US'）"}
        description:  {type: "string", size: 255, description: "説明"}
    }
    indexes: [
		{
			unique: true
			keys: [
				{
					column: columns.country_code.name
				},
			]
		},
	]
}
