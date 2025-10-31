// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	// "errors"
	// "fmt"
	// "strconv"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/infra/mq"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	checkout "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/checkout"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/email"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/payment"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"

	// "github.com/cloudwego/kitex/pkg/klog"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type CheckoutService struct {
	ctx context.Context
} // NewCheckoutService new CheckoutService
func NewCheckoutService(ctx context.Context) *CheckoutService {
	return &CheckoutService{ctx: ctx}
}

/*
	Run

1. get cart
2. calculate cart
3. create order
4. empty cart
5. pay
6. change order result
7. finish
*/
func (s *CheckoutService) Run(req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	// TODO 1.get cart (使用RPC调用Cart服务以获得购物车信息)
	cartResp, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{UserId: req.UserId})
	if err != nil {
		return nil, err
	}
	if cartResp == nil || cartResp.Cart == nil || len(cartResp.Cart.Items) == 0 {
		return &checkout.CheckoutResp{
			OrderId:       "",
			TransactionId: "",
		}, nil
	}

	// TODO 2.calc cart（根据第1步的购物车信息，计算总价和订单项信息）
	var total float64
	var orderItems []*order.OrderItem

	for _, cartItem := range cartResp.Cart.Items {
		productResp, err := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{Id: cartItem.ProductId})
		if err != nil {
			return nil, err
		}
		if productResp == nil || productResp.Product == nil {
			continue
		}
		cost := productResp.Product.Price * float32(cartItem.Quantity)
		total += float64(cost)
		orderItems = append(orderItems, &order.OrderItem{
			Item: cartItem,
			Cost: cost,
		})
	}

	// TODO 3.create order（根据第1步和第2步的信息，创建order.PlaceOrderReq，并使用RPC调用Order服务创建订单）
	orderResp, err := rpc.OrderClient.PlaceOrder(s.ctx, &order.PlaceOrderReq{
		UserId:       req.UserId,
		UserCurrency: "dollar",
		Address: &order.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       10002,
		},
		Email:      req.Email,
		OrderItems: orderItems,
	})
	if err != nil {
		return nil, err
	}
	if orderResp == nil || orderResp.Order == nil {
		return nil, err
	}

	// TODO 4.empty cart（使用RPC调用Cart服务清空购物车）
	_, err = rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	// TODO 5.pay（使用RPC调用Payment服务进行支付）
	paymentResp, err := rpc.PaymentClient.Charge(s.ctx, &payment.ChargeReq{
		Amount:     float32(total),
		CreditCard: req.CreditCard,
		OrderId:    orderResp.Order.OrderId,
		UserId:     req.UserId,
	})
	if err != nil {
		return nil, err
	}

	// TODO 6.send email（使用MQ发送邮件通知）
	data, _ := proto.Marshal(&email.EmailReq{
		From:        "from@example.com",
		To:          req.Email,
		ContentType: "text/plain",
		Subject:     "You just created an order in CloudWeGo shop",
		Content:     "You just created an order in CloudWeGo shop",
	})
	msg := &nats.Msg{Subject: "email", Data: data}
	_ = mq.Nc.PublishMsg(msg)

	// TODO 7.finish（返回订单ID和支付结果）
	resp = &checkout.CheckoutResp{
		OrderId:       orderResp.Order.OrderId,
		TransactionId: paymentResp.TransactionId,
	}
	return
}
