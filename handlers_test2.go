package main

import (
	"device-store/handlers"
	"device-store/models"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	dsn := "host=localhost user=postgres password=2001 dbname=devicestore port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	_ = db.Exec("TRUNCATE TABLE devices, brands, categories RESTART IDENTITY CASCADE")
	db.AutoMigrate(&models.Brand{}, &models.Category{}, &models.Device{})
	os.Exit(m.Run())
}

func createTestDevice(t *testing.T) models.Device {
	brand := models.Brand{Name: "Test Brand"}
	category := models.Category{Name: "Test Category"}
	if err := db.Create(&brand).Error; err != nil {
		t.Fatalf("failed to create brand: %v", err)
	}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("failed to create category: %v", err)
	}
	device := models.Device{
		Name:       "Test Device",
		BrandID:    brand.ID,
		CategoryID: category.ID,
		Price:      100.0,
	}
	if err := db.Create(&device).Error; err != nil {
		t.Fatalf("failed to create device: %v", err)
	}
	return device
}

func setupRouter() *mux.Router {
	h := &handlers.Handler{DB: db}
	r := mux.NewRouter()
	r.HandleFunc("/devices", h.GetDevices).Methods("GET")
	r.HandleFunc("/devices", h.CreateDevice).Methods("POST")
	r.HandleFunc("/devices/{id:[0-9]+}", h.GetDeviceByID).Methods("GET")
	r.HandleFunc("/devices/{id:[0-9]+}", h.UpdateDevice).Methods("PUT")
	r.HandleFunc("/devices/{id:[0-9]+}", h.DeleteDevice).Methods("DELETE")
	return r
}

// ---------- GET /devices ----------
func TestGetDevices(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/devices", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %v", w.Code)
	}
}

func TestGetDevicesWithPagination(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/devices?page=1&limit=2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %v", w.Code)
	}
}

// ---------- GET /devices/{id} ----------
func TestGetDeviceByID(t *testing.T) {
	device := createTestDevice(t)
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/devices/"+strconv.Itoa(int(device.ID)), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %v", w.Code)
	}
}

func TestDeviceNotFound(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("GET", "/devices/999999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %v", w.Code)
	}
}

// ---------- POST /devices ----------
func TestCreateDevice(t *testing.T) {
	brand := models.Brand{Name: "CreateBrand"}
	category := models.Category{Name: "CreateCategory"}
	db.Create(&brand)
	db.Create(&category)

	r := setupRouter()
	body := fmt.Sprintf(`{"name":"Created Device","brand_id":%d,"category_id":%d,"price":250.0}`, brand.ID, category.ID)
	req, _ := http.NewRequest("POST", "/devices", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %v", w.Code)
	}
}

func TestCreateInvalidDevice(t *testing.T) {
	r := setupRouter()
	body := `{"name":"","brand_id":0,"category_id":0,"price":0}`
	req, _ := http.NewRequest("POST", "/devices", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %v", w.Code)
	}
}

func TestCreateDeviceWithMissingFields(t *testing.T) {
	r := setupRouter()
	body := `{"brand_id":1,"category_id":1}`
	req, _ := http.NewRequest("POST", "/devices", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %v", w.Code)
	}
}

// ---------- PUT /devices/{id} ----------
func TestUpdateDevice(t *testing.T) {
	device := createTestDevice(t)
	r := setupRouter()
	body := fmt.Sprintf(`{"name":"Updated Device","brand_id":%d,"category_id":%d,"price":200.0}`, device.BrandID, device.CategoryID)
	req, _ := http.NewRequest("PUT", "/devices/"+strconv.Itoa(int(device.ID)), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %v", w.Code)
	}
}

func TestUpdateInvalidDevice(t *testing.T) {
	device := createTestDevice(t)
	r := setupRouter()
	body := `{"name":"","brand_id":0,"category_id":0,"price":0}`
	req, _ := http.NewRequest("PUT", "/devices/"+strconv.Itoa(int(device.ID)), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %v", w.Code)
	}
}

func TestUpdateNonExistentDevice(t *testing.T) {
	r := setupRouter()
	body := `{"name":"Doesn't exist","brand_id":1,"category_id":1,"price":999}`
	req, _ := http.NewRequest("PUT", "/devices/999999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %v", w.Code)
	}
}

// ---------- DELETE /devices/{id} ----------
func TestDeleteDevice(t *testing.T) {
	device := createTestDevice(t)
	r := setupRouter()
	req, _ := http.NewRequest("DELETE", "/devices/"+strconv.Itoa(int(device.ID)), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected 204, got %v", w.Code)
	}
}

func TestDeleteNonExistentDevice(t *testing.T) {
	r := setupRouter()
	req, _ := http.NewRequest("DELETE", "/devices/999999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %v", w.Code)
	}
}
