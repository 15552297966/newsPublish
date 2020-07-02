package models

import (
	"database/sql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//传统方法操作数据库
func init1() {
	//操作数据库的代码

	//第一个参数：数据库驱动
	//第二个参数：连接数据库参数
	db, err := sql.Open("mysql", "root:11111111@tcp(127.0.0.1:3306)/stu?charset=utf8")
	if err != nil {
		beego.Error("连接数据库错误", err)
		return
	}
	//关闭数据库
	defer db.Close()
	//创建表
	//_,err=db.Exec("create table itcast(name varchar(20),password varchar(20));")
	//if err!=nil{
	//	beego.Error("创建表失败：",err)
	//	return
	//}

	//插入数据
	//db.Exec("insert into itcast(name,password) values (?,?)","chuanzhi","heima")

	//查询
	res, err := db.Query("select name from itcast")
	var name string
	for res.Next() {
		res.Scan(&name)
		beego.Info(name)
	}

}

//表的设计
//定义一个结构体
type User struct {
	Id       int //在没有设置主键时，ORM会自动选名为id，类型为int的列为主键
	Name     string
	PassWord string //这里的修改某一列会新增一列
	//Pass_Word 通过ORM在数据库里面创建为Pass__word，这样是不对的，__在ORM有特殊含义。

	Articles []*Article `orm:"reverse(many)"`
}
type Article struct {
	Id       int       `orm:"pk;auto"`
	ArtiName string    `orm:"size(20)"`
	Atime    time.Time `orm:"auto_now"`
	Acount   int       `orm:"default(0);null"`
	Acontent string    `orm:"size(500)"`
	Aimg     string    `orm:"size(100)"`

	ArticleType *ArticleType `orm:"rel(fk)"` //多的一方设置外键 rel正向设置 与reverse相对应
	                                         //在mysql中，一对多关系中，多的一方表中增加一个属性article_type_id,一的一方不变。D



	Users []*User `orm:"rel(m2m)"`   //多对多关系中，会创建一个关系表article_users
}

//文章类型表
type ArticleType struct {
	Id       int
	TypeName string `orm:"size(20)"`

	Articles []*Article `orm:"reverse(many)"` //一的一方设置少的一方的切片 reverse反向设置
}

//在ORM里面__是有特殊含义的

//ORM方法操作数据库
func init() {
	//获取连接对象 第一个参数是数据库别名
	orm.RegisterDataBase("default", "mysql", "root:11111111@tcp(127.0.0.1:3306)/stu?charset=utf8")
	//注册表
	orm.RegisterModel(new(User), new(Article), new(ArticleType))

	//生成表
	//第一个参数是数据库别名。
	//第二个参数是是否强制更新这个表（所谓强制更新，每次生成表的时候先把之前的表drop掉。为false时，会在你表
	//变动比较大的时候，也会给你重新生成表，改变表结构）。
	//第三个参数，生成表的SQL语句是否可见。
	orm.RunSyncdb("default", false, true)

}
