package v1

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"test/internal/config"
	"test/internal/service/dto"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauth2Config = config.GetConfig().Oauth2Config

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  oauth2Config.RedirectURL,
	ClientID:     oauth2Config.ClientID,
	ClientSecret: oauth2Config.ClientSecret,
	Scopes:       oauth2Config.Scopes,
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v3/userinfo?access_token="
const coockieOauth2StateName = "oauthstate"

func OauthGoogleLogin(c *gin.Context) {

	oauthState := generateStateOauthCookie(c.Writer)
	u := googleOauthConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: coockieOauth2StateName, Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}


func OauthGoogleCallback(ctx *gin.Context) {
	oauthState, _ := ctx.Request.Cookie(coockieOauth2StateName)
	if ctx.Request.FormValue("state") != oauthState.Value {
		newResponse(ctx, http.StatusBadRequest, "invalid oauth google state")
		return
	}

	code := ctx.Request.FormValue("code")
	token, err := googleOauthConfig.Exchange(ctx.Request.Context(), code)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "failed to get google access token throught code")
		return
	}

	resp, err := http.Get(oauthGoogleUrlAPI + url.QueryEscape(token.AccessToken))
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "failed to requst google userInfo")
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		newResponse(ctx, http.StatusBadRequest, "failed to read body")
		return
	}

	// TODO: do something with access token
	var userDTO dto.CreateUserDTO
	if err = json.Unmarshal(response, &userDTO); err != nil {
		newResponse(ctx, http.StatusBadRequest, "failed to unmarshal google user")
		return
	}
	userDTO.Password = token.AccessToken
}
