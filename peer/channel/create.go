/*
Copyright IBM Corp. 2017 All Rights Reserved.

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

package channel

import (
	"fmt"
	"io/ioutil"

	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/udo/common/configtx"
	localsigner "github.com/hyperledger/udo/common/localmsp"
	"github.com/hyperledger/udo/common/tools/configtxgen/encoder"
	genesisconfig "github.com/hyperledger/udo/common/tools/configtxgen/localconfig"
	"github.com/hyperledger/udo/common/util"
	"github.com/hyperledger/udo/peer/common"
	cb "github.com/hyperledger/udo/protos/common"
	"github.com/hyperledger/udo/protos/utils"
	"github.com/spf13/cobra"
)

//ConfigTxFileNotFound channel create configuration tx file not found
type ConfigTxFileNotFound string

const createCmdDescription = "Create a channel"

func (e ConfigTxFileNotFound) Error() string {
	return fmt.Sprintf("channel create configuration tx file not found %s", string(e))
}

//InvalidCreateTx invalid channel create transaction
type InvalidCreateTx string

func (e InvalidCreateTx) Error() string {
	return fmt.Sprintf("Invalid channel create transaction : %s", string(e))
}

func createCmd(cf *ChannelCmdFactory) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: createCmdDescription,
		Long:  createCmdDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd, args, cf)
		},
	}
	flagList := []string{
		"channelID",
		"file",
		"timeout",
	}
	attachFlags(createCmd, flagList)

	return createCmd
}

func createChannelFromDefaults(cf *ChannelCmdFactory) (*cb.Envelope, error) {
	chCrtEnv, err := encoder.MakeChannelCreationTransaction(channelID, localsigner.NewSigner(), nil, genesisconfig.Load(genesisconfig.SampleSingleMSPChannelProfile))
	if err != nil {
		return nil, err
	}

	return chCrtEnv, nil
}

func createChannelFromConfigTx(configTxFileName string) (*cb.Envelope, error) {
	cftx, err := ioutil.ReadFile(configTxFileName)
	if err != nil {
		return nil, ConfigTxFileNotFound(err.Error())
	}

	return utils.UnmarshalEnvelope(cftx)
}

func sanityCheckAndSignConfigTx(envConfigUpdate *cb.Envelope) (*cb.Envelope, error) {
	payload, err := utils.ExtractPayload(envConfigUpdate)
	if err != nil {
		return nil, InvalidCreateTx("bad payload")
	}

	if payload.Header == nil || payload.Header.ChannelHeader == nil {
		return nil, InvalidCreateTx("bad header")
	}

	ch, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, InvalidCreateTx("could not unmarshall channel header")
	}

	if ch.Type != int32(cb.HeaderType_CONFIG_UPDATE) {
		return nil, InvalidCreateTx("bad type")
	}

	if ch.ChannelId == "" {
		return nil, InvalidCreateTx("empty channel id")
	}

	// Specifying the chainID on the CLI is usually redundant, as a hack, set it
	// here if it has not been set explicitly
	if channelID == "" {
		channelID = ch.ChannelId
	}

	if ch.ChannelId != channelID {
		return nil, InvalidCreateTx(fmt.Sprintf("mismatched channel ID %s != %s", ch.ChannelId, channelID))
	}

	configUpdateEnv, err := configtx.UnmarshalConfigUpdateEnvelope(payload.Data)
	if err != nil {
		return nil, InvalidCreateTx("Bad config update env")
	}

	signer := localsigner.NewSigner()
	sigHeader, err := signer.NewSignatureHeader()
	if err != nil {
		return nil, err
	}

	configSig := &cb.ConfigSignature{
		SignatureHeader: utils.MarshalOrPanic(sigHeader),
	}

	configSig.Signature, err = signer.Sign(util.ConcatenateBytes(configSig.SignatureHeader, configUpdateEnv.ConfigUpdate))

	configUpdateEnv.Signatures = append(configUpdateEnv.Signatures, configSig)

	return utils.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, channelID, signer, configUpdateEnv, 0, 0)
}

func sendCreateChainTransaction(cf *ChannelCmdFactory) error {
	var err error
	var chCrtEnv *cb.Envelope

	if channelTxFile != "" {
		if chCrtEnv, err = createChannelFromConfigTx(channelTxFile); err != nil {
			return err
		}
	} else {
		if chCrtEnv, err = createChannelFromDefaults(cf); err != nil {
			return err
		}
	}

	if chCrtEnv, err = sanityCheckAndSignConfigTx(chCrtEnv); err != nil {
		return err
	}

	var broadcastClient common.BroadcastClient
	broadcastClient, err = cf.BroadcastFactory()
	if err != nil {
		return fmt.Errorf("Error getting broadcast client: %s", err)
	}

	defer broadcastClient.Close()
	err = broadcastClient.Send(chCrtEnv)

	return err
}

func executeCreate(cf *ChannelCmdFactory) error {
	var err error

	if err = sendCreateChainTransaction(cf); err != nil {
		return err
	}

	var block *cb.Block
	if block, err = getGenesisBlock(cf); err != nil {
		return err
	}

	b, err := proto.Marshal(block)
	if err != nil {
		return err
	}

	file := channelID + ".block"
	if err = ioutil.WriteFile(file, b, 0644); err != nil {
		return err
	}

	return nil
}

func create(cmd *cobra.Command, args []string, cf *ChannelCmdFactory) error {
	//the global chainID filled by the "-c" command
	if channelID == common.UndefinedParamValue {
		return errors.New("Must supply channel ID")
	}

	var err error
	if cf == nil {
		cf, err = InitCmdFactory(EndorserNotRequired, OrdererRequired)
		if err != nil {
			return err
		}
	}
	return executeCreate(cf)
}
