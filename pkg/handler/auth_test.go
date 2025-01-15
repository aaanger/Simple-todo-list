package handler

import (
	"bytes"
	"github.com/aaanger/p1/pkg/models"
	services2 "github.com/aaanger/p1/pkg/services"
	mock_services "github.com/aaanger/p1/pkg/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"net/http/httptest"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_services.MockAuthorization, user models.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           models.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test", "password":"test"}`,
			inputUser: models.User{
				Username: "test",
				Password: "test",
			},
			mockBehavior: func(s *mock_services.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty fields",
			inputBody:           `{"password":"test"}`,
			mockBehavior:        func(s *mock_services.MockAuthorization, user models.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error": "Invalid input body"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_services.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &services2.Service{
				Authorization: auth,
			}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/signup", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signup", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}
