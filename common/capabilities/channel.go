/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	"github.com/hyperledger/udo/msp"
	cb "github.com/hyperledger/udo/protos/common"
)

const (
	channelTypeName = "Channel"

	// ChannelV1_1 is the capabilties string for standard new non-backwards compatible udo v1.1 channel capabilities.
	ChannelV1_1 = "V1_1"
)

// ChannelProvider provides capabilities information for channel level config.
type ChannelProvider struct {
	*registry
	v11 bool
}

// NewChannelProvider creates a channel capabilities provider.
func NewChannelProvider(capabilities map[string]*cb.Capability) *ChannelProvider {
	cp := &ChannelProvider{}
	cp.registry = newRegistry(cp, capabilities)
	_, cp.v11 = capabilities[ChannelV1_1]
	return cp
}

// Type returns a descriptive string for logging purposes.
func (cp *ChannelProvider) Type() string {
	return channelTypeName
}

// HasCapability returns true if the capability is supported by this binary.
func (cp *ChannelProvider) HasCapability(capability string) bool {
	switch capability {
	// Add new capability names here
	case ChannelV1_1:
		return true
	default:
		return false
	}
}

// MSPVersion returns the level of MSP support required by this channel.
func (cp *ChannelProvider) MSPVersion() msp.MSPVersion {
	switch {
	case cp.v11:
		return msp.MSPv1_1
	default:
		return msp.MSPv1_0
	}
}
