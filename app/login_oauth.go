package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

func (a *App) AuthenticateTokenForLogin(token string, service string) (user *model.User, err *model.AppError) {
	if len(token) == 0 {
		return nil, model.NewAppError("createSessionForUserAccessToken", "app.user_access_token.invalid_or_missing", nil, "", http.StatusBadRequest)
	}

	sso := *a.Config().GetOAuthProvider(service)
	var url = *(sso.UserApiEndpoint)
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + token

	// Create a new request using http
	req, _ := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, httperr := client.Do(req)

	if httperr != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, model.NewAppError("createSessionForUserAccessToken", "app.user_access_token.invalid_or_missing", nil, "", http.StatusForbidden)
	}

	jsonDataFromHTTP, _ := ioutil.ReadAll(resp.Body)
	log.Println(string([]byte(jsonDataFromHTTP)))

	defer resp.Body.Close()

	jsonData := map[string]string{}
	json.Unmarshal([]byte(jsonDataFromHTTP), &jsonData)

	// if jsonerr != nil {
	// 	log.Println("Error on unmarshall.\n[ERRO] -", err)
	// }
	email := jsonData["email"]
	if sso.EmailField != nil {
		emailField := *(sso.EmailField)
		email = jsonData[emailField]
	}
	name := jsonData["name"]
	if sso.NameField != nil {
		nameField := *(sso.NameField)
		name = jsonData[nameField]
	}

	user, err = a.GetUserForLogin("", email)
	split := strings.Split(email, "@")
	trusted := split[1] == *(sso.TrustedDomain)
	fmt.Println(name, email)
	if err != nil && trusted {
		user := model.User{Email: email, Nickname: name, Username: model.NewId(), Roles: model.SYSTEM_USER_ROLE_ID} //Username: GenerateTestUsername(), Roles: model.SYSTEM_USER_ROLE_ID}
		fmt.Println(user)
		var ruser *model.User
		ruser, err = a.CreateUserFromSignup(&user)
		fmt.Println(ruser)
		fmt.Println(err)
		return ruser, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}
