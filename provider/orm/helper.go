package orm

import (
	"cvgo/kit/castkit"
	"database/sql"
)

// 扫描获取结果集
func ScanSqlRows(rows *sql.Rows) []map[string]*castkit.GoodleVal {
	if rows == nil {
		return nil
	}
	defer rows.Close()           // 扫描完后一定要退出，否则goroutine占用过多内存
	columns, _ := rows.Columns() // 所有查到的字段名，slice
	columnLength := len(columns)
	// 临时存储每行数据
	cache := make([]interface{}, columnLength)
	// 为每个字段初始化一个指针
	for index, _ := range cache {
		var t3interface interface{}
		cache[index] = &t3interface
	}
	// 返回的切片
	var list []map[string]*castkit.GoodleVal
	for rows.Next() {
		_ = rows.Scan(cache...)
		item := make(map[string]*castkit.GoodleVal)
		for i, data := range cache {
			val := *data.(*interface{}) // 取实际类型
			switch val.(type) {
			// 使用了指针得到的中文类型是字[]uint8(节切片)，要转string，否则json序列化后中文是字节切片
			// 这是ASCII码值与英文ASCII字符的互转,ASCII的码值本质上就是uint8类型
			case []uint8:
				val = string(val.([]uint8))
			case nil:
				val = ""
			}
			// 格式化create_time, update_time字段, 将datatime类型格式化为2020-10-10 20:20:20
			//if columns[i] == "create_time" || columns[i] == "update_time" {
			//	val = timekit.DateTimeFormat(val)
			//}
			item[columns[i]] = &castkit.GoodleVal{val}

		}
		list = append(list, item)
	}
	return list
}
