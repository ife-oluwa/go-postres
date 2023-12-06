package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ife-oluwa/go-postres/models"
	"github.com/joho/godotenv"
)

type response struct {
	ID int64 `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}


func createConnection() *sql.DB {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connection established.")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatal("Unable to decode stock body: ", err)
	}

	insertID, err := insertStock(stock)

	res := response{
		ID: insertID,
		Message: "stock created successfully.",
	}

	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string to int: %v", err)
	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Unable to get stock: %v", err)
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err	:= getAllStocks()

	if err != nil {
		log.Fatalf("Unable to get all stock: %v", err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatal("Unable to convert the string to int. Got: ", err)
	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode reqeust body: %v", err)
	}

	updatedRows, err := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stock updated successfully. Total rows/records affected: %v", updatedRows)

	res := response {
		ID: int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert id: %v", err)
	}

	deletedRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock deleted successfully. Total rows/records %v", deletedRows)

	res := response{
		ID: int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func getStock(id int64) (models.Stock, error) {
	db := createConnection()
	defer db.Close()

	sqlStatement := "SELECT * FROM stock WHERE stockid = $1"

	var stock models.Stock

	err := db.QueryRow(sqlStatement, id).Scan(
		&stock.StockID,
		&stock.Name,
		&stock.Price,
		&stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan rows. Error: %v", err)
	}

	return stock, err
}

func getAllStocks() ([]models.Stock, error){
	db := createConnection()

	defer db.Close()

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute query. Error: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatalf("Unable to scan the row %v", err)
		}

		stocks = append(stocks, stock)
	}
	return stocks, err
}

func insertStock(stock models.Stock) (int64, error){
	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING stockid`

	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. Error: %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)
	return id, err
}

func updateStock(id int64, stock models.Stock) (int64, error){
	db := createConnection()

	defer db.Close()

	sqlStatement := `UPDATE stocks SET name = $2, price = $3, company = $4 WHERE stockid = $1`

	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("Unable to execute the query. Error: %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. Error: %v", err)
	}

	fmt.Printf("Total rows/records affected: %v", rowsAffected)

	return rowsAffected, err

}

func deleteStock(id int64) int64{
	db := createConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("unable to delete row. Error: %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/records affected %v", rowsAffected)

	return rowsAffected
}