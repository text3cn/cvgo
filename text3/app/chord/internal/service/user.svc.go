package service

import (
    "text3/entity/mysql"
    "text3/app"
	"github.com/textthree/cvgoweb"
	"sync"
)

var userServiceInstance *UserService
var userServiceOnce sync.Once

type UserService struct {
	ctx *httpserver.Context
	uid int64
}

func UserSvc(ctx *httpserver.Context) *UserService {
	userServiceOnce.Do(func() {
		userServiceInstance = &UserService{
			ctx: ctx,
			uid: ctx.GetVal("uid").ToInt64(),
		}
	})
	return userServiceInstance
}

// ListUserinfo
func (self *UserService) ListUserinfo() (total int64, list []mysql.UserEntity) {
    page := 1
	rows := 20
	where := "1=?"
	bindValues := []any{1}
	skip := (page - 1) * rows
	err := app.Db.Model(&mysql.UserEntity{}).
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
	err = app.Db.Model(&mysql.UserEntity{}).Where(where, bindValues...).Count(&total).Error
	if err != nil {
		app.Log.Error(err)
		return
	}
	return
}
