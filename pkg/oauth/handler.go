package oauth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ophum/tfstate-manager/pkg/models"
	"github.com/ophum/tfstate-manager/pkg/utils"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type Github struct {
	db              *gorm.DB
	config          *oauth2.Config
	frontendBaseURL *url.URL
}

func NewGithub(db *gorm.DB, config *oauth2.Config, frontendBaseURL *url.URL) *Github {
	return &Github{
		db:              db,
		config:          config,
		frontendBaseURL: frontendBaseURL,
	}
}

func (g *Github) RegisterHandlers(router gin.IRouter) {
	r := router.Group("/oauth/github")
	r.GET("", g.begin)
	r.GET("/callback", g.callback)
}

func (g *Github) begin(ctx *gin.Context) {
	session := sessions.Default(ctx)

	state, err := utils.RandomString()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session.Set("state", state)
	if err := session.Save(); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	url := g.config.AuthCodeURL(state)
	log.Println(url)

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (g *Github) callback(ctx *gin.Context) {
	session := sessions.Default(ctx)

	state, ok := session.Get("state").(string)
	if !ok {
		log.Println("stateが存在しません")
		ctx.AbortWithStatus(http.StatusPreconditionFailed)
		return
	}

	var queries struct {
		Code  string `form:"code"`
		State string `form:"state"`
	}
	if err := ctx.ShouldBindQuery(&queries); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if queries.State != state {
		ctx.AbortWithStatus(http.StatusPreconditionFailed)
		return
	}

	token, err := g.config.Exchange(ctx, queries.Code)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		_ = ctx.AbortWithError(http.StatusForbidden, err)
		return
	}
	defer resp.Body.Close()

	var body struct {
		ID    uint   `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user models.User
	if err := g.db.Where("github_id = ?", body.ID).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		user = models.User{
			Name:     body.Login,
			Email:    body.Email,
			GithubID: body.ID,
		}
		if err := g.db.Create(&user).Error; err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	session.Set("user_id", user.ID)
	session.Set("token", token)
	if err := session.Save(); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	u := g.frontendBaseURL
	ctx.Redirect(http.StatusFound, u.String())
}
