package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Message string `json:"message"`
}

type ShowResponse struct {
	Data []showBookData `json:"data"`
}

type deleteBookData struct {
	ID int `json:"id"`
}

type showBookData struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  int    `json:"price"`
}

type UpdateBookData struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  int    `json:"price"`
}

type BookData struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  int    `json:"price"`
}

type Storage struct {
	Store map[int]BookData
}

func (s *Storage) AddNewBook(data BookData) {
	NewID := len(s.Store) + 1
	s.Store[NewID] = data
}

func (s *Storage) getUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the HTML page
	http.ServeFile(w, r, "static/update.html")
}

func (s *Storage) postUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request into the BookData struct
	var bookData UpdateBookData
	err := json.NewDecoder(r.Body).Decode(&bookData)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	newData := BookData{
		Title:  bookData.Title,
		Author: bookData.Author,
		Price:  bookData.Price,
	}
	s.Store[bookData.ID] = newData

	// Log the received data
	fmt.Printf("Updated Book Data: %+v\n", bookData)

	// Construct a response message
	response := Response{
		Message: fmt.Sprintf("Book '%s' by %s priced at %d added successfully!", bookData.Title, bookData.Author, bookData.Price),
	}

	// Set response headers and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Storage) getCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the HTML page
	http.ServeFile(w, r, "static/index.html")
}

func (s *Storage) postCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request into the BookData struct
	var bookData BookData
	err := json.NewDecoder(r.Body).Decode(&bookData)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	s.AddNewBook(bookData)

	// Log the received data
	fmt.Printf("Received Book Data: %+v\n", bookData)

	// Construct a response message
	response := Response{
		Message: fmt.Sprintf("Book '%s' by %s priced at %d added successfully!", bookData.Title, bookData.Author, bookData.Price),
	}

	// Set response headers and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Storage) getDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the HTML page
	http.ServeFile(w, r, "static/delete.html")
}

func (s *Storage) postDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON request into the BookData struct
	var bookData deleteBookData
	err := json.NewDecoder(r.Body).Decode(&bookData)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// s.AddNewBook(bookData)
	delete(s.Store, bookData.ID)

	// Log the received data
	fmt.Printf("Deleted Book Data: %+v\n", bookData)

	// Construct a response message
	response := Response{
		Message: fmt.Sprint("Book is successfully Deleted"),
	}

	// Set response headers and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (s *Storage) getShowHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the HTML page
	http.ServeFile(w, r, "static/show.html")
}

func (s *Storage) showStorage(w http.ResponseWriter, r *http.Request) {
	var response []showBookData

	for key, value := range s.Store {
		temp := showBookData{
			ID:     key,
			Title:  value.Title,
			Author: value.Author,
			Price:  value.Price,
		}
		response = append(response, temp)
	}

	data := ShowResponse{
		Data: response,
	}
	// Set response headers and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func CreateNewStorage() *Storage {
	return &Storage{
		Store: make(map[int]BookData),
	}
}

func main() {

	r := mux.NewRouter()
	storage := CreateNewStorage()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/create", storage.getCreateHandler).Methods(http.MethodGet)
	r.HandleFunc("/create", storage.postCreateHandler).Methods(http.MethodPost)
	r.HandleFunc("/show", storage.getShowHandler).Methods(http.MethodGet)
	r.HandleFunc("/getShowData", storage.showStorage).Methods(http.MethodGet)
	r.HandleFunc("/update", storage.getUpdateHandler).Methods(http.MethodGet)
	r.HandleFunc("/update", storage.postUpdateHandler).Methods(http.MethodPost)
	r.HandleFunc("/delete", storage.getDeleteHandler).Methods(http.MethodGet)
	r.HandleFunc("/delete", storage.postDeleteHandler).Methods(http.MethodPost)
	// Start the server
	fmt.Println("Server is running on http://localhost:3000")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
