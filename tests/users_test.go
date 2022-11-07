package tests

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"test/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
func (s *ApiTestSuite) TestUserCreate() {
	router := s.handler.Init()
	r := s.Require()
	email, password := "test@test.com", "qwerty123"
	usersData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBuffer([]byte(usersData)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	var user domain.User
	err := s.db.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	s.NoError(err)

	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	r.Equal(email, user.Email)
	r.Equal(passwordHash, user.PasswordHash)
}

func (s *ApiTestSuite) TestUserFindOne() {
	router := s.handler.Init()
	r := s.Require()
	email, password := "test@test.com", "qwerty123"
	id := primitive.NewObjectID()

	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	_, err = s.db.Collection("users").InsertOne(context.Background(), domain.User{
		Id:           id,
		PasswordHash: passwordHash,
		Email:        email,
		Session:      domain.Session{},
	})
	s.NoError(err)

	req, _ := http.NewRequest("GET", "/api/v1/users/"+id.Hex(), &bytes.Reader{})
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var user domain.User
	err = s.db.Collection("users").FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	s.NoError(err)

	r.Equal(id, user.Id)
	r.Equal(email, user.Email)
	r.Equal(passwordHash, user.PasswordHash)

}

func (s *ApiTestSuite) TestUserFindAll() {
	router := s.handler.Init()
	r := s.Require()

	email1, email2 := "test@test.com", "test2@test.com"
	password1, password2 := "qwerty123", "qwerty1234"
	id1, id2 := primitive.NewObjectID(), primitive.NewObjectID()

	passwordHash1, err := s.hasher.Hash(password1)
	s.NoError(err)
	passwordHash2, err := s.hasher.Hash(password2)
	s.NoError(err)

	testUsers := []domain.User{
		{
			Id:           id1,
			PasswordHash: passwordHash1,
			Email:        email1,
		},
		{
			Id:           id2,
			PasswordHash: passwordHash2,
			Email:        email2,
		},
	}
	s.db.Collection("users").DeleteMany(context.Background(), bson.D{})
	insertUsers := make([]interface{}, len(testUsers))
	for i := 0; i < len(insertUsers); i++ {
		insertUsers[i] = testUsers[i]
	}

	_, err = s.db.Collection("users").InsertMany(context.Background(), insertUsers)
	s.NoError(err)

	req, _ := http.NewRequest("GET", "/api/v1/users/", &bytes.Reader{})
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var users []domain.User
	cursor, err := s.db.Collection("users").Find(context.Background(), bson.D{})
	s.NoError(err)

	err = cursor.All(context.Background(), &users)
	s.NoError(err)

	for _, testUser := range testUsers {
		r.Contains(users, testUser)
	}
}

func (s *ApiTestSuite) TestUserUpdate() {
	router := s.handler.Init()
	r := s.Require()

	email, password := "test@test.com", "qwerty123"
	id := primitive.NewObjectID()

	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	_, err = s.db.Collection("users").InsertOne(context.Background(), domain.User{
		Id:           id,
		PasswordHash: passwordHash,
		Email:        email,
		Session:      domain.Session{},
	})
	s.NoError(err)

	updateEmail, updatePassword := "test@test.ru", "qwerty1234"
	usersData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, updateEmail, updatePassword)

	req, _ := http.NewRequest("PUT", "/api/v1/users/"+id.Hex(), bytes.NewBuffer([]byte(usersData)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var user domain.User
	err = s.db.Collection("users").FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	s.NoError(err)

	updatePasswordHash, err := s.hasher.Hash(updatePassword)
	s.NoError(err)

	r.Equal(updateEmail, user.Email)
	r.Equal(updatePasswordHash, user.PasswordHash)
}

func (s *ApiTestSuite) TestUserDelete() {
	router := s.handler.Init()
	r := s.Require()

	email, password := "test@test.com", "qwerty123"
	id := primitive.NewObjectID()

	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	_, err = s.db.Collection("users").InsertOne(context.Background(), domain.User{
		Id:           id,
		PasswordHash: passwordHash,
		Email:        email,
		Session:      domain.Session{},
	})
	s.NoError(err)

	req, _ := http.NewRequest("DELETE", "/api/v1/users/"+id.Hex(), &bytes.Reader{})
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	result := s.db.Collection("users").FindOne(context.Background(), bson.M{"_id": id})

	s.ErrorIs(result.Err(), mongo.ErrNoDocuments)
}

func (s *ApiTestSuite) TestUserCreateSetSession() {
	router := s.handler.Init()
	r := s.Require()
	email, password := "test@test.com", "qwerty123"
	usersData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req, _ := http.NewRequest("POST", "/api/v1/users/", bytes.NewBuffer([]byte(usersData)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	var user domain.User
	err := s.db.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	s.NoError(err)

	r.NotEmpty(user.Session.RefreshToken, user.Session.ExpiresAt)

}

func (s *ApiTestSuite) TestUserGetUserByRefreshToken() {
	router := s.handler.Init()
	r := s.Require()

	email, password := "test@test.com", "qwerty123"
	id := primitive.NewObjectID()

	session := domain.Session{
		RefreshToken: "Test token",
		ExpiresAt:    time.Now().Add(time.Minute * 2),
	}
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	_, err = s.db.Collection("users").InsertOne(context.Background(), domain.User{
		Id:           id,
		PasswordHash: passwordHash,
		Email:        email,
		Session:      session,
	})
	s.NoError(err)

	req, _ := http.NewRequest("GET", "/api/v1/users/"+id.Hex()+"/auth/refresh", &bytes.Reader{})
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var user domain.User
	err = s.db.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	s.NoError(err)

	r.NotEqual(user.Session.RefreshToken, session.RefreshToken)
	r.NotEqual(user.Session.ExpiresAt, session.ExpiresAt)

}
