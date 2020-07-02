package controllers

import "C"
import (
	"bj3q/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller  //继承了beego.Controller
}

func (c *MainController) Get() {   //beego默认的get请求访问get方法，post请求访问post方法
	c.Data["Website"] = "beego.me"  //传递数据给视图
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["data"]="china"
	c.TplName = "test.html"  //指定视图页面
}
func (c *MainController)Post(){

	c.Data["data"]="post请求页面"
	c.TplName="test.html"
}


func(c *MainController)ShowGet(){
	//获取ORM对象
     o:=orm.NewOrm()
	//执行某个操作函数 增删改查

	//插入操作 （插入结构体对象）


     var user models.User
     user.Name="heima"
     user.PassWord="chuanzhi"

     //插入操作
     count,err:=o.Insert(&user)
     if err!=nil{
     	beego.Error("插入失败")
	 }
	 beego.Info(count)



	//查询操作
	/*
	//创建查询对象
	var user models.User
	user.Id=1

	err:=o.Read(&user,"Id")
	if err!=nil{
		beego.Error("查询失败")

	}


	//返回结果
	beego.Info(user)
*/

	//更新操作
	/*
	//获取更新对象
	var user models.User
	user.Id=1
	err:=o.Read(&user,"Id")
	if err!=nil{
		beego.Error("要更新的数据不存在")
	}
	user.Name="shanghaigengxin"
	count,err:=o.Update(&user)
	if err!=nil{
		beego.Error("更新失败")
	}
	beego.Info(count)
*/

	//删除操作
	/*
	//创建删除对象
	var user models.User
	user.Id=1
	user.Name="heima"
//如果不查询，直接删除，删除对象的主键要有值
	count,err:=o.Delete(&user,"Name")
	if err!=nil{
		beego.Error("删除失败")
	}
	beego.Info(count)
*/
	c.Data["data"]="上海"
	c.TplName="test.html"
}

