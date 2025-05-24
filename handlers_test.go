package main

import (
	"encoding/json"
	"fmt"
	deviceHandlers "device-store/device-service/handlers"
	deviceModels "device-store/device-service/models"
	userHandlers "device-store/user-service/handlers"
	userModels "device-store/user-service/models"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupDeviceDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=2001 dbname=devices port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("failed to connect to Postgres device DB:", err)
	}
	if err := db.AutoMigrate(&deviceModels.Device{}, &deviceModels.Brand{}, &deviceModels.Category{}); err != nil {
		t.Fatal(err)
	}
	// Optional: clean DB for tests
	db.Exec("DELETE FROM devices")
	db.Exec("DELETE FROM brands")
	db.Exec("DELETE FROM categories")
	return db
}

func setupUserDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=2001 dbname=devices port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("failed to connect to Postgres user DB:", err)
	}
	if err := db.AutoMigrate(&userModels.User{}); err != nil {
		t.Fatal(err)
	}
	db.Exec("DELETE FROM users")
	return db
}

// -------- User tests ---------

func TestRegister_Success(t *testing.T) {
	db := setupUserDB(t)
	h := &userHandlers.Handler{DB: db}

	payload := `{"username":"testuser","password":"testpass"}`
	req := httptest.NewRequest("POST", "/register", strings.NewReader(payload))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var user userModels.User
	if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
		t.Fatal(err)
	}
	if user.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %s", user.Username)
	}
}

// ... other user tests unchanged except using userHandlers.Handler, userModels.User

// -------- Device tests ---------

func TestGetDevices_Empty(t *testing.T) {
	db := setupDeviceDB(t)
	h := &deviceHandlers.Handler{DB: db}

	req := httptest.NewRequest("GET", "/devices", nil)
	w := httptest.NewRecorder()

	h.GetDevices(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var devices []deviceModels.Device
	if err := json.NewDecoder(w.Body).Decode(&devices); err != nil {
		t.Fatal(err)
	}

	if len(devices) != 0 {
		t.Errorf("expected empty devices list, got %d", len(devices))
	}
}

func TestCreateDevice_Success(t *testing.T) {
	db := setupDeviceDB(t)
	h := &deviceHandlers.Handler{DB: db}

	payload := `{"name":"iPhone","price":999}`
	req := httptest.NewRequest("POST", "/devices", strings.NewReader(payload))
	w := httptest.NewRecorder()

	h.CreateDevice(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var device deviceModels.Device
	if err := json.NewDecoder(w.Body).Decode(&device); err != nil {
		t.Fatal(err)
	}

	if device.Name != "iPhone" || device.Price != 999 {
		t.Errorf("device data mismatch, got %+v", device)
	}
}

func TestGetDeviceByID_NotFound(t *testing.T) {
	db := setupDeviceDB(t)
	h := &deviceHandlers.Handler{DB: db}

	req := httptest.NewRequest("GET", "/devices/123", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "123"})
	w := httptest.NewRecorder()

	h.GetDeviceByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateDevice_Success(t *testing.T) {
	db := setupDeviceDB(t)
	h := &deviceHandlers.Handler{DB: db}

	device := deviceModels.Device{Name: "OldName", Price: 100}
	db.Create(&device)

	payload := `{"name":"NewName","price":200}`
	req := httptest.NewRequest("PUT", "/devices/"+fmt.Sprintf("%d", device.ID), strings.NewReader(payload))
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", device.ID)})
	w := httptest.NewRecorder()

	h.UpdateDevice(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var updatedDevice deviceModels.Device
	if err := json.NewDecoder(w.Body).Decode(&updatedDevice); err != nil {
		t.Fatal(err)
	}

	if updatedDevice.Name != "NewName" || updatedDevice.Price != 200 {
		t.Errorf("device not updated correctly, got %+v", updatedDevice)
	}
}

func TestDeleteDevice_Success(t *testing.T) {
	db := setupDeviceDB(t)
	h := &deviceHandlers.Handler{DB: db}

	device := deviceModels.Device{Name: "ToDelete", Price: 10}
	db.Create(&device)

	req := httptest.NewRequest("DELETE", "/devices/"+fmt.Sprintf("%d", device.ID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", device.ID)})
	w := httptest.NewRecorder()

	h.DeleteDevice(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	var check deviceModels.Device
	if err := db.First(&check, device.ID).Error; err == nil {
		t.Errorf("device was not deleted")
	}
}
