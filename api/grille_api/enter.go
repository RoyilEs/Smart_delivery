package grille_api

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"Smart_delivery_locker/models/ctype"
	"Smart_delivery_locker/models/res"
	"fmt"
	"github.com/gin-gonic/gin"
	"sort"
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

// GrilleConfigDTO 定义格口基础信息结构体（存储总体箱子+位置的尺寸映射）
type GrilleConfigDTO struct {
	BoxCode   string  // 总体箱子编号
	BoxLength float64 // 总体箱子长度
	BoxWidth  float64 // 总体箱子宽度
	BoxHeight float64 // 总体箱子高度
}

// GenerateGrilleIDs 生成格口ID列表
func GenerateGrilleIDs(matrix []string, size int, grilleCount int) (dto []GrilleConfigDTO) {
	var (
		boxCode string
	)
	for _, m := range matrix {
		for i := range grilleCount {
			boxCode = fmt.Sprintf("%s_%d", m, i)
			switch ctype.Size(size) {
			case ctype.SizeSmall:
				dto = append(dto, GrilleConfigDTO{
					BoxCode:   boxCode,
					BoxLength: 30,
					BoxWidth:  20,
					BoxHeight: 15,
				})
			case ctype.SizeMedium:
				dto = append(dto, GrilleConfigDTO{
					BoxCode:   boxCode,
					BoxLength: 50,
					BoxWidth:  35,
					BoxHeight: 25,
				})
			case ctype.SizeLarge:
				dto = append(dto, GrilleConfigDTO{
					BoxCode:   boxCode,
					BoxLength: 80,
					BoxWidth:  60,
					BoxHeight: 40,
				})
			case ctype.SizeXLarge:
				dto = append(dto, GrilleConfigDTO{
					BoxCode:   boxCode,
					BoxLength: 120,
					BoxWidth:  80,
					BoxHeight: 60,
				})
			}
		}
	}
	return
}

type GrilleFormItemCreateRequest struct {
	LogisticsIds []string `json:"logistics_ids"`
}

// GrilleFormItemCreateView 通过订单ID创建格口
func (GrilleApi) GrilleFormItemCreateView(c *gin.Context) {
	var (
		cr    GrilleFormItemCreateRequest
		count int
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
		items = append(items, item)
	}

	// 获取空格口
	var grilles []models.Grille
	if err := global.DB.Find(&grilles, "logistics_id = ?", "").Error; err != nil {
		global.Log.Error("[error] 获取空格口失败", err)
		return
	}

	if len(grilles) < len(items) {
		res.ResultFailWithMsg("格口数量不足 请管理员添加格口", c)
		return
	}

	for _, item := range items {
		for _, grille := range grilles {
			flag := IsItemFitGrille(item, grille)
			// 成功则适配 检索下一个 放入表中
			if flag {
				global.DB.Model(&grilles[count]).Update("logistics_id", item.LogisticsId)
				count++
				break
			}
		}
	}
	res.ResultOkWithMsg(fmt.Sprintf("成功放入 %d 个订单", count), c)
}

type GrilleCreateRequest struct {
	Matrix []string `json:"matrix"`
	Size   int      `json:"size"`
	Count  int      `json:"count"`
}

func (GrilleApi) GrilleCreateView(c *gin.Context) {
	var cr GrilleCreateRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.ResultFailWithError(err, &cr, c)
		return
	}

	grilles := GenerateGrilleIDs(cr.Matrix, cr.Size, cr.Count)
	for _, grille := range grilles {
		grilleModel := models.Grille{
			GrilleId: grille.BoxCode,
			X:        grille.BoxLength,
			Y:        grille.BoxWidth,
			Z:        grille.BoxHeight,
			Size:     ctype.Size(cr.Size),
		}
		global.DB.Create(&grilleModel)
	}
	res.ResultOkWithMsg("格口创建成功!", c)
}
