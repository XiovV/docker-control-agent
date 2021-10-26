package app

import (
	"encoding/json"
	"errors"
	"github.com/XiovV/dokkup-agent/config"
	"github.com/XiovV/dokkup-agent/controller"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func (m *mockDockerController) FindContainerIDByName(containerName string) (string, bool) {
	args := m.Called(containerName)

	return args.String(0), args.Bool(1)
}

func (m *mockDockerController) UpdateContainer(containerName, image string, keep bool) error {
	args := m.Called(containerName, image, keep)

	return args.Error(0)
}

func (m *mockDockerController) RollbackContainer(containerName string) error {
	args := m.Called(containerName)

	return args.Error(0)
}

func TestUpdateContainer(t *testing.T) {
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

	t.Run("Valid update request", func(t *testing.T) {
		mockController.On("UpdateContainer", "validContainer", "imageName:latest", true).
			Return(nil).Once()

		w := sendRequest(router, "PUT", "/v1/containers/update?container=validContainer&image=imageName:latest&keep=true", apiKey)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.NewDecoder(w.Body).Decode(&success)
		assert.Nil(t, err)

		assert.Equal(t, "container updated successfully", success.Message)
	})

	t.Run("Empty container value", func(t *testing.T) {
		w := sendRequest(router, "PUT", "/v1/containers/update?container=&image=imageName:latest&keep=true", apiKey)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "container value must not be empty", errorResponse.Error)
	})

	t.Run("Without container query parameter", func(t *testing.T) {
		w := sendRequest(router, "PUT", "/v1/containers/update?image=imageName:latest&keep=true", apiKey)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "container value must not be empty", errorResponse.Error)
	})

	t.Run("Empty image value", func(t *testing.T) {
		w := sendRequest(router, "PUT", "/v1/containers/update?container=validContainer&image=&keep=true", apiKey)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "image value must not be empty", errorResponse.Error)
	})

	t.Run("Without image query parameter", func(t *testing.T) {
		w := sendRequest(router, "PUT", "/v1/containers/update?container=validContainer&keep=true", apiKey)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "image value must not be empty", errorResponse.Error)
	})

	t.Run("Invalid keep value", func(t *testing.T) {
		w := sendRequest(router, "PUT", "/v1/containers/update?container=validContainer&image=imageName:latest&keep=abc", apiKey)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "keep value must be either true or false", errorResponse.Error)
	})

	t.Run("Empty keep value", func(t *testing.T) {
		w := sendRequest(router, "PUT", "/v1/containers/update?container=validContainer&image=imageName:latest&keep=", apiKey)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "keep value must be either true or false", errorResponse.Error)
	})

	t.Run("Without keep parameter", func(t *testing.T) {
		w := sendRequest(router, "PUT", "/v1/containers/update?container=validContainer&image=imageName:latest", apiKey)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "keep value must be either true or false", errorResponse.Error)
	})

	t.Run("Image without name", func(t *testing.T) {
		mockController.On("UpdateContainer", "validContainer", ":latest", true).
			Return(controller.ErrImageFormatInvalid).Once()

		w := sendRequest(router, "PUT", "/v1/containers/update?container=validContainer&image=:latest&keep=true", apiKey)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "image format is invalid", errorResponse.Error)
	})

	t.Run("Image without tag", func(t *testing.T) {
		mockController.On("UpdateContainer", "validContainer", "imageName:", true).
			Return(controller.ErrImageFormatInvalid).Once()

		w := sendRequest(router, "PUT", "/v1/containers/update?container=validContainer&image=imageName:&keep=true", apiKey)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "image format is invalid", errorResponse.Error)
	})

	t.Run("Non-existent container name", func(t *testing.T) {
		mockController.On("UpdateContainer", "invalidContainer", "imageName:latest", true).
			Return(controller.ErrContainerNotFound).Once()

		w := sendRequest(router, "PUT", "/v1/containers/update?container=invalidContainer&image=imageName:latest&keep=true", apiKey)

		assert.Equal(t, http.StatusNotFound, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "the requested container could not be found", errorResponse.Error)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockController.On("UpdateContainer", "validContainer", "imageName:latest", true).
			Return(errors.New("some unknown error")).Once()

		w := sendRequest(router, "PUT", "/v1/containers/update?container=validContainer&image=imageName:latest&keep=true", apiKey)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "some unknown error", errorResponse.Error)
	})
}

func TestRollbackContainer(t *testing.T) {
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

	t.Run("Valid rollback request", func(t *testing.T) {
		mockController.On("RollbackContainer", "containerName").
			Return(nil).Once()

		w := sendRequest(router, "PUT", "/v1/containers/rollback?container=containerName", apiKey)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.NewDecoder(w.Body).Decode(&success)
		assert.Nil(t, err)

		assert.Equal(t, "successfully restored container", success.Message)
	})

	t.Run("Empty container value", func(t *testing.T) {
		w := sendRequest(router, "PUT", "/v1/containers/rollback?container=", apiKey)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "container value must not be empty", errorResponse.Error)
	})

	t.Run("Without container parameter", func(t *testing.T) {
		w := sendRequest(router, "PUT", "/v1/containers/rollback", apiKey)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "container value must not be empty", errorResponse.Error)
	})

	t.Run("Non-existent container", func(t *testing.T) {
		mockController.On("RollbackContainer", "invalidContainer").
			Return(controller.ErrContainerNotFound).Once()

		w := sendRequest(router, "PUT", "/v1/containers/rollback?container=invalidContainer", apiKey)

		assert.Equal(t, http.StatusNotFound, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "the requested container does not exist", errorResponse.Error)
	})

	t.Run("Non-existent rollback container", func(t *testing.T) {
		mockController.On("RollbackContainer", "invalidContainer").
			Return(controller.ErrRollbackContainerNotFound).Once()

		w := sendRequest(router, "PUT", "/v1/containers/rollback?container=invalidContainer", apiKey)

		assert.Equal(t, http.StatusNotFound, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "the requested container does not have a rollback container", errorResponse.Error)
	})

	t.Run("Container not running", func(t *testing.T) {
		mockController.On("RollbackContainer", "containerName").
			Return(controller.ErrContainerNotRunning).Once()

		w := sendRequest(router, "PUT", "/v1/containers/rollback?container=containerName", apiKey)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "the container failed to start", errorResponse.Error)
	})

	t.Run("Container failed to start", func(t *testing.T) {
		mockController.On("RollbackContainer", "containerName").
			Return(controller.ErrContainerStartFailed{Reason: errors.New("some random reason")}).Once()

		w := sendRequest(router, "PUT", "/v1/containers/rollback?container=containerName", apiKey)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "some random reason", errorResponse.Error)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockController.On("RollbackContainer", "containerName").
			Return(errors.New("some unknown error")).Once()

		w := sendRequest(router, "PUT", "/v1/containers/rollback?container=containerName", apiKey)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		err = json.NewDecoder(w.Body).Decode(&errorResponse)
		assert.Nil(t, err)

		assert.Equal(t, "some unknown error", errorResponse.Error)
	})
}