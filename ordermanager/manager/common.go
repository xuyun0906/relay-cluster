/*

  Copyright 2017 Loopring Project Ltd (Loopring Foundation).

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package manager

import (
	"fmt"
	"github.com/Loopring/relay-cluster/dao"
	"github.com/Loopring/relay-lib/eth/loopringaccessor"
	"github.com/Loopring/relay-lib/types"
	"math/big"
)

func NewOrderEntity(state *types.OrderState, blockNumber *big.Int) (*dao.Order, error) {
	state.DealtAmountS = big.NewInt(0)
	state.DealtAmountB = big.NewInt(0)
	state.SplitAmountS = big.NewInt(0)
	state.SplitAmountB = big.NewInt(0)
	state.CancelledAmountB = big.NewInt(0)
	state.CancelledAmountS = big.NewInt(0)

	if blockNumber == nil {
		state.UpdatedBlock = big.NewInt(0)
	} else {
		state.UpdatedBlock = blockNumber
	}

	// calculate order amount and settled
	SettleOrderAmountOnChain(state)

	// check order finished status
	SettleOrderStatus(state, false)

	// convert order
	model := &dao.Order{}
	model.ConvertDown(state)

	return model, nil
}

func SettleOrderAmountOnChain(state *types.OrderState) error {
	var (
		cancelled, cancelOrFilled, dealt *big.Int
		err                              error
	)

	protocol := state.RawOrder.DelegateAddress
	orderhash := state.RawOrder.Hash

	// get order cancelled amount from chain
	if cancelled, err = loopringaccessor.GetCancelled(protocol, orderhash, "latest"); err != nil {
		return fmt.Errorf("order manager,handle gateway order,order %s getCancelled error:%s", orderhash.Hex(), err.Error())
	}

	// get order cancelledOrFilled amount from chain
	if cancelOrFilled, err = loopringaccessor.GetCancelledOrFilled(protocol, orderhash, "latest"); err != nil {
		return fmt.Errorf("order manager,handle gateway order,order %s getCancelledOrFilled error:%s", orderhash.Hex(), err.Error())
	}

	if cancelOrFilled.Cmp(cancelled) < 0 {
		return fmt.Errorf("order manager,handle gateway order,order %s cancelOrFilledAmount:%s < cancelledAmount:%s", orderhash.Hex(), cancelOrFilled.String(), cancelled.String())
	}

	dealt = big.NewInt(0).Sub(cancelOrFilled, cancelled)

	if state.RawOrder.BuyNoMoreThanAmountB {
		state.DealtAmountB = dealt
		state.CancelledAmountB = cancelled
	} else {
		state.DealtAmountS = dealt
		state.CancelledAmountS = cancelled
	}

	return nil
}

func SettleOrderStatus(state *types.OrderState, isCancel bool) {
	zero := big.NewInt(0)
	finishAmountS := big.NewInt(0).Add(state.CancelledAmountS, state.DealtAmountS)
	totalAmountS := big.NewInt(0).Add(finishAmountS, state.SplitAmountS)
	finishAmountB := big.NewInt(0).Add(state.CancelledAmountB, state.DealtAmountB)
	totalAmountB := big.NewInt(0).Add(finishAmountB, state.SplitAmountB)
	totalAmount := big.NewInt(0).Add(totalAmountS, totalAmountB)

	if totalAmount.Cmp(zero) <= 0 {
		state.Status = types.ORDER_NEW
	} else if !marketCapProvider.IsOrderValueDust(state) {
		state.Status = types.ORDER_PARTIAL
	} else if isCancel {
		state.Status = types.ORDER_CANCEL
	} else {
		state.Status = types.ORDER_FINISHED
	}
}

// 记录已经存在的情况下 判断是否重复
func IsEventDuplicate(eventStatus types.TxStatus, modelStatus uint8) bool {
	if eventStatus == types.TX_STATUS_PENDING {
		return true
	}
	if eventStatus != types.TX_STATUS_PENDING && modelStatus != uint8(types.TX_STATUS_PENDING) {
		return true
	}
	return false
}

func IsDustyOrder(state *types.OrderState) bool {
	return marketCapProvider.IsOrderValueDust(state)
}
