package qq

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
)

var (
	clientId = "101413897"
	clientSecret = "aba01a5a90a22ae64c9d0e04ce048ec3"
	grantType = "authorization_code"
	redirectUri = "admin.window.360baige.com"
)

type UserInfoData struct {
	ClientId string
	Openid   string
}

func GetUserInfo(code string) (UserInfoData, error) {
	var userInfo UserInfoData
	getAccessTokenUrl := "https://graph.qq.com/oauth2.0/token?grant_type=" + grantType + "&client_id=" + clientId + "&client_secret=" + clientSecret + "&code=" + code + "&redirect_uri=" + redirectUri
	tokenResponse, err := http.Get(getAccessTokenUrl)
	if err != nil {
		return userInfo, err
	}

	defer tokenResponse.Body.Close()
	tokenBody, err := ioutil.ReadAll(tokenResponse.Body)
	if err != nil {
		return userInfo, err
	}

	var accessTokenData = strings.Split(string(tokenBody), "&")
	if accessTokenData[0] == "" {
		return userInfo, err
	}
	var accessToken = strings.Split(accessTokenData[0], "=")
	if accessToken[1] == "" {
		return userInfo, err
	}
	userInfoUrl := "https://graph.qq.com/oauth2.0/me?access_token=" + accessToken[1]
	infoResponse, err := http.Get(userInfoUrl)
	if err != nil {
		return userInfo, err
	}

	defer infoResponse.Body.Close()
	infoBody, err := ioutil.ReadAll(infoResponse.Body)
	if err != nil {
		return userInfo, err
	}
	var infoData = strings.Fields(string(infoBody))
	err = json.Unmarshal([]byte(infoData[1]), &userInfo)
	if err != nil {
		return userInfo, err
	}
	return userInfo, err
}
