package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"strings"
	"fmt"
	"strconv"
)

//文章
type Post struct {
	Id int
	//用户id
	Userid int
	//作者
	Author string `orm:"size(15)"`
	//标题
	Title string `orm:"size(100)"`
	//标题颜色
	Color string `orm:"size(7)"`
	//文章内容
	Content string `orm:"type(text)"`
	//标签名称
	Tags string `orm:"size(100)"`
	//浏览量
	Views int
	//状态
	Status int
	//发表时间
	Posttime time.Time `orm:"type(datetime)"`
	//更新时间
	Updated time.Time `orm:"type(datetime)"`
	//是否置顶
	Istop int
	//封面
	Cover string `orm:"size(70)"`
}


func (post *Post) TableName() string {
	//从配置文件中获取表的前缀
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "post"
}

//插入
func (post *Post) Insert() error {
	if _, err := orm.NewOrm().Insert(post); err != nil {
		return err
	}
	return nil
}


//读取
func (post *Post) Read(fields ...string) error {
	if err := orm.NewOrm().Read(post, fields...); err != nil {
		return err
	}
	return nil
}

//更新
func (post *Post) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(post, fields...); err != nil {
		return err
	}
	return nil
}


func (post *Post) TagsLink() string {
	if post.Tags == "" {
		return ""
	}
	tagslink := strings.Trim(post.Tags, ",")
	return tagslink
}

func (post *Post) Link() string {
	// /index:page:int.html
	//   /article/2
	return "/article/" + strconv.Itoa(post.Id)
}

//返回带颜色的标题
func (post *Post) ColorTitle() string {
	if post.Color != "" {
		return fmt.Sprintf("<span style='color:%s'>%s</span>", post.Color, post.Title)
	}
	return post.Title
}

func (post *Post) Excerpt() string {
	return post.Content
}


//根据当前文章获取上一篇文章和下一篇文章
func (this *Post) GetPreAndNext() (pre, next *Post) {
	pre = &Post{}
	//上一篇文章
	err := orm.NewOrm().QueryTable(new(Post)).OrderBy("-id").Filter("id__lt", this.Id).Filter("status", 0).Limit(1).One(pre)
	if err != nil {
		pre = nil
	}

	next = &Post{}
	//下一篇文章
	err = orm.NewOrm().QueryTable(new(Post)).OrderBy("id").Filter("id__gt", this.Id).Filter("status", 0).Limit(1).One(next)
	if err != nil {
		pre = nil
	}
	return
}

//删除文章
func (post *Post) Delete() error {
	//判断文章的标签是否为空
	if post.Tags != "" {
		//获取标签文章表的句柄并通过文章id进行过滤
		query := orm.NewOrm().QueryTable(&TagPost{}).Filter("postid", post.Id)
		//创建切片，用于存储查询结果
		var tagpost []*TagPost
		if n, err := query.All(&tagpost); n > 0 && err == nil {
			//遍历查询结果
			for i := 0; i < len(tagpost); i++ {
				//创建标签对象并初始化标签id
				var tag = &Tag{Id:tagpost[i].Tagid}
				//根据标签id查询标签表，当查询没有出现错误并且文章数量大于0的时候才跟新count字段
				if err = tag.Read(); err == nil && tag.Count > 0 {
					tag.Count--
					//更新count字段
					tag.Update("count")
				}
			}
		}
		//删除标签文章表中的相关记录
		query.Delete()
	}
	//删除文章
	if _, err := orm.NewOrm().Delete(post); err != nil {
		return err
	}
	return nil
}