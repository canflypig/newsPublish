package controllers

import (
	"fmt"
	//"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"math"
	"path"
	"shanghaiyiqi/models"
	"time"
)

type ArticleController struct {
	beego.Controller
}

// 展示文章列表页
func (this *ArticleController) ShowArticleList() {

	//session 判断
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	//获取数据
	//高级查询
	//指定表
	o := orm.NewOrm()
	qs := o.QueryTable("Article") //queryseter
	var articles []models.Article
	//总数据
	//_, err := qs.All(&articles)
	//if err != nil {
	//	fmt.Println("查询数据错误")
	//}
	//查询总记录数
	typeName := this.GetString("select")
	var count int64

	//获取一页显示第三条数据
	pageSize := 2
	//获取总页数

	//获取传递过来的第一页页码数
	pageIndex, err1 := this.GetInt("pageIndex")
	fmt.Println("pageCount:", pageIndex)
	if err1 != nil {
		fmt.Println("err1:", err1)
		pageIndex = 1
	}
	//获取数据
	/*作用就三获取数据库部分数据，第一个参数 获取几条
	  第二个参数 从那条数据开始获取 返回值还是querySeter
	  起始位置计算

	RelatedSel("ArticleType")  关联文章类型
	*/

	start := (pageIndex - 1) * pageSize

	if typeName == "" {
		count, _ = qs.Count()
	} else {
		count, _ = qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
	}
	pageCount := math.Ceil(float64(count) / float64(pageSize))
	//获取文章类型
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types

	//根据选中类型查询相应类型文章

	fmt.Println("typeName:", typeName)
	if typeName == "" {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
	} else {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	}

	//传递数据
	this.Data["userName"] = userName.(string)
	this.Data["typeName"] = typeName
	this.Data["pageIndex"] = pageIndex
	this.Data["pageCount"] = int(pageCount)
	this.Data["count"] = count
	this.Data["articles"] = articles

	//指定视图布局
	this.Layout = "layout.html"
	this.TplName = "index.html"
}

// 展示添加文章页面
func (this *ArticleController) ShowAddArticle() {
	//查询所有类型数据
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	//传递数据
	this.Data["types"] = types
	this.TplName = "add.html"
}

// 获取添加文章数据
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
	//处理文件上传
	file, head, err := this.GetFile("uploadname")
	defer file.Close()
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return
	}
	//1.文件大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "文件太大，请重新上传"
		this.TplName = "add.html"
		return
	}
	//2.文件格式
	//a.jpg 获取图片名后缀.jpg
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "文件格式错误。请重新上传"
		this.TplName = "add.html"
		return
	}
	//3.防止重名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext
	//存储
	this.SaveToFile("uploadname", "./static/img/"+fileName)

	//fmt.Println("fileName:" + fileName)

	//处理数据
	//插入操作
	o := orm.NewOrm()

	var articles models.Article
	articles.ArtiName = articleName
	articles.Acontent = content
	articles.Aimg = "/static/img/" + fileName
	//给文章添加类型
	//获取类型数据
	typeName := this.GetString("select")
	//fmt.Println(typeName)
	//根据名称查询类型
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType, "TypeName")
	articles.ArticleType = &articleType
	o.Insert(&articles)

	//4.返回页面
	this.Redirect("/article/showArticleList", 302)
}

// 展示文章详情页面
func (this *ArticleController) ShowArticleDetail() {

	//获取数据
	id, er := this.GetInt("articleId")
	//数据校验
	if er != nil {
		fmt.Println("传递连接错误")
	}
	//操作数据
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	//o.Read(&article)
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id", id).One(&article)

	//修改阅读量
	article.Acount += 1
	o.Update(&article)

	//多对多插入浏览记录
	/*
		1获取orm对象
		2获取操作对象
		3获取多对多操作对象
		4获取插入对象
		5插入
	*/
	m2m := o.QueryM2M(&article, "Users")
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	var user  models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")
	//插入操作
	m2m.Add(user)

	//查询
	//o.LoadRelated(&article,"Users")  这种没办法过滤
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id", id).Distinct().All(&users)

	//返回试图
	this.Data["users"] = users
	this.Data["article"] = article

	userNameLayout := this.GetSession("userName")
	this.Data["userName"] = userNameLayout.(string)
	this.Layout = "layout.html"
	this.TplName = "content.html"
}

// 显示编辑页面
func (this *ArticleController) ShowUpdateArticle() {
	//获取数据
	id, err := this.GetInt("articleId")
	//校验数据
	if err != nil {
		fmt.Println("请求文章错误")
		return
	}
	//数据处理
	//查询相应文章
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)

	//返回视图
	this.Data["article"] = article
	this.TplName = "update.html"
}

// 封装上传文件函数
func UploadFile(this *beego.Controller, filePath string) string {
	//处理文件上传
	file, head, err := this.GetFile(filePath)
	if head.Filename == "" {
		return "NoImg"
	}
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return ""
	}
	defer file.Close()
	//1.文件大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "文件太大，请重新上传"
		this.TplName = "add.html"
		return ""
	}

	//2.文件格式
	//a.jpg
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "文件格式错误。请重新上传"
		this.TplName = "add.html"
		return ""
	}
	//3.防止重名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext
	//存储
	this.SaveToFile(filePath, "./static/img/"+fileName)

	return "/static/img/" + fileName
}

// 处理编辑页面数据
func (this *ArticleController) HandleUpdateArticle() {
	//获取数据
	id, err := this.GetInt("articleId")
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	filePath := UploadFile(&this.Controller, "uploadname")
	//数据校验
	if err != nil || articleName == "" || content == "" || filePath == "" {
		fmt.Println("请求错误")
		return
	}

	//数据处理
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		fmt.Println("更新的文章不存在")
		return
	}
	article.ArtiName = articleName
	article.Acontent = content
	if filePath != "NoImg" {
		article.Aimg = filePath
	}
	o.Update(&article)
	//返回视图
	this.Redirect("/article/showArticleList", 302)
}

// 删除文章
func (this *ArticleController) DeleteArticle() {
	//获取数据
	id, err := this.GetInt("articleId")
	//校验数据
	if err != nil {
		fmt.Println("删除文章请求路径错误")
		return
	}
	//数据处理
	//删除操作
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Delete(&article)

	//返回视图
	this.Redirect("/article/showArticleList", 302)

}

// 展示添加类型
func (this *ArticleController) ShowAddType() {

	//查询
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	//传递数据
	this.Data["types"] = types

	userNameLayout := this.GetSession("userName")
	this.Data["userName"] = userNameLayout.(string)
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

// 处理添加类型数据
func (this *ArticleController) HandleAddType() {

	//获取数据
	typeName := this.GetString("typeName")
	//校验数据
	if typeName == "" {
		fmt.Println("信息不完整,请重新输入")
		return
	}
	//处理数据
	//插入数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Insert(&articleType)

	//返回视图
	this.Redirect("/article/addType", 302)
}

// 删除类型
func (this *ArticleController) DeleteType() {
	//获取数据
	id, err := this.GetInt("id")
	if err != nil {
		fmt.Println("删除类型错误", err)
		return
	}
	//处理数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	o.Delete(&articleType)
	//返回视图
	this.Redirect("/article/addType", 302)

}
