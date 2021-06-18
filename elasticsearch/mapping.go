package elasticsearch

const PutMapping = `
	{
	  "properties":{
			"title":{
				"type": "text",
				"analyzer":"russian",
				"fields": {
          			"keyword": { 
            			"type": "keyword"
					}
        		}
			},
			"author": {
				"type":"keyword"
			},
			"link":{
				"type":"keyword"
			},
			"description":{
				"type":"text",
				"fields": {
          			"keyword": { 
            			"type": "keyword"
					}
        	},
				"analyzer":"russian"
			},
			"date":{
				"type":"date"
			},
			"article":{
				"type":"text",
				"analyzer":"russian"
			}
	 }
}`
