package main

import (
	_ "bj3q/models"
	_ "bj3q/routers" //下划线的作用：执行包里的init函数
	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("prepage",ShowPrePage)
	beego.AddFuncMap("nextpage",ShowNextPage)

	beego.Run()
}

//后台定义一个函数
func ShowPrePage(pageIndex int)int{
	if pageIndex==1{
		return pageIndex
	}
return pageIndex-1
}

func ShowNextPage(pageIndex int,pageCount int)int{
	if pageIndex==pageCount{
		return pageIndex
	}
	return pageIndex+1
}