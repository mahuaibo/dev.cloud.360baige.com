package utils

import (
	"github.com/astaxie/beego"
	"dev.model.360baige.com/action"
	"dev.model.360baige.com/models/user"
	"dev.cloud.360baige.com/rpc/client"
	"errors"
)

func UserPosition(accessToken string, currentTimestamp int64) (*user.UserPosition, error) {
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_token", Val: accessToken },
			action.CondValue{Type: "Or", Key: "transit_token", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type", "access_token", "expire_in", "transit_token", "transit_expire_in"},
	}, &replyUserPosition)

	if err != nil {
		return nil, errors.New("系统异常，请稍后重试")
	}

	if accessToken == replyUserPosition.AccessToken && currentTimestamp <= replyUserPosition.ExpireIn {
		return &replyUserPosition, nil
	} else if accessToken == replyUserPosition.TransitToken && replyUserPosition.TransitExpireIn-user.UserPositionTransitExpireIn <= currentTimestamp && currentTimestamp <= replyUserPosition.TransitExpireIn+user.UserPositionTransitExpireIn {
		return &replyUserPosition, nil
	}
	return nil, errors.New("访问令牌无效，请重新登录")
}
