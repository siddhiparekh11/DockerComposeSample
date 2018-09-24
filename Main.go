package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)


// This struct is for application level object - database and mux router
type App struct {
	Router *mux.Router
	DB *sql.DB
}

//Book struct model
type Book struct {

	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	IDAuthor string `json:"idauthor"`
}

//Author struct model
type Author struct {

	ID string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}


var a App


func initialize(username string, password string,dbname string){

	connectionString := fmt.Sprintf("%s:%s@/%s", username, password, dbname)
	db,err := sql.Open("mysql",connectionString)
	if err!=nil {
		fmt.Print(err.Error())
	}
	a.DB = db
}



//Get all books
func getBooks(w http.ResponseWriter, r *http.Request){


	var books []Book
	var book Book
	w.Header().Set("Content-Type","application/json")
	rows, err := a.DB.Query("select idBook,bookTitle, bookISBN, idAuthor from Books;")
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&book.ID,&book.Title, &book.Isbn, &book.IDAuthor)
		books = append(books, book)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	defer rows.Close()
	json.NewEncoder(w).Encode(books)

}

//Get single book acc to title/id
func getBook(w http.ResponseWriter, r *http.Request){

	var books []Book
	var book Book
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	query := "select idBook,bookTitle, bookISBN, idAuthor from Books where idBook=" + params["id"] + ";"
	rows, err := a.DB.Query(query)
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&book.ID,&book.Title, &book.Isbn, &book.IDAuthor)
		books = append(books, book)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	defer rows.Close()
	json.NewEncoder(w).Encode(books)

	}



//create a new book
func createBook(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type","application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	insert, err := a.DB.Prepare("insert into Books(bookTitle,bookISBN,idAuthor) values(?,?,?);")
	if err != nil {
		fmt.Print(err.Error())
	}
	id, err := strconv.Atoi(book.IDAuthor)
	if err != nil {
		fmt.Println(err)
	}
	_, err = insert.Exec(book.Title,book.Isbn,id)
	if err != nil {
		fmt.Print(err.Error())
	}
	json.NewEncoder(w).Encode(book)

}


//create a new Author

func createAuthor(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type","application/json")
	var author Author
	_ = json.NewDecoder(r.Body).Decode(&author)
	insert, err := a.DB.Prepare("insert into AUTHORS(firstName,lastName) values(?,?);")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = insert.Exec(author.Firstname,author.Lastname)
	if err != nil {
		fmt.Print(err.Error())
	}
	json.NewEncoder(w).Encode(author)

}

// get Authors

func getAuthors(w http.ResponseWriter, r *http.Request){


	var authors []Author
	var author Author
	w.Header().Set("Content-Type","application/json")
	rows, err := a.DB.Query("select idAuthor,firstName, lastName from Authors;")
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&author.ID,&author.Firstname, &author.Lastname)
		authors = append(authors, author)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	defer rows.Close()
	json.NewEncoder(w).Encode(authors)

}


//update a book acc to id
func updateBook(w http.ResponseWriter, r *http.Request){

	var books []Book
	var book Book
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	_ = json.NewDecoder(r.Body).Decode(&book)
	query := "update Books set bookTitle='" + book.Title + "', bookISBN='" + book.Isbn + "', idAuthor='" + book.IDAuthor + "', idBook='" + params["id"] + "' where idBook=" + params["id"] + ";"
	fmt.Println(query)
	update, err:= a.DB.Query(query)
	if err != nil {
		fmt.Print(err.Error())
	}
	update.Close()
	query1 := "select bookTitle, bookISBN, idAuthor from Books where idBook=" + params["id"]
	rows, err := a.DB.Query(query1)
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&book.Title, &book.Isbn, &book.IDAuthor)
		books = append(books, book)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	defer rows.Close()

	json.NewEncoder(w).Encode(books)

}

//delete a book acc to authorId
func deleteBook(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	fmt.Println(params["id"])
	query := "delete from Books where idAuthor=" + params["id"]
	delete, err := a.DB.Query(query)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer delete.Close()

}




func main() {

	//initializing router and database
	a := App{}
	a.Router = mux.NewRouter()
	initialize("root","12356","library")
	fmt.Printf("%T",a)

	//endpoints and functions
	a.Router.HandleFunc("/api/books",getBooks).Methods("GET")
	a.Router.HandleFunc("/api/books/{id}",getBook).Methods("GET")
	a.Router.HandleFunc("/api/book",createBook).Methods("POST")
	a.Router.HandleFunc("/api/author",createAuthor).Methods("POST")
	a.Router.HandleFunc("/api/books/{id}",updateBook).Methods("PUT")
	a.Router.HandleFunc("/api/books/{id}",deleteBook).Methods("DELETE")
	a.Router.HandleFunc("/api/authors",getAuthors).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000",a.Router))
	
}
