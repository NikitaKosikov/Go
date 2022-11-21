package service

import (
	"context"
	"fmt"
	"test/internal/domain"
	db_mocks "test/internal/repository/mocks"
	"test/internal/service/dto"
	"test/pkg/api/auth"
	"test/pkg/hash"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mockUserService(t *testing.T) (*UserService, *db_mocks.MockUserRepository) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	userRepoMock := db_mocks.NewMockUserRepository(mockCtl)

	userService := NewUserService(
		userRepoMock,
		&auth.Manager{},
		&hash.SHA1Hasher{},
		1*time.Minute,
		1*time.Minute,
	)

	return userService, userRepoMock
}

func TestUserRepository_Create(t *testing.T) {
	userService, userRepoMock := mockUserService(t)
	type mockRepoBehavior func(dbmock *db_mocks.MockUserRepository)

	testTable := []struct {
		name               string
		userDTO            dto.CreateUserDTO
		expectedResult     dto.TokenDTO
		mockRepoBehavior   mockRepoBehavior
		assertServiceTests func() []func(t *testing.T, err error, i ...interface{})
	}{
		{
			name:    "OK",
			userDTO: dto.CreateUserDTO{Email: "test@test.ru", Password: "test1234"},
			expectedResult: dto.TokenDTO{
				AccessToken:  "Rand string",
				RefreshToken: "Rand stirng",
			},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().Create(context.Background(), gomock.Any()).Return(primitive.NewObjectID(), nil)
				dbmock.EXPECT().SetSession(context.Background(), gomock.Any(), gomock.Any()).Return(nil)
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.NotEmpty(t, i)
					},
					func(t *testing.T, err error, i ...interface{}) {
						assert.Nil(t, err)
					},
				}
			},
		},
		{
			name:    "User already exists",
			userDTO: dto.CreateUserDTO{Email: "test@test.ru", Password: "test1234"},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, domain.ErrUserAlreadyExists)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
					},
				}
			},
		},
		{
			name:             "Email Invalid",
			userDTO:          dto.CreateUserDTO{Email: "testest.ru", Password: "test1234"},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "Invalid userDTO parameters")
					},
				}
			},
		},
		{
			name:             "Password Invalid",
			userDTO:          dto.CreateUserDTO{Email: "test@test.ru", Password: "test"},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "Invalid userDTO parameters")
					},
				}
			},
		},
		{
			name:             "Hash Error",
			userDTO:          dto.CreateUserDTO{Email: "test@test.ru", Password: ""},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "Failed to hash password")
					},
				}
			},
		},
		{
			name:    "Create session service Failure",
			userDTO: dto.CreateUserDTO{Email: "test@test.ru", Password: "test1234"},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().Create(context.Background(), gomock.Any()).Return(primitive.NewObjectID(), nil)
				dbmock.EXPECT().SetSession(context.Background(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("create session service failure"))
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "create session service failure")
					},
				}
			},
		},
		{
			name:           "Repository Failure",
			expectedResult: dto.TokenDTO{},
			userDTO:        dto.CreateUserDTO{Email: "test@test.ru", Password: "test1234"},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().Create(context.Background(), gomock.Any()).Return(primitive.ObjectID{}, fmt.Errorf("repository failure"))
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "repository failure")
					},
				}
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			testCase.mockRepoBehavior(userRepoMock)

			actualToken, err := userService.Create(context.Background(), testCase.userDTO)

			for _, assert := range testCase.assertServiceTests() {
				assert(t, err, testCase.expectedResult, actualToken)
			}
		})
	}
}

func TestUserRepository_FindOne(t *testing.T) {
	userService, userRepoMock := mockUserService(t)
	type mockRepoBehavior func(dbmock *db_mocks.MockUserRepository)

	testTable := []struct {
		name               string
		id                 string
		expectedResult     domain.User
		mockRepoBehavior   mockRepoBehavior
		assertServiceTests func() []func(t *testing.T, err error, i ...interface{})
	}{
		{
			name: "OK",
			id:   primitive.NewObjectID().Hex(),
			expectedResult: domain.User{
				Email: "test@test.ru",
			},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindOne(context.Background(), gomock.Any()).Return(domain.User{Email: "test@test.ru"}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Equal(t, i[0], i[1])
					},
					func(t *testing.T, err error, i ...interface{}) {
						assert.Nil(t, err)
					},
				}
			},
		},
		{
			name:             "Id Invalid",
			id:               "0000000000000",
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "invalid id param")
					},
				}
			},
		},
		{
			name: "Repository Failure",
			id:   "000000000000",
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindOne(context.Background(), gomock.Any()).Return(domain.User{}, fmt.Errorf("repository failure"))
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "repository failure")
					},
				}
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			testCase.mockRepoBehavior(userRepoMock)

			actualToken, err := userService.FindOne(context.Background(), testCase.id)

			for _, assert := range testCase.assertServiceTests() {
				assert(t, err, testCase.expectedResult, actualToken)
			}
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	userService, userRepoMock := mockUserService(t)
	type mockRepoBehavior func(dbmock *db_mocks.MockUserRepository)

	testTable := []struct {
		name               string
		email              string
		expectedResult     domain.User
		mockRepoBehavior   mockRepoBehavior
		assertServiceTests func() []func(t *testing.T, err error, i ...interface{})
	}{
		{
			name:  "OK",
			email: "test@test.ru",
			expectedResult: domain.User{
				Email: "test@test.ru",
			},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{Email: "test@test.ru"}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Equal(t, i[0], i[1])
					},
					func(t *testing.T, err error, i ...interface{}) {
						assert.Nil(t, err)
					},
				}
			},
		},

		{
			name:           "Repository Failure",
			email:          "test@test.ru",
			expectedResult: domain.User{},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, fmt.Errorf("repository failure"))
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Equal(t, i[0], i[1])
					},
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "repository failure")
					},
				}
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			testCase.mockRepoBehavior(userRepoMock)

			user, err := userService.FindByEmail(context.Background(), testCase.email)

			for _, assert := range testCase.assertServiceTests() {
				assert(t, err, testCase.expectedResult, user)
			}
		})
	}
}

func TestUserRepository_FindAll(t *testing.T) {
	userService, userRepoMock := mockUserService(t)
	type mockRepoBehavior func(dbmock *db_mocks.MockUserRepository)
	testTable := []struct {
		name               string
		limit              string
		offset             string
		filter             string
		sortBy             string
		expectedResult     []domain.User
		mockRepoBehavior   mockRepoBehavior
		assertServiceTests func() []func(t *testing.T, err error, i ...interface{})
	}{
		{
			name:   "OK. Without sorting, filtres, pagination",
			limit:  "",
			offset: "",
			filter: "",
			sortBy: "",
			expectedResult: []domain.User{
				{
					Email: "email1",
				},
				{
					Email: "email2",
				},
			},

			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindAll(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.User{
					{
						Email: "email1",
					},
					{
						Email: "email2",
					},
				}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Equal(t, i[0], i[1])
					},
					func(t *testing.T, err error, i ...interface{}) {
						assert.Nil(t, err)
					},
				}
			},
		},
		{
			name:   "OK. With sorting, filtres, pagination",
			limit:  "2",
			offset: "1",
			filter: "email[eq]=email",
			sortBy: "email.desc",
			expectedResult: []domain.User{
				{
					Email: "email1",
				},
				{
					Email: "email2",
				},
			},

			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindAll(context.Background(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.User{
					{
						Email: "email1",
					},
					{
						Email: "email2",
					},
				}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Equal(t, i[0], i[1])
					},
					func(t *testing.T, err error, i ...interface{}) {
						assert.Nil(t, err)
					},
				}
			},
		},
		{
			name:             "Limit Invalid",
			limit:            "-2",
			offset:           "1",
			filter:           "email[eq]=email",
			sortBy:           "email.desc",
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "limit query parameter is no valid number")
					},
				}
			},
		},
		{
			name:             "Offset Invalid",
			limit:            "2",
			offset:           "-1",
			filter:           "email[eq]=email,password[eq]=password",
			sortBy:           "email.desc",
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "offset query parameter is no valid number")
					},
				}
			},
		},

		{
			name:             "Filter Signutare Invalid",
			limit:            "2",
			offset:           "1",
			filter:           "emaileq]=email,password[eq]password",
			sortBy:           "email.desc",
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "malformed filter query parameter, should be field[operator]=value")
					},
				}
			},
		},
		{
			name:             "Filter Operation Invalid",
			limit:            "2",
			offset:           "1",
			filter:           "email[ewqeq]=email,password[eq]=password",
			sortBy:           "email.desc",
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "invalid operator")
					},
				}
			},
		},
		{
			name:             "Sorting Invalid",
			limit:            "2",
			offset:           "1",
			filter:           "email[eq]=email,password[eq]=password",
			sortBy:           "email.qdesc,",
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "sortBy query parameter is no valid number")
					},
				}
			},
		},
		{
			name:   "Repository Failure",
			limit:  "2",
			offset: "1",
			filter: "email[eq]=email,password[eq]=password",
			sortBy: "email.desc",
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().Create(context.Background(), gomock.Any()).Return(primitive.ObjectID{}, fmt.Errorf("repository failure"))
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "repository failure")
					},
				}
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			testCase.mockRepoBehavior(userRepoMock)

			actualToken, err := userService.FindAll(context.Background(), testCase.limit, testCase.offset, testCase.filter, testCase.sortBy)

			for _, assert := range testCase.assertServiceTests() {
				assert(t, err, testCase.expectedResult, actualToken)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	userService, userRepoMock := mockUserService(t)
	type mockRepoBehavior func(dbmock *db_mocks.MockUserRepository)
	testTable := []struct {
		name               string
		userDTO            dto.UpdateUserDTO
		mockRepoBehavior   mockRepoBehavior
		assertServiceTests func() []func(t *testing.T, err error, i ...interface{})
	}{
		{
			name: "OK",

			userDTO: dto.UpdateUserDTO{
				Id:       primitive.NewObjectID().Hex(),
				Email:    "test@test.ru",
				Password: "test1234",
			},

			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().Update(context.Background(), gomock.Any()).Return(nil)
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Nil(t, err)
					},
				}
			},
		},
		{
			name: "User with this email already exist",

			userDTO: dto.UpdateUserDTO{
				Id:       primitive.NewObjectID().Hex(),
				Email:    "test@test.ru",
				Password: "test1234",
			},

			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, domain.ErrUserAlreadyExists)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
					},
				}
			},
		},
		{
			name: "Email Invalid",

			userDTO: dto.UpdateUserDTO{
				Id:       primitive.NewObjectID().Hex(),
				Email:    "testest.ru",
				Password: "test1234",
			},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "Invalid userDTO parameters")
					},
				}
			},
		},
		{
			name: "Password Invalid",
			userDTO: dto.UpdateUserDTO{
				Id:       primitive.NewObjectID().Hex(),
				Email:    "test@test.ru",
				Password: "test",
			},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "Invalid userDTO parameters")
					},
				}
			},
		},
		{
			name: "Hash Error",
			userDTO: dto.UpdateUserDTO{
				Id:       primitive.NewObjectID().Hex(),
				Email:    "test@test.ru",
				Password: "",
			},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "Failed to hash password")
					},
				}
			},
		},
		{
			name: "Convert userDTO to user Error",
			userDTO: dto.UpdateUserDTO{
				Id:       "000000000000",
				Email:    "test@test.ru",
				Password: "test1234",
			},
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "invalid id param")
					},
				}
			},
		},
		{
			name: "Repository Failure",
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().Create(context.Background(), gomock.Any()).Return(primitive.ObjectID{}, fmt.Errorf("repository failure"))
				dbmock.EXPECT().FindByEmail(context.Background(), gomock.Any()).Return(domain.User{}, nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "repository failure")
					},
				}
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			testCase.mockRepoBehavior(userRepoMock)

			err := userService.Update(context.Background(), testCase.userDTO)

			for _, assert := range testCase.assertServiceTests() {
				assert(t, err)
			}
		})
	}
}
func TestUserRepository_Delete(t *testing.T) {
	userService, userRepoMock := mockUserService(t)
	type mockRepoBehavior func(dbmock *db_mocks.MockUserRepository)
	testTable := []struct {
		name               string
		id                 string
		mockRepoBehavior   mockRepoBehavior
		assertServiceTests func() []func(t *testing.T, err error, i ...interface{})
	}{
		{
			name: "OK",
			id:   primitive.NewObjectID().Hex(),
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().Delete(context.Background(), gomock.Any()).Return(nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Nil(t, err)
					},
				}
			},
		},
		{
			name:             "Id Invalid",
			id:               "0000000000000",
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "invalid id param")
					},
				}
			},
		},
		{
			name: "Repository Failure",
			id:   primitive.NewObjectID().Hex(),
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().Delete(context.Background(), gomock.Any()).Return(fmt.Errorf("repository failure"))
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "repository failure")
					},
				}
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			testCase.mockRepoBehavior(userRepoMock)

			err := userService.Delete(context.Background(), testCase.id)

			for _, assert := range testCase.assertServiceTests() {
				assert(t, err)
			}
		})
	}
}

func TestUserRepository_RefreshUserToken(t *testing.T) {
	userService, userRepoMock := mockUserService(t)
	type mockRepoBehavior func(dbmock *db_mocks.MockUserRepository)

	testTable := []struct {
		name               string
		id                 string
		mockRepoBehavior   mockRepoBehavior
		assertServiceTests func() []func(t *testing.T, err error, i ...interface{})
	}{
		{
			name: "OK",
			id:   primitive.NewObjectID().Hex(),
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().GetUserByRefreshToken(context.Background(), gomock.Any()).Return(domain.User{
					Id:           [12]byte{},
					PasswordHash: "123456781234",
					Email:        "test@test.ru",
					Session:      domain.Session{},
				}, nil)
				dbmock.EXPECT().SetSession(context.Background(), gomock.Any(), gomock.Any()).Return(nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.NotEmpty(t, i)
					},
					func(t *testing.T, err error, i ...interface{}) {
						assert.Nil(t, err)
					},
				}
			},
		},
		{
			name:             "Id Invalid",
			id:               "0000000000000",
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "invalid id param")
					},
				}
			},
		},
		{
			name: "Create session service Failure",
			id:   primitive.NewObjectID().Hex(),
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().GetUserByRefreshToken(context.Background(), gomock.Any()).Return(domain.User{
					Id:           [12]byte{},
					PasswordHash: "123456781234",
					Email:        "test@test.ru",
					Session:      domain.Session{},
				}, nil)
				dbmock.EXPECT().SetSession(context.Background(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("create session service failure"))
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "create session service failure")
					},
				}
			},
		},
		{
			name: "Repository Failure",
			id:   primitive.NewObjectID().Hex(),
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().GetUserByRefreshToken(context.Background(), gomock.Any()).Return(
					domain.User{},
					fmt.Errorf("repository failure"))
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "repository failure")
					},
				}
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			testCase.mockRepoBehavior(userRepoMock)

			actualToken, err := userService.RefreshUserToken(context.Background(), testCase.id)

			for _, assert := range testCase.assertServiceTests() {
				assert(t, err, actualToken)
			}
		})
	}
}

func TestUserRepository_CreateSession(t *testing.T) {
	userService, userRepoMock := mockUserService(t)
	type mockRepoBehavior func(dbmock *db_mocks.MockUserRepository)

	testTable := []struct {
		name               string
		id                 primitive.ObjectID
		mockRepoBehavior   mockRepoBehavior
		assertServiceTests func() []func(t *testing.T, err error, i ...interface{})
	}{
		{
			name: "OK",
			id:   primitive.NewObjectID(),
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().SetSession(context.Background(), gomock.Any(), gomock.Any()).Return(nil)
			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.NotEmpty(t, i)
					},
					func(t *testing.T, err error, i ...interface{}) {
						assert.Nil(t, err)
					},
				}
			},
		},
		{
			name: "Repository Failure",
			id:   primitive.NewObjectID(),
			mockRepoBehavior: func(dbmock *db_mocks.MockUserRepository) {
				dbmock.EXPECT().SetSession(context.Background(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("repository failure"))

			},
			assertServiceTests: func() []func(t *testing.T, err error, i ...interface{}) {
				return []func(t *testing.T, err error, i ...interface{}){
					func(t *testing.T, err error, i ...interface{}) {
						assert.Error(t, err, "repository failure")
					},
				}
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			testCase.mockRepoBehavior(userRepoMock)

			actualToken, err := userService.CreateSession(context.Background(), testCase.id)

			for _, assert := range testCase.assertServiceTests() {
				assert(t, err, actualToken)
			}
		})
	}
}
