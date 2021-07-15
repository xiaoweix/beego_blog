package admin

import (
	"os"
	"runtime"
	"github.com/astaxie/beego/orm"
	"beego_blog/models"
)

type IndexController struct {
	baseController

}


func (this *IndexController) Index() {
	//hostname  gover  os  arch  cpunum  postnum  tagnum  usernum
	this.Data["hostname"], _ = os.Hostname()
	this.Data["gover"] = runtime.Version()
	this.Data["os"] = runtime.GOOS
	this.Data["arch"] = runtime.GOARCH
	this.Data["cpunum"] = runtime.NumCPU()
	this.Data["postnum"], _ = orm.NewOrm().QueryTable(new(models.Post)).Count()
	this.Data["tagnum"], _ = orm.NewOrm().QueryTable(new(models.Tag)).Count()
	this.Data["usernum"], _ = orm.NewOrm().QueryTable(new(models.User)).Count()
	this.display()
}