本文章是使用 Go 來寫一個 Repository Restful API 的留言板，並且會使用 gin 以及 gorm (使用 Mysql)套件。

建議可以先觀看 [Go 介紹](https://pin-yi.me/blog/golang/go-introduce/) 文章來簡單學習 Go 語言。

版本資訊

* macOS：11.6
* Go：go version go1.18 darwin/amd64
* Mysql：mysql  Ver 8.0.28 for macos11.6 on x86_64 (Homebrew)
<br>

## 實作

### 檔案結構

```
.
├── controller
│   └── controller.go
├── go.mod
├── go.sum
├── main.go
├── model
│   └── model.go
├── repository
│   └── repository.go
├── router
│   └── router.go
└── sql
    ├── connect.yaml
    └── sql.go
```

<br>

我們來說明一下上面的資料夾個別功能與作用

* sql：放置連線資料庫檔案。
* controller：商用邏輯控制。
* model：定義資料表資料型態。
* repository：處理與資料庫進行交握。
* router：設定網站網址路由。

### go.mod

一開始我們創好資料夾後，要先來設定 go.mod 的 module

```sh
$ go mod init message
```

* go.mod 檔案
```sh
module message

go 1.18
```

<br>

接著使用 `go get` 來引入 `gin`、`gorm`、`mysql`、`yaml` 套件
```sh
$ go get -u github.com/gin-gonic/gin
$ go get -u gorm.io/gorm
$ go get -u gorm.io/driver/mysql
$ go get -u gopkg.in/yaml.v2
```
可以在查看一下 go.mod 檔案是否多了很多 indirect

<br>

### main.go

```go
package main

import (
	"fmt"
	"message/model"
	"message/router"
	"message/sql"
)

func main() {
	//連線資料庫
	if err := sql.InitMySql(); err != nil {
		panic(err)
	}

	//連結模型
	sql.Connect.AutoMigrate(&model.Message{})
	//sql.Connect.Table("message") //也可以使用連線已有資料表方式

	//註冊路由
	r := router.SetRouter()

	//啟動埠為8081的專案
	fmt.Println("開啟127.0.0.0.1:8081...")
	r.Run("127.0.0.1:8081")
}
```
引入我們 Repository 架構，將 config、model、router 導入，先測試是否可以連線資料庫，使用 `AutoMigrate` 來新增資料表(如果沒有才新增)，或是使用 Table 來連線已有資料表，註冊網址路由，最後啟動專案，我們將 Port 設定成 8081。

<br>

### sql

我們剛剛有引入 `yaml` 套件，因為我們設定檔案會使用 yaml 來編輯

* connect.yaml
```yaml
host: 127.0.0.1
username: root
password: "密碼"
dbname: "資料庫名稱"
port: 3306
```
我們把 mysql 連線的資訊寫在此處。

<br>

* sql.go (下面為一個檔案，但長度有點長，分開說明)

```go
package sql

import (
	"io/ioutil"
	"fmt"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)
```
import 會使用到的套件。

<br>

```go
var Connect *gorm.DB

type conf struct {
	Host     string `yaml:"host"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
	Port     string `yaml:"port"`
}

func (c *conf) getConf() *conf {
	//讀取config/connect.yaml檔案
	yamlFile, err := ioutil.ReadFile("sql/connect.yaml")

	//若出現錯誤，列印錯誤訊息
	if err != nil {
		fmt.Println(err.Error())
	}

	//將讀取的字串轉換成結構體conf
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}
```
設定資料庫連線的 conf 來讀取 yaml 檔案。

<br>

```go
//初始化連線資料庫
func InitMySql() (err error) {
	var c conf

	//獲取yaml配置引數
	conf := c.getConf()

	//將yaml配置引數拼接成連線資料庫的url
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName,
	)

	//連線資料庫
	Connect, err = gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{})
	return
}
```
初始化資料庫，會把剛剛讀取 yaml 的 conf  串接成可以連接資料庫的 url ，最後連線資料庫。

<br>

### router.go

```go
package router

import (
	"message/controller"
	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	//顯示 debug 模式
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	v1 := r.Group("api/v1")
	{
		//新增留言
		v1.POST("/message", controller.Create)
		//查詢全部留言
		v1.GET("/message", controller.GetAll)
		//查詢 {id} 留言
		v1.GET("/message/:id", controller.Get)
		//修改 {id} 留言
		v1.PATCH("/message/:id", controller.Update)
		//刪除 {id} 留言
		v1.DELETE("/message/:id", controller.Delete)
	}
	return r
}
```
設定路由，版本 v1 網址是 `api/v1` ，分別是新增留言、查詢全部留言、查詢 {id} 留言、修改 {id} 留言、刪除 {id} 留言，連接到不同的 `controller` function 。

<br>

### model.go

```go
package model

import 	"gorm.io/gorm"

func (Message) TableName() string {
	return "message"
}

type Message struct {
	Id        int    `gorm:"primary_key,type:INT;not null;AUTO_INCREMENT"`
	User_Id   int    `json:"User_Id"  binding:"required"`
	Content   string `json:"Content"  binding:"required"`
	Version   int    `gorm:"default:0"`
	// 包含 CreatedAt 和 UpdatedAt 和 DeletedAt 欄位
    gorm.Model
}
```
設定資料表的結構，使用 gorm.Model 預設裡面會包含 CreatedAt 和 UpdatedAt 和 DeletedAt 欄位。

<br>

### controller.go

**(下面為一個檔案，但長度有點長，分開說明)**

```go
package controller

import (
	"message/model"
	"message/repository"
	"net/http"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)
```
import 會使用到的套件。

<br>

**查詢留言功能**

```go
func GetAll(c *gin.Context) {
	message, err := repository.GetAllMessage()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func Get(c *gin.Context) {
	var message model.Message

	if err := repository.GetMessage(&message, c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "找不到留言"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}
```
`GetAll()` 會使用到 `repository.GetAllMessage()` 查詢並回傳顯示查詢的資料。

`c.Param("id")` 是網址讀入後的 id，網址是`http://127.0.0.1:8081/api/v1/message/{id}` ，將輸入的 id 透過 `repository.GetMessage()` 查詢並回傳顯示查詢的資料。


<br>

**新增留言功能**

```go
func Create(c *gin.Context) {
	var message model.Message

	if c.PostForm("Content") == "" || utf8.RuneCountInString(c.PostForm("Content")) >= 20 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "沒有輸入內容或長度超過20個字元"})
		return
	}

	c.Bind(&message)
	repository.CreateMessage(&message)
	c.JSON(http.StatusCreated, gin.H{"message": message})
}
```

使用 Gin 框架中的 `Bind 函數`，可以將 url 的查詢參數 query parameter，http 的 Header、body 中提交的數據給取出，透過 `repository.CreateMessage()` 將要新增的資料帶入，如果失敗就顯示 `http.StatusBadRequest`，如果成功就顯示 `http.StatusCreated` 以及新增的資料。

<br>

**修改留言功能**

```go
func Update(c *gin.Context) {
	var message model.Message

	if c.PostForm("Content") == "" || utf8.RuneCountInString(c.PostForm("Content")) >= 20 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "沒有輸入內容或長度超過20個字元"})
		return
	}

	if err := repository.UpdateMessage(&message, c.PostForm("Content"), c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "找不到留言"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}
```

先使用 `repository.GetMessage()` 以及 `c.Param("id")` 來查詢此 id 是否存在，再帶入要修改的 `Content` ，透過 `repository.UpdateMessage()` 將資料修改，，如果失敗就顯示 `http.StatusNotFound` 以及找不到留言，如果成功就顯示 `http.StatusOK` 以及修改的資料。

<br>

**刪除留言功能**

```go
func Delete(c *gin.Context) {
	var message model.Message

	if err := repository.DeleteMessage(&message, c.Param("id")); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "找不到留言"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "刪除留言成功"})
}
```

透過 `repository.DeleteMessage()` 將資料刪除，如果失敗就顯示 `http.StatusNotFound ` 以及找不到留言，如果成功就顯示 `http.StatusNoContent`。

<br>

### repository.go

**(下面為一個檔案，但長度有點長，分開說明)**

所有的邏輯判斷都要在 controller 處理，所以 repository.go 就單純對資料庫就 CRUD：

```go
package repository

import (
	"message/model"
	"message/sql"
)
```
import 會使用到的套件。

<br>

**查詢留言資料讀取**

```go
//查詢全部留言
func GetAllMessage() (message []*model.Message, err error) {
	err = sql.Connect.Find(&message).Error
	return
}

//查詢 {id} 留言
func GetMessage(message *model.Message, id string) (err error) {
	err = sql.Connect.Where("id=?", id).First(&message).Error
	return
}
```

<br>

**新增留言資料讀取**

```go
//新增留言
func CreateMessage(message *model.Message) (err error) {
	err = sql.Connect.Create(&message).Error
	return
}
```

<br>

**修改留言資料讀取**

```go
//更新 {id} 留言
func UpdateMessage(message *model.Message, content, id string) (err error) {
	err = sql.Connect.Where("id=?", id).First(&message).Update("content", content).Error
	return
}
```

<br>

**刪除留言資料讀取**

```go
//刪除 {id} 留言
func DeleteMessage(message *model.Message, id string) (err error) {
	err = sql.Connect.Where("id=?", id).First(&message).Delete(&message).Error
	return
}
```

<br>

## Postman 測試

### 查詢全部留言 - 成功(無資料)

![圖片](https://raw.githubusercontent.com/880831ian/go-restful-api-repository-messageboard/master/images/get-success-1.png)

<br>

### 查詢全部留言 - 成功(有資料)

![圖片](https://raw.githubusercontent.com/880831ian/go-restful-api-repository-messageboard/master/images/get-success-2.png)

<br>

### 查詢{id}留言 - 成功

![圖片](https://raw.githubusercontent.com/880831ian/go-restful-api-repository-messageboard/master/images/get-id-succes.png)

<br>

### 查詢{id}留言 - 失敗

![圖片](https://raw.githubusercontent.com/880831ian/go-restful-api-repository-messageboard/master/images/get-error.png)

<br>

### 新增留言 - 成功

![圖片](https://raw.githubusercontent.com/880831ian/go-restful-api-repository-messageboard/master/images/create.png)

<br>

### 修改{id}留言 - 成功

![圖片](https://raw.githubusercontent.com/880831ian/go-restful-api-repository-messageboard/master/images/patch-success.png)

<br>

### 修改{id}留言 - 失敗

![圖片](https://raw.githubusercontent.com/880831ian/go-restful-api-repository-messageboard/master/images/patch-error.png)

<br>

### 刪除{id}留言 - 成功

![圖片](https://raw.githubusercontent.com/880831ian/go-restful-api-repository-messageboard/master/images/delete.png)

<br>

### 執行結果

![圖片](https://raw.githubusercontent.com/880831ian/go-restful-api-repository-messageboard/master/images/gin-cli.png)

<br>

## 參考資料

[基於Gin+Gorm框架搭建MVC模式的Go語言後端系統](https://iter01.com/609571.html)
