package wechat

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var (
	appid = "wxe8d941078f9472af"
	secret = "f6b34820fde984be4a0cba8e3315c3ce"
	grantType = "authorization_code"
)

type AccessTokenData struct {
	Access_token  string
	Expires_in    int64
	Refresh_token string
	Openid        string
	Scope         string
	Unionid       string
}

type UserInfoData struct {
	Openid     string
	Nickname   string
	Sex        int
	Province   string
	City       string
	Country    string
	Headimgurl string
	Privilege  interface{}
	Unionid    string
}

func GetUserInfo(code string) (UserInfoData, error) {
	var userInfo UserInfoData
	getAccessTokenUrl := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + appid + "&secret=" + secret + "&code=" + code + "&grant_type=" + grantType
	tokenResponse, err := http.Get(getAccessTokenUrl)
	if err != nil {
		return userInfo, err
	}

	defer tokenResponse.Body.Close()
	tokenBody, err := ioutil.ReadAll(tokenResponse.Body)
	if err != nil {
		return userInfo, err
	}

	var accessTokenData AccessTokenData
	err = json.Unmarshal(tokenBody, &accessTokenData)
	if err != nil {
		return userInfo, err
	}

	userInfoUrl := "https://api.weixin.qq.com/sns/userinfo?access_token=" + accessTokenData.Access_token + "&openid=" + accessTokenData.Openid
	infoResponse, err := http.Get(userInfoUrl)
	if err != nil {
		return userInfo, err
	}

	defer infoResponse.Body.Close()
	infoBody, err := ioutil.ReadAll(infoResponse.Body)
	if err != nil {
		return userInfo, err
	}

	err = json.Unmarshal(infoBody, &userInfo)
	if err != nil {
		return userInfo, err
	}
	return userInfo, err
}
