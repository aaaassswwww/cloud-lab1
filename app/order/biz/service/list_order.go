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

	"github.com/cloudwego/biz-demo/gomall/app/order/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/order/biz/model"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	order "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/kitex/pkg/kerrors"
	// "github.com/cloudwego/kitex/pkg/klog"
)

type ListOrderService struct {
	ctx context.Context
} // NewListOrderService new ListOrderService
func NewListOrderService(ctx context.Context) *ListOrderService {
	return &ListOrderService{ctx: ctx}
}

// Run create note info
func (s *ListOrderService) Run(req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	// TODO 请实现ListOrder的业务逻辑，从数据库中的order表和order_item表中查询数据
	// 可以参考其他服务的源代码实现这个函数
	orderlist, err := model.ListOrder(mysql.DB, s.ctx, req.UserId)
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, err.Error())
	}
	var orders []*order.Order
	for _, v := range orderlist {
		var items []*order.OrderItem
		for _, item := range v.OrderItems {
			items = append(items, &order.OrderItem{
				Item: &cart.CartItem{
					ProductId: item.ProductId,
					Quantity: item.Quantity,
				},
				Cost: item.Cost,
			})
		}
		orders = append(orders, &order.Order{
			OrderItems: items,
			OrderId: v.OrderId,
			UserId: v.UserId,
			UserCurrency: v.UserCurrency,
			Address: &order.Address{
				StreetAddress: v.Consignee.StreetAddress,
				City: v.Consignee.City,
				Country: v.Consignee.Country,
				ZipCode: v.Consignee.ZipCode,
			},
			Email: v.Consignee.Email,
		})
	}
	resp = &order.ListOrderResp{
		Orders: orders,
	}
	return resp, nil
}
