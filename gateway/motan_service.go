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

package gateway

import (
	"github.com/Loopring/relay-cluster/accountmanager"
	"github.com/Loopring/relay-cluster/ordermanager/manager"
	"github.com/Loopring/relay-cluster/ordermanager/viewer"
	"github.com/Loopring/relay-lib/motan"
	"math/big"
	"time"
)

type MotanService struct {
	accountManager accountmanager.AccountManager
}

func (s *MotanService) GetBalanceAndAllowance(req *motan.AccountBalanceAndAllowanceReq) *motan.AccountBalanceAndAllowanceRes {
	//start := msecNow()

	res := &motan.AccountBalanceAndAllowanceRes{}
	if balance, allowance, err := accountmanager.GetBalanceAndAllowance(req.Owner, req.Token, req.Spender); nil != err {
		res.Allowance = big.NewInt(int64(0))
		res.Balance = big.NewInt(int64(0))
		res.Err = err.Error()
	} else {
		err := ""
		if nil == balance {
			balance = big.NewInt(int64(0))
			err = err + "balance is nil,"
		}
		if nil == allowance {
			allowance = big.NewInt(int64(0))
			err = err + " allowance is nil"
		}
		res.Balance = new(big.Int).Set(balance)
		res.Allowance = new(big.Int).Set(allowance)
		res.Err = err
	}

	//stop := msecNow()
	//log.Debugf("GetBalanceAndAllowance,owner:%s, token:%s, spender:%s, execute time:%d(msec)", req.Owner.Hex(), req.Token.Hex(), req.Spender.Hex(), stop-start)

	return res
}

func (s *MotanService) GetMinerOrders(req *motan.MinerOrdersReq) *motan.MinerOrdersRes {
	//start := msecNow()

	res := &motan.MinerOrdersRes{}
	res.List = manager.MinerOrders(req.Delegate, req.TokenS, req.TokenB, req.Length, req.ReservedTime, req.StartBlockNumber, req.EndBlockNumber, req.FilterOrderHashLists...)

	//stop := msecNow()
	//log.Debugf("motan service, GetMinerOrders list length:%d, execute time:%d(msec)", len(res.List), stop-start)

	return res
}

func StartMotanService(options motan.MotanServerOptions, accountManager accountmanager.AccountManager, orderViewer viewer.OrderViewer) {
	service := &MotanService{}
	service.accountManager = accountManager
	options.ServerInstance = service
	go motan.RunServer(options)
}

func msecNow() int64 {
	return time.Now().UnixNano() / 1e6
}
