package packages_api

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"Smart_delivery_locker/models/res"
	"github.com/gin-gonic/gin"
)

// PackageLogVO 包裹日志响应结构
type PackageLogVO struct {
	ID          uint   `json:"id"`
	LogisticsId string `json:"logistics_id"`
	Action      string `json:"action"`
	ActionText  string `json:"action_text"`
	Status      string `json:"status"`
	Detail      string `json:"detail"`
	Operator    string `json:"operator"`
	OperatorId  string `json:"operator_id"`
	CreatedAt   string `json:"created_at"`
}

// GetPackageLogsResponse 获取包裹日志响应
type GetPackageLogsResponse struct {
	Count int            `json:"count"`
	List  []PackageLogVO `json:"list"`
}

// PackageLogsQuery 日志查询参数
type PackageLogsQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Action   string `form:"action"`
}

// GetPackageLogs 获取包裹日志
func (PackageApi) GetPackageLogs(c *gin.Context) {
	// 获取包裹ID
	logisticsId := c.Query("id")
	if logisticsId == "" {
		res.ResultFailWithMsg("包裹ID不能为空", c)
		return
	}

	// 获取查询参数
	var query PackageLogsQuery
	if err := c.BindQuery(&query); err != nil {
		res.ResultFailWithError(err, &query, c)
		return
	}

	// 设置默认分页参数
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 20
	}
	offset := (query.Page - 1) * query.PageSize

	// 查询日志列表
	var logs []models.Package
	var total int64

	db := global.DB.Model(&models.Package{}).Where("logistics_id = ?", logisticsId)

	// 操作类型筛选
	if query.Action != "" {
		db = db.Where("action = ?", query.Action)
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		global.Log.Errorf("查询包裹日志总数失败: %v", err)
		res.ResultFailWithMsg("查询失败", c)
		return
	}

	// 查询日志列表
	if err := db.Order("created_at DESC, id DESC").
		Limit(query.PageSize).
		Offset(offset).
		Find(&logs).Error; err != nil {
		global.Log.Errorf("查询包裹日志列表失败: %v", err)
		res.ResultFailWithMsg("查询失败", c)
		return
	}

	// 转换为VO
	list := make([]PackageLogVO, 0, len(logs))
	for _, log := range logs {
		vo := PackageLogVO{
			ID:          log.ID,
			LogisticsId: log.LogisticsId,
			Action:      log.Action,
			ActionText:  getActionText(log.Action),
			Status:      log.Status,
			Detail:      log.Detail,
			Operator:    log.Operator,
			OperatorId:  log.OperatorId,
			CreatedAt:   log.CreatedAt,
		}
		list = append(list, vo)
	}

	res.ResultOkWithData(GetPackageLogsResponse{
		Count: int(total),
		List:  list,
	}, c)
}

// GetPackageLogsAll 获取包裹所有日志（不分页）
func (PackageApi) GetPackageLogsAll(c *gin.Context) {
	logisticsId := c.Param("id")
	if logisticsId == "" {
		res.ResultFailWithMsg("包裹ID不能为空", c)
		return
	}

	var logs []models.Package
	if err := global.DB.Where("logistics_id = ?", logisticsId).
		Order("created_at ASC, id ASC").
		Find(&logs).Error; err != nil {
		global.Log.Errorf("查询包裹日志失败: %v", err)
		res.ResultFailWithMsg("查询失败", c)
		return
	}

	list := make([]PackageLogVO, 0, len(logs))
	for _, log := range logs {
		vo := PackageLogVO{
			ID:          log.ID,
			LogisticsId: log.LogisticsId,
			Action:      log.Action,
			ActionText:  getActionText(log.Action),
			Status:      log.Status,
			Detail:      log.Detail,
			Operator:    log.Operator,
			OperatorId:  log.OperatorId,
			CreatedAt:   log.CreatedAt,
		}
		list = append(list, vo)
	}

	res.ResultOkWithData(GetPackageLogsResponse{
		Count: len(list),
		List:  list,
	}, c)
}

func (PackageApi) GetLogsList(c *gin.Context) {
	var logs []models.Package
	global.DB.Find(&logs)

	list := make([]PackageLogVO, 0, len(logs))
	for _, log := range logs {
		vo := PackageLogVO{
			ID:          log.ID,
			LogisticsId: log.LogisticsId,
			Action:      log.Action,
			ActionText:  getActionText(log.Action),
			Status:      log.Status,
			Detail:      log.Detail,
			Operator:    log.Operator,
			OperatorId:  log.OperatorId,
			CreatedAt:   log.CreatedAt,
		}
		list = append(list, vo)
	}
	res.ResultOkWithData(GetPackageLogsResponse{
		Count: len(list),
		List:  list,
	}, c)
}

// getActionText 获取操作类型中文描述
func getActionText(action string) string {
	switch action {
	case models.PackageActionCreate:
		return "创建订单"
	case models.PackageActionUpdate:
		return "更新订单"
	case models.PackageActionStore:
		return "存入柜子"
	case models.PackageActionPickup:
		return "取件"
	case models.PackageActionExpire:
		return "过期"
	case models.PackageActionAutoOutbound:
		return "自动出库"
	case models.PackageActionManualOutbound:
		return "手动出库"
	case models.PackageActionDelete:
		return "删除订单"
	case models.PackageActionNotify:
		return "发送通知"
	default:
		return action
	}
}
