package admin

import (
	"github.com/astaxie/beego/orm"
	"beego_blog/models"
	"strings"
)

type LinkController struct {
	baseController
}

//友链列表
func (this *LinkController) List() {
	var list []*models.Link
	orm.NewOrm().QueryTable(new(models.Link)).OrderBy("-rank").All(&list)
	this.Data["list"] = list
	this.display()
}

//编辑友链
func (this *LinkController) Edit() {
	//获取id
	id, _ := this.GetInt("id")
	link := &models.Link{Id:id}
	if err := link.Read(); err != nil {
		this.showmsg("友链不存在!")
	}
	if this.Ctx.Request.Method == "POST" {
		sitename := strings.TrimSpace(this.GetString("sitename"))
		url := strings.TrimSpace(this.GetString("url"))
		/*rank, err := this.GetInt("rank")
		if err != nil {
			rank = 0
		}*/
		rank, _ := this.GetInt("rank", 0)
		link.Sitename = sitename
		link.Url = url
		link.Rank = rank
		link.Update()
		this.Redirect("/admin/link/list", 302)
	}
	this.Data["link"] = link
	this.display()
}

//删除友情链接
func (this *LinkController) Delete() {
	id, err := this.GetInt("id")
	if err != nil {
		this.showmsg("删除失败!")
	}
	link := &models.Link{Id:id}
	if err = link.Read(); err == nil {
		link.Delete()
	}
	this.Redirect("/admin/link/list", 302)
}

//添加友链
func (this *LinkController) Add() {
	if this.Ctx.Request.Method == "POST" {
		//sitename  url  rank
		sitename := this.GetString("sitename")
		url := this.GetString("url")
		rank, err := this.GetInt("rank")
		if err != nil {
			rank = 0
		}
		var link = &models.Link{Sitename:sitename, Url:url, Rank:rank}
		if err = link.Insert(); err != nil {
			this.showmsg("添加失败!")
		}
		this.Redirect("/admin/link/list", 302)
	}
	this.display()
}




