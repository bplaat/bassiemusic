package controllers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/bplaat/bassiemusic/routes"
	"github.com/gofiber/fiber/v2"
)

func TestHome(t *testing.T) {
	api := routes.Api()
	request := httptest.NewRequest("GET", "/", nil)
	response, _ := api.Test(request)

	if response.StatusCode != fiber.StatusOK {
		t.Errorf("TestHome(): response status is not OK")
	}

	body, _ := ioutil.ReadAll(response.Body)
	if _, err := json.Marshal(body); err != nil {
		t.Errorf("TestHome(): response is not valid JSON")
	}
}
