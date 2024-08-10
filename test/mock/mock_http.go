package mock

import (
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock type for the HTTP client
type MockHTTPClient struct {
	mock.Mock
}

// Do is a mock implementation of the HTTP client's Do method
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// MockHTTPResponse is a helper function to create a mock HTTP response
func MockHTTPResponse(statusCode int, body string) *http.Response {
	rec := httptest.NewRecorder()
	rec.WriteHeader(statusCode)
	rec.WriteString(body)

	return rec.Result()
}
