package admin

import (
	"github.com/astaxie/beego/orm"
	"beego_blog/models"
	"fmt"
	"strings"
	"time"
	"math/rand"
	"strconv"
)

type ArticleController struct {
	baseController
}

//文章列表
func (this *ArticleController) List() {
	//获取搜索状态
	status, _ := this.GetInt("status")
	//获取搜索类型
	searchtype := this.GetString("searchtype")
	//获取搜索关键字
	keyword := this.GetString("keyword")
	var list []*models.Post

	//得到文章表的句柄并设置过滤条件
	query := orm.NewOrm().QueryTable(new(models.Post)).Filter("status", status)
	//搜索关键字不为空
	if keyword != "" {
		switch searchtype {
		//搜索类型为标题
		case "title":
			//select * from tb_post where title like '%keyword%'
			query = query.Filter("title__icontains", keyword)
		//搜索类型为作者
		case "author":
			query = query.Filter("author__icontains", keyword)
		//搜索类型为标签
		case "tag":
			query = query.Filter("tag__icontains", keyword)
		}
	}
	count, _ := query.Count()
	//判断count是否大于0
	if count > 0 {
		//设置偏移量
		offset := (this.pager.Page - 1) * this.pager.Pagesize
		//分页查询
		query.Limit(this.pager.Pagesize, offset).All(&list)
	}
	this.Data["list"] = list
	this.Data["status"] = status
	//获取草稿箱中文章的数量
	//0:已发布  1：草稿箱  2：回收站
	this.Data["count_1"], _ = orm.NewOrm().QueryTable(&models.Post{}).Filter("status", 1).Count()
	//获取回收站中文章的数量
	this.Data["count_2"], _ = orm.NewOrm().QueryTable(&models.Post{}).Filter("status", 2).Count()
	//搜索类型
	this.Data["searchtype"] = searchtype
	//搜索的关键字
	this.Data["keyword"] = keyword

	this.pager.SetTotalnum(int(count))
	this.pager.SetUrlpath(fmt.Sprintf("/admin/article/list?searchtype=%s&keyword=%s&status=%d&page=%s", searchtype, keyword, status, "%d"))
	this.Data["pagebar"] = this.pager.ToString()

	this.display()
}

//跳转到文章添加页面
func (this *ArticleController) Add() {
	this.display()
}


//添加文章
/*
思路：
1.获取用户输入的文章信息，插入数据库
2.用户有可能输入多个标签，需要检测每个标签的合法性，例如标签与标签之间的空格需要去除，并判断是否有重复标签
3.第二步处理完成之后，判断标签表中是否存在这些标签，如果存在需要更新count，如果不存在则创建
4.在标签文章表中插入相应的字段

*/
func (this *ArticleController) Save() {
	//title  color   istop   tags   posttime  status  content
	var post models.Post
	//一.获取用户输入
	//1.获取前台传递过来的数据
	post.Title = strings.TrimSpace(this.GetString("title"))
	if post.Title == "" {
		this.showmsg("标题不能为空!")
	}
	post.Color = strings.TrimSpace(this.GetString("color"))
	post.Istop, _ = this.GetInt("istop", 0)
	tags :=  strings.TrimSpace(this.GetString("tags"))
	timestr := strings.TrimSpace(this.GetString("posttime"))
	post.Status, _ = this.GetInt("status", 0)
	post.Content = this.GetString("content")
	//2.补全其他字段的信息
	post.Userid = this.userid
	post.Author = this.username
	post.Updated = time.Now()
	//3.设置随机数种子
	rand.Seed(time.Now().Unix())
	//生成[0,10)之间的随机数
	var r = rand.Intn(10)
	//拼接图片的路径
	post.Cover = "/static/upload/blog" + fmt.Sprintf("%d", r) + ".jpg"
	//Mon Jan 2 15:04:05 -0700 MST 2006
	//将时间转换为time类型
	posttime, err := time.Parse("2006-01-02 15:04:05", timestr)
	//判断转换是否成功，如果不成功，取当前时间为文章发布时间
	if err == nil {
		post.Posttime = posttime
	}else {
		post.Posttime = time.Now()
	}
	//插入文章
	if err = post.Insert(); err != nil {
		this.showmsg("文章添加失败!")
	}
	//  0   1     2    3  4
	//  C  C++  JAVA  GO  C
	//  C  C++  JAVA  GO
	//去重之后的结果切片
	addtags := make([]string, 0)
	//二.处理标签
	if tags != "" {
		//将中文的逗号全部替换为英文的逗号
		tags = strings.Replace(tags, "，", ",", -1)
		//通过英文逗号切割标签
		tagarr := strings.Split(tags, ",")
		//遍历存储标签的切片
		for _, v := range tagarr {
			if tag := strings.TrimSpace(v); tag != "" {
				//定义标志，默认没有重复标签
				exists := false
				for _, vv := range addtags {
					//有重复标签
					if vv == tag {
						exists = true
						//退出循环
						break
					}
				}
				//没有重复标签，则将tag追加到结果切片中
				if !exists {
					addtags = append(addtags, tag)
				}
			}
		}
	}

	//三.将结果切片中的标签插入标签表中
	if len(addtags) > 0 {
		//遍历结果标签
		for _, v := range addtags {
			//创建标签对象并初始化变迁名称
			tag := &models.Tag{Name: v}
			//根据名称查询标签
			if err := tag.Read("Name"); err == orm.ErrNoRows {
				tag.Count = 1
				tag.Insert()
			}else {
				//该标签下的文章数量加一
				tag.Count += 1
				//更新count字段
				tag.Update("Count")
			}
			//创建标签文章对象，并初始化各个字段
			tp := &models.TagPost{Tagid:tag.Id, Postid:post.Id, Poststatus:post.Status, Posttime:post.Posttime}
			//插入标签文章对象
			tp.Insert()
		}
		//拼接标签名称
		post.Tags = "," + strings.Join(addtags, ",") + ","
	}
	post.Updated = time.Now()
	//跟新
	post.Update("tags", "updated")
	this.Redirect("/admin/article/list", 302)
}


//删除文章
func (this *ArticleController) Delete() {
	//获取文章id
	id, _ := this.GetInt("id")
	//创建文章结构体并初始化id
	post := &models.Post{Id:id}
	if post.Read() == nil {
		post.Delete()
	}
	this.Redirect("/admin/article/list", 302)
}

//批量操作
func (this *ArticleController) Batch() {
	//获取用户所选择的文章id
	ids := this.GetStrings("ids[]")
	//创建切片，用于存储转换之后的结果
	idarr := make([]int, 0)
	//遍历获取到的文章id
	for _, v := range ids {
		//将字符串id转换为整形
		if id, _ := strconv.Atoi(v); id > 0 {
			//将转换之后的结果追加到结果切片中
			idarr =append(idarr, id)
		}
	}

	//获取用户所选择的操作
	op := this.GetString("op")
	query := orm.NewOrm().QueryTable(new(models.Post))
	switch op {
	//移至已发布
	case "topub":
		query.Filter("id__in", ids).Update(orm.Params{"status": 0})
	//移至草稿箱
	case "todrafts":
		query.Filter("id__in", ids).Update(orm.Params{"status": 1})
	//移至回收站
	case "totrash":
		query.Filter("id__in", ids).Update(orm.Params{"status": 2})
	//删除
	case "delete":
		for _, id := range idarr {
			//创建文章结构体，并初始化id
			obj := models.Post{Id:id}
			//查询
			if obj.Read() == nil {
				//删除
				obj.Delete()
			}
		}
	}
	//重定向到上一个页面
	this.Redirect(this.Ctx.Request.Referer(), 302)
}

//编辑文章(跳转到文章编辑页面)
func (this *ArticleController) Edit() {
	//获取需要被编辑的文章的id
	id, _ := this.GetInt("id")
	post := &models.Post{Id:id}
	//查询出现错误
	if post.Read() != nil {
		this.showmsg("未找到该篇文章")
	}
	//去除标签前后的英文逗号
	post.Tags = strings.Trim(post.Tags, ",")
	this.Data["post"] = post
	//将文章发布时间转换为字符串
	this.Data["posttime"] = post.Posttime.Format("2006-01-02 15:04:05")
	this.display()
}

//更新文章
/*
思路：
第一种情况：需要判断用户是否修改了标签，如果没有修改直接更新文章表
第二种情况：如果修改了标签，则需要更新标签文章表和标签表，首先我们
需要判断修改之前该文章的标签是否为空，如果不为空，则将标签文章中的相关记录删除，
更新标签表中对应的count字段，然后需要处理用户输入的新标签，也就是去除新标签两边的空格，
去除重复标签等等，在标签中如果没有对应的标签则创建该标签，如果已经存在该标签则更新count字段，
最后在标签文章中插入对应的记录
*/
func (this *ArticleController) Update() {
	//创建文章结构体
	var post models.Post
	//获取文章id
	id, err := this.GetInt("id")
	//处理错误
	if err != nil {
		this.showmsg("文章不存在!")
	}
	post.Id = id
	//通过id查询文章
	if post.Read() != nil {
		this.Redirect("/admin/article/list", 302)
	}



	//获取文章标题
	post.Title = strings.TrimSpace(this.GetString("title"))
	if post.Title == "" {
		this.showmsg("标题不能为空!")
	}
	post.Color = strings.TrimSpace(this.GetString("color"))
	post.Istop, _ = this.GetInt("istop")
	tags := strings.TrimSpace(this.GetString("tags"))
	timestr := strings.TrimSpace(this.GetString("posttime"))
	//将字符串时间转换为time类型
	if posttime, err := time.Parse("2006-01-02 15:04:05", timestr); err == nil {
		post.Posttime = posttime
	}
	post.Status, _ = this.GetInt("status")
	post.Content = this.GetString("content")
	//修改文章的修改时间
	post.Updated = time.Now()

	//------第一种情况-------------
	//用户没有修改文章所属的标签
	if strings.Trim(post.Tags, ",") == tags {
		post.Update("title", "color", "istop", "posttime", "status", "content", "updated")
		this.Redirect("/admin/article/list", 302)
	}


	//------第二种情况-------------
	if post.Tags != "" {
		var tagpost models.TagPost
		//获得标签文章表的句柄并通过文章id进行过滤
		query := orm.NewOrm().QueryTable(&tagpost).Filter("postid", post.Id)
		//用于存储查询结果
		var tagpostarr []*models.TagPost
		if n, err := query.All(&tagpostarr); n > 0 && err == nil {
			//遍历tagpostarr
			for i := 0; i < len(tagpostarr); i++ {
				//去除tagid对新创建的标签对象赋值
				var tag = &models.Tag{Id:tagpostarr[i].Tagid}
				//通过id查询标签，如果没有出现错误且count字段大于0
				if err =tag.Read(); err == nil && tag.Count > 0 {
					tag.Count--
					tag.Update("count")
				}
			}
		}
		//在标签文章表中删除对应的记录
		query.Delete()
	}




	//  0   1     2    3  4
	//  C  C++  JAVA  GO  C
	//  C  C++  JAVA  GO
	//去重之后的结果切片
	addtags := make([]string, 0)
	//二.处理标签
	if tags != "" {
		//将中文的逗号全部替换为英文的逗号
		tags = strings.Replace(tags, "，", ",", -1)
		//通过英文逗号切割标签
		tagarr := strings.Split(tags, ",")
		//遍历存储标签的切片
		for _, v := range tagarr {
			if tag := strings.TrimSpace(v); tag != "" {
				//定义标志，默认没有重复标签
				exists := false
				for _, vv := range addtags {
					//有重复标签
					if vv == tag {
						exists = true
						//退出循环
						break
					}
				}
				//没有重复标签，则将tag追加到结果切片中
				if !exists {
					addtags = append(addtags, tag)
				}
			}
		}
	}

	//三.将结果切片中的标签插入标签表中
	if len(addtags) > 0 {
		//遍历结果标签
		for _, v := range addtags {
			//创建标签对象并初始化变迁名称
			tag := &models.Tag{Name: v}
			//根据名称查询标签
			if err := tag.Read("Name"); err == orm.ErrNoRows {
				tag.Count = 1
				tag.Insert()
			}else {
				//该标签下的文章数量加一
				tag.Count += 1
				//更新count字段
				tag.Update("Count")
			}
			//创建标签文章对象，并初始化各个字段
			tp := &models.TagPost{Tagid:tag.Id, Postid:post.Id, Poststatus:post.Status, Posttime:post.Posttime}
			//插入标签文章对象
			tp.Insert()
		}
		//拼接标签名称
		post.Tags = "," + strings.Join(addtags, ",") + ","
	}

	post.Update("title", "color", "istop", "posttime", "status", "content", "updated", "tags")
	this.Redirect("/admin/article/list", 302)
}




















