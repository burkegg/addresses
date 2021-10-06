package addresses

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"

	//"github.com/gocarina/gocsv"
	"io/fs"
	"log"
	"net/http"
)

//go:embed assets
var content embed.FS

type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
	DBConn   *gorm.DB
	Version  string
}

type House struct {
	gorm.Model
	Address string `csv:"ADDRESS"`
	PropType string `csv:"PROPERTY TYPE"`
	Price int `csv:"PRICE"`
	SqFeet int `csv:"SQUARE FEET"`
	City string `csv:"LOCATION"`
	URL string `csv:"URL (SEE http://www.redfin.com/buy-a-home/comparative-market-analysis FOR INFO ON PRICING)"`
}

func BuildDBConfig(host string, port int, user string, dbName string, password string, version string) (dbConfig *DBConfig, err error) {
	dbConfig = &DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		DBName:   dbName,
		Password: password,
		DBConn:    nil,
		Version:  version,
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
	DBConn, err := gorm.Open( "postgres", dsn)
	if err != nil {
		log.Printf("Failed to open connection to database using: %s\n", dsn)
		return dbConfig, err
	}
	dbConfig.DBConn = DBConn
	err = dbConfig.DBConn.AutoMigrate(&House{}).Error
	if err != nil {
		log.Println("Failed to migrate House type to database table")
		return dbConfig, err
	}
	return dbConfig, err
}

func Serve(urlPrefix string, efs embed.FS) gin.HandlerFunc {
	// the embedded filesystem has a 'views/' at the top level.  We wanna strip this so we can treat the root of the views directory as the web root.
	fmt.Printf("urlPrefix: %s\n", urlPrefix)
	fsys, err := fs.Sub(efs, "assets")
	if err != nil {
		log.Fatalf(err.Error())
	}

	fileserver := http.FileServer(http.FS(fsys))
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}

	return func(c *gin.Context) {
		fileserver.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func (db *DBConfig) FetchAddresses(c *gin.Context) {
	var Addresses []House
	db.DBConn.Find(&Addresses)
	c.JSON(http.StatusOK, Addresses)
}

func (db *DBConfig) SetUpRouter(address string, port int) (router *gin.Engine) {
	router = gin.Default()

	api := router.Group("/api")

	// Routes
	router.Use(Serve("/", content))
	api.GET("/addresses", db.FetchAddresses)
	addr := fmt.Sprintf("%s:%d", address, port)
	fmt.Printf("Server starting on %s\n", addr)
	return router
}

// For now take in the csv locally and read it into psql
func (db *DBConfig) ImportData() {
	in, err := os.Open("/go/src/addresses-challenge/pkg/addresses/assets/addressdata.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()
	var data []*House
	if err := gocsv.UnmarshalFile(in, &data); err != nil {
		panic(err)
	}
	for _, home := range data {
		fmt.Println("Address, ", home.Address)
	}
}

func RunServer(dbHost string, dbPort int, dbUser string, dbPassword string, dbName string, address string, port int, version string) (err error) {
	// Set up DB
	db, err := BuildDBConfig(dbHost, dbPort, dbUser, dbName, dbPassword, version)
	if err != nil {
		log.Println("Database init failed.")
		log.Printf("dbHost: %v, dbport: %v, dbUser %v, dbName %v, dbPassword: %v", dbHost, dbPort, dbUser, dbName, dbPassword )
		return err
	}

	db.ImportData()
	r := db.SetUpRouter(address, port)
	addr := fmt.Sprintf("%s:%d", address, port)
	err = r.Run(addr)
	return err
}
