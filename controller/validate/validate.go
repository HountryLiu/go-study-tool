package validate

import (
	"github.com/HountryLiu/go-study-tool/model"
	"github.com/HountryLiu/go-study-tool/utils"
	"github.com/gin-gonic/gin"
)

// @Tags Validate校验
// @Summary 数据校验
// @Description 数据校验
// @accept application/json
// @Produce application/json
// @Param data body model.ValidateData true "ValidateData object"
// @Failure 200 {object} object{no=int,data=string}
// @Router /api/validate/create [post]
func Create(ctx *gin.Context) {
	var req model.ValidateData
	err := ctx.ShouldBind(&req)

	if err != nil {
		utils.Error(ctx, utils.UnmarshalError, err)
		return
	}

	if err := utils.Validate.Struct(&req); err != nil {
		utils.Error(ctx, utils.ParamsRequired, err)
		return
	}

	if err := req.Create(); err != nil {
		utils.Error(ctx, utils.DatabaseInsertError, err)
		return
	}
	utils.Success(ctx)
}
