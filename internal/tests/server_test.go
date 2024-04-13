package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/randnull/banner-service/pkg/models"
)


var TokenResponse models.Token


func TestRegister(t *testing.T) {
	register_model := models.Register{Username: "user_test", Password: "password_test"}

	register_json, err := json.Marshal(register_model)

	if err != nil {
		t.Fatal(err)
	}

	register_response, err := http.Post("http://127.0.0.1:6050/register/admin", "application/json", bytes.NewBuffer(register_json))

	if err != nil {
		t.Fatal(err)
	}

	if register_response.StatusCode != http.StatusCreated {
		t.Errorf("Except status 201, but get %v", register_response.StatusCode)
	}
}


func TestLogin(t *testing.T) {
	login_model := models.Register{Username: "user_test", Password: "password_test"}

	login_json, err := json.Marshal(login_model)

	if err != nil {
		t.Fatal(err)
	}

	login_response, err := http.Post("http://127.0.0.1:6050/login", "application/json", bytes.NewBuffer(login_json))

	if err != nil {
		t.Fatal(err)
	}

	if login_response.StatusCode != http.StatusOK {
		t.Errorf("Except status 200, but get %v", login_response.StatusCode)
	}

	err = json.NewDecoder(login_response.Body).Decode(&TokenResponse)

	if err != nil {
		t.Fatal(err)
	}
}


func TestAddBanner(t *testing.T) {
	client := &http.Client{}

	banner := models.Banner{
		TagIds:	[]int{1, 2, 3},
		FeatureId: 1,
		Content: models.Content{
			Title: "title",
			Text:  "text",
			Url:   "example.com",
		},
		IsActive: true,
	}

	banner_json, err := json.Marshal(banner)

	if err != nil {
		t.Fatal(err)
	}

	create_request, err := http.NewRequest("POST", "http://127.0.0.1:6050/banner", bytes.NewBuffer(banner_json))

	if err != nil {
		t.Fatal(err)
	}

	create_request.Header.Set("token", TokenResponse.Token)
	create_request.Header.Set("Content-Type", "application/json")

	create_response, err := client.Do(create_request)

	if err != nil {
		t.Fatal(err)
	}

	if create_response.StatusCode != http.StatusCreated {
		t.Errorf("Except status 201, but get %v", create_response.StatusCode)
	}
}

func TestGetBanner(t *testing.T) {
	client := &http.Client{}

	get_request, err := http.NewRequest("GET", "http://127.0.0.1:6050/user_banner?tag_id=1&feature_id=1", nil)
	
	if err != nil {
		t.Fatal(err)
	}
	
	get_request.Header.Set("token", TokenResponse.Token)
	get_request.Header.Set("Content-Type", "application/json")

	get_response, err := client.Do(get_request)

	if err != nil {
		t.Fatal(err)
	}
	
	if get_response.StatusCode != http.StatusOK {
		t.Errorf("Except status 200, but get %v", get_response.StatusCode)
	}
}
