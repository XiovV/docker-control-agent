package app

import (
	"encoding/json"
	"fmt"
	"github.com/XiovV/docker_control/config"
	"github.com/XiovV/docker_control/controller"
	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const testConfigFilename = "../config_test.json"

type mockDockerController struct {
	mock.Mock
}

func (m *mockDockerController) PullImage(image string) error {
	args := m.Called(image)

	return args.Error(0)
}

func (m *mockDockerController) FindContainerByName(containerName string) (types.Container, bool) {
	args := m.Called(containerName)

	container := args.Get(0)
	return container.(types.Container), args.Bool(1)
}

func (m *mockDockerController) FindContainerIDByName(containerName string) (string, bool) {
	args := m.Called(containerName)

	return args.String(0), args.Bool(1)
}

func (m *mockDockerController) UpdateContainer(containerName, image string, keep bool) error {
	args := m.Called()

	return args.Error(0)
}

func (m *mockDockerController) RollbackContainer(containerName string) error {
	args := m.Called()

	return args.Error(0)
}

func TestPullImage(t *testing.T) {
	defer removeConfig(t)
	cfg, apiKey, err := config.New(testConfigFilename)
	assert.Nil(t, err)

	mockController := new(mockDockerController)

	app := New(mockController, cfg)

	router := app.Router()

	var success struct{
		Message string `json:"message"`
	}

	var errorResponse struct {
		Error string `json:"error"`
	}

	t.Run("Valid image name", func(t *testing.T) {
		mockController.On("PullImage", "imageName:latest").Return(nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/images/pull?image=imageName:latest", nil)
		req.Header.Add("key", apiKey)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.NewDecoder(w.Body).Decode(&success)
		assert.Nil(t, err)

		assert.Equal(t, "image pulled successfully", success.Message)
	})

	t.Run("Image without name", func(t *testing.T) {
		mockController.On("PullImage", ":latest").Return(controller.ErrImageFormatInvalid).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/images/pull?image=:latest", nil)
		req.Header.Add("key", apiKey)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "image format is invalid", errorResponse.Error)
	})

	t.Run("Image without tag", func(t *testing.T) {
		mockController.On("PullImage", "imageName:").Return(controller.ErrImageFormatInvalid).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/images/pull?image=imageName:", nil)
		req.Header.Add("key", apiKey)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "image format is invalid", errorResponse.Error)
	})


	t.Run("Empty image name", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/images/pull?image=", nil)
		req.Header.Add("key", apiKey)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "image value must not be empty", errorResponse.Error)
	})

	t.Run("Without image query parameter", func(t *testing.T) {
		mockController.On("PullImage", "imageName:latest").Return(fmt.Errorf("some unknown error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/images/pull?image=imageName:latest", nil)
		req.Header.Add("key", apiKey)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "some unknown error", errorResponse.Error)
	})

	t.Run("Internal server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/v1/images/pull", nil)
		req.Header.Add("key", apiKey)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "image value must not be empty", errorResponse.Error)
	})
}

func TestGetContainerImage(t *testing.T) {
	defer removeConfig(t)
	cfg, apiKey, err := config.New(testConfigFilename)
	assert.Nil(t, err)

	mockController := new(mockDockerController)

	app := New(mockController, cfg)

	router := app.Router()

	var success struct {
		Image string `json:"image"`
	}

	var errorResponse struct {
		Error string `json:"error"`
	}

	t.Run("Valid container name", func(t *testing.T) {
		mockController.On("FindContainerByName", "containerName").Return(types.Container{Image: "imageName:latest"}, true).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/containers/image/containerName", nil)
		req.Header.Add("key", apiKey)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.NewDecoder(w.Body).Decode(&success)
		assert.Nil(t, err)

		assert.Equal(t, "imageName:latest", success.Image)
	})

	t.Run("Empty container name", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/containers/image/", nil)
		req.Header.Add("key", apiKey)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "container not found", errorResponse.Error)
	})

	t.Run("Non-existent container name", func(t *testing.T) {
		mockController.On("FindContainerByName", "doesntExist").Return(types.Container{Image: ""}, false).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/containers/image/doesntExist", nil)
		req.Header.Add("key", apiKey)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "container not found", errorResponse.Error)
	})
}

func removeConfig(t *testing.T) {
	err := os.Remove(testConfigFilename)
	assert.Nil(t, err)
}
