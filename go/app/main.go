package main


import (
	"encoding/json" 
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
    "crypto/sha256"
    "io"
    "encoding/hex"
	"strconv"

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
    ID            int    `json:"id"`
    Name          string `json:"name"`
    Category      string `json:"category"`
    ImageFilename string `json:"image_filename"`
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

//最後に使用されたIDを追跡するためのグローバル変数
var lastID int = 0


func addItem(c echo.Context) error {
	// データの取得
	name := c.FormValue("name")
	category := c.FormValue("category")  
	c.Logger().Infof("Receive item: %s, category: %s", name, category)  
    
	//c.FormValue("name")とc.FormValue("category")が空文字列でないことを確認
	if name == "" || category == "" {
		return c.JSON(http.StatusBadRequest, &Response{Message: "name and category are required"})
	}

	//画像が存在しない場合のファイル作成処理
	if _, err := os.Stat("./images"); os.IsNotExist(err) {
		os.Mkdir("./images", 0755)
	}

  //画像が存在する場合のファイル作成処理
	
	// ファイルアップロードの処理
    file, err := c.FormFile("image")
    if err != nil {
        return err
    }
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    // ファイルのSHA256ハッシュを計算
    h := sha256.New()
    if _, err := io.Copy(h, src); err != nil {
        log.Fatal(err)
    }
    hash := hex.EncodeToString(h.Sum(nil))

    // ハッシュをファイル名としてファイルを保存
    dst, err := os.Create(fmt.Sprintf("images/%s.jpg", hash))
    if err != nil {
        return err
    }
    defer dst.Close()

    src.Seek(0, 0) // ファイルの読み取り位置をリセット

    if _, err = io.Copy(dst, src); err != nil {
        return err
    }

	lastID++
	item := Item{ID: lastID, Name: name, Category: category, ImageFilename: fmt.Sprintf("%s.jpg", hash)}
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


//特定の商品の詳細情報を取得するための新しいエンドポイントを作成
func getItem(c echo.Context) error {
    // item_idをパスパラメータから取得
    id, err := strconv.Atoi(c.Param("item_id"))
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
    }

    // item_idに対応する商品を検索
    for _, item := range items.Items {
        if item.ID == id {
            // 商品が見つかったらJSONとして返す
            return c.JSON(http.StatusOK, item)
        }
    }

    // 商品が見つからなかった場合は404エラーを返す
    return echo.NewHTTPError(http.StatusNotFound, "Item not found")
}

func getImg(c echo.Context) error {
	c.Logger().Debugf("Current log level: %v", c.Echo().Logger.Level())

	//イメージパスの作成
	imgPath := path.Join(ImgDir, c.Param("imageFilename"))

	if !strings.HasSuffix(imgPath, ".jpg") {
		res := Response{Message: "Image path does not end with .jpg"}
		return c.JSON(http.StatusBadRequest, res)
	}

	// ファイルが存在するかどうかを確認する
	_, err := os.Stat(imgPath)
	if os.IsNotExist(err) {
		c.Logger().Debugf("Image not found: %s", imgPath)
		imgPath = path.Join(ImgDir, "default.jpg")
	}

	//  default.jpgも表示
	if imgPath == path.Join(ImgDir, "default.jpg") {
		c.Logger().Debugf("Default image displayed: %s", imgPath)
	}

	return c.File(imgPath)
}




//位置変更

//リスト取るためのコード
func getItems(c echo.Context) error {
    return c.JSON(http.StatusOK, items)
}


func handler(c echo.Context) error {
    c.Logger().Info("This is an info message")
    c.Logger().Error("This is an error message")
    return nil
}


func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(log.DEBUG)

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
	e.GET("/items/:item_id", getItem) //エンドポイントルート追加
	e.GET("/image/:imageFilename", getImg)
	e.POST("/items", addItem)  




	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
