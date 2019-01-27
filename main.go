package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
)

/*
// Create a GORM-backend model
type User struct {
	gorm.Model
	Email     string
	Password  string
	Name      sql.NullString
	Gender    string
	Role      string
	Addresses []Address
}
type Address struct {
	gorm.Model
	UserID   uint
	Address1 string
	Address2 string
}
*/
// Create another GORM-backend model
type Product struct {
	gorm.Model
	Type        int
	Name        string
	Color       int
	Description string
}
type ProductType struct {
	gorm.Model
	ID   int
	Name string
}
type ProductColor struct {
	gorm.Model
	Color string
}

func main() {
	DB, _ := gorm.Open("sqlite3", "demo.db")
	DB.AutoMigrate(&Product{}, &ProductType{}, &ProductColor{})
	DB.LogMode(true)
	// Initalize
	Admin := admin.New(&admin.AdminConfig{DB: DB, SiteName: "Qor Example"})

	// Allow to use Admin to manage User, Product
	product := Admin.AddResource(&Product{})
	Admin.AddResource(&ProductType{})
	Admin.AddResource(&ProductColor{})
	product.Meta(&admin.Meta{
		Name:  "Type",
		Label: "Тип",
		Type:  "select_one",
		Collection: func(value interface{}, context *qor.Context) [][]string {
			var productTypeL []ProductType
			var collectionValues [][]string
			context.GetDB().Find(&productTypeL)
			for _, data := range productTypeL {
				collectionValues = append(collectionValues, []string{fmt.Sprintf("%d", data.ID), data.Name})
			}
			return collectionValues
		},
		Setter: func(record interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			//ylog.Debug("admin", fmt.Sprintf("Setter, record = %+v", record))
			//ylog.Debug("admin", fmt.Sprintf("Setter, resource = %+v", *metaValue))
			//ylog.Debug("admin", fmt.Sprintf("Setter, context = %+v", *context))
			//if newStatusID := utils.ToInt(metaValue.Value); newStatusID != 0 {
			//	record.(*asr.Payment).StatusID = int(newStatusID)
			//}
		},
		Valuer: func(record interface{}, context *qor.Context) (result interface{}) {
			if rec, ok := record.(*Product); ok {
				var productTypeL ProductType
				context.DB.Where(&ProductType{ID: rec.Type}).First(&productTypeL)
				fmt.Println(productTypeL)
				//mapPaymentStatus := asr.BookListMap("paymentStatus")
				return productTypeL.Name //asr.BookListMap("paymentStatus")[rec.StatusID] //mapPaymentStatus[rec.TypeID]
			}
			return ""
		}})

	// initalize an HTTP request multiplexer
	mux := http.NewServeMux()

	// Mount admin interface to mux
	Admin.MountTo("/", mux)

	fmt.Println("Listening on: 9000")
	http.ListenAndServe(":9000", mux)
}
