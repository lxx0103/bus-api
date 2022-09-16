package response

import "github.com/gin-gonic/gin"

type ErrorRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SuccessRes struct {
	Data interface{} `json:"data"`
}

type ListRes struct {
	PageID   int         `json:"page_id"`
	PageSize int         `json:"page_size"`
	Count    int         `json:"count"`
	Data     interface{} `json:"data"`
}

func ResponseList(c *gin.Context, page int, page_size int, count int, data interface{}) {
	var res ListRes
	res.PageID = page
	res.PageSize = page_size
	res.Count = count
	res.Data = data
	c.AbortWithStatusJSON(200, res)
}

func Response(c *gin.Context, data interface{}) {
	var res SuccessRes
	res.Data = data
	c.AbortWithStatusJSON(200, res)
}

func ResponseError(c *gin.Context, code string, err error) {
	var res ErrorRes
	res.Code = code
	res.Message = err.Error()
	if res.Message == "sql: no rows in result set" {
		res.Message = "Data Not Exist"
	}
	c.AbortWithStatusJSON(400, res)
}

func ResponseUnauthorized(c *gin.Context, code string, err error) {
	var res ErrorRes
	res.Code = code
	res.Message = err.Error()
	c.AbortWithStatusJSON(401, res)
}
