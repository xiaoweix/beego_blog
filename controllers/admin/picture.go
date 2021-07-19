package admin

import "log"

type PictureController struct {
	baseController
}

func (this *PictureController) List() {
	if this.Ctx.Request.Method == "POST" {
		f, h, err := this.GetFile("picture")
		if err != nil {
			log.Fatal("getfile err ", err)
			this.Abort("403")
		}
		defer f.Close()
		err = this.SaveToFile("picture", "static/upload/"+h.Filename) // 保存位置在 static/upload, 没有文件夹要先创建
		if err != nil {
			this.Abort("403")
		}
	}
	this.display()
}
