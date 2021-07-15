package admin

import (
	"github.com/astaxie/beego"
	"strings"
	"strconv"
	"beego_blog/models"
)

type baseController struct {
	beego.Controller
	userid int
	username string
	controllerName string
	actionName string
	pager *models.Pager
}
func (this *baseController) Prepare() {
	//获取控制器名称和方法名称  IndexController_Index.html
	//IndexController  LinkController TagController
	controllerName, actionName := this.GetControllerAndAction()
	//去除控制器名称尾部的Controller并将结果转换为小写
	this.controllerName = strings.ToLower(controllerName[:len(controllerName)-10])
	//将方法名称转换为小写
	this.actionName = strings.ToLower(actionName)

	this.auth()
	page, err := this.GetInt("page")
	if err != nil {
		page = 1
	}
	pagesize := 2
	this.pager = models.NewPager(page, pagesize, 0, "")
}

//身份验证
func (this *baseController) auth() {
	//用户请求的是前台页面或者是登录页面不需要身份验证
	if this.controllerName == "main" || (this.controllerName == "account" && this.actionName == "login") {
		return
	}
	//获取cookie并通过|进行切割
	arr := strings.Split(this.Ctx.GetCookie("auth"), "|")
	if len(arr) == 2 {
		idstr, password := arr[0], arr[1]
		//将id转换为整数
		id, _ := strconv.Atoi(idstr)
		if id > 0 {
			user := new(models.User)
			user.Id = id
			if user.Read() == nil && user.Password == password {
				this.userid = user.Id
				this.username = user.Username
			}
		}
	}
	//验证失败
	if this.userid == 0 {
		this.Redirect("/admin/login", 302)
	}
}

//  views/admin/tag_list.html
func (this *baseController) display(tplname ...string) {
	modileName := "admin/"
	this.Layout = modileName + "layout.html"

	this.Data["version"] = beego.AppConfig.String("version")
	this.Data["adminname"] = this.username

	if len(tplname) == 1 {
		this.TplName = modileName + tplname[0] + ".html"
	}else {
		this.TplName = modileName + this.controllerName + "_" + this.actionName + ".html"
	}
}

func (this *baseController) showmsg(msg ...string) {
	if len(msg) == 0 {
		msg = append(msg, "出错了!")
	}
	//拼接上一个页面的地址
	msg = append(msg, this.Ctx.Request.Referer())
	this.Data["msg"] = msg[0]
	this.Data["redirect"] = msg[1]
	this.display("showmsg")
	this.Render()
	this.StopRun()
}