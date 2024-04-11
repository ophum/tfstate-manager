package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"net/http/cgi"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	statev1 "github.com/ophum/tfstate-manager/api/state/v1"
	userv1 "github.com/ophum/tfstate-manager/api/user/v1"
	"github.com/ophum/tfstate-manager/gen/api/state/v1/statev1connect"
	"github.com/ophum/tfstate-manager/gen/api/user/v1/userv1connect"
	httpbackend "github.com/ophum/tfstate-manager/pkg/http_backend"
	"github.com/ophum/tfstate-manager/pkg/middlewares"
	"github.com/ophum/tfstate-manager/pkg/models"
	"github.com/ophum/tfstate-manager/pkg/oauth"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
)

type Config struct {
	SessionStoreKeyPairs string `yaml:"sessionStoreKeyPairs"`
	Database             string `yaml:"database"`
	ClientID             string `yaml:"clientID"`
	ClientSecret         string `yaml:"clientSecret"`
	RedirectURL          string `yaml:"redirectURL"`
}

var (
	config Config
)

func init() {
	if err := initialize(); err != nil {
		log.Fatal(err)
	}
}

func initialize() error {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.yaml"
	}

	f, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return err
	}

	gob.RegisterName("token", &oauth2.Token{})
	return nil
}

func main() {
	log.SetOutput(os.Stderr)
	gin.SetMode(gin.ReleaseMode)

	dbLog, err := os.Create("/tmp/db.log")
	if err != nil {
		log.Fatal(err)
	}
	defer dbLog.Close()

	db, err := gorm.Open(sqlite.Open(config.Database), &gorm.Config{
		Logger: logger.New(log.New(dbLog, "", log.LstdFlags), logger.Config{
			Colorful: true,
			LogLevel: logger.Info,
		}),
	})
	if err != nil {
		log.Fatal(err)
	}

	_ = db.AutoMigrate(&models.User{}, &models.State{})

	store := gormsessions.NewStore(db, true, []byte(config.SessionStoreKeyPairs))

	conf := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes: []string{
			"user:email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://github.com/login/oauth/authorize",
			TokenURL:  "https://github.com/login/oauth/access_token",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}
	router := gin.New()
	router.ContextWithFallback = true
	router.Use(gin.LoggerWithWriter(os.Stderr), gin.Recovery())

	scriptName := os.Getenv("SCRIPT_NAME")
	r := router.Group(scriptName)

	httpBackend := httpbackend.NewHTTPBackendServer(db)
	httpBackend.RegisterHandlers(r)

	r.Use(sessions.Sessions("sessions", store))

	frontendBaseURL, err := url.Parse("http://localhost:5173")
	if err != nil {
		log.Fatal(err)
	}

	oauthGithub := oauth.NewGithub(db, conf, frontendBaseURL)
	oauthGithub.RegisterHandlers(r)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowOrigins = append(corsConfig.AllowOrigins, "http://localhost:5173")
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "connect-protocol-version")
	r.Use(cors.New(corsConfig))
	r.Use(middlewares.SessionAuth())
	r.Use(middlewares.StripPrefix(scriptName))
	{
		s := statev1.NewStateServer(db)
		p, h := statev1connect.NewStateServiceHandler(s)
		registerHandler(r, p, h)
	}
	{
		s := userv1.NewUserServer(db)
		p, h := userv1connect.NewUserServiceHandler(s)
		registerHandler(r, p, h)
	}

	for _, r := range router.Routes() {
		dbLog.WriteString(fmt.Sprintf("%s %s\n", r.Method, r.Path))
	}
	if err := cgi.Serve(router); err != nil {
		log.Fatal(err)
	}
}

func registerHandler(r gin.IRouter, path string, handler http.Handler) {
	r.Any(path+"*any", gin.WrapH(handler))
}
