package main

import (
	"encoding/json" 
	"fmt"
	"io"
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

func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}

func addItem(c echo.Context) error {
	// Get form data
	name := c.FormValue("name")
	category := c.FormValue("category")  // added line
	c.Logger().Infof("Receive item: %s, category: %s", name, category)  // modified line

	item := Item{Name: name, Category: category}
	items.Items = append(items.Items, item)

	// Open the file in append mode, or create it if it doesn't exist
	f, err := os.OpenFile("items.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Make sure to close the file
	defer f.Close()

	// Marshal the items to indented JSON
	data, err := json.MarshalIndent(items, "", "    ")
	if err != nil {
		return err
	}

	// Write to the file
	_, err = f.Write(data)
	if err != nil {
		return err
	}

	// Write a newline (optional, but makes the file easier to read)
	_, err = f.WriteString("\n")
	if err != nil {
		return err
	}

	message := fmt.Sprintf("item received: %s, category: %s", name, category)  // modified line
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

//Define a new structure to store product information
type Item struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

type Items struct {
	Items []Item `json:"items"`
}


func main() {
	e := echo.New()

	// Middleware
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
	e.GET("/items", addItem)
	e.POST("/items", addItem)
	e.GET("/image/:imageFilename", getImg)


	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
