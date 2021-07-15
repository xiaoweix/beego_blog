package admin

import (
	"github.com/astaxie/beego/orm"
	"beego_blog/models"
	"strings"
	"github.com/astaxie/beego/validation"
)

type UserController struct {
	baseController
}

//用户列表
func (this *UserController) List() {
	var list []*models.User
	query := orm.NewOrm().QueryTable(new(models.User))
	count, _ := query.Count()
	if count > 0 {
		offset := (this.pager.Page - 1) * this.pager.Pagesize
		query.OrderBy("-id").Limit(this.pager.Pagesize, offset).All(&list)
	}
	this.Data["list"] = list
	this.pager.SetTotalnum(int(count))
	this.pager.SetUrlpath("/admin/user/list?page=%d")
	this.Data["pagebar"] = this.pager.ToString()
	this.display()
}

func (this *UserController) Delete() {
	id, _ := this.GetInt("id")
	if id == 7 {
		this.showmsg("不能删除超级管理员!")
	}
	user := &models.User{Id:id}
	if user.Read() == nil {
		user.Delete()
	}
	this.Redirect("/admin/user/list", 302)
}


//编辑用户
func (this *UserController) Edit() {
	//获取用户id
	id, _ := this.GetInt("id")
	user := &models.User{Id:id}
	if err := user.Read(); err != nil {
		this.showmsg("用户不存在!")
	}
	//用户存储错误提示
	errmsg := make(map[string]string)
	if this.Ctx.Input.Method() == "POST" {
		password := strings.TrimSpace(this.GetString("password"))
		password2 := strings.TrimSpace(this.GetString("password2"))
		email := strings.TrimSpace(this.GetString("email"))
		active, _ := this.GetInt("active")
		valid := validation.Validation{}

		if password != "" {
			if result := valid.Required(password2, "password2"); !result.Ok {
				errmsg["password2"] = "确认密码不能为空!"
			}else if password != password2 {
				errmsg["password2"] = "两次输入的密码不一致!"
			}else {
				user.Password = models.Md5([]byte(password))
			}
		}
		if result := valid.Required(email, "email"); !result.Ok {
			errmsg["email"] = "邮箱不能为空!"
		}else if result := valid.Email(email, "email"); !result.Ok {
			errmsg["email"] = "邮箱不合法!"
		}else {
			user.Email = email
		}

		if active > 0 {
			user.Active = 1
		}else {
			user.Active = 0
		}
		if len(errmsg) == 0 {
			user.Update()
			this.Redirect("/admin/user/list", 302)
		}
	}
	this.Data["user"] = user
	this.Data["errmsg"] = errmsg
	this.display()
}

//添加用户
func (this *UserController) Add() {
	//创建map，用于错误回显
	input := make(map[string]string)
	//创建map，用于回显错误
	errmsg := make(map[string]string)
	if this.Ctx.Request.Method == "POST" {
		//username  password   password2  email  active
		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		password2 := strings.TrimSpace(this.GetString("password2"))
		email := strings.TrimSpace(this.GetString("email"))
		active, _ := this.GetInt("active")


		input["username"] = username
		input["password"] = password
		input["password2"] = password2
		input["email"] = email

		valid := validation.Validation{}


		if result := valid.Required(username, "username"); !result.Ok {
			errmsg["username"] = "用户名不能为空!"
		}else if result := valid.MaxSize(username, 15, "username"); !result.Ok {
			errmsg["username"] = "用户名长度不能大于15个字符!"
		}

		if result := valid.Required(password, "password"); !result.Ok {
			errmsg["password"] = "密码不能为空!"
		}

		if result := valid.Required(password2, "password2"); !result.Ok {
			errmsg["password2"] = "确认密码不能为空!"
		}else if password != password2 {
			errmsg["password2"] = "两次输入的密码不一致!"
		}

		if result := valid.Required(email, "email"); !result.Ok {
			errmsg["email"] = "邮箱不能为空!"
		}else if result := valid.Email(email, "email"); !result.Ok{
			errmsg["email"] = "邮箱不合法!"
		}

		if active > 0 {
			active = 1
		}else {
			active = 0
		}

		//数据校验没有出现错误
		if len(errmsg) == 0 {
			var user = &models.User{}
			user.Username = username
			user.Password = models.Md5([]byte(password))
			user.Email = email
			user.Active = active
			if err := user.Insert(); err != nil {
				this.showmsg(err.Error())
			}
			this.Redirect("/admin/user/list", 302)
		}
	}
	this.Data["input"] = input
	this.Data["errmsg"] = errmsg
	this.display()
}






































