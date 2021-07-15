package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type User struct {
	Id int
	Username string `orm:"size(15)"`
	Password string `orm:"size(32)"`
	Email string `orm:"size(50)"`
	Logincount int//登录次数
	Authkey string `orm:"size(10)"`
	Active int//是否激活
}


func (user *User) TableName() string {
	//从配置文件中获取表的前缀
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "user"
}

//插入
func (user *User) Insert() error {
	if _, err := orm.NewOrm().Insert(user); err != nil {
		return err
	}
	return nil
}

//删除
func (user *User) Delete() error {
	if _, err := orm.NewOrm().Delete(user); err != nil {
		return err
	}
	return nil
}

//读取
func (user *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(user, fields...); err != nil {
		return err
	}
	return nil
}

//更新
func (user *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(user, fields...); err != nil {
		return err
	}
	return nil
}
