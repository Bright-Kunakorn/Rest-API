package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type SKU_branch struct {
	SKUID      *string    `json:"skuid"`
	MerchantID *string    `json:"merchantid"`
	BranchID   *string    `json:"branchid"`
	Price      *float64   `json:"price"`
	StartDate  *time.Time `json:"startdate"`
	EndDate    *time.Time `json:"enddate"`
	IsActive   *int32     `json:"isactive"`
}

func main() {
	connStr := "postgresql://root:secret@localhost:5433?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to PostgreSQL database!")
	r := gin.Default()
	r.GET("/skus_branch", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM backendposdata_sku_branch_price")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var skus_branch []SKU_branch
		for rows.Next() {
			var sku_branch SKU_branch

			err := rows.Scan(&sku_branch.SKUID, &sku_branch.MerchantID, &sku_branch.BranchID, &sku_branch.Price, &sku_branch.StartDate, &sku_branch.EndDate, &sku_branch.IsActive)
			if err != nil {
				log.Fatal(err)
			}
			skus_branch = append(skus_branch, sku_branch)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		file, err := os.Create("skus_branch.json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		jsonData, err := json.Marshal(skus_branch)
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.Write(jsonData)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, skus_branch)
	})
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
