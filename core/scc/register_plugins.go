// +build pluginsenabled,go1.9,linux,cgo
// +build !ppc64le

/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package scc

//RegisterSysCCs is the hook for system chaincodes where system chaincodes are registered with the udo
//note the chaincode must still be deployed and launched like a user chaincode will be
func RegisterSysCCs() {
	systemChaincodes = append(systemChaincodes, loadSysCCs()...)

	for _, sysCC := range systemChaincodes {
		registerSysCC(sysCC)
	}
}
