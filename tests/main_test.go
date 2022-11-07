package tests

import (
	"context"
	"os"
	"test/internal/config"
	v1 "test/internal/delivery/http/v1"
	"test/internal/repository"
	"test/internal/service"
	"test/pkg/api/auth"
	"test/pkg/client/mongodb"
	"test/pkg/hash"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dbURI  = "mongodb://localhost:27019"
	dbName = "testDb"
)

type ApiTestSuite struct {
	suite.Suite

	db       *mongo.Database
	handler  *v1.Handler
	services *service.Services
	repos    *repository.Repository

	tokenManager auth.TokenManager
	hasher       *hash.SHA1Hasher
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(ApiTestSuite))
}

func (s *ApiTestSuite) SetupSuite() {
	cfg := config.MongodbConfig{
		URI: dbURI,
	}
	if client, err := mongodb.NewClient(cfg); err != nil {
		s.FailNow("Failed to connect to mongo", err)
	} else {
		s.db = client.Database(dbName)
	}

	s.initDeps()
}

func (s *ApiTestSuite) TearDownSuite() {
	s.db.Client().Disconnect(context.Background()) //nolint:errcheck
}

func (s *ApiTestSuite) BeforeTest(suiteName, testName string) {
	s.db.Collection("users").DeleteMany(context.Background(), bson.D{})
}

func (s *ApiTestSuite) initDeps() {
	// Init domain deps
	repos := repository.NewRepository(s.db)
	hasher := hash.NewSHA1Hasher("salt")

	tokenManager, err := auth.NewManager("signing_key")
	if err != nil {
		s.FailNow("Failed to initialize token manager", err)
	}

	services := service.NewServices(service.Deps{

		Repos:        repos,
		Hasher:       hasher,
		TokenManager: tokenManager,

		AccessTokenTTL:  time.Minute * 15,
		RefreshTokenTTL: time.Minute * 15,
	})

	s.repos = repos
	s.services = services
	s.handler = v1.NewHandler(services, tokenManager)
	s.hasher = hasher
	s.tokenManager = tokenManager
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}
