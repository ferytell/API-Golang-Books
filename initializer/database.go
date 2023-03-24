package initializer

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

//var DB *gorm.DB

var (
	db  *sql.DB
	err error
)

func ConnectToDB() {
	var (
		host     = os.Getenv("DB_HOST")
		port     = 5432
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	log.Println("connected to database")
}

type Employee struct {
	ID        int
	Full_name string
	Email     string
	Age       int
	Division  string
}

// func CreateEmployee() {
// 	var employee = Employee{}

// 	sqlStatement := `
// 	INSERT INTO employees (full_name, email, age, division)
// 	VALUES ($1, $2, $3, $4)
// 	Returning`

// 	err == DB.QueryRow()
// }
