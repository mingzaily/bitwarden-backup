package model

import "time"

// BackupTask 备份任务配置
type BackupTask struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	SourceServerID uint      `gorm:"not null" json:"source_server_id"`
	CronExpression string    `gorm:"size:100" json:"cron_expression"`
	Enabled        bool      `gorm:"default:true" json:"enabled"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// 关联
	SourceServer ServerConfig        `json:"source_server"`
	Destinations []BackupDestination `gorm:"many2many:task_destinations;" json:"destinations"`
}

// TaskRequest 任务请求 DTO
type TaskRequest struct {
	Name           string `json:"name"`
	SourceServerID uint   `json:"source_server_id"`
	CronExpression string `json:"cron_expression"`
	Enabled        *bool  `json:"enabled"`
	DestinationIDs []uint `json:"destination_ids"`
}

// TaskResponse 任务响应 DTO（隐藏敏感数据）
type TaskResponse struct {
	ID             uint                    `json:"id"`
	Name           string                  `json:"name"`
	SourceServerID uint                    `json:"source_server_id"`
	CronExpression string                  `json:"cron_expression"`
	Enabled        bool                    `json:"enabled"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
	SourceServer   ServerResponse          `json:"source_server"`
	Destinations   []DestinationResponse   `json:"destinations"`
}

// ToResponse 转换为响应结构
func (t *BackupTask) ToResponse() TaskResponse {
	dests := make([]DestinationResponse, len(t.Destinations))
	for i, d := range t.Destinations {
		dests[i] = d.ToResponse()
	}
	return TaskResponse{
		ID:             t.ID,
		Name:           t.Name,
		SourceServerID: t.SourceServerID,
		CronExpression: t.CronExpression,
		Enabled:        t.Enabled,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
		SourceServer:   t.SourceServer.ToResponse(),
		Destinations:   dests,
	}
}
