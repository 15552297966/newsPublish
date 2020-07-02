package controllers

import (
	"bj3q/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"time"
)

type ArticleController struct {
	beego.Controller
}

//展示文章列表页
func (this *ArticleController) ShowArticleList() {
	//session判断
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/login",302)
		return
	}

	//从数据库取数据
	//高级查询
	//指定表
	o:=orm.NewOrm()
	qs:=o.QueryTable("Article")  //queryseter查询集对象
	var articles []models.Article
	//_,err:=qs.All(&articles)
	//if err!=nil{
	//	beego.Info("查询数据错误")
	//}


	//获取总页数   天花板函数：两个浮点数相除，得到一个浮点数，向上取整
	//            地板函数：两个浮点数相除，得到一个浮点数，向下取整
	pageSize:=2

	//根据选中类型查询相应的类型文章
	typeName:=this.GetString("select")
	//beego.Info(typeName)

	//查询总记录数
	var count int64

	////获取页码
	pageIndex,err:=this.GetInt("pageIndex")
	if err!=nil{
		pageIndex=1
	}

	beego.Info(pageIndex)

	//获取数据
	//作用就是获取数据库部分数据，第一个参数，获取几条,第二个参数，从哪条数据开始获取
	//返回值还是querySeter
	//qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)

	//起始位置
	start:=(pageIndex-1)*pageSize
	if typeName==""{
		count,_=qs.Count()
	}else{
		count,_=qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()
	}
	//共多少页
	pageCount:=math.Ceil(float64(count)/ float64(pageSize))

	//获取文章类型
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"]=types



	if typeName==""{ //刚登录的时候，没有指定类型，应该能查到所有文章
		qs.Limit(pageSize,start).All(&articles)  //从start开始取，取2个数据，存到articles对象中
	}else {
		//拿到typeName后，RelatedSel是多表联查，Article表与ArticleType表联查，通过Filter指定ArticleType的TypeName字段==typeName的数据查出来放在articles中。
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)

	}



	//传递数据
	this.Data["session"]=userName
	this.Data["typeName"]=typeName  //当前选中的文章类型，是通过下拉框获得的
	this.Data["count"]=count
	this.Data["articles"]=articles
	this.Data["pageCount"]=int(pageCount)
	this.Data["pageIndex"]=pageIndex


	//指定视图布局
	this.Layout="layout.html"
	this.TplName = "index.html"


}

//展示添加文章页面
func (this *ArticleController) ShowAddArticle() {
	//查询所有类型数据，并展示
	o:=orm.NewOrm()

	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	//传递数据
	this.Data["types"]=types
	userName:=this.GetSession("userName").(string)

	this.Data["session"]=userName
	this.Layout="layout.html"
	this.TplName = "add.html"
}

//获取添加文章数据
func (this *ArticleController) HandleAddArticle() {
	//1.获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")

	//2.校验数据
	if articleName == "" || content == "" {
		this.Data["errmsg"] = "添加数据不完整"
		this.TplName = "add.html"
		return
	}
	beego.Info(articleName, content)


	//3.处理数据
    filePath:=UploadFile(&this.Controller,"uploadname")
	//插入操作
	o:=orm.NewOrm()
	var article models.Article
	article.ArtiName=articleName
	article.Acontent=content
	article.Aimg=filePath
	article.Atime=time.Now()
	//给文章添加类型
	//获取类型数据
	typeName:=this.GetString("select")
	//根据名称查询类型
	var articleType models.ArticleType
	articleType.TypeName=typeName
	o.Read(&articleType,"TypeName")
	article.ArticleType=&articleType
	_,err:=o.Insert(&article)

	//之前一直插入失败是因为数据没填全。。。。。我草了
	if err!=nil{
		this.Data["errmsg"] = "插入数据库失败"
		beego.Error("插入数据库失败")
		this.TplName = "add.html"
		return
	}
	//4.返回页面
	this.Redirect("/article/showArticleList",302)



}

//展示文章详情页面
func(this*ArticleController)ShowArticleDetail(){
	//获取数据
	id,err:=this.GetInt("articleId")
	//校验数据
	if err!=nil{
		beego.Info("传递的链接有误")
	}
	//操作数据
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id

	//o.Read(&article)
	//高级查询
	qs:=o.QueryTable("Article")
	qs.RelatedSel("ArticleType").Filter("Id",id).One(&article) //过滤的时候如果是本张表的字段不需要加表__字段
	                                                                    //一条用One，存对象里。多条用all，存切片里。

	                                                                    //修改阅读量

	article.Acount+=1
	o.Update(&article)

	//多对多插入浏览记录
	//获取多对多操作对象   表    字段
	m2m:=o.QueryM2M(&article,"Users")
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/login",302)
		return
	}
	var user models.User
	user.Name=userName.(string)
	o.Read(&user,"Name")

	//多对多插入
	m2m.Add(user)

	//多对多查询
	//查询方法一： 一般查询
	//o.LoadRelated(&article,"Users")
	//查询方法二： 高级查询
	var users []models.User
	//我们插入的时候是向article中插user，查询的时候是从user中取获取。 Filter的参数是User做多表的字段值__表名__表的主键
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)  //Distinct()去重



	//返回视图页面
	this.Data["title"]="文章详情"
	this.Data["users"]=users
	this.Data["article"]=article
	//传给layout的数据
	this.Data["session"]=userName
	this.Layout="layout.html"
	this.TplName="content.html"
}

//显示编辑页面
func(this*ArticleController)ShowUpdateArticle(){
	//获取数据
	id,err:=this.GetInt("articleId")

	//校验
	if err!=nil{
		beego.Info("请求文章错误")
		return
	}
	//处理
	//查询相应文章
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	o.Read(&article)
	//返回视图
	this.Data["article"]=article
	userName:=this.GetSession("userName").(string)

	this.Data["session"]=userName
	this.Layout="layout.html"
	this.TplName="update.html"
}

//封装上传文件函数
func UploadFile(this *beego.Controller,filePath string)string{
	//处理文件上传
	//后台能获取到文件上传的两个函数

	//获取文件 head是文件头
	file,head,err := this.GetFile(filePath)  //bug：前端form未加enctype="multipart/form-data" 导致文件上传失败

	//defer file.Close()  //这句加在这会报错 但我不懂为啥放在末尾就不会
	                      //我懂了 不是要在末尾 而是在 if err的后面，因为如果err了就没有file了
	                      //beego.Info(head.Filename)
	if head==nil{
		return "NoImg"
	}
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return ""
	}
	defer file.Close()
	//存储文件 第一个参数是前端的文件上传的name，第二个参数是保存的路径+文件名
	/*
		err=this.SaveToFile("uploadname","./static/img/"+head.Filename) //第二个参数中.一定要加
		if err != nil {
			this.Data["errmsg"] = "文件存储失败"
			this.TplName = "add.html"
			return
		}
	*/

	//文件上传后台需要考虑的三个方面
	//1.文件大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "文件太大，请重新上传"
		this.TplName = "add.html"
		return ""
	}
	//2.文件格式
	//path包下专门用来获取后缀名的方法
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "文件格式错误，请重新上传"
		this.TplName = "add.html"
		return ""
	}
	//3.防止重名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext //time格式转成string格式
	//存储
	this.SaveToFile(filePath, "./static/img/"+fileName)
	//defer file.Close()
	return "/static/img/"+fileName


}

//处理编辑页面数据
func (this *ArticleController)HandleUpdateArticle(){
	//获取数据
	id,err:=this.GetInt("articleId")
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	filePath:=UploadFile(&this.Controller,"uploadname")


	//校验
	if err!=nil||articleName==""||content==""||filePath==""{
		beego.Info("请求错误")
		return
	}
	//数据处理
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	err=o.Read(&article)
	if err!=nil{
		beego.Info("更新的文章不存在")
		return
	}
	article.ArtiName=articleName
	article.Acontent=content
	if filePath!="NoImg"{
		article.Aimg=filePath
	}

	o.Update(&article)



	//返回视图
	this.Redirect("/article/showArticleList",302)




}

//删除文章处理
func (this *ArticleController)DeleteArticle () {
	//获取数据
	id,err:=this.GetInt("articleId")
	//校验
	if err!=nil{
		beego.Info("删除文章请求路径错误")
		return
	}
	//删除操作
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	o.Delete(&article)

	//返回页面 这里不能用tplname渲染，因为首页需要数据
	this.Redirect("/article/showArticleList",302)
}

//展示添加分类
func (this *ArticleController)ShowAddType (){
	//查询
	o:=orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	//传递数据
	this.Data["types"]=types
	userName:=this.GetSession("userName").(string)

	this.Data["session"]=userName
	this.Layout="layout.html"
	this.TplName="addType.html"
}

//处理添加类型数据
func (this *ArticleController)HandleAddType (){
	//获取数据
	typeName:=this.GetString("typeName")
	//校验数据
	if typeName==""{
		beego.Info("信息不完整，请重新输入")
		return
	}
	//处理数据
	//插入操作
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName=typeName
	o.Insert(&articleType)
	//返回视图
	this.Redirect("/article/addType",302)
}

//删除类型
func (this *ArticleController)DeleteType (){
	//获取数据
	id,err:=this.GetInt("id")
	//校验数据
	if err!=nil{
		beego.Error("删除类型错误",err)
		return
	}
	//处理数据
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id=id
	o.Delete(&articleType)
	//返回视图
	this.Redirect("/article/addType",302)
}