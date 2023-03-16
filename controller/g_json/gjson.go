package g_json

import (
	"fmt"

	"github.com/HountryLiu/go-study-tool/utils"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

// @Tags GJson
// @Summary GJson使用
// @Description GJson使用
// @accept application/json
// @Produce application/json
// @Failure 200 {object} object{no=int,data=string}
// @Router /api/gjson [get]
func Index(ctx *gin.Context) {
	str_json := `{
		"name": "money_page_v5_de",
		"alias": "12-29",
		"frame": {
			"cssPath": "/",
			"id": 25
		},
		"version": "temporary",
		"tags": [
			{
				"id": 68,
				"created_at": "2022-12-29T07:52:36.167214Z",
				"updated_at": "2022-12-29T07:52:36.167214Z",
				"tag_status": 1,
				"tag_name": "money_page_v5_de",
				"operator": "孙艺岭",
				"operator_id": 1651,
				"kind_id": 36
			},
			{
				"id": 83,
				"created_at": "2023-01-06T03:16:50.064809Z",
				"updated_at": "2023-01-06T03:16:50.064809Z",
				"tag_status": 1,
				"tag_name": "MoneyPage(prelander)",
				"operator": "廖玲",
				"operator_id": 1666,
				"kind_id": 65
			}
		],
		"settings": {
			"theme": {
				"id": 65,
				"name": "money_page_v5_de"
			}
		},
		"custom": {
			"theme": "default",
			"sections": [
				{
					"slots": [
						{
							"slot": "header",
							"components": [
								{
									"id": 812,
									"name": "money_page_v5_de_header",
									"type": "component",
									"cacheKey": "b7c30b6f-ddb2-4965-8e9e-0b9db210f124"
								}
							]
						},
						{
							"slot": "footer",
							"components": [
								{
									"id": 846,
									"name": "money_page_v5_de_footer",
									"type": "component",
									"cacheKey": "894b7b5d-dcb0-4bf2-b0fe-ded346c0e969"
								}
							]
						}
					],
					"type": "layout",
					"name": "money_page_v5_de_layout",
					"id": 523,
					"alias": "money_page_v5_de_layout",
					"cacheKey": "42462431-af9c-4ed5-851d-99bf0c1a132a"
				}
			],
			"plugin": [
				{}
			]
		}
	}`

	fmt.Println("frame.id ====>", gjson.Get(str_json, "frame.id"))

	for _, v := range gjson.Get(str_json, "custom.sections").Array() {
		for kk, vv := range v.Get("slots").Array() {
			fmt.Println("slot", kk, " ====>", vv.Get("slot"))
		}
	}
	utils.Success(ctx, gjson.Get(str_json, "frame.id"))
}
