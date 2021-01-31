package web

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"GoTuringCoffee/internal/service/lib"
	repo "GoTuringCoffee/internal/service/web/model/repository"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
)

func initRepositoryManager(dbConf repo.MongoDBConfig) (*repo.RepositoryManager, error) {
	ctx := context.TODO()
	repoManager, err := repo.NewRepositoryManager(ctx, &dbConf)
	if err != nil {
		return nil, err
	}

	// Delete all cookbooks before testing
	repoManager.Cookbook.DeleteAll(ctx)

	// Create a cookbook
	newCookbook := lib.Cookbook{
		Name:        "new-cookbook",
		Description: "new cookbook",
		Tags:        []string{},
		Notes:       []string{},
		Processes: []lib.Process{
			lib.Process{
				ID:        "1",
				Name:      "Circle",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Impl:      &lib.Circle{},
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repoManager.Cookbook.Create(ctx, newCookbook)

	return repoManager, nil
}

func initWebService() (*echo.Echo, *Service, error) {
	dbConf := repo.MongoDBConfig{
		URL:      "mongodb+srv://turingcoffee:D56pNXo9bKlosH8W@cluster0.m5idb.gcp.mongodb.net/testturingcoffee?retryWrites=true&w=majority",
		Database: "turingcoffee",
	}
	dbConf.Collections.Cookbook = "cookbooks"

	webConf := WebConfig{
		StaticFilePath: ".",
		Port:           5000,
	}

	s := Service{
		DBConfig:  dbConf,
		WebConfig: webConf,
	}

	ctx := context.TODO()

	repoManager, err := initRepositoryManager(dbConf)
	if err != nil {
		return nil, nil, err
	}

	e := s.InitWebServer(
		ctx,
		repoManager,
		nil,
		nil,
	)

	return e, &s, nil
}

func setup(t *testing.T) *echo.Echo {
	t.Log("Starting web service")
	e, service, err := initWebService()
	assert.NoError(t, err)
	assert.NotNil(t, e)
	assert.NotNil(t, service)

	return e
}

func sendRequest(t *testing.T, e *echo.Echo, url string, method string, params interface{}) *httptest.ResponseRecorder {
	t.Log(fmt.Sprintf("Sending request(%s) to %s", method, url))

	var req *http.Request
	if params != nil {
		jsonParams, err := json.Marshal(params)
		assert.NoError(t, err)
		t.Log(fmt.Sprintf("Params: %s", jsonParams))
		req = httptest.NewRequest(method, url, strings.NewReader(string(jsonParams)))
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	req.Header.Add("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	t.Log("Response:")
	t.Log(rec.Body.String())

	return rec
}

func TestListAllCookbooks(t *testing.T) {
	t.Log("Starting List all cookbooks testcase")
	e := setup(t)

	// Create a fake http request and context
	rec := sendRequest(t, e, "/api/cookbooks", http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	var mapResult map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &mapResult)

	assert.Equal(t, float64(200), mapResult["status"])
}

func generateSampleCookbook() map[string]interface{} {
	return map[string]interface{}{
		"name":        "Cretae New Cookbook - 1",
		"description": "Create new Cookbook - 1 Description",
		"tags":        []string{},
		"notes":       []string{},
		"processes": []map[string]interface{}{
			{
				"id":         "1",
				"name":       "Home",
				"created_at": time.Now().UTC().Unix(),
				"updated_at": time.Now().UTC().Unix(),
				"impl":       map[string]string{},
			},
		},
	}
}

func TestCreateCookbook(t *testing.T) {
	t.Log("Starting create new cookbook")
	e := setup(t)

	// Send request to create cookbook api
	params := generateSampleCookbook()
	rec := sendRequest(t, e, "/api/cookbooks", http.MethodPost, params)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Try to get the new cookbook
	rec = sendRequest(t, e, "/api/cookbooks", http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	var mapResult map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &mapResult)
	assert.Equal(t, float64(200), mapResult["status"])

	payload := mapResult["payload"].([]interface{})
	assert.Equal(t, 2, len(payload))
}

func TestUpdateCookbook(t *testing.T) {
	var mapResult map[string]interface{}

	t.Log("Starting update existing cookbook")
	e := setup(t)

	// Get a cookbook
	rec := sendRequest(t, e, "/api/cookbooks", http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &mapResult)

	payload := mapResult["payload"].([]interface{})
	cookbookJson := payload[0].(map[string]interface{})

	assert.Equal(t, float64(200), mapResult["status"])
	assert.Equal(t, "new-cookbook", cookbookJson["name"])

	// Update the cookbook
	cookbookJson["name"] = "new-cookbook updated"
	rec = sendRequest(t, e, fmt.Sprintf("/api/cookbooks/%s", cookbookJson["id"]), http.MethodPut, cookbookJson)

	// Get the cookbook again
	rec = sendRequest(t, e, fmt.Sprintf("/api/cookbooks/%s", cookbookJson["id"]), http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &mapResult)
	updatedCookbookJson := mapResult["payload"].(map[string]interface{})
	assert.Equal(t, "new-cookbook updated", updatedCookbookJson["name"])
}

func TestDeleteCookbook(t *testing.T) {
	t.Log("Starting delete existing cookbook")
	e := setup(t)

	// Get cookbook
	rec := sendRequest(t, e, "/api/cookbooks", http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	var mapResult map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &mapResult)
	assert.Equal(t, float64(200), mapResult["status"])
	payload := mapResult["payload"].([]interface{})
	cookbookJson := payload[0].(map[string]interface{})
	assert.Equal(t, "new-cookbook", cookbookJson["name"])

	// Delete cookbook
	rec = sendRequest(t, e, fmt.Sprintf("/api/cookbooks/%s", cookbookJson["id"]), http.MethodDelete, nil)

	// Get the cookbook again
	rec = sendRequest(t, e, fmt.Sprintf("/api/cookbooks/%s", cookbookJson["id"]), http.MethodGet, nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetDefaultProcesses(t *testing.T) {
	t.Log("Starting getting all default processes")
	e := setup(t)

	rec := sendRequest(t, e, "/api/default/processes", http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	var mapResult map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &mapResult)
	assert.Equal(t, float64(200), mapResult["status"])

	libDefaultCookbook := lib.GenerateDefaultCookbook()
	defaultProcesses := mapResult["payload"].(map[string]interface{})
	assert.Equal(t, len(libDefaultCookbook.Processes), len(defaultProcesses))
}

func TestGetDefaultProcess(t *testing.T) {
	t.Log("Starting getting default processes by name")
	e := setup(t)

	rec := sendRequest(t, e, "/api/default/processes/Circle", http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	var mapResult map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &mapResult)
	assert.Equal(t, float64(200), mapResult["status"])

	defaultProcesses := mapResult["payload"].(map[string]interface{})
	assert.Equal(t, "Circle", defaultProcesses["name"])

}
