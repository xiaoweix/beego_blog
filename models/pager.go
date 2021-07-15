package models

import (
	"bytes"
	"fmt"
)

type Pager struct {
	Page int //当前页码
	Pagesize int //每页大小
	Totalnum int //总数量
	urlpath string//每页对应的url
}

//创建pager对象
func NewPager(page, pagesize, totaonum int, urlpath string) *Pager {
	pager := new(Pager)
	pager.Page = page
	pager.Pagesize = pagesize
	pager.Totalnum = totaonum
	pager.urlpath = urlpath
	return pager
}

func (this *Pager) SePage(page int) {
	this.Page = page
}

func (this *Pager) SetPagesize(page int) {
	this.Page = page
}

func (this *Pager) SetTotalnum(totalnum int) {
	this.Totalnum = totalnum
}

func (this *Pager) SetUrlpath(urlpath string) {
	this.urlpath = urlpath
}

//this.Pager.SetUrlpath("/index%d.html")
func (this *Pager) url(page int) string {
	return fmt.Sprintf(this.urlpath, page)
}

//根据当前页码计算现在需要显示的10页，将者10页的页码放在a标签中一字符创的形式返回
func (this *Pager) ToString() string {
	//文章的总数量小于等于每页显示的文章的数量
	if this.Totalnum <= this.Pagesize {
		return ""
	}
	var totalpage int//总页码
	linknum := 10//需要显示的页码
	var from int //从那一页开始显示
	var to int //显示到那一页
	offset := 5//偏移量
	if this.Totalnum % this.Pagesize != 0 {
		totalpage = this.Totalnum/this.Pagesize + 1
	}else {
		totalpage = this.Totalnum/this.Pagesize
	}

	//总的页码小于10.直接从第一页显示到最后一页
	if totalpage < linknum {
		from = 1
		to = totalpage
	}else {
		//开始页码
		from = this.Page - offset
		//最后一页
		to = from + linknum
		//判断起始页是否小于1
		if from < 1 {
			from = 1
			to = from + linknum - 1
		}else if to > totalpage {
			to = totalpage
			from = to - linknum + 1
		}
	}
   	var buf bytes.Buffer
   	buf.WriteString("<div class='page'>")
	//上一页
	if this.Page > 1 {
		//this.Pager.SetUrlpath("/index%d.html")
		buf.WriteString(fmt.Sprintf("<a href='%s'>&laquo;</a>", this.url(this.Page-1)))
	}

	for i := from; i <= to; i++ {
		if i == this.Page {
			buf.WriteString(fmt.Sprintf("<b>%d</b>", i))
		}else {
			buf.WriteString(fmt.Sprintf("<a href='%s'>%d</a>", this.url(i), i))
		}
	}

	//设置下一页
	if this.Page < totalpage {
		buf.WriteString(fmt.Sprintf("<a href='%s'>&raquo;</a>", this.url(this.Page+1)))
	}

	buf.WriteString("</div>")

	str := buf.String()
	return str
}


