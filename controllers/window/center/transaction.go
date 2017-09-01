package center

import (
	"github.com/astaxie/beego"
	//"dev.model.360baige.com/action"
	"dev.model.360baige.com/models/account"
	"dev.cloud.360baige.com/utils"
	"dev.cloud.360baige.com/rpc/client"
	"fmt"
)

type TransactionController struct {
	beego.Controller
}

// 充值
func Recharge() {
	// 运营商账号- ￥
	// 充值账号+ ￥
}

// 消费
func Consume() {
	// 购买账号+ ￥
	// 运营商账号+ ￥
}

// 交易
func Transaction() {

}

func AddTransaction(fromAccountId, toAccountId, amount int64, orderCode, remark string) (account.Transaction, error) {
	currentTimestamp := utils.CurrentTimestamp()
	var replyTransaction account.Transaction
	args := account.Transaction{
		CreateTime:    currentTimestamp,
		UpdateTime:    currentTimestamp,
		FromAccountId: fromAccountId,
		ToAccountId:   toAccountId,
		Amount:        amount,
		OrderCode:     orderCode,
		Remark:        remark,
		Status:        0,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Transaction", "Add", args, &replyTransaction)
	return replyTransaction, err
}

func AddAccountItem() {
	currentTimestamp := utils.CurrentTimestamp()
	var replyAccountItem account.AccountItem
	args := account.AccountItem{
		CreateTime: currentTimestamp,
		UpdateTime: currentTimestamp,
		//Amount:     amount,
		//Remark:     remark,
		Status:     0,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "Add", args, &replyAccountItem)
	fmt.Println(err)
}
