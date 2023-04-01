package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/adapter/orm"
	beego "github.com/beego/beego/v2/server/web"
	"shanghaiyiqi/models"
)

type UserController struct {
	beego.Controller
}

// 显示注册页面
func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

// 处理注册数据
func (this *UserController) HandlePost() {
	//获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	//fmt.Println(userName, pwd)
	//2.校验数据
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "注册数据不完整,请重新注册"
		fmt.Println("注册数据不完整,请重新注册")
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
	//插入
	o.Insert(&user)
	//返回结果

	//4.返回页面
	//this.Ctx.WriteString("注册页面")
	this.Redirect("/login", 302)
	//this.TplName = "login.html"
}

// 展示登录页面
func (this *UserController) ShowLogin() {
	//this.Data["data"] = "hehehe"
	userName := this.Ctx.GetCookie("userName")
	if userName == "" {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	} else {
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	}

	this.TplName = "login.html"
}

// 登录页面
func (this *UserController) HandleLogin() {
	//1.获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	//2.校验数据
	if userName == "" || pwd == "" {
		//this.Data["errmsg"] = "呵呵呵"
		this.TplName = "login.html"
		return
	}
	//3.操作数据
	//1.获取orm对象
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		this.Data["errmsg"] = "用户不存在"
		this.TplName = "login.html"
		return
	}
	if user.PassWord != pwd {
		this.Data["errmsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}
	//4.返回页面
	//this.Ctx.WriteString("登录成功")

	data := this.GetString("remember")
	fmt.Println(data)
	if data == "on" {
		this.Ctx.SetCookie("userName", userName, 100)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}
	this.SetSession("userName", userName)

	this.Redirect("/article/showArticleList", 302)

}

// 退出登录
func (this *UserController) Logout() {
	//删除session
	this.DelSession("userName")
	//跳转登录页面
	this.Redirect("/login", 302)
}
