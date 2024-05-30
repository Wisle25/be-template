package server_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/wisle25/be-template/commons"
	"github.com/wisle25/be-template/infrastructures/database/db_helper"
	"github.com/wisle25/be-template/infrastructures/server"
	"io"
	"net/http"

	"net/http/httptest"
	"testing"
)

func TestAddUser(t *testing.T) {
	config := commons.LoadConfig("../../")

	t.Run("Should add a new user and return 201", func(t *testing.T) {
		defer db_helper.CleanUserDB()

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
	})
}
