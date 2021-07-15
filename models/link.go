package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//友情链接
type Link struct {
	Id int
	Sitename string `orm:"size(80)"`//网站名称
	Url string `orm:"size(200)"`//网址
	Rank int //排序
}

func (link *Link) TableName() string {
	//从配置文件中获取表的前缀
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "link"
}

//插入
func (link *Link) Insert() error {
	if _, err := orm.NewOrm().Insert(link); err != nil {
		return err
	}
	return nil
}

//删除
func (link *Link) Delete() error {
	if _, err := orm.NewOrm().Delete(link); err != nil {
		return err
	}
	return nil
}

//读取
func (link *Link) Read(fields ...string) error {
	if err := orm.NewOrm().Read(link, fields...); err != nil {
		return err
	}
	return nil
}

//更新
func (link *Link) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(link, fields...); err != nil {
		return err
	}
	return nil
}