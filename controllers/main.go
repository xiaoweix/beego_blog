package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"beego_blog/models"
	"fmt"
	"strconv"
)

type MainController struct {
	beego.Controller
	Pager *models.Pager
}

func (this *MainController) Prepare() {
	var page int
	var err error

	if page, err = strconv.Atoi(this.Ctx.Input.Param(":page")); err != nil {
		page = 1
	}
	this.Pager = models.NewPager(page, 2, 0, "")
}


//首页
func (this *MainController) Index() {
	//创建文章切片，用于存储查询结果
	var list []*models.Post
	post := models.Post{}
	//获得文章表的句柄，并设置过滤条件(正常状态的文章)
	query := orm.NewOrm().QueryTable(&post).Filter("status", 0)
	//获得符合条件的记录数
	count, _ := query.Count()
	//设置总的数量
	this.Pager.SetTotalnum(int(count))
	//设置每页对应的路径
	this.Pager.SetUrlpath("/index%d.html")

	if count > 0 {
		offset := (this.Pager.Page-1)*this.Pager.Pagesize
		_, err := query.OrderBy("-istop", "-views").Limit(this.Pager.Pagesize, offset).All(&list)
		if err != nil {
			fmt.Println("err = ", err)
		}
	}
	this.Data["list"] = list
	this.Data["pagebar"] = this.Pager.ToString()

	this.setRight()
	this.setHeadMeater()
	this.display("index")
}

//http://localhost:8080/
//http://localhost:8080/index2.html
//http://localhost:8080/index3.html
//http://localhost:8080/index4.html

//设置首页右侧部分
func (this *MainController) setRight() {
	//最新文章
	this.Data["latestblog"] = models.GetLatestBlog()
	//浏览量最多的4篇文章
	this.Data["hotblog"] = models.GetHotBlog()
	//友情链接
	this.Data["links"] = models.GetLinks()
}

func (this *MainController) display(tplname string) {
	theme := "double"
	this.Layout = theme + "/layout.html"
	this.TplName = theme + "/" + tplname + ".html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["head"] = theme + "/head.html"
	this.LayoutSections["foot"] = theme + "/foot.html"

	if tplname == "index" {
		this.LayoutSections["right"] = theme + "/right.html"
		this.LayoutSections["banner"] = theme + "/banner.html"
		this.LayoutSections["middle"] = theme + "/middle.html"
	}else if tplname == "life" {
		this.LayoutSections["right"] = theme + "/right.html"
	}
}

//设置头部信息
func (this *MainController) setHeadMeater() {
	this.Data["title"] = beego.AppConfig.String("title")
	this.Data["keywords"] = beego.AppConfig.String("keywords")
	this.Data["description"] = beego.AppConfig.String("description")
}


//通过文章id查看文章详情
func (this *MainController) Show() {
	//获取文章id并转换为整数
	id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err != nil {
		this.Redirect("/404", 302)
	}
	//创建文章结构体
	post := new(models.Post)
	post.Id = id
	//查询文章
	err = post.Read()
	if err != nil {
		this.Redirect("/404", 302)
	}
	//浏览量加一
	post.Views++
	//更新浏览量
	post.Update("Views")

	this.Data["post"] = post
	//获取上一篇文章和下一篇文章
	pre, next := post.GetPreAndNext()
	this.Data["pre"] = pre
	this.Data["next"] = next
	this.Data["smalltitle"] = "文章详情"

	this.display("article")
}

//关于我
func (this *MainController) About() {
	this.setHeadMeater()
	this.display("about")
}

//成长录
func (this *MainController) BlogList() {
	query := orm.NewOrm().QueryTable(new(models.Post)).Filter("status", 0)
	count, _ := query.Count()
	this.Pager.SetTotalnum(int(count))
	var list []*models.Post
	if count > 0 {
		offset := (this.Pager.Page - 1) * this.Pager.Pagesize
		query.OrderBy("-istop", "posttime").Limit(this.Pager.Pagesize, offset).All(&list)
	}
	this.Pager.SetUrlpath("/life%d.html")//"  /life:page:int.html"
	this.Data["list"] = list
	this.Data["pagebar"] = this.Pager.ToString()
	this.setRight()
	this.setHeadMeater()
	this.display("life")

}

//碎言碎语
func (this *MainController) Mood() {
	var list []*models.Mood
	//获得tb_post表的句柄
	query := orm.NewOrm().QueryTable(new(models.Mood))
	//查询数量
	//select count(*) form tb_mood;
	count, _ := query.Count()
	if count > 0 {
		offset := (this.Pager.Page - 1) * this.Pager.Pagesize
		query.OrderBy("-posttime").Limit(this.Pager.Pagesize, offset).All(&list)
	}
	this.Data["list"] = list
	//设置总数量
	this.Pager.SetTotalnum(int(count))
	//设置urlpath
	//   /mood:page:int.html
	this.Pager.SetUrlpath("/mood%d.html")
	this.setHeadMeater()
	this.Data["pagebar"] = this.Pager.ToString()
	this.display("mood")
}








































