package server_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/database"
	"github.com/wisle25/be-template/infrastructures/database/db_helper"
	"github.com/wisle25/be-template/infrastructures/server"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var accessTokenCookie string
var refreshTokenCookie string

func TestAddUser(t *testing.T) {
	config := commons.LoadConfig("../../")

	// Arrange
	payload := map[string]string{
		"username":        "user",
		"email":           "user@example.com",
		"password":        "password",
		"confirmPassword": "password",
	}
	body, _ := json.Marshal(payload)

	app := server.CreateServer(config)

	// Action
	req := httptest.NewRequest(
		"POST",
		"/users",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res, _ := app.Test(req)

	// Assert
	resBody, _ := io.ReadAll(res.Body)
	var resMap map[string]string
	_ = json.Unmarshal(resBody, &resMap)

	assert.Equal(t, res.StatusCode, http.StatusCreated)
	assert.NotNil(t, resMap["id"])
}

func TestLoginUser(t *testing.T) {
	config := commons.LoadConfig("../../")
	db := database.ConnectDB(config)
	userHelperDB := &db_helper.UserHelperDB{
		DB: db,
	}

	defer userHelperDB.CleanUserDB()

	// Arrange
	payload := map[string]string{
		"identity": "user",
		"password": "password",
	}
	body, _ := json.Marshal(payload)

	app := server.CreateServer(config)

	req := httptest.NewRequest(
		"POST",
		"/auths",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	// Action
	res, _ := app.Test(req)

	// Save cookie for refresh token test
	for _, cookie := range res.Cookies() {
		if cookie.Name == "refresh_token" {
			refreshTokenCookie = cookie.Value
		} else if cookie.Name == "access_token" {
			accessTokenCookie = cookie.Value
		}
	}

	// Assert
	resBody, _ := io.ReadAll(res.Body)
	var resMap map[string]string
	_ = json.Unmarshal(resBody, &resMap)

	assert.Equal(t, 3, len(res.Cookies()))
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, resMap["status"], "success")
	assert.NotEmpty(t, refreshTokenCookie)
}

func TestRefreshToken(t *testing.T) {
	// Arrange
	config := commons.LoadConfig("../../")

	app := server.CreateServer(config)

	req := httptest.NewRequest(
		"PUT",
		"/auths",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshTokenCookie,
	})

	// Action
	res, _ := app.Test(req)

	// Assert
	resBody, _ := io.ReadAll(res.Body)
	var resMap map[string]string
	_ = json.Unmarshal(resBody, &resMap)

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, resMap["status"], "success")
}

func TestLogoutUser(t *testing.T) {
	config := commons.LoadConfig("../../")

	// Arrange
	app := server.CreateServer(config)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/auths",
		nil,
	)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: accessTokenCookie,
	})
	req.AddCookie(&http.Cookie{
		Name:  "refresh_token",
		Value: refreshTokenCookie,
	})

	// Action
	res, _ := app.Test(req)

	// Assert
	resBody, _ := io.ReadAll(res.Body)
	var resMap map[string]string
	_ = json.Unmarshal(resBody, &resMap)

	assert.Equal(t, 3, len(res.Cookies()))
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, resMap["status"], "success")
	assert.NotEmpty(t, refreshTokenCookie)
}
