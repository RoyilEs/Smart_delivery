package grille_api

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"Smart_delivery_locker/models/ctype"
	"Smart_delivery_locker/models/ctype/status"
	"Smart_delivery_locker/models/res"
	CODE "Smart_delivery_locker/models/res/code"
	"Smart_delivery_locker/service/common"
	"Smart_delivery_locker/service/email_ser"
	"Smart_delivery_locker/utils"
	"Smart_delivery_locker/utils/jwts"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"math/big"
	"sort"
	"strings"
	"time"
)

type GrilleApi struct{}

// sortDimensions 将三维尺寸按从小到大排序，返回切片
func sortDimensions(x, y, z float64) []float64 {
	dims := []float64{x, y, z}
	sort.Float64s(dims) // 升序排序
	return dims
}

// IsItemFitGrille 判断单个包裹是否能放入单个格口（支持维度旋转）
// 返回值：true=适配，false=不适配
func IsItemFitGrille(item models.Item, grille models.Grille) bool {
	// 对包裹和格口的尺寸分别排序
	pkgDims := sortDimensions(item.X, item.Y, item.Z)
	gridDims := sortDimensions(grille.X, grille.Y, grille.Z)

	// 逐维度对比：包裹的每个维度必须 ≤ 格口对应维度
	return pkgDims[0] <= gridDims[0] && pkgDims[1] <= gridDims[1] && pkgDims[2] <= gridDims[2]
}

// FindAllFitGrilles 批量匹配：返回包裹可适配的所有格口列表
func FindAllFitGrilles(item models.Item, grilles []models.Grille) []models.Grille {
	var fitGrids []models.Grille
	for _, grid := range grilles {
		if IsItemFitGrille(item, grid) {
			fitGrids = append(fitGrids, grid)
		}
	}
	return fitGrids
}

type SequenceGenerator struct {
	used map[string]bool // 存储已使用的序号
}

func NewSequenceGenerator() *SequenceGenerator {
	return &SequenceGenerator{
		used: make(map[string]bool),
	}
}

// MarkUsed 标记已经使用的序号
func (sg *SequenceGenerator) MarkUsed(seq string) {
	sg.used[strings.ToUpper(seq)] = true
}

// GenerateNext 生成下一个可用的序号
func (sg *SequenceGenerator) GenerateNext() string {
	// 从A开始查找第一个未使用的序号
	current := 0 // 0代表A，1代表B...25代表Z，26代表AA，27代表AB...
	for {
		seq := numberToSequence(current)
		if !sg.used[seq] {
			sg.used[seq] = true // 标记为已使用
			return seq
		}
		current++
	}
}

// numberToSequence 将数字转换为序号（0->A, 25->Z, 26->AA, 27->AB...）
func numberToSequence(num int) string {
	var result strings.Builder

	for {
		// 计算当前位的字符：0->A, 1->B...25->Z
		remainder := num % 26
		result.WriteByte('A' + byte(remainder))

		num = (num-remainder)/26 - 1
		if num < 0 {
			break
		}
	}

	return reverseString(result.String())
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// getSizeDimensions 根据尺寸类型返回格口的长宽高 (X, Y, Z)
func getSizeDimensions(size ctype.Size) (x, y, z float64) {
	switch size {
	case ctype.SizeSmall:
		return 30, 20, 15
	case ctype.SizeMedium:
		return 50, 35, 25
	case ctype.SizeLarge:
		return 80, 60, 40
	case ctype.SizeXLarge:
		return 120, 80, 60
	default:
		return 50, 35, 25
	}
}

// getNextCabinetLetter 获取下一个可用的柜子字母（A, B, C, ...）
// 查询数据库中已存在的 CabinetId，提取首字母，取最大值后 +1
func getNextCabinetLetter(db *gorm.DB) (string, error) {
	var letters []string
	// 假设 CabinetId 格式为 "{字母}X-1"，提取第一个字符
	if err := db.Model(&models.Grille{}).
		Distinct("SUBSTRING(cabinet_id, 1, 1)").
		Pluck("SUBSTRING(cabinet_id, 1, 1)", &letters).Error; err != nil {
		return "A", err // 出错时默认返回A
	}

	maxLetter := 'A' - 1
	for _, l := range letters {
		if len(l) > 0 {
			c := rune(l[0])
			if c >= 'A' && c <= 'Z' && c > maxLetter {
				maxLetter = c
			}
		}
	}
	next := maxLetter + 1
	if next > 'Z' {
		return "", fmt.Errorf("柜子字母已用尽（A-Z）")
	}
	return string(next), nil
}

// GenerateGrilleIDs 根据 matrix(层数)、size(尺寸类型)、count(每层格口数) 生成格口列表
// 返回生成的 Grille 对象切片（未入库），并自动计算行列位置、分配柜子字母和格口序号
func GenerateGrilleIDs(matrix int, size ctype.Size, count int) ([]models.Grille, error) {
	var grilles []models.Grille

	cabinetLetter, err := getNextCabinetLetter(global.DB)
	if err != nil {
		return nil, err
	}
	cabinetId := fmt.Sprintf("%sX-1", cabinetLetter)   // 例如 "AX-1"
	cabinetCode := fmt.Sprintf("%s区主柜", cabinetLetter) // 例如 "A区主柜"

	cols := int(math.Ceil(math.Sqrt(float64(count))))
	rows := int(math.Ceil(float64(count) / float64(cols)))

	generator := NewSequenceGenerator()
	// 标记已使用的序号（从数据库中已有的 GrilleId 提取前缀）
	var existingGrilles []models.Grille
	if err := global.DB.Find(&existingGrilles).Error; err != nil {
		return nil, err
	}
	for _, g := range existingGrilles {
		parts := strings.Split(g.GrilleId, "_")
		if len(parts) > 0 {
			generator.MarkUsed(parts[0])
		}
	}

	xDim, yDim, zDim := getSizeDimensions(size)

	for layer := 1; layer <= matrix; layer++ {
		// 每层使用一个新的序号前缀（A, B, C...）
		seq := generator.GenerateNext()

		for i := 0; i < count; i++ {
			// 计算在当前层中的行号和列号（从1开始）
			row := i/cols + 1
			col := i%cols + 1
			// 如果 row 超过了实际行数（最后一行的列数可能不足），跳过（实际不会超过，因为 rows 已按 ceil 计算）
			if row > rows {
				continue
			}

			grilleId := fmt.Sprintf("%s_%d", seq, i)

			grille := models.Grille{
				GrilleId:     grilleId,
				X:            xDim,
				Y:            yDim,
				Z:            zDim,
				Size:         size,
				CabinetId:    cabinetId,
				CabinetCode:  cabinetCode,
				MatrixRow:    row,
				MatrixColumn: col,
				Layer:        layer,
				Status:       "idle",
				Remark:       "",
			}
			grilles = append(grilles, grille)
		}
	}

	return grilles, nil
}

// GeneratePickupCode 生成纯数字取件码
// length: 取件码长度，建议6-8位
func GeneratePickupCode(length int) string {
	maxNum := big.NewInt(1)
	for i := 0; i < length; i++ {
		maxNum.Mul(maxNum, big.NewInt(10))
	}

	n, err := rand.Int(rand.Reader, maxNum)
	if err != nil {
		// 降级方案：使用时间戳
		return fmt.Sprintf("%0*d", length, time.Now().UnixNano()%1000000)
	}

	return fmt.Sprintf("%0*d", length, n)
}

type GrilleFormItemCreateRequest struct {
	LogisticsIds []string `json:"logistics_ids"`
}

type GrilleFormItemCreateResponse struct {
	Count int           `json:"count"`
	Items []models.Item `json:"list"`
}

// GrilleFormItemCreateView 通过订单ID创建格口
func (GrilleApi) GrilleFormItemCreateView(c *gin.Context) {
	var (
		cr GrilleFormItemCreateRequest
	)
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}

	// 获取需要的放入的订单
	var items []models.Item
	for _, logisticsId := range cr.LogisticsIds {
		item := models.Item{}
		global.DB.Find(&item, "logistics_id = ?", logisticsId)
		if item.LogisticsId != "" { // 确保订单存在
			items = append(items, item)
		}
	}

	// 剔除已存在格口的订单
	filteredItems := make([]models.Item, 0)
	for _, item := range items {
		if item.GrilleId == "" {
			filteredItems = append(filteredItems, item)
		} else {
			global.Log.Printf("%s 已存在格口中已剔除", item.LogisticsId)
		}
	}
	items = filteredItems

	if len(items) == 0 {
		res.ResultFailWithMsg("没有需要放入的有效订单", c)
		return
	}

	// 获取空格口
	var grilles []models.Grille
	if err := global.DB.Find(&grilles, "logistics_id = ? AND status = ?", "", status.Idle.String()).Error; err != nil {
		global.Log.Error("[error] 获取空格口失败", err)
		res.ResultFailWithMsg("获取格口失败", c)
		return
	}

	if len(grilles) < len(items) {
		res.ResultFailWithMsg(fmt.Sprintf("格口数量不足，需要%d个，可用%d个", len(items), len(grilles)), c)
		return
	}

	// 计算过期时间
	timeoutHours := global.Config.Basic.PickupTimeoutHours
	if timeoutHours <= 0 {
		timeoutHours = 72 // 默认72小时
	}
	expireTime := time.Now().Add(time.Duration(timeoutHours) * time.Hour)
	expireTimeStr := utils.ToISO8601(expireTime)

	// 创建邮件服务实例
	emailService := &email_ser.EmailService{}

	// 记录成功处理的订单
	successCount := 0
	processedItems := make([]models.Item, 0)

	for i, item := range items {
		for j, grille := range grilles {
			flag := IsItemFitGrille(item, grille)
			if flag && grille.Status == status.Idle.String() {
				pickupCode := GeneratePickupCode(global.Config.Pickup.CodeLength)
				receiverToken := generateReceiverToken(item.LogisticsId)

				// 开启事务
				tx := global.DB.Begin()

				// 更新格口信息
				if err := tx.Model(&grilles[j]).
					Updates(map[string]interface{}{
						"logistics_id": item.LogisticsId,
						"status":       status.Occupied.String(),
					}).Error; err != nil {
					tx.Rollback()
					global.Log.Errorf("更新格口失败: %v", err)
					continue
				}

				iso8601 := utils.ToISO8601(time.Now())

				// 更新物品信息 - 注意：使用指针类型需要特殊处理
				updates := map[string]interface{}{
					"grille_id":      grilles[j].GrilleId,
					"cabinet_id":     grilles[j].CabinetId,
					"cabinet_code":   grilles[j].CabinetCode,
					"grille_status":  grilles[j].Status,
					"status":         status.Stored.String(),
					"pickup_code":    pickupCode,
					"receiver_token": receiverToken,
					"expire_at":      expireTimeStr,
					"inbound_at":     iso8601,
				}

				if err := tx.Model(&items[i]).Updates(updates).Error; err != nil {
					tx.Rollback()
					global.Log.Errorf("更新物品失败: %v", err)
					continue
				}

				// 存储到Redis
				if err := storeTokenToRedis(receiverToken, item.LogisticsId, expireTime); err != nil {
					global.Log.Errorf("存储token到Redis失败: %v", err)
				}

				// 提交事务
				if err := tx.Commit().Error; err != nil {
					global.Log.Errorf("提交事务失败: %v", err)
					continue
				}

				// 重新查询更新后的订单信息
				var updatedItem models.Item
				if err := global.DB.Where("logistics_id = ?", item.LogisticsId).First(&updatedItem).Error; err == nil {
					processedItems = append(processedItems, updatedItem)
				}

				// 异步发送邮件（需要确保收件人邮箱不为空）
				if item.ReceiverEmail != "" {
					go func(itemCopy models.Item, grilleCopy models.Grille, code, expire string) {
						defer func() {
							if r := recover(); r != nil {
								global.Log.Errorf("发送邮件时发生panic: %v", r)
							}
						}()

						// 安全地获取过期时间字符串
						var expireAtStr string
						if itemCopy.ExpireAt != nil {
							expireAtStr = *itemCopy.ExpireAt
						} else {
							expireAtStr = expire
						}

						err := emailService.SendPickupNotification(
							itemCopy.ReceiverEmail,
							itemCopy.ReceiverName,
							code,
							grilleCopy.CabinetCode,
							grilleCopy.GrilleId,
							expireAtStr,
						)

						if err != nil {
							global.Log.Errorf("发送取件通知失败 [订单: %s]: %v", itemCopy.LogisticsId, err)
						} else {
							global.Log.Infof("取件通知已发送 [订单: %s, 收件人: %s]", itemCopy.LogisticsId, itemCopy.ReceiverEmail)
						}
					}(item, grille, pickupCode, expireTimeStr)
				} else {
					global.Log.Warnf("订单 %s 没有收件人邮箱，无法发送通知", item.LogisticsId)
				}

				successCount++
				break
			}
		}
	}

	// 返回结果
	if successCount == 0 {
		res.ResultFailWithMsg("没有成功放入任何订单", c)
		return
	}

	res.ResultOK(GrilleFormItemCreateResponse{
		Count: successCount,
		Items: processedItems,
	}, fmt.Sprintf("成功放入 %d 个订单，已发送邮件通知", successCount), c)
}

// generateReceiverToken 生成接收者令牌
func generateReceiverToken(logisticsId string) string {
	data := fmt.Sprintf("%s_%d", logisticsId, time.Now().UnixNano())
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// storeTokenToRedis 存储token到Redis
func storeTokenToRedis(token, logisticsId string, expireTime time.Time) error {
	key := fmt.Sprintf("receiver_token:%s", token)
	expiration := time.Until(expireTime)
	if expiration <= 0 {
		return fmt.Errorf("过期时间无效")
	}
	return global.Redis.Set(key, logisticsId, expiration).Err()
}

type GrilleCreateRequest struct {
	Matrix int    `json:"matrix"`
	Size   string `json:"size"`
	Count  int    `json:"count"`
	Remark string `json:"remark"`
}

type GrilleCreateResponse struct {
	Count   int             `json:"count"`
	Grilles []models.Grille `json:"list"`
}

// GrilleCreateView 创建格口
func (GrilleApi) GrilleCreateView(c *gin.Context) {
	var (
		cr        GrilleCreateRequest
		count     int
		newGrille []models.Grille
	)
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}
	var size ctype.Size
	switch cr.Size {
	case ctype.SizeSmall.String():
		size = ctype.SizeSmall
	case ctype.SizeLarge.String():
		size = ctype.SizeLarge
	case ctype.SizeMedium.String():
		size = ctype.SizeMedium
	case ctype.SizeXLarge.String():
		size = ctype.SizeXLarge
	}

	grilles, err := GenerateGrilleIDs(cr.Matrix, size, cr.Count)
	if err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}

	for _, grille := range grilles {
		var existing models.Grille
		result := global.DB.Where("grille_id = ?", grille.GrilleId).First(&existing)
		if result.Error == nil {
			continue
		}
		fmt.Println(grille)
		global.DB.Create(&grille)
	}

	// 重新计数并获取新创建的记录
	count = 0
	newGrille = []models.Grille{}
	for _, grille := range grilles {
		grilleModel := models.Grille{}
		result := global.DB.Where("grille_id = ?", grille.GrilleId).First(&grilleModel)
		if result.Error == nil {
			count++
			newGrille = append(newGrille, grilleModel)
		}
	}

	res.ResultOK(GrilleCreateResponse{
		Count:   count,
		Grilles: newGrille,
	}, fmt.Sprintf("成功创建 %d 个格口", count), c)
}

type ItemOutGrilleRequest struct {
	LogisticsIds []string `json:"logistics_ids"`
}

// ItemOutGrilleView 订单出格口
func (GrilleApi) ItemOutGrilleView(c *gin.Context) {
	var (
		cr      ItemOutGrilleRequest
		items   []models.Item
		grilles []models.Grille
	)
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}
	for _, id := range cr.LogisticsIds {
		item := models.Item{}
		grille := models.Grille{}
		if err := global.DB.Find(&item, "logistics_id = ?", id).Error; err != nil {
			res.ResultFailWithMsg("订单不存在", c)
			return
		}
		if err := global.DB.Find(&grille, "logistics_id = ?", id).Error; err != nil {
			res.ResultFailWithMsg("订单不在格口中", c)
			return
		}
		items = append(items, item)
		grilles = append(grilles, grille)
	}

	// 出库操作
	for i, item := range items {
		iso8601 := utils.ToISO8601(time.Now())
		global.DB.Model(&item).Update("grille_id", "").Update("status", "picked_up").Update("outbound_at", iso8601)
		global.DB.Model(&grilles[i]).Update("logistics_id", "").Update("status", "idle")
	}

	// TODO 测试阶段不做删除
	//err := global.DB.Delete(&items).Error
	//if err != nil {
	//	global.Log.Error("[error] 出库失败", err)
	//	res.ResultFailWithMsg("出库失败", c)\
	//	return
	//}

	res.ResultOK(cr, "出库成功", c)
}

// PhoneUri 获取包裹信息
type PhoneUri struct {
	Phone string `uri:"phone"`
}

type ItemResponse struct {
	models.Item
}

type ItemListRequest struct {
	models.PageInfo
}

// PhoneGetItemsView 通过手机号获取已入库的订单
func (GrilleApi) PhoneGetItemsView(c *gin.Context) {
	var cr PhoneUri
	if err := c.ShouldBindUri(&cr); err != nil {
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}

	var page ItemListRequest
	if err := c.ShouldBind(&page); err != nil {
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}

	var (
		userModel models.User
		total     int64
	)
	global.DB.Where("phone = ?", cr.Phone).Find(&userModel).Count(&total)
	if total == 0 {
		res.ResultFailWithMsg("用户不存在", c)
		return
	}
	fmt.Println(userModel)
	var (
		items        []ItemResponse
		count        int64
		setGrilleNum int64
	)

	list, _, _ := common.ComList(models.Item{ReceiverPhone: userModel.Phone}, common.Option{
		PageInfo: page.PageInfo,
	})
	fmt.Println(len(list))
	for _, item := range list {
		if item.ReceiverPhone == userModel.Phone {
			items = append(items, ItemResponse{
				Item: item,
			})
			count++
		}
	}

	for _, item := range items {
		if item.GrilleId == "" {
			setGrilleNum++
		}
	}

	res.ResultOkWithListMsg(items, count, fmt.Sprintf("%d个包裹未放入格口", setGrilleNum), c)
}

type GrilleListRequest struct {
	Count   int         `json:"count"`
	Grilles []GrilleDTO `json:"list"`
}

type GrilleDTO struct {
	models.Grille
	PickupCode string `json:"pickupCode"`
}

func (GrilleApi) GrilleListView(c *gin.Context) {
	var grillesList []models.Grille

	err := global.DB.Find(&grillesList).Error
	if err != nil {
		res.ResultFailWithError(err, nil, c)
		return
	}

	grilles := GrilleListRequest{
		Count:   len(grillesList),
		Grilles: make([]GrilleDTO, len(grillesList)),
	}

	for i, grille := range grillesList {
		grilles.Grilles[i].Grille = grille

		var item models.Item
		err := global.DB.Where("logistics_id = ?", grille.LogisticsId).First(&item).Error
		if err == nil {
			grilles.Grilles[i].PickupCode = item.PickupCode
		}
	}

	res.ResultOkWithData(grilles, c)
}

type GrilleUpdateOneRequest struct {
	models.Grille
}

// GrilleUpdateOneView 修改单个格口信息 此处id为grille_id
func (GrilleApi) GrilleUpdateOneView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	if ctype.Role(claims.Role) == ctype.PermissionUser {
		res.ResultFailWithMsg("权限不足", c)
		return
	}

	id := c.Param("id")
	var cr GrilleUpdateOneRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}

	var grille models.Grille
	err = global.DB.Take(&grille, "grille_id", cr.GrilleId).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		global.Log.Warn("格口已存在")
		res.ResultFailWithMsg("格口已存在 请重新申请格口ID", c)
		return
	}
	global.DB.Find(&grille, "grille_id = ?", id)
	// 通过尺寸重写xyz
	cr.X, cr.Y, cr.Z = getSizeDimensions(cr.Size)
	err = global.DB.Model(&grille).Where("grille_id = ?", id).Updates(cr).Error
	if err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}
	global.DB.Find(&grille, "grille_id = ?", id)
	res.ResultOkWithData(grille, c)
}

type GrilleBatchUpdateRequest struct {
	Ids    []string `json:"ids"`
	Status string   `json:"status"`
}

func (GrilleApi) GrilleUpdateBatchView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*jwts.CustomClaims)
	if ctype.Role(claims.Role) == ctype.PermissionUser || ctype.Role(claims.Role) == ctype.PermissionCourier {
		res.ResultFailWithMsg("权限不足", c)
		return
	}

	var cr GrilleBatchUpdateRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.ResultFailWithCode(CODE.ArgumentError, c)
		return
	}

	result := global.DB.Model(&models.Grille{}).Where("grille_id in ?", cr.Ids).Update("status", cr.Status)
	if result.Error != nil {
		res.ResultFailWithError(result.Error, &cr, c)
		return
	}
	res.ResultOkWithMsg(fmt.Sprintf("%d个格口状态更新成功", result.RowsAffected), c)
}
