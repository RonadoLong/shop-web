package logic

import (
	"shop-web/common/commonUtils"
	"shop-web/common/dbutil"
	"shop-web/module/service/model"

	"github.com/gin-gonic/gin"
)

func FindAreatListByParentId(c *gin.Context) {

	areaname := c.Param("areaname")

	var areaId int

	if areaname == "normal" {
		areaId = 1
	} else {
		var area model.TohAreaUsa
		err := dbutil.DB.Table("toh_areas_usa").Where("area_name = ?", areaname).Find(&area).Error
		if err != nil {
			Logger.Error(err.Error())
			commonUtils.CreateError(c)
			return
		}
		areaId = area.AreaId
	}

	var areaList []model.TohAreaUsa
	err := dbutil.DB.Table("toh_areas_usa").
		Where("parent_id = ?", areaId).
		Find(&areaList).Error

	if err != nil {
		Logger.Error(err.Error())
		commonUtils.CreateError(c)
		return
	}

	var areaRespList = make([]model.TohAreaUsaResp, len(areaList))
	for _, area := range areaList {
		areaResp := model.TohAreaUsaResp{}
		areaResp.Label = area.AreaName
		areaResp.Value = area.AreaName
		areaRespList = append(areaRespList, areaResp)
	}
	commonUtils.CreateSuccess(c, areaRespList)
}
