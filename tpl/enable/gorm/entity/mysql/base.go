package mysql

import "time"

// 一张表除主键外，有必要再建一个非主键索引，非聚簇索引对 count(*) 扫表友好
type CommonField struct {
	Id        int64     `gorm:"primarykey; autoIncrement"`
	CreatedAt int64     `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
