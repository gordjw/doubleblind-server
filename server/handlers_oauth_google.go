package server

import (
	"context"
	"crypto/rand"
	"doubleblind/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
const oauthGoogleTokenEndpoint = "https://oauth2.googleapis.com/token"

func (env *Env) oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := env.generateStateOauthCookie(w)

	// TODO - These are not set on initialisation, investigate why
	googleOauthConfig.ClientID = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	googleOauthConfig.ClientSecret = os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")

	url := googleOauthConfig.AuthCodeURL(oauthState, oauth2.ApprovalForce, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (env *Env) oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	user, err := env.getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	env.users.GetOrCreateUser(user)

	cookie := http.Cookie{Name: "user_id", Value: user.OauthId, Path: "/"}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (env *Env) generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (env *Env) getUserDataFromGoogle(code string) (models.User, error) {

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return models.User{}, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return models.User{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	var user models.User
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("failed decoding user info: %s", err.Error())
	}
	user.OauthProvider = "google"

	// If a refresh token was provided, update the session storage
	if token.RefreshToken != "" {
		env.sessions[user.OauthId] = token.RefreshToken
	}

	return user, nil
}

func (env *Env) middlewareOauth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie("user_id")
		if err != nil {
			fmt.Printf("error reading auth cookie: %v\n", err.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}
		var refreshToken = env.sessions[authCookie.Value]
		if refreshToken == "" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}

		postData := url.Values{}
		postData.Add("client_id", os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
		postData.Add("client_secret", os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
		postData.Add("refresh_token", refreshToken)
		postData.Add("grant_type", "refresh_token")

		response, err := http.PostForm(
			oauthGoogleTokenEndpoint,
			postData,
		)
		if err != nil {
			fmt.Printf("failed getting user info: %s\n", err.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}
		defer response.Body.Close()

		token := struct {
			AccessToken string `json:"access_token"`
		}{}
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&token)
		if err != nil {
			fmt.Printf("failed decoding user info: %s\n", err.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// fmt.Printf("bearerToken %v returned accessToken %v\n", authCookie.Value, token.AccessToken)

		// If we made it here, the user is authenticated!
		// Set the userId in the request context for other middleware and handlers to use
		ctx := context.WithValue(r.Context(), ContextKeyAuthToken, authCookie.Value)
		r = r.WithContext(ctx)

		// And update the client-side cookie
		cookie := http.Cookie{Name: "user_id", Value: authCookie.Value, Path: "/"}
		http.SetCookie(w, &cookie)

		next.ServeHTTP(w, r)
	})
}
