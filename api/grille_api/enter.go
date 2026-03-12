package grille_api

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"Smart_delivery_locker/models/ctype"
	"Smart_delivery_locker/models/res"
	"Smart_delivery_locker/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"sort"
	"strings"
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

// GenerateGrilleIDs 生成格口ID列表
// 参数: matrix  矩阵总体箱子首位编号
// 尺寸: size 尺寸大小参考 ctype.grille 中的尺寸
// 格口数量: grilleCount  矩阵中的箱子数量
func GenerateGrilleIDs(matrix int, size int, grilleCount int) (dto []GrilleConfigDTO) {
	var (
		boxCode string
		grilles []models.Grille
		seqs    []string
	)
	// 创建序号生成器 依旧无脑责任链模式 到处拉屎
	generator := NewSequenceGenerator()

	if err := global.DB.Find(&grilles).Error; err != nil {
		global.Log.Error("[error] 获取格口失败", err)
		return
	}
	// 获取已使用的序号
	for _, m := range grilles {
		split := strings.Split(m.GrilleId, "_")
		seqs = append(seqs, split[0])
	}

	for _, seq := range seqs {
		generator.MarkUsed(seq)
	}

	for range matrix {
		seq := generator.GenerateNext()
		for i := range grilleCount {
			boxCode = fmt.Sprintf("%s_%d", seq, i)
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

	for _, item := range items {
		if item.GrilleId != "" {
			newItem := utils.DeleteByValue(items, item)
			items = newItem
			global.Log.Printf("%s存在在格口中已剔除", item.LogisticsId)
		}
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
				global.DB.Model(&items[count]).Update("grille_id", grilles[count].GrilleId)
				count++
				break
			}
		}
	}
	res.ResultOkWithMsg(fmt.Sprintf("成功放入 %d 个订单", count), c)
}

type GrilleCreateRequest struct {
	Matrix int `json:"matrix"`
	Size   int `json:"size"`
	Count  int `json:"count"`
}

// GrilleCreateView 创建格口
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
		global.DB.Model(&item).Update("grille_id", "")
		global.DB.Model(&grilles[i]).Update("logistics_id", "")
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
