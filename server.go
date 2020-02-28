package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	
	"github.com/labstack/echo"
	"github.com/go-resty/resty/v2"
	"github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type TblItemModel struct {
	gorm.Model
	Name string `gorm:"size:255"`
}

var db *gorm.DB

func main() {
	populateDB()

	e := echo.New()
	e.GET("/", getData)
	e.Logger.Fatal(e.Start(":8080"))
}

func populateDB() {
	fmt.Println("Connecting DB...")

	// Create a Resty Client
	client := resty.New()

	// Pull
	resp, err := client.R().
		Get("https://demo.ckan.org/api/3/action/package_list")

	if resp.StatusCode() != http.StatusOK {
		fmt.Println(err)
		return
	}

	// Connect DB
	db, err := gorm.Open("sqlite3", "./sqlite_dummy.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	fmt.Println("Connection Established")

	// Migrate DB
	db.AutoMigrate(&TblItemModel{})

	// Parse Data
	values, _, _, err := jsonparser.Get(resp.Body(), "result")
	var items []string
	_ = json.Unmarshal([]byte(values), &items)
	
	// Populate DB
	for _, item := range items {
		db.Create(&TblItemModel{Name: item})
	}
}

func getData(c echo.Context) error {
	db, err := gorm.Open("sqlite3", "sqlite_dummy.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	
	var items []TblItemModel
	if err := db.Find(&items).Error; err != nil {
		panic(err)
	}
	
	var dist []string
	for _, item := range items {
		dist = append(dist, string(item.Name))
	}

	return c.JSONPretty(http.StatusCreated, dist, "  ")
}
