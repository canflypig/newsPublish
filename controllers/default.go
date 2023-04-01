package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["data"] = "China"
	//c.TplName = "index.tpl"
	c.TplName = "test.html"
}

func (c *MainController) Post() {
	c.Data["data"] = "上海一期最棒"
	c.TplName = "test.html"

}
func (c *MainController) ShowGet() {
	//获取ORM对象
	//o := orm.NewOrm()
	//执行某个操作行署 增删改查
	//插入
	/*
		var user models.User
		user.Name = "heima"
		user.PassWord = "chuanzhi"
		//插入操作
		count, err := o.Insert(&user)
		if err != nil {
			fmt.Println("插入失败")
		}
		fmt.Println("################################################################################################")
		fmt.Println(count)
	*/

	//查询
	/*
		var user models.User
		user.Id = 1
		//err := o.Read(&user, "Id")
		//默认主键查询可以省略
		err := o.Read(&user)
		if err != nil {
			fmt.Println("查询失败")
		}
		fmt.Println(user)
	*/

	//更新操作
	/*
		var user models.User
		user.Id = 1
		err := o.Read(&user)
		if err != nil {
			fmt.Println("更新的数据不存在")
		}
		user.Name = "shanghiayiqi"
		count, err1 := o.Update(&user)
		if err1 != nil {
			fmt.Println("更新失败")
		}
		fmt.Println(count)
	*/
	//删除操作
	/*var user models.User
	user.Id = 1
	//如果不查询 直接删除 删除对象的主键要有值
	count, err := o.Delete(&user)
	if err != nil {
		fmt.Println("删除失败")
	}
	fmt.Println(count)



	*/
	c.Data["data"] = "上海"
	c.TplName = "test.html"
}
