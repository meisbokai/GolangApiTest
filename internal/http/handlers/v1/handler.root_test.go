package v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	V1Handlers "github.com/meisbokai/GolangApiTest/internal/http/handlers/v1"
	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	setup(t)

	ginMock.POST("/", V1Handlers.Root)

	t.Run("Test 1 | Root Test", func(t *testing.T) {

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", nil)

		r.Header.Set("Content-Type", "application/json")

		// Perform request
		ginMock.ServeHTTP(w, r)

		// parsing json to raw text
		response := w.Body.String()

		var responseJson V1Handlers.BaseResponse

		json.Unmarshal([]byte(response), &responseJson)

		// Assertions
		// Assert status code
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		// assert.Equal(t, true, responseJson.Status)
		assert.Equal(t, "v1 online...", responseJson.Message)
	})

}
