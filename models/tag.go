package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

//标签
type Tag struct {
	Id int
	Name string `orm:size(20)`//标签名称
	Count int //文章数量
}



func (tag *Tag) TableName() string {
	//从配置文件中获取表的前缀
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "tag"
}

//插入
func (tag *Tag) Insert() error {
	if _, err := orm.NewOrm().Insert(tag); err != nil {
		return err
	}
	return nil
}



//读取
func (tag *Tag) Read(fields ...string) error {
	if err := orm.NewOrm().Read(tag, fields...); err != nil {
		return err
	}
	return nil
}

//更新
func (tag *Tag) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(tag, fields...); err != nil {
		return err
	}
	return nil
}

//删除
/*
思路：根据需要被删除的标签的id在标签文章表中查找相关记录，从这些记录中获取到文章id，
然后根据文章id在文章表中找到对应的记录，将标签名替换掉，最后删除标签文章中的相关记录。
*/
func (tag *Tag) Delete() error {
	var list []*TagPost
	//得到标签文章表的句柄并根据标签id进行过滤
	orm.NewOrm().QueryTable(&TagPost{}).Filter("tagid", tag.Id).All(&list)
	if len(list) > 0 {
		ids := make([]string, 0, len(list))
		//遍历list
		for _, tagpost := range list {
			//将文章id转换为字符串拼接到ids中
			ids = append(ids, strconv.Itoa(tagpost.Postid))
		}
		//UPDATE tb_post SET tags = REPLACE(tags, ':', ',') WHERE id IN(11, 12, 16);
		table := new(Post).TableName()
		//将文章表中的tags替换为逗号
		orm.NewOrm().Raw("update " + table +
			" set tags = REPLACE(tags, ?, ',') where id in (" +
			strings.Join(ids, ",") + ")", ","+tag.Name+",").Exec()
		//删除标签文章表中的相关记录
		orm.NewOrm().QueryTable(&TagPost{}).Filter("tagid", tag.Id).Delete()
	}
	//删除标签
	if  _, err := orm.NewOrm().Delete(tag); err != nil {
		return err
	}
	return nil
}

//2.根据需要被合并的标签id在标签文章表中查找对应的记录，将对应记录的标签id替换目标标签的id
//3.在文章表中将原标签的名称替换为目标标签的名称
func (tag *Tag) MergeTo(to *Tag) {
	var list []*TagPost
	var tp TagPost
	//得到标签文章的句柄
	query:= orm.NewOrm().QueryTable(&tp).Filter("tagid", tag.Id)
	//根据标签id过滤
	query.All(&list)
	if len(list) > 0 {
		//在标签文章表将需要被合并的标签id替换为目标标签的id
		query.Update(orm.Params{"tagid": to.Id})

		ids := make([]string, 0, len(list))
		//将文章id拼接到切片中
		for _, v :=  range list {
			ids = append(ids, strconv.Itoa(v.Postid))
		}

		//UPDATE tb_post SET tags = REPLACE(tags, ':', ',') WHERE id IN(11, 12, 16);
		orm.NewOrm().Raw("update " + new(Post).TableName() +
			" set tags = REPLACE(tags, ?, ?) where id in (" +
			strings.Join(ids, ",") + ")", ","+tag.Name+",", ","+to.Name+",").Exec()
	}
}


//更新标签的count字段
func (tag *Tag) UpCount() {
	//获取标签文章表的句柄并根据tagid进行过滤，得到符合条件的记录的数量
	count, err := orm.NewOrm().QueryTable(&TagPost{}).Filter("tagid", tag.Id).Count()
	newcount := int(count)
	//查询没有出现错误并且所得到的数量和标签中原来的count不一样
	if err == nil && newcount != tag.Count {
		tag.Count = newcount
		tag.Update("count")
	}
}