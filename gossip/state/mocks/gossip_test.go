/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	"testing"

	"github.com/hyperledger/udo/gossip/api"
	"github.com/hyperledger/udo/gossip/common"
	"github.com/hyperledger/udo/gossip/discovery"
	proto "github.com/hyperledger/udo/protos/gossip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGossipMock(t *testing.T) {
	g := GossipMock{}
	mkChan := func() <-chan *proto.GossipMessage {
		c := make(chan *proto.GossipMessage, 1)
		c <- &proto.GossipMessage{}
		return c
	}
	g.On("Accept", mock.Anything, false).Return(mkChan(), nil)
	a, b := g.Accept(func(o interface{}) bool {
		return true
	}, false)
	assert.Nil(t, b)
	assert.NotNil(t, a)
	assert.Panics(t, func() {
		g.SuspectPeers(func(identity api.PeerIdentityType) bool { return false })
	})
	assert.Panics(t, func() {
		g.Send(nil, nil)
	})
	assert.Panics(t, func() {
		g.Peers()
	})
	g.On("PeersOfChannel", mock.Anything).Return([]discovery.NetworkMember{})
	assert.Empty(t, g.PeersOfChannel(common.ChainID("A")))

	assert.Panics(t, func() {
		g.UpdateMetadata([]byte{})
	})
	assert.Panics(t, func() {
		g.Gossip(nil)
	})
	assert.NotPanics(t, func() {
		g.UpdateChannelMetadata([]byte{}, common.ChainID("A"))
		g.Stop()
		g.JoinChan(nil, common.ChainID("A"))
	})
}
