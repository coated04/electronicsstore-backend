package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "bookstore/models"
)

var books = []models.Book{}
var nextBookID = 1

func GetBooks(w http.ResponseWriter, r *http.Request) {
    categoryID := r.URL.Query().Get("category")
    pageStr := r.URL.Query().Get("page")
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }
    perPage := 5


    filteredBooks := []models.Book{}
    for _, book := range books {
        if categoryID == "" || strconv.Itoa(book.CategoryID) == categoryID {
            filteredBooks = append(filteredBooks, book)
        }
    }

    start := (page - 1) * perPage
    end := start + perPage

    if start >= len(filteredBooks) {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode([]models.Book{})
        return
    }
    if end > len(filteredBooks) {
        end = len(filteredBooks)
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(filteredBooks[start:end])
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
    var book models.Book
    if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
        http.Error(w, "Invalid JSON input", http.StatusBadRequest)
        return
    }

    if book.Title == "" || book.AuthorID == 0 || book.CategoryID == 0 || book.Price <= 0 {
        http.Error(w, "missing required fields or invalid values", http.StatusBadRequest)
        return
    }

    book.ID = nextBookID
    nextBookID++
    books = append(books, book)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(book)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "book ID is invalid", http.StatusBadRequest)
        return
    }

    for _, book := range books {
        if book.ID == id {
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(book)
            return
        }
    }
    http.Error(w, "book not found", http.StatusNotFound)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid book ID", http.StatusBadRequest)
        return
    }

    for i, book := range books {
        if book.ID == id {
            books = append(books[:i], books[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    http.Error(w, "Book not found", http.StatusNotFound)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid book ID", http.StatusBadRequest)
        return
    }

    var updatedBook models.Book
    if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
        http.Error(w, "Invalid JSON input", http.StatusBadRequest)
        return
    }


    if updatedBook.Title == "" || updatedBook.AuthorID == 0 || updatedBook.CategoryID == 0 || updatedBook.Price <= 0 {
        http.Error(w, "Missing required fields or invalid values", http.StatusBadRequest)
        return
    }

    for i, book := range books {
        if book.ID == id {
            updatedBook.ID = id
            books[i] = updatedBook
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(updatedBook)
            return
        }
    }
    http.Error(w, "Book not found", http.StatusNotFound)
}
