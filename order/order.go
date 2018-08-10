package order

import (
	"github.com/pkg/errors"
	"fotomxq/gobase/log"
	"fotomxq/gobase/distribution"
)

//订单服务
// 用户自由创建订单，自动关闭过期的订单
// 每个用户每天只能创建30个订单，超过该数字后拒绝生成订单

var(
	//状态
	ORDER_STATUS_WAIT string = "wait" //等待用户支付的订单
	ORDER_STATUS_FINISH string = "finish" //已经完成的订单
	ORDER_STATUS_TRASH string = "trash" //作废的订单
)

//自动维护工具
// 关闭超过3小时的订单
func Auto(){
	//更新子服务项目 OrderType.Run
	err := distribution.UpdateSubTaskBySelf("OrderType.Run")
	if err != nil{
		log.SendError("0.0.0.0","","OrderType.Run",err)
	}
}

//查询订单列表
func Find() ([]FieldsOrder,error){
	res := []FieldsOrder{}
	return res,errors.New("unable to query order list.")
}

//生成订单并展示付款二维码
//param createUserInfo FieldsUser.UserFields 创建订单的用户信息组
//param userInfo FieldsUser 订单需求的用户信息组
//param groupInfo FieldsUserGroup 服务信息
//param serviceUnit int64
//return FieldsOrder.OrderFields 创建订单
//return error 错误代码

//取消订单