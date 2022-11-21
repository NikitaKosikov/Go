package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"test/internal/domain"
	"test/internal/service"
	"test/internal/service/dto"
	"test/internal/service/mocks"
	apierrors "test/pkg/api/api_errors"
	"test/pkg/api/auth"
	"test/pkg/api/params"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Create(t *testing.T) {
	type mockBehavior func(s *mocks.MockUsers, user dto.CreateUserDTO)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           dto.CreateUserDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"Test","password":"qwerty"}`,
			inputUser: dto.CreateUserDTO{
				Email:    "Test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mocks.MockUsers, userDTO dto.CreateUserDTO) {
				s.EXPECT().Create(context.Background(), userDTO).Return(dto.TokenDTO{
					AccessToken:  "Rand string",
					RefreshToken: "Rand string",
				}, nil)
			},
			expectedStatusCode: 201,
		},
		{
			name:      "User already exist",
			inputBody: `{"email":"Test","password":"qwerty"}`,
			inputUser: dto.CreateUserDTO{
				Email:    "Test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mocks.MockUsers, userDTO dto.CreateUserDTO) {
				s.EXPECT().Create(context.Background(), userDTO).Return(dto.TokenDTO{}, domain.ErrUserAlreadyExists)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"user with such email already exists"}`,
		},
		{
			name:                "Empty fields",
			mockBehavior:        func(s *mocks.MockUsers, userDTO dto.CreateUserDTO) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"failed to bind user and json"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"email":"Test","password":"qwerty"}`,
			inputUser: dto.CreateUserDTO{
				Email:    "Test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mocks.MockUsers, userDTO dto.CreateUserDTO) {
				s.EXPECT().Create(context.Background(), userDTO).Return(dto.TokenDTO{}, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userMockService := mocks.NewMockUsers(c)
			testCase.mockBehavior(userMockService, testCase.inputUser)

			services := &service.Services{Users: userMockService}
			handler := NewHandler(services, &auth.Manager{})

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			w := httptest.NewRecorder()

			r.POST("/users", handler.Create)
			req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

func TestHandler_FindOne(t *testing.T) {
	type mockBehavior func(s *mocks.MockUsers, id string)

	testTable := []struct {
		name                string
		id                  string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			id:   "000000000000",
			mockBehavior: func(s *mocks.MockUsers, id string) {
				s.EXPECT().FindOne(context.Background(), id).Return(domain.User{
					Id:           [12]byte{1},
					PasswordHash: "password",
					Email:        "email",
					Session:      domain.Session{},
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"email":"email"}`,
		},
		{
			name: "Empty id",
			id:   "000000000000",
			mockBehavior: func(s *mocks.MockUsers, id string) {
				s.EXPECT().FindOne(context.Background(), id).Return(domain.User{}, params.ErrInvalidIdParam)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid id param"}`,
		},
		{
			name: "User not found",
			id:   "000000000000",
			mockBehavior: func(s *mocks.MockUsers, id string) {
				s.EXPECT().FindOne(context.Background(), id).Return(domain.User{}, domain.ErrUserNotFound)
			},
			expectedStatusCode:  404,
			expectedRequestBody: `{"message":"user doesn't exists"}`,
		},
		{
			name: "Service Failure",
			id:   "000000000000",
			mockBehavior: func(s *mocks.MockUsers, id string) {
				s.EXPECT().FindOne(context.Background(), id).Return(domain.User{}, fmt.Errorf("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userMockService := mocks.NewMockUsers(c)
			testCase.mockBehavior(userMockService, testCase.id)

			services := &service.Services{Users: userMockService}
			handler := NewHandler(services, &auth.Manager{})

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			w := httptest.NewRecorder()

			r.GET("/users/:id", handler.FindOne)
			req := httptest.NewRequest("GET", "/users/"+testCase.id, &bytes.Reader{})

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

func TestHandler_FindAll(t *testing.T) {
	type mockBehavior func(s *mocks.MockUsers, limit, offset, filter, sortBy string)

	testTable := []struct {
		name                string
		limit               string
		offset              string
		filter              string
		sortBy              string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:   "OK. Without sorting, filtres, pagination",
			limit:  "",
			offset: "",
			filter: "",
			sortBy: "",
			mockBehavior: func(s *mocks.MockUsers, limit, offset, filter, sortBy string) {
				s.EXPECT().FindAll(context.Background(), limit, offset, filter, sortBy).Return([]domain.User{
					{
						Id:           [12]byte{1},
						PasswordHash: "password1",
						Email:        "email1",
						Session:      domain.Session{},
					},
					{
						Id:           [12]byte{2},
						PasswordHash: "password2",
						Email:        "email2",
						Session:      domain.Session{},
					},
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[{"email":"email1"},{"email":"email2"}]`,
		},
		{
			name:   "OK. With sorting, filtres, pagination",
			limit:  "2",
			offset: "1",
			filter: "email[eq]=email,password[eq]=password",
			sortBy: "email.desc,password.asc",
			mockBehavior: func(s *mocks.MockUsers, limit, offset, filter, sortBy string) {
				s.EXPECT().FindAll(context.Background(), limit, offset, filter, sortBy).Return([]domain.User{
					{
						Id:           [12]byte{1},
						PasswordHash: "password1",
						Email:        "email1",
						Session:      domain.Session{},
					},
					{
						Id:           [12]byte{2},
						PasswordHash: "password2",
						Email:        "email2",
						Session:      domain.Session{},
					},
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[{"email":"email1"},{"email":"email2"}]`,
		},
		{
			name:   "Limit Invalid",
			limit:  "-2",
			offset: "1",
			filter: "email[eq]=email,password[eq]=password",
			sortBy: "email.desc,password.asc",
			mockBehavior: func(s *mocks.MockUsers, limit, offset, filter, sortBy string) {
				s.EXPECT().FindAll(context.Background(), limit, offset, filter, sortBy).Return([]domain.User{}, apierrors.ErrLimitInvalid)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"limit query parameter is no valid number"}`,
		},
		{
			name:   "Offset Invalid",
			limit:  "2",
			offset: "-1",
			filter: "email[eq]=email,password[eq]=password",
			sortBy: "email.desc,password.asc",
			mockBehavior: func(s *mocks.MockUsers, limit, offset, filter, sortBy string) {
				s.EXPECT().FindAll(context.Background(), limit, offset, filter, sortBy).Return([]domain.User{}, apierrors.ErrOffsetInvalid)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"offset query parameter is no valid number"}`,
		},
		{
			name:   "Filter Signutare Invalid",
			limit:  "2",
			offset: "1",
			filter: "emaileq]=email,password[eq]password",
			sortBy: "email.desc,password.asc",
			mockBehavior: func(s *mocks.MockUsers, limit, offset, filter, sortBy string) {
				s.EXPECT().FindAll(context.Background(), limit, offset, filter, sortBy).Return([]domain.User{}, apierrors.ErrFilterInvalid)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"malformed filter query parameter, should be field[operator]=value"}`,
		},
		{
			name:   "Filter Operation Invalid",
			limit:  "2",
			offset: "1",
			filter: "email[ewqeq]=email,password[eq]=password",
			sortBy: "email.desc,password.asc",
			mockBehavior: func(s *mocks.MockUsers, limit, offset, filter, sortBy string) {
				s.EXPECT().FindAll(context.Background(), limit, offset, filter, sortBy).Return([]domain.User{}, apierrors.ErrFilterOperatorInvalid)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid filter operator"}`,
		},
		{
			name:   "Sorting Invalid",
			limit:  "2",
			offset: "1",
			filter: "email[eq]=email,password[eq]=password",
			sortBy: "email.qdesc,password..asc",
			mockBehavior: func(s *mocks.MockUsers, limit, offset, filter, sortBy string) {
				s.EXPECT().FindAll(context.Background(), limit, offset, filter, sortBy).Return([]domain.User{}, apierrors.ErrSortByInvalid)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"sortBy query parameter is no valid number"}`,
		},
		{
			name:   "Service Failure",
			limit:  "2",
			offset: "1",
			filter: "email[eq]=email,password[eq]=password",
			sortBy: "email.desc,password.asc",
			mockBehavior: func(s *mocks.MockUsers, limit, offset, filter, sortBy string) {
				s.EXPECT().FindAll(context.Background(), limit, offset, filter, sortBy).Return([]domain.User{}, fmt.Errorf("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userMockService := mocks.NewMockUsers(c)
			testCase.mockBehavior(userMockService, testCase.limit, testCase.offset, testCase.filter, testCase.sortBy)

			services := &service.Services{Users: userMockService}
			handler := NewHandler(services, &auth.Manager{})

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			w := httptest.NewRecorder()

			r.GET("/users", handler.FindAll)
			req := httptest.NewRequest("GET", "/users", &bytes.Reader{})

			q := req.URL.Query()
			q.Add("limit", testCase.limit)
			q.Add("offset", testCase.offset)
			q.Add("sortBy", testCase.sortBy)
			q.Add("filter", testCase.filter)
			req.URL.RawQuery = q.Encode()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

func TestHandler_Update(t *testing.T) {
	type mockBehavior func(s *mocks.MockUsers, user dto.UpdateUserDTO)

	testTable := []struct {
		name                string
		id                  string
		inputBody           string
		inputUser           dto.UpdateUserDTO
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"Test","password":"qwerty"}`,
			id:        "000000000000",
			inputUser: dto.UpdateUserDTO{
				Id:       "000000000000",
				Email:    "Test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mocks.MockUsers, userDTO dto.UpdateUserDTO) {
				s.EXPECT().Update(context.Background(), userDTO).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:      "User with this email already exist",
			inputBody: `{"email":"Test","password":"qwerty"}`,
			id:        "000000000000",
			inputUser: dto.UpdateUserDTO{
				Id:       "000000000000",
				Email:    "Test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mocks.MockUsers, userDTO dto.UpdateUserDTO) {
				s.EXPECT().Update(context.Background(), userDTO).Return(domain.ErrUserAlreadyExists)
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"user with such email already exists"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"email":"Test","password":"qwerty"}`,
			id:        "000000000000",
			inputUser: dto.UpdateUserDTO{
				Id:       "000000000000",
				Email:    "Test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mocks.MockUsers, userDTO dto.UpdateUserDTO) {
				s.EXPECT().Update(context.Background(), userDTO).Return(fmt.Errorf("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userMockService := mocks.NewMockUsers(c)
			testCase.mockBehavior(userMockService, testCase.inputUser)

			services := &service.Services{Users: userMockService}
			handler := NewHandler(services, &auth.Manager{})

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			w := httptest.NewRecorder()

			r.PUT("/users/:id", handler.Update)
			req := httptest.NewRequest("PUT", "/users/"+testCase.id, bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

func TestHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mocks.MockUsers, id string)

	testTable := []struct {
		name                string
		id                  string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			id:   "000000000000",
			mockBehavior: func(s *mocks.MockUsers, id string) {
				s.EXPECT().Delete(context.Background(), id).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name: "Service Failure",
			id:   "000000000000",
			mockBehavior: func(s *mocks.MockUsers, id string) {
				s.EXPECT().Delete(context.Background(), id).Return(fmt.Errorf("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userMockService := mocks.NewMockUsers(c)
			testCase.mockBehavior(userMockService, testCase.id)

			services := &service.Services{Users: userMockService}
			handler := NewHandler(services, &auth.Manager{})

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			w := httptest.NewRecorder()

			r.DELETE("/users/:id", handler.Delete)
			req := httptest.NewRequest("DELETE", "/users/000000000000", &bytes.Reader{})

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_RefreshToken(t *testing.T) {
	type mockBehavior func(s *mocks.MockUsers, id string)

	testTable := []struct {
		name                string
		userId              string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:   "OK",
			userId: "000000000000",
			mockBehavior: func(s *mocks.MockUsers, userId string) {
				s.EXPECT().RefreshUserToken(context.Background(), userId).Return(dto.TokenDTO{
					AccessToken:  "Rand string",
					RefreshToken: "Rand string",
				}, nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:   "Service Failure",
			userId: "000000000000",
			mockBehavior: func(s *mocks.MockUsers, userId string) {
				s.EXPECT().RefreshUserToken(context.Background(), userId).Return(dto.TokenDTO{}, fmt.Errorf("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userMockService := mocks.NewMockUsers(c)
			testCase.mockBehavior(userMockService, testCase.userId)

			services := &service.Services{Users: userMockService}
			handler := NewHandler(services, &auth.Manager{})

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			w := httptest.NewRecorder()

			r.GET("/users/:id/auth/refresh", handler.RefreshToken)
			req := httptest.NewRequest("GET", "/users/000000000000/auth/refresh", &bytes.Reader{})

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
