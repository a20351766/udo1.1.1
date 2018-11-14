/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"fmt"

	"github.com/pkg/errors"

	"encoding/json"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/udo/core/scc/qscc"
	"github.com/hyperledger/udo/peer/common"
	cb "github.com/hyperledger/udo/protos/common"
	pb "github.com/hyperledger/udo/protos/peer"
	"github.com/hyperledger/udo/protos/utils"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func getinfoCmd(cf *ChannelCmdFactory) *cobra.Command {
	getinfoCmd := &cobra.Command{
		Use:   "getinfo",
		Short: "get blockchain information of a specified channel.",
		Long:  "get blockchain information of a specified channel. Requires '-c'.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getinfo(cf)
		},
	}
	flagList := []string{
		"channelID",
	}
	attachFlags(getinfoCmd, flagList)

	return getinfoCmd
}
func (cc *endorserClient) getBlockChainInfo() (*cb.BlockchainInfo, error) {
	var err error

	invocation := &pb.ChaincodeInvocationSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type:        pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]),
			ChaincodeId: &pb.ChaincodeID{Name: "qscc"},
			Input:       &pb.ChaincodeInput{Args: [][]byte{[]byte(qscc.GetChainInfo), []byte(channelID)}},
		},
	}

	var prop *pb.Proposal
	c, _ := cc.cf.Signer.Serialize()
	prop, _, err = utils.CreateProposalFromCIS(cb.HeaderType_ENDORSER_TRANSACTION, "", invocation, c)
	if err != nil {
		return nil, errors.WithMessage(err, "cannot create proposal")
	}

	var signedProp *pb.SignedProposal
	signedProp, err = utils.GetSignedProposal(prop, cc.cf.Signer)
	if err != nil {
		return nil, errors.WithMessage(err, "cannot create signed proposal")
	}

	proposalResp, err := cc.cf.EndorserClient.ProcessProposal(context.Background(), signedProp)
	if err != nil {
		return nil, errors.WithMessage(err, "failed sending proposal")
	}

	if proposalResp.Response == nil || proposalResp.Response.Status != 200 {
		return nil, errors.Errorf("received bad response, status %d", proposalResp.Response.Status)
	}

	blockChainInfo := &cb.BlockchainInfo{}
	err = proto.Unmarshal(proposalResp.Response.Payload, blockChainInfo)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read qscc response")
	}

	return blockChainInfo, nil

}

func getinfo(cf *ChannelCmdFactory) error {
	//the global chainID filled by the "-c" command
	if channelID == common.UndefinedParamValue {
		return errors.New("Must supply channel ID")
	}

	var err error
	if cf == nil {
		cf, err = InitCmdFactory(EndorserRequired, OrdererNotRequired)
		if err != nil {
			return err
		}
	}

	client := &endorserClient{cf}

	blockChainInfo, err := client.getBlockChainInfo()
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(blockChainInfo)
	if err != nil {
		return err
	}

	fmt.Printf("Blockchain info: %s\n", string(jsonBytes))

	return nil
}
