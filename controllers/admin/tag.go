package admin

import (
	"github.com/astaxie/beego/orm"
	"beego_blog/models"
	"strconv"
	"strings"
)

type TagController struct {
	baseController
}

//标签列表
func (this *TagController) List() {
	//接收对应的操作
	act := this.GetString("act")
	switch act {
	case "batch":
		this.batch()
	default:
		this.tagList()
	}
}

func (this *TagController) batch() {
	//获取用户所选择的id
	ids := this.GetStrings("ids[]")
	//获取用户所选择的操作
	op := this.GetString("op")
	idarr := make([]int, 0)
	for _, v := range ids {
		if id, _ := strconv.Atoi(v); id > 0 {
			idarr = append(idarr, id)
		}
	}
	switch op {
	//合并
	/*
	1.获取用户输入的新的标签名称，并去除两边的空格，根据该名称去标签表中查找记录，
	如果没有查找到则需要插入新的记录。
	2.根据需要被合并的标签id在标签文章表中查找对应的记录，将对应记录的标签id替换目标标签的id
	3.在文章表中将原标签的名称替换为目标标签的名称
	4.删除原始标签
	5.更新目标标签的count字段

	*/
	case "merge":
		//获取目标标签的名称
		toname := strings.TrimSpace(this.GetString("toname"))
		//目标标签的名称不为空并且idarr也不为空
		if toname != "" && len(idarr) > 0 {
			//创建标签对象
			tag := new(models.Tag)
			//赋值
			tag.Name = toname
			//根据标签名称查询
			if tag.Read("name") != nil {
				tag.Count = 0
				tag.Insert()
			}
			for _, id := range idarr {
				//创建标签结构体，并初始化id
				obj := models.Tag{Id:id}
				if obj.Read() == nil {
					obj.MergeTo(tag)
					obj.Delete()
				}
			}
			tag.UpCount()
		}

	//删除
	case "delete":
		//遍历id切片
		for _, id := range idarr {
			obj := models.Tag{Id:id}
			if obj.Read() == nil {
				obj.Delete()
			}
		}
	}
	this.Redirect("/admin/tag", 302)
}

//标签列表
func (this *TagController) tagList() {
	var list []*models.Tag
	//获得标签表的句柄
	query := orm.NewOrm().QueryTable(new(models.Tag))
	//获得标签的数量
	count, _ := query.Count()
	if count > 0 {
		//计算偏移量
		offset := (this.pager.Page - 1) * this.pager.Pagesize
		//分页查询
		query.Limit(this.pager.Pagesize, offset).All(&list)
	}
	this.Data["list"] = list
	this.pager.SetTotalnum(int(count))
	this.pager.SetUrlpath("/admin/tag?page=%d")
	this.Data["pagebar"] = this.pager.ToString()
	this.display("tag_list")
}