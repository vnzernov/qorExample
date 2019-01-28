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

// Create a GORM-backend model
type User struct {
	//gorm.Model
	ID       int
	Email    string
	Password string
	Name     string
	//Gender    string
	//Role      string
	//Addresses []Address
}

/*
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
	//ID          uint
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

var DB *gorm.DB

func main() {
	DB, _ = gorm.Open("sqlite3", "demo.db")
	DB.AutoMigrate(&Product{}, &ProductType{}, &ProductColor{}, User{})
	DB.LogMode(true)
	// Initalize
	Admin := admin.New(&admin.AdminConfig{DB: DB, SiteName: "Qor Example"})

	// Allow to use Admin to manage User, Product
	Admin.AddResource(&User{})
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
			//if newStatusID := utils.ToInt(metaValue.Value); newStatusID != 0 {
			//	record.(*asr.Payment).StatusID = int(newStatusID)
			//}
		},
		Valuer: func(record interface{}, context *qor.Context) (result interface{}) {
			if rec, ok := record.(*Product); ok {
				var productTypeL ProductType
				context.DB.Where(&ProductType{ID: rec.Type}).First(&productTypeL)
				fmt.Println(productTypeL)
				return productTypeL.Name
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

func (User) ConfigureQorResourceBeforeInitialize(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		//fmt.Printf("resAfter = %+v\n", *res)
		//res.Meta(&admin.Meta{
		//	Name: "ScheduledStartAt",
		//	Valuer: func(interface{}, *qor.Context) interface{} {
		//		return ""
		//	},
		//})
		// do something before initialize
		//res.
		res.FindManyHandler = func(result interface{}, context *qor.Context) error {

			//if res.HasPermission(roles.Read, context) {
			fmt.Printf("result = %+v\n", result)
			fmt.Printf("context = %+v\n", context)
			db := context.GetDB()
			data, ok := db.Get("qor:getting_total_count")
			fmt.Printf("data = %+v\n ok = %+v\n", data, ok)
			if ok {
				//context.GetDB().Count(result)
				fmt.Printf("resFindManyHandler results1 = %+v\n Type %T\n", result.(*int), result)
				return nil
			}
			//context.GetDB().Set("gorm:order_by_primary_key", "DESC").Find(result)
			//userL := make([]*User, 2)
			DB.Set("gorm:order_by_primary_key", "DESC").Find(result)
			//fmt.Printf("resFindManyHandler results = %+v\n Type %T", result, result)
			/*
				result = 2
				return nil
				userL := make([]*User, 2)
				userL[0] = &User{ID: 1, Name: "user1", Password: "qwerty", Email: "aaa@aaa"}
				userL[1] = &User{ID: 2, Name: "user2", Password: "qwerty", Email: "aaa@aaa"}
				result = &userL
			*/
			for _, data := range *result.(*[]*User) {
				//fmt.Printf("resFindManyHandler results = %+v\n Type %T\n", *result.(*[]*User), result)
				fmt.Printf("resFindManyHandler results = %+v\n Type %T\n", data, data)
			}
			return nil
			//}
			//return roles.ErrPermissionDenied
			//return nil
		}
	}
}

/*
func (User) ConfigureQorResource(res resource.Resourcer) {
	fmt.Printf("res = %+v\n", res)
	// do something after initialize
}
*/
