package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbDriver  = "mysql"
	dbUser    = os.Getenv("DB_USER")
	dbPass    = os.Getenv("DB_PASS")
	dbName    = os.Getenv("DB_NAME")
	dbHost    = os.Getenv("DB_HOST")
	dbPort    = os.Getenv("DB_PORT")
	query     = os.Getenv("QUERY")
	dbMaxCon  = os.Getenv("DB_MAXCON")
	dbMaxIdle = os.Getenv("DB_MAXIDLE")
)

func main() {
	// Set up the MySQL database connection
	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	maxCon, err := strconv.Atoi(dbMaxCon)
	if err != nil {
		maxCon = 10
	}

	maxIdle, err := strconv.Atoi(dbMaxIdle)
	if err != nil {
		maxIdle = 5
	}
	db.SetMaxOpenConns(maxCon) // Adjust based on your application's needs

	// Set the maximum number of idle connections
	db.SetMaxIdleConns(maxIdle) // Adjust based on your application's needs

	// Initialize the Gin router
	router := gin.Default()

	// Define a simple route that queries the database
	router.Any("/", func(c *gin.Context) {
		// Query the database (assuming you have a 'users' table)

		conn, err := db.Conn(c)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error getting a database connection"})
			return
		}
		defer conn.Close()
		rows, err := conn.QueryContext(c, query)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error querying the database"})
			return
		}
		defer rows.Close()
		// Get column names dynamically
		columns, err := rows.Columns()
		if err != nil {
			c.JSON(500, gin.H{"error": "Error retrieving column names"})
			return
		}

		// Create a slice to hold the column values dynamically
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		// Process the query results and create a response
		var data []map[string]interface{}
		for rows.Next() {
			err := rows.Scan(values...)
			if err != nil {
				c.JSON(500, gin.H{"error": "Error scanning rows"})
				return
			}

			// Create a map to store column names and values
			rowData := make(map[string]interface{})
			for i, column := range columns {
				val := *(values[i].(*interface{}))
				rowData[column] = val
			}

			data = append(data, rowData)
		}

		// Return the response
		c.JSON(200, data)
	})

	// Run the Gin server
	err = router.Run(":5000")
	if err != nil {
		log.Fatal("Error starting Gin server:", err)
	}
}
