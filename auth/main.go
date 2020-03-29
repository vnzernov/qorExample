package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qor/admin"
	"github.com/qor/auth"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth_themes/clean"
	"github.com/qor/qor"
	"github.com/qor/roles"
	"github.com/qor/session/manager"
	"golang.org/x/crypto/bcrypt"
)

// Define another GORM-backend model
type Product struct {
	gorm.Model
	Name        string
	Description string
}

var (
	// Initialize gorm DB
	gormDB, _ = gorm.Open("sqlite3", "sample.db")

	Auth = clean.New(&auth.Config{
		DB: gormDB,
		//Mailer:    config.Mailer,
		UserModel: User{},
	})
	// Authority initialize Authority for Authorization
	//Authority = authority.New(&authority.Config{
	//	Auth: Auth,
	//})
)

type AdminAuth struct{}

func (AdminAuth) LoginURL(c *admin.Context) string {
	return "/auth/login"
}

func (AdminAuth) LogoutURL(c *admin.Context) string {
	return "/auth/logout"
}

func (AdminAuth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	currentUser := Auth.GetCurrentUser(c.Request)
	if currentUser != nil {
		qorCurrentUser, ok := currentUser.(qor.CurrentUser)
		if !ok {
			fmt.Printf("User %#v haven't implement qor.CurrentUser interface\n", currentUser)
		}
		return qorCurrentUser
	}
	return nil
}

func init() {
	roles.Register("admin", func(req *http.Request, currentUser interface{}) bool {
		return currentUser != nil && currentUser.(*User).Role == "Admin"
	})
}

func main() {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	fmt.Printf("%s\n", string(bcryptPassword))
	gormDB.LogMode(true)
	gormDB.AutoMigrate(&User{}, &Product{}, &auth_identity.AuthIdentity{})
	Admin := admin.New(&admin.AdminConfig{DB: gormDB, Auth: &AdminAuth{}})

	Admin.AddResource(&User{})
	Admin.AddResource(&Product{})

	// Initalize an HTTP request multiplexer
	mux := http.NewServeMux()

	// Mount admin to the mux
	Admin.MountTo("/", mux)

	// Mount Auth to Router
	mux.Handle("/auth/", Auth.NewServeMux())

	fmt.Println("Listening on: 7000")
	http.ListenAndServe(":7000", manager.SessionManager.Middleware(mux))
}

type User struct {
	gorm.Model
	Email                  string `form:"email"`
	Password               string
	Name                   string `form:"name"`
	Gender                 string
	Role                   string
	Birthday               *time.Time
	Balance                float32
	DefaultBillingAddress  uint `form:"default-billing-address"`
	DefaultShippingAddress uint `form:"default-shipping-address"`
	//	Addresses              []Address
	//Avatar AvatarImageStorage

	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Recover
	RecoverToken       string
	RecoverTokenExpiry *time.Time

	// Accepts
	AcceptPrivate bool `form:"accept-private"`
	AcceptLicense bool `form:"accept-license"`
	AcceptNews    bool `form:"accept-news"`
}

func (user User) DisplayName() string {
	return user.Email
}
