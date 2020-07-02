package controllers

import (
	"bj3q/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

//显示注册页面
func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

//处理注册数据
func (this *UserController) HandlePost() {
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("passWord")

	//beego.Info(userName,pwd)
	//2.校验数据
	//判空校验
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "注册数据不完整，请重新注册"
		beego.Info("注册数据不完整，请重新注册")
		this.TplName = "register.html"
		return
	}
	//3.操作数据

	//获取ORM对象
	o := orm.NewOrm()
	//获取插入对象
	var user models.User
	//给插入对象赋值
	user.Name = userName
	user.PassWord = pwd
	//执行插入操作
	o.Insert(&user)

	//返回结果
	//4.返回页面

	//this.Ctx.WriteString("注册成功")

	//重定向
	//重定向用到的方法是this.Redirect()函数，第一个参数是请求路径，第二个参数是http状态码。
	//请求路径就不说了，就是和超链接一样的路径。
	//重定向函数（跳转函数） 第一个参数请求路径 第二个参数状态码
	//HTTP状态码的内容
	//1xx 服务端已经接收到了客户端请求，但是客户端应当继续发送请求。100
	//2xx 200 请求成功
	//3xx 300 302 请求的资源转换路径了，重定向（请求被跳转）
	//4xx 404 请求端的错误
	//5xx 500 服务端错误

	//重定向的工作流程：
	//1. 当服务端向客户端响应redirect后，并没有提供任何view数据进行渲染，仅仅是告诉浏览器响应为redirect，以及重定向的目标地址。
	//2. 浏览器受到服务端redirect过来的响应，会再次发起一个http请求。
	//3. 由于是浏览器再次发起了一个新的http请求，所以浏览器地址栏中的url会发生变化。

	this.Data["tplname"]="服务器直接渲染"
	this.Redirect("/login",302) //因为是浏览器再次发送了请求，所以浏览器url会改变。
	//因为服务器向客户端发送了302状态码，不能向客户端传递数据，但是能够获取到加载页面时的数据。
	//（这里说的获取到加载页面时的数据应该是说显示视图的函数中的数据）。
	//使用场景：页面跳转的时候。


	//渲染方式：直接给浏览器返回视图
	//this.Data["tplname"]="服务器直接渲染"
	//this.TplName="login.html"  //直接向客户端返回了页面，是可以传递数据this.Data的，但如果要获取加载页面时数据，需要再次写相关代码。
	//（这里说的获取到加载页面时的数据应该是说显示视图的函数中的数据）。
    //使用场景：页面加载或者是登录注册失败的时候。

}

//显示登陆页面
func (this *UserController) ShowLogin() {
	userName:=this.Ctx.GetCookie("userName")
	if userName==""{
		this.Data["userName"]=""
		this.Data["checked"]=""
	}else{
		this.Data["userName"]=userName
		this.Data["checked"]="checked"
	}
	this.TplName = "login.html"
}

//处理登陆数据
func (this *UserController) HandleLogin() {
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("passWord")
	//2.校验数据
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "登录数据不完整，请重新登录"
		beego.Info("注册数据不完整，请重新注册")
		this.TplName = "login.html"
		return
	}
	//3.操作数据（登录是用来查询的）
	//获取orm对象
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err:=o.Read(&user, "Name")
	if err!=nil{
		this.Data["errmsg"]="用户不存在"
		this.TplName="login.html"
		return
	}
	if user.PassWord!=pwd{
		this.Data["errmsg"]="密码不正确"
		this.TplName="login.html"
		return
	}
	//4.返回页面
	//this.Ctx.WriteString("登录成功")

	data:=this.GetString("remember")
	beego.Info(data)
	if data=="on"{
		//登录成功后，设置一个cookie
		this.Ctx.SetCookie("userName",userName,100)
	}else{

		this.Ctx.SetCookie("userName",userName,-1)
	}

	//设置session
	this.SetSession("userName",userName)

	this.Redirect("/article/showArticleList",302)
}

//退出登录
func (this *UserController) Logout(){
	//删除session
	this.DelSession("userName")
	//跳转登录页面
	this.Redirect("/login",302)
}
