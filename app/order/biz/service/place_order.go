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
	// "fmt"

	"github.com/cloudwego/biz-demo/gomall/app/order/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/order/biz/model"
	order "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"

	// "github.com/cloudwego/kitex/pkg/klog"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlaceOrderService struct {
	ctx context.Context
} // NewPlaceOrderService new PlaceOrderService
func NewPlaceOrderService(ctx context.Context) *PlaceOrderService {
	return &PlaceOrderService{ctx: ctx}
}

// Run create note info
func (s *PlaceOrderService) Run(req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	// TODO 请实现PlaceOrder的业务逻辑，插入数据到数据库中的order表和order_item表，生成一个随机的uuid作为订单号
	// 可以参考其他服务的源代码实现这个函数

	orderId := uuid.NewString()

	// 使用事务确保订单和订单项的原子性
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 先创建订单
		newOrder := &model.Order{
			OrderId:      orderId,
			UserId:       req.UserId,
			UserCurrency: req.UserCurrency,
			Consignee: model.Consignee{
				Email:         req.Email,
				StreetAddress: req.Address.StreetAddress,
				City:          req.Address.City,
				Country:       req.Address.Country,
				ZipCode:       req.Address.ZipCode,
			},
		}

		if err := tx.WithContext(s.ctx).Create(newOrder).Error; err != nil {
			return err
		}

		// 再创建订单项
		for _, v := range req.OrderItems {
			item := &model.OrderItem{
				ProductId:    v.Item.ProductId,
				Cost:         v.Cost,
				Quantity:     v.Item.Quantity,
				OrderIdRefer: orderId,
			}
			if err := tx.WithContext(s.ctx).Create(item).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return &order.PlaceOrderResp{}, err
	}

	resp = &order.PlaceOrderResp{
		Order: &order.OrderResult{
			OrderId: orderId,
		},
	}

	return resp, nil
}
