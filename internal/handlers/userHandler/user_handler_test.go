package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ArsHighway/Tasks-PSQL/internal/handlers/mocks"
	user "github.com/ArsHighway/Tasks-PSQL/internal/handlers/userHandler"
	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func testGinContext(w *httptest.ResponseRecorder, method, path string, body *bytes.Buffer, params ...gin.Param) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	var r *http.Request
	if body != nil {
		r = httptest.NewRequest(method, path, body)
	} else {
		r = httptest.NewRequest(method, path, nil)
	}
	c.Request = r
	if len(params) > 0 {
		c.Params = params
	}
	return c
}

func TestUserHandler_CreateUser(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockUserServ{
		UserToReturn: &models.User{ID: 1, Name: "User 1", Email: "testuser@gmail.com"},
		ErrToReturn:  nil,
	}
	h := user.NewUserHandler(mockServ)
	body := bytes.NewBufferString(`{"Name":"User 1","Email":"testuser@gmail.com"}`)
	rec := httptest.NewRecorder()
	c := testGinContext(rec, http.MethodPost, "/", body)
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateUser(c)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	var got models.User
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	want := mockServ.UserToReturn
	if !reflect.DeepEqual(&got, want) {
		t.Fatalf("unexpected body: got %+v, want %+v", got, want)
	}
}

func TestUserHandler_GetUserWithID(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockUserServ{
		UserToReturn: &models.User{ID: 1, Name: "User 1", Email: "testuser@gmail.com"},
		ErrToReturn:  nil,
	}
	h := user.NewUserHandler(mockServ)
	rec := httptest.NewRecorder()
	c := testGinContext(rec, http.MethodGet, "/users/1", nil, gin.Param{Key: "id", Value: "1"})
	h.GetUserWithID(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got models.User
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	want := mockServ.UserToReturn
	if !reflect.DeepEqual(&got, want) {
		t.Fatalf("unexpected body: got %+v, want %+v", got, want)
	}
}

func TestUserHandler_PatchUser(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockUserServ{
		UserToReturn: &models.User{ID: 1, Name: "User 1", Email: "testuser@gmail.com"},
		ErrToReturn:  nil,
	}
	h := user.NewUserHandler(mockServ)
	body := bytes.NewBufferString(`{"Name":"User 1","Email":"testuser@gmail.com"}`)
	rec := httptest.NewRecorder()
	c := testGinContext(rec, http.MethodPatch, "/users/1", body, gin.Param{Key: "id", Value: "1"})
	c.Request.Header.Set("Content-Type", "application/json")
	h.PatchUser(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got models.User
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	want := mockServ.UserToReturn
	if !reflect.DeepEqual(&got, want) {
		t.Fatalf("unexpected body: got %+v, want %+v", got, want)
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockUserServ{
		UserToReturn: &models.User{ID: 1, Name: "User 1", Email: "testuser@gmail.com"},
		ErrToReturn:  nil,
	}
	h := user.NewUserHandler(mockServ)
	rec := httptest.NewRecorder()
	c := testGinContext(rec, http.MethodDelete, "/users/1", nil, gin.Param{Key: "id", Value: "1"})
	h.DeleteUser(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
