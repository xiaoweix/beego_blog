package admin

import (
	"strings"
	"beego_blog/models"
	"strconv"
)

type AccountController struct {
	baseController
}


//登录
func (this *AccountController) Login() {
	//判断请求方式是否是post请求
	if this.GetString("dosubmit") == "yes" {
		//account   password   remember
		//获取账号并去除两边的空格
		account := strings.TrimSpace(this.GetString("account"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := strings.TrimSpace(this.GetString("remember"))
		//判断账号密码是否都不为空
		if account != "" && password != "" {
			var user = &models.User{}
			user.Username = account
			//根据账号查询用户，并且判断查询到的密码和用户输入的密码通过md5哈希之后的结果是否一致
			if user.Read("username") != nil || user.Password != models.Md5([]byte(password)){
				this.Data["errmsg"] = "账号或密码错误!"
			}else if user.Active == 0 {//判断该账号是否激活
				this.Data["errmsg"] = "账号尚未激活"
			}else {
				//登录次数加一
				user.Logincount += 1
				//更新Logincount
				user.Update("logincount")
				authkey := models.Md5([]byte(password))
				if remember == "yes" {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id) + "|" + authkey, 60*60*24*7)
				}else {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id) + "|" + authkey)
				}
				//重定向到后台首页
				this.Redirect("/admin", 302)
			}
		}
	}
	this.TplName = "admin/" + this.controllerName + "_" + this.actionName + ".html"
}

//退出登录
func (this *AccountController) Logout() {
	//清空cookie
	this.Ctx.SetCookie("auth", "")
	this.Redirect("/admin/login", 302)
}

func (this *AccountController) Profile() {
	user := &models.User{Id:this.userid}
	if err := user.Read(); err != nil {
		this.showmsg(err.Error())
	}
	updated := false
	errmsg := make(map[string]string)

	if this.Ctx.Request.Method == "POST" {
		//password   newpassword   newpassword2
		password := strings.TrimSpace(this.GetString("password"))
		newpassword := strings.TrimSpace(this.GetString("newpassword"))
		newpassword2 := strings.TrimSpace(this.GetString("newpassword2"))

		if newpassword != "" {
			if password == "" || models.Md5([]byte(password)) != user.Password {
				errmsg["password"] = "当前密码错误!"
			}else if len(newpassword) < 6 {
				errmsg["newpassword"] = "新密码不能少于6个字符"
			}else if newpassword != newpassword2 {
				errmsg["newpassword2"] = "两次输入的密码不一致!"
			}
		}

		if len(errmsg) == 0{
			user.Password = models.Md5([]byte(newpassword))
			user.Update("password")
			updated = true
		}
	}
	this.Data["updated"] = updated
	this.Data["errmsg"] = errmsg

	this.Data["user"] = user
	this.display()
}


/*
crontab -e创建定时任务
crontab -l查看定时任务
crontab -r删除定时任务

* * * * * echo "hello world" >> /tmp/test.log
(1).表示分钟,可以取0-59
(2).表示小时，可以取0-23
(3).表示日期，可以取1-31
(4).表示月份，可以取1-12
(5).表示星期几，可以取0-7,0和7表示的都是星期日
可以简记为：分时日月周

(1).星号(*):代表所有可能的值
(2).逗号(,):指定一个列表范围，例如:"1, 2, 5, 9"
(3).中杆(-):表示整数范围，例如"2-6",("2, 3, 4, 5, 6")
(4).正斜杆(/):指的是时间的间隔频率，"0-23/2"

每小时的第2和第20分钟执行
2,20 * * * * command

在上午7点到10点的第2和第20分钟执行
2,20 7-10 * * * command

每个星期一的上午7点到10点的第2和第20分钟执行
2,20 7-10 * * 1 command

每隔两天的上午7到10点的第2和第20分钟执行


*/

















































