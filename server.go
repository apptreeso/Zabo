package main

import (
	"fmt"
	"net/http"
	// "encoding/json"
	"database/sql"
	"time"
	
	"github.com/labstack/echo"
	"github.com/go-resty/resty/v2"
	// "github.com/buger/jsonparser"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// type TblItemModel struct {
// 	gorm.Model
// 	Name string `gorm:"size:255"`
// }

// type TblItemModel struct {
// 	gorm.Model
// 	Station_locator_url string `gorm:"size:255"`
// 	Total_results int
// 	Station_counts StationCountsModel `gorm:"foreignkey:Station_counts_id"`
// 	Station_counts_id int
// 	Fuel_stations FuelStationsModel `gorm:"foreignkey:Fuel_stations_id"`
// 	Fuel_stations_id int
// }

type TblItemModel struct {
	gorm.Model
	Station_locator_url string `gorm:"size:255"`
	Total_results int
	Total_model TotalModel `gorm:"foreignkey:Total_model_id"`
	Total_model_id uint
}

type StationCountsModel struct {
	gorm.Model
	Total int
	Fuels []FuelsModel `gorm:"foreignkey:Fuels_id"`
	Fuels_id int
}

type FuelsModel struct {
	gorm.Model
	BD TotalModel `gorm:"foreignkey:BD_id"`
	BD_id int
	E85 TotalModel `gorm:"foreignkey:E85_id"`
	E85_id int
	ELEC ELECModel `gorm:"foreignkey:ELEC_id"`
	ELEC_id int
	HY TotalModel `gorm:"foreignkey:HY_id"`
	HY_id int
	LNG TotalModel `gorm:"foreignkey:LNG_id"`
	LNG_id int
	CNG TotalModel `gorm:"foreignkey:CNG_id"`
	CNG_id int
	LPG TotalModel `gorm:"foreignkey:LPG_id"`
	LPG_id int
}

type ELECModel struct {
	gorm.Model
	Total int
	Stations TotalModel
	Stations_id int `gorm:"foreignkey:Stations_id"`
}

type TotalModel struct {
	gorm.Model
	Total int 
}

type FuelStationsModel struct {
	gorm.Model
	Access_code string `gorm:"size:255"`
	Access_days_time sql.NullString
	Access_detail_code sql.NullString
	Cards_accepted sql.NullString
	Date_last_confirmed time.Time
	Expected_date sql.NullString
	Fuel_type_code string `gorm:"size:255"`
	Groups_with_access_code string `gorm:"size:255"`
	Id int
	Open_date time.Time
	Owner_type_code string `gorm:"size:255"`
	Status_code string `gorm:"size:255"`
	Station_name string `gorm:"size:255"`
	Station_phone sql.NullString
	Updated_at time.Time
	Facility_type string `gorm:"size:255"`
	Geocode_status string `gorm:"size:255"`
	Latitude float64
	Longitude float64
	City string `gorm:"size:255"`
	Intersection_directions sql.NullString
	Plus4 sql.NullString
	State string `gorm:"size:255"`
	Street_address string `gorm:"size:255"`
	Zip string `gorm:"size:255"`
	Country string `gorm:"size:255"`
	Bd_blends sql.NullString
	Cng_dispenser_num sql.NullString
	Cng_fill_type_code string `gorm:"size:255"`
	Cng_psi string `gorm:"size:255"`
	Cng_renewable_source sql.NullString
	Cng_total_compression sql.NullString
	Cng_total_storage sql.NullString
	Cng_vehicle_class string `gorm:"size:255"`
	E85_blender_pump sql.NullString
	E85_other_ethanol_blends sql.NullString
	Ev_connector_types sql.NullString
	Ev_dc_fast_num sql.NullString
	Ev_level1_evse_num sql.NullString
	Ev_level2_evse_num sql.NullString
	Ev_network sql.NullString
	Ev_network_web sql.NullString
	Ev_other_evse sql.NullString
	Ev_pricing sql.NullString
	Ev_renewable_source sql.NullString
	Hy_is_retail sql.NullString
	Hy_pressures sql.NullString
	Hy_standards sql.NullString
	Hy_status_link sql.NullString
	Lng_renewable_source sql.NullString
	Lng_vehicle_class sql.NullString
	Lpg_primary sql.NullString
	Lpg_nozzle_types sql.NullString
	Ng_fill_type_code string `gorm:"size:255"`
	Ng_psi string `gorm:"size:255"`
	Ng_vehicle_class string `gorm:"size:255"`
	Access_days_time_fr sql.NullString
	Intersection_directions_fr sql.NullString
	Bd_blends_fr sql.NullString
	Groups_with_access_code_fr string `gorm:"size:255"`
	Ev_pricing_fr sql.NullString
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
	// Get("https://demo.ckan.org/api/3/action/package_list")
	//https://developer.nrel.gov/api/alt-fuel-stations/v1.json?limit=1&api_key=UFtSv13cjRgZYDpq9tXZet9bz7YGbEhIE0RKkXsp
	resp, err := client.R().
		 Get("https://demo.ckan.org/api/3/action/package_list")
		// Get("https://developer.nrel.gov/api/alt-fuel-stations/v1.json?limit=1&api_key=UFtSv13cjRgZYDpq9tXZet9bz7YGbEhIE0RKkXsp")

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
	db.AutoMigrate(&TblItemModel{}, &TotalModel{})

	bus := TblItemModel{
		Station_locator_url: "Test",
		Total_results: 1010,
		Total_model_id: TotalModel{
			TotalModel{Total: 2222},
		},
	}
	db.Create(&bus)

	// // Parse Data
	// values, _, _, err := jsonparser.Get(resp.Body(), "result")
	// var items []string
	// _ = json.Unmarshal([]byte(values), &items)
	
	// // Populate DB
	// for _, item := range items {
	// 	db.Create(&TblItemModel{Station_locator_url: item})
	// }
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
		dist = append(dist, string(item.Station_locator_url))
	}

	return c.JSONPretty(http.StatusOK, dist, "  ")
}
