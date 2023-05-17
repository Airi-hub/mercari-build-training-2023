package main


import (
	"encoding/json" 
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const (
	ImgDir = "images"
)

type Response struct {
	Message string `json:"message"`
}

type Item struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

type Items struct {
	Items []Item `json:"items"`
}

var items Items

func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}


//ファイルから既存のアイテムの読み込み
func init() {
    f, err := os.Open("items.json")
    if err == nil {
        defer f.Close()
        err = json.NewDecoder(f).Decode(&items)
        if err != nil {
            log.Error(err)
        }
    } else if !os.IsNotExist(err) {
        log.Error(err)
    }
}



//リスト取るためのコード
func getItems(c echo.Context) error {
    return c.JSON(http.StatusOK, items)
}


func addItem(c echo.Context) error {
	// データの取得
	name := c.FormValue("name")
	category := c.FormValue("category")  
	c.Logger().Infof("Receive item: %s, category: %s", name, category)  
    
	//c.FormValue("name")とc.FormValue("category")が空文字列でないことを確認
	if name == "" || category == "" {
		return c.JSON(http.StatusBadRequest, &Response{Message: "name and category are required"})
	}

	item := Item{Name: name, Category: category}
	items.Items = append(items.Items, item)

    //ファイルを追記モードで開くか、存在しない場合は作成する
	f, err := os.OpenFile("items.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
	  log.Errorf("Failed to open file: %v", err)
	  return err
	}

	// ファイルを閉じる
	defer f.Close()

	// itemsをインデント付きのJSONに変換
	data, err := json.MarshalIndent(items, "", "    ")
	 if err != nil {
	   log.Errorf("Failed to marshal JSON: %v", err)
	   return err
	}

	// ファイルを書く
	_, err = f.Write(data)
	if err != nil {
	  log.Errorf("Failed to write to file: %v", err)
	  return err
	}

	// 改行を書く
	_, err = f.WriteString("\n")
	if err != nil {
	  log.Errorf("Failed to write newline: %v", err)
	  return err
	}

	message := fmt.Sprintf("item received: %s, category: %s", name, category)
	res := Response{Message: message}

	return c.JSON(http.StatusOK, res)
}



func getImg(c echo.Context) error {
	// Create image path
	imgPath := path.Join(ImgDir, c.Param("imageFilename"))

	if !strings.HasSuffix(imgPath, ".jpg") {
		res := Response{Message: "Image path does not end with .jpg"}
		return c.JSON(http.StatusBadRequest, res)
	}
	if _, err := os.Stat(imgPath); err != nil {
		c.Logger().Debugf("Image not found: %s", imgPath)
		imgPath = path.Join(ImgDir, "default.jpg")
	}
	return c.File(imgPath)
}




func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.INFO)

	front_url := os.Getenv("FRONT_URL")
	if front_url == "" {
		front_url = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{front_url},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	

	// Routes
	e.GET("/", root)
	e.GET("/items", getItems)
	e.POST("/items", addItem)
	e.GET("/image/:imageFilename", getImg)


	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
