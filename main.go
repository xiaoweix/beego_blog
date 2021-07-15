package main

import (
	_ "beego_blog/routers"
	"github.com/astaxie/beego"
	_ "beego_blog/models"
)

func main() {
	beego.Run()
}

