package task_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ArsHighway/Tasks-PSQL/internal/handlers/mocks"
	task "github.com/ArsHighway/Tasks-PSQL/internal/handlers/taskHandler"
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

func TestTaskHandler_CreateTask(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		TaskToReturn: &models.Task{ID: 1, Title: "Task 1"},
		ErrToReturn:  nil,
	}
	h := task.NewTaskHandler(mockServ)
	body := bytes.NewBufferString(`{"title":"Task 1","user_id":1}`)
	rec := httptest.NewRecorder()
	c := testGinContext(rec, http.MethodPost, "/", body)
	c.Request.Header.Set("Content-Type", "application/json")
	h.CreateTask(c)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	var got models.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ID != mockServ.TaskToReturn.ID || got.Title != mockServ.TaskToReturn.Title {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestTaskHandler_GetTaskWithID(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		TaskToReturn: &models.Task{ID: 1, Title: "Task1 "},
		ErrToReturn:  nil,
	}
	h := task.NewTaskHandler(mockServ)
	rec := httptest.NewRecorder()
	c := testGinContext(rec, http.MethodGet, "/tasks/1", nil, gin.Param{Key: "id", Value: "1"})
	h.GetTaskWithID(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got models.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ID != mockServ.TaskToReturn.ID || got.Title != mockServ.TaskToReturn.Title {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestTaskHandler_GetTasks(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		TasksToReturn: []models.Task{{ID: 1, Title: "Task1 "}},
		ErrToReturn:   nil,
	}
	h := task.NewTaskHandler(mockServ)
	rec := httptest.NewRecorder()
	c := testGinContext(rec, http.MethodGet, "/tasks?status=open", nil)
	h.GetTasks(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got []models.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(got, mockServ.TasksToReturn) {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestTaskHandler_UpdateTask(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		TaskToReturn: &models.Task{ID: 1, Title: "Task1 ", Description: "Description1", Status: "open"},
		ErrToReturn:  nil,
	}
	h := task.NewTaskHandler(mockServ)
	body := bytes.NewBufferString(`{"title":"Task1 ","description":"Description1","status":"open","user_id":1}`)
	rec := httptest.NewRecorder()
	c := testGinContext(rec, http.MethodPut, "/tasks/1", body, gin.Param{Key: "id", Value: "1"})
	c.Request.Header.Set("Content-Type", "application/json")
	h.UpdateTask(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got models.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ID != mockServ.TaskToReturn.ID || got.Title != mockServ.TaskToReturn.Title {
		t.Fatalf("unexpected body: %+v", got)
	}
	if got.Description != mockServ.TaskToReturn.Description || got.Status != mockServ.TaskToReturn.Status {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestTaskHandler_PatchTask(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		TaskToReturn: &models.Task{ID: 1, Title: "Task1 ", Description: "Description1", Status: "open"},
		ErrToReturn:  nil,
	}
	h := task.NewTaskHandler(mockServ)
	body := bytes.NewBufferString(`{"title":"Task1 ","description":"Description1","status":"open","user_id":1}`)
	rec := httptest.NewRecorder()
	c := testGinContext(rec, http.MethodPatch, "/tasks/1", body, gin.Param{Key: "id", Value: "1"})
	c.Request.Header.Set("Content-Type", "application/json")
	h.PatchTask(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got models.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ID != mockServ.TaskToReturn.ID || got.Title != mockServ.TaskToReturn.Title {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		ErrToReturn: nil,
	}
	h := task.NewTaskHandler(mockServ)
	rec := httptest.NewRecorder()
	c := testGinContext(rec, http.MethodDelete, "/tasks/1", nil, gin.Param{Key: "id", Value: "1"})
	h.DeleteTask(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got struct {
		Message string `json:"message"`
		TaskID  int    `json:"taskID"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Message != "Task deleted successfully" || got.TaskID != 1 {
		t.Fatalf("unexpected body: %+v", got)
	}
}
