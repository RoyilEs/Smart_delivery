package tasks

import (
	"Smart_delivery_locker/global"
	"Smart_delivery_locker/models"
	"Smart_delivery_locker/models/ctype/action"
	"Smart_delivery_locker/models/ctype/status"
	"Smart_delivery_locker/utils"
	"fmt"
	"time"
)

// StartTokenExpiryChecker 启动token过期检测定时任务
func StartTokenExpiryChecker() {
	// 立即执行一次
	go checkAndProcessExpiredTokens()

	// 每10分钟检查一次
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			checkAndProcessExpiredTokens()
		}
	}()
}

// checkAndProcessExpiredTokens 检查并处理过期的token
func checkAndProcessExpiredTokens() {
	defer func() {
		if r := recover(); r != nil {
			global.Log.Errorf("过期检测任务panic: %v", r)
		}
	}()

	var expiredItems []models.Item

	// 查询所有已过期但未出库的物品
	now := time.Now()
	nowStr := utils.ToISO8601(now)

	if err := global.DB.Where("receiver_token IS NOT NULL AND expire_at <= ? AND status != ?",
		nowStr, status.PickedUp).
		Find(&expiredItems).Error; err != nil {
		global.Log.Errorf("查询过期物品失败: %v", err)
		return
	}

	if len(expiredItems) == 0 {
		return
	}

	global.Log.Infof("发现 %d 个过期待处理的订单，开始自动出库", len(expiredItems))

	for _, item := range expiredItems {
		// 处理过期订单出库
		if err := processExpiredItem(item); err != nil {
			global.Log.Errorf("处理过期订单失败 [LogisticsId: %s]: %v", item.LogisticsId, err)
		}
	}
}

// processExpiredItem 处理过期的物品（自动出库）
func processExpiredItem(item models.Item) error {
	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新物品状态为已出库
	outboundTime := utils.ToISO8601(time.Now())
	if err := tx.Model(&item).Updates(map[string]interface{}{
		"status":         status.PickedUp,
		"outbound_at":    outboundTime,
		"receiver_token": nil, // 清空token
		"expire_at":      nil,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新物品状态失败: %v", err)
	}

	// 释放格口
	if item.GrilleId != "" {
		if err := tx.Model(&models.Grille{}).
			Where("grille_id = ?", item.GrilleId).
			Updates(map[string]interface{}{
				"logistics_id": "",
				"status":       status.Idle.String(),
			}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("释放格口失败: %v", err)
		}
	}

	// 从Redis中删除token
	if item.ReceiverToken != nil {
		key := fmt.Sprintf("receiver_token:%s", *item.ReceiverToken)
		if err := global.Redis.Del(key).Err(); err != nil {
			global.Log.Warnf("删除Redis token失败: %v", err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	global.Log.Infof("订单自动出库成功 [LogisticsId: %s, GrilleId: %s, 存放时长: %s]",
		item.LogisticsId,
		item.GrilleId,
		calculateDuration(item.InboundAt, outboundTime))

	// 记录操作日志
	recordAutoOutboundLog(item)

	return nil
}

// calculateDuration 计算存放时长
func calculateDuration(inboundAt, outboundAt string) string {
	inTime, err1 := time.Parse(time.RFC3339, inboundAt)
	outTime, err2 := time.Parse(time.RFC3339, outboundAt)
	if err1 != nil || err2 != nil {
		return "未知"
	}
	duration := outTime.Sub(inTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	return fmt.Sprintf("%d小时%d分钟", hours, minutes)
}

// recordAutoOutboundLog 记录自动出库日志
func recordAutoOutboundLog(item models.Item) {
	log := models.Package{
		LogisticsId: item.LogisticsId,
		Action:      action.PickedUp.String(),
		Operator:    item.SenderEmail,
		CreatedAt:   utils.ToISO8601(time.Now()),
		Detail:      "超过取件时限自动出库",
	}

	if err := global.DB.Create(&log).Error; err != nil {
		global.Log.Errorf("记录自动出库日志失败: %v", err)
	}
}
