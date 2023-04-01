package models

//mysql 操作数据库
//import (
//	"database/sql"
//	"fmt"
//	_ "github.com/beego/beego/v2/adapter/session/mysql"
//)

//ORM 操作数据库
import (
	"github.com/beego/beego/v2/client/orm"
	//"github.com/beego/beego/v2/adapter/orm"
	//_ "github.com/beego/beego/v2/adapter/session/mysql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//表的设计

// 定义一个结构体
type User struct {
	Id       int
	Name     string
	PassWord string
	//Pass_Word
	Articles []*Article `orm:"rel(m2m)"`
}

type Article struct {
	Id       int       `orm:"pk;auto"`
	ArtiName string    `orm:"size(20)"`
	Atime    time.Time `orm:"auto_now"`
	Acount   int       `orm:"default(0);null"`
	Acontent string    `orm:"size(500)"`
	Aimg     string    `orm:"size(100)"`

	ArticleType *ArticleType `orm:"rel(fk)"`
	Users       []*User      `orm:"reverse(many)"`
}

// 类型表
type ArticleType struct {
	Id       int
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"`
}

// 在orm 中双下划线__是有特殊含义的

func init() {
	//
	//conn, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	//defer conn.Close() //随手关闭数据库是个好习惯
	//if err != nil {
	//	//beego.info("连接错误", err)
	//	fmt.Println("连接错误", err)
	//	//beego.Error("连接错误", err)
	//	return
	//}

	//创建表
	//res, err1 := conn.Exec("create table itcast(name VARCHAR(40),password VARCHAR(40))")
	//if err1 != nil {
	//	//beego.info("连接错误", err)
	//	//beego.error("连接错误", err)
	//	fmt.Println("创建表错误", err)
	//	return
	//}
	//fmt.Println("创建表信息", res)

	//插入数据
	//conn.Exec("insert into itcast(name,password) values(?,?)", "chuanzhi", "heima")
	//查询
	//res, err2 := conn.Query("select name from itcast")
	//if err2 != nil {
	//	fmt.Println("查询错误")
	//}
	//var name string
	//for res.Next() {
	//	res.Scan(&name)
	//	fmt.Println(name)
	//}

	//ORM操作数据库
	//获取连接对象
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	//创建表

	orm.RegisterModel(new(User), new(Article), new(ArticleType))

	//生成表
	//第一个参数是数据库别名，第二个参数三 是否强制更新
	orm.RunSyncdb("default", false, true)

	//err := orm.RunSyncdb("default", false, true)

	//if err != nil {
	//	fmt.Println("orm.RunSyncdb err:", err)
	//}
	//操作表
}
