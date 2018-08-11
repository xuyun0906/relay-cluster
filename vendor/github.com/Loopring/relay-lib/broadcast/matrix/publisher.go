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

package matrix

import (
	"fmt"
	"github.com/Loopring/relay-lib/broadcast"
	"github.com/Loopring/relay-lib/log"
)

type MatrixPublisherOption struct {
	MatrixClientOptions
	Rooms []string
}

type MatrixPublisher struct {
	matrixClient *MatrixClient
	Rooms        []string
}

func (publisher *MatrixPublisher) PubOrder(hash string, orderData []byte) error {
	var err error
	for _, room := range publisher.Rooms {
		if eventId, err1 := publisher.matrixClient.SendMessages(room, LoopringOrderType, hash, LoopringOrderType, string(orderData)); nil != err1 {
			if nil == err {
				err = fmt.Errorf("%s:%s", publisher.matrixClient.HSUrl, err1.Error())
			} else {
				err = fmt.Errorf("%s,%s:%s", err.Error(), publisher.matrixClient.HSUrl, err1.Error())
			}
		} else {
			log.Infof("broadcast order:%s in room:%s with eventId:%s", hash, room, eventId)
		}
	}
	return err
}

func (publisher *MatrixPublisher) Name() string {
	return "matrixPublisher"
}

func NewPublishers(options []MatrixPublisherOption) ([]broadcast.Publisher, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("matrixPublisher.options can't be %d", len(options))
	}
	publishers := []broadcast.Publisher{}
	for _, option := range options {
		if matrixClient, err := NewMatrixClient(option.MatrixClientOptions); nil != err {
			return nil, fmt.Errorf("client:%s, err:%s", matrixClient.HSUrl, err.Error())
		} else {
			rooms, err := matrixClient.CheckAndJoinRoom(option.Rooms)
			if nil != err {
				return publishers, err
			}
			joinedRooms := []string{}
			for room, err1 := range rooms {
				if nil == err1 {
					joinedRooms = append(joinedRooms, room)
				} else {
					log.Errorf("a room can't join with error:%s", err1.Error())
				}
			}
			publishers = append(publishers, &MatrixPublisher{
				matrixClient: matrixClient,
				Rooms:        joinedRooms,
			})
		}
	}
	return publishers, nil
}
