package main

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/textthree/cvgokit/strkit"
	"github.com/textthree/provider/clog"
	"os"
	"reflect"
	"text3/scripts/cvgo/codegen/entityregistry"
	"time"
)

// go build -o $GOPATH/bin/cvg
// go run scripts/cvgo/codegen/curdl.go c
func main() {
	curdType := os.Args[1]
	tableName := os.Args[2]
	funcName := os.Args[3]
	fileNamePascalCase := os.Args[4]
	cursorPaging := cast.ToBool(os.Args[5])
	switch curdType {
	case "c":
		fmt.Println(CurdCreate(tableName, funcName, fileNamePascalCase))
	case "u":
		fmt.Printf(CurdUpdate(tableName, funcName, fileNamePascalCase))
	case "r":
		fmt.Printf(CurdGet(tableName, funcName, fileNamePascalCase))
	case "d":
		fmt.Printf(CurdDelete(tableName, funcName, fileNamePascalCase))
	case "l":
		fmt.Printf(CurdList(tableName, funcName, fileNamePascalCase, cursorPaging))
	}
}

// Create 基础代码
func CurdCreate(tableName, funcName, fileNamePascalCase string) string {
	entityType, found := entityregistry.EntityRegistry[tableName]
	if !found {
		clog.RedPrintln("类型", strkit.SnakeToPascalCase(tableName)+"Entity", "不存在")
		return ""
	}
	// 创建实体实例
	entityValue := reflect.New(entityType).Elem()
	packageName := "mysql"
	structName := strkit.SnakeToPascalCase(tableName) + "Entity"
	content := `
// ` + funcName + `
func (self *` + fileNamePascalCase + `Service) ` + funcName + `() {
`
	content += fmt.Sprintf("    result := %s.Db.Create(&%s.%s{\n", "app", packageName, structName)
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		value := entityValue.Field(i)

		// 过滤嵌套结构体
		// 判断字段类型是否为结构体，但不是时间类型（time.Time也是一个结构体，需要排除）
		if field.Type.Kind() == reflect.Struct && field.Type.Name() != "Time" {
			continue
		}

		// 获取字段默认值的字符串表示形式
		var fieldValue string
		switch value.Kind() {
		case reflect.String:
			fieldValue = `""`
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldValue = "0"
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldValue = "0"
		case reflect.Float32, reflect.Float64:
			fieldValue = "0.0"
		case reflect.Bool:
			fieldValue = "false"
		case reflect.Struct:
			// 检查是否为 time.Time 类型
			if field.Type == reflect.TypeOf(time.Time{}) {
				fieldValue = "time.Time{}"
			}
		default:
			fieldValue = "nil"
		}

		content += fmt.Sprintf("\t\t%s: %s,\n", field.Name, fieldValue)
	}

	content += "    })"
	content += `
	if result.Error != nil {
		app.Log.Error(result.Error.Error())
	}
}`

	return content
}

// 根据主键修改记录
func CurdUpdate(tableName, funcName, fileNamePascalCase string) string {
	structName := strkit.SnakeToPascalCase(tableName) + "Entity"
	content := `
// ` + funcName + `
func (self *` + fileNamePascalCase + `Service) ` + funcName + `() {
`
	content += `    result := app.Db.Model(&mysql.` + structName + `{}).
		Where("id = ?", 0).
		Updates(mysql.` + structName + `{})
	if result.Error != nil {
		app.Log.Error(result.Error.Error())
	}
}`
	return content
}

// 根据主键删除记录
func CurdDelete(tableName, funcName, fileNamePascalCase string) string {
	structName := strkit.SnakeToPascalCase(tableName) + "Entity"
	content := `
// ` + funcName + `
func (self *` + fileNamePascalCase + `Service) ` + funcName + `() {
`
	content += `    result := app.Db.Where("id = ?", 0).Delete(&mysql.` + structName + `{}, 0)
	if result.Error != nil {
		app.Log.Error(result.Error.Error())
	}
}`
	return content
}

// 根据主键获取记录
func CurdGet(tableName, funcName, fileNamePascalCase string) string {
	structName := strkit.SnakeToPascalCase(tableName) + "Entity"
	content := `
// ` + funcName + `
func (self *` + fileNamePascalCase + `Service) ` + funcName + `() (rows mysql.` + structName + `) {
`
	content += `    result := app.Db.Where("id = ?", 0).Take(&rows)
	if result.Error != nil {
		app.Log.Error(result.Error.Error())
	}
	return
}`
	return content
}

// 获取列表
func CurdList(tableName string, funcName, fileNamePascalCase string, cursorPaging bool) (content string) {
	entityStruct := strkit.SnakeToPascalCase(tableName) + "Entity"
	if cursorPaging {
		// 游标分页
		content = `
// ` + funcName + `
func (self *` + fileNamePascalCase + `Service) ` + funcName + `(cursor int64, rows ...int) (list []mysql.` + entityStruct + `, total int64) {
`
		content += `    limitRows := 20
	if len(rows) > 0 && rows[0] < 50 {
		limitRows = rows[0]
	}
	where := "1=?"
	bindValues := []any{1}
	tx := app.Db.Model(&mysql.` + entityStruct + `{}).Where(where, bindValues...)
	if cursor > 0 {
		// 适用于按 id 降序列表，升序列表请改为 id > ?
		tx.Where("id < ?", cursor)
	}
	err := tx.Limit(limitRows + 1).Find(&list).Error
	if err != nil {
		app.Log.Error(err)
		return
	}
	// 分页处理
	if len(list) == 0 {
		return
	}
	rowsNum := len(list)
	lastIndex := rowsNum - 1
	if rowsNum > limitRows {
		lastIndex--
	}
	return
}
`
	} else {
		// 传统分页
		content = `
// ` + funcName + `
func (self *` + fileNamePascalCase + `Service) ` + funcName + `() (total int64, list []mysql.` + entityStruct + `) {
`
		content += `    page := 1
	rows := 20
	where := "1=?"
	bindValues := []any{1}
	skip := (page - 1) * rows
	err := app.Db.Model(&mysql.` + entityStruct + `{}).
	Where(where, bindValues...).
	Offset(skip).
	Order("id DESC").
	Limit(rows).
	Find(&list).Error
	if err != nil {
		app.Log.Error(err)
		return
	}
	// 统计总数
	err = app.Db.Model(&mysql.` + entityStruct + `{}).Where(where, bindValues...).Count(&total).Error
	if err != nil {
		app.Log.Error(err)
		return
	}
	return
}
`
	}
	return content
}
