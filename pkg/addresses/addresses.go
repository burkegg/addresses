package addresses

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"io/fs"
	"log"
	"net/http"
	"os"
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

type Search struct {
	Term string`json:"Term"`
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
	return dbConfig, err
}

func (db *DBConfig) ConnectConfigToDB() (err error){
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, db.User, db.Password, db.DBName)
	DBConn, err := gorm.Open( "postgres", dsn)
	if err != nil {
		log.Printf("Failed to open connection to database using: %s\n", dsn)
		return err
	}
	db.DBConn = DBConn
	err = db.DBConn.AutoMigrate(&House{}).Error
	if err != nil {
		log.Println("Failed to migrate House type to database table")
		return err
	}
	return nil
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

// Take a search condition and do the PSQL query based on that.
func (db *DBConfig) FetchAddresses(c *gin.Context) {
	var Addresses []House
	fmt.Printf("Req body: %+v\n", c.Request.Body)
	var s Search
	err := c.BindJSON(&s)
	if err != nil{
		log.Fatalf(err.Error())
	}

	db.DBConn.Where("Address ILIKE ?", fmt.Sprintf("%%%s%%", s.Term)).Find(&Addresses)

	c.JSON(http.StatusOK, Addresses)
}

// InsertHouse inserts one house into the db.  Could speed things up w/ batch insert if we cared.
func (db *DBConfig) InsertHouse(h *House) (err error){
	err = db.DBConn.Create(h).Error
	return err
}

// SetUpRouter gets the engine ready to serve static files and to handle routes
func (db *DBConfig) SetUpRouter(address string, port int) (router *gin.Engine) {
	router = gin.Default()

	api := router.Group("/api")

	// Routes
	router.Use(Serve("/", content))
	api.POST("/addresses", db.FetchAddresses)
	addr := fmt.Sprintf("%s:%d", address, port)
	fmt.Printf("Server starting on %s\n", addr)
	return router
}

// ImportData takes in the csv locally and reads it into psql
func (db *DBConfig) ImportData(fileLoc string) (data []*House){
	in, err := os.Open(fileLoc)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer in.Close()
	if err := gocsv.UnmarshalFile(in, &data); err != nil {
		log.Fatalf(err.Error())
	}
	return data
}

// InsertData - this should have been a batch insert, but here we are.q
func (db *DBConfig) InsertData(data []*House) {
	for _, home := range data {
		err := db.InsertHouse(home)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func RunServer(dbHost string, dbPort int, dbUser string, dbPassword string, dbName string, address string, port int, version string) (err error) {
	// Set up DB
	db, err := BuildDBConfig(dbHost, dbPort, dbUser, dbName, dbPassword, version)
	if err != nil {
		log.Println("Database config setup failed.")
		return err
	}
	err = db.ConnectConfigToDB()
	if err != nil {
		log.Println("Database init failed.")
		log.Printf("dbHost: %v, dbport: %v, dbUser %v, dbName %v, dbPassword: %v", dbHost, dbPort, dbUser, dbName, dbPassword )
		return err
	}

	data := db.ImportData("/go/src/addresses-challenge/pkg/addresses/assets/addressdata.csv")
	db.InsertData(data)
	r := db.SetUpRouter(address, port)
	addr := fmt.Sprintf("%s:%d", address, port)
	err = r.Run(addr)

	return err
}
