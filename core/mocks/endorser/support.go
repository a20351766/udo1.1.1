/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package endorser

import (
	"github.com/hyperledger/udo/common/channelconfig"
	"github.com/hyperledger/udo/common/resourcesconfig"
	"github.com/hyperledger/udo/core/ledger"
	mc "github.com/hyperledger/udo/core/mocks/ccprovider"
	"github.com/hyperledger/udo/protos/common"
	pb "github.com/hyperledger/udo/protos/peer"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type MockSupport struct {
	*mock.Mock
	IsSysCCAndNotInvokableExternalRv bool
	IsSysCCRv                        bool
	ExecuteCDSResp                   *pb.Response
	ExecuteCDSEvent                  *pb.ChaincodeEvent
	ExecuteCDSError                  error
	ExecuteResp                      *pb.Response
	ExecuteEvent                     *pb.ChaincodeEvent
	ExecuteError                     error
	ChaincodeDefinitionRv            resourcesconfig.ChaincodeDefinition
	ChaincodeDefinitionError         error
	GetTxSimulatorRv                 *mc.MockTxSim
	GetTxSimulatorErr                error
	CheckInstantiationPolicyError    error
	GetTransactionByIDErr            error
	CheckACLErr                      error
	SysCCMap                         map[string]struct{}
	IsJavaRV                         bool
	IsJavaErr                        error
	GetApplicationConfigRv           channelconfig.Application
	GetApplicationConfigBoolRv       bool
}

func (s *MockSupport) IsSysCCAndNotInvokableExternal(name string) bool {
	return s.IsSysCCAndNotInvokableExternalRv
}

func (s *MockSupport) GetTxSimulator(ledgername string, txid string) (ledger.TxSimulator, error) {
	if s.Mock == nil {
		return s.GetTxSimulatorRv, s.GetTxSimulatorErr
	}

	args := s.Called(ledgername, txid)
	return args.Get(0).(ledger.TxSimulator), args.Error(1)
}

func (s *MockSupport) GetHistoryQueryExecutor(ledgername string) (ledger.HistoryQueryExecutor, error) {
	return nil, nil
}

func (s *MockSupport) GetTransactionByID(chid, txID string) (*pb.ProcessedTransaction, error) {
	return nil, s.GetTransactionByIDErr
}

func (s *MockSupport) IsSysCC(name string) bool {
	if s.SysCCMap != nil {
		_, in := s.SysCCMap[name]
		return in
	}
	return s.IsSysCCRv
}

func (s *MockSupport) Execute(ctxt context.Context, cid, name, version, txid string, syscc bool, signedProp *pb.SignedProposal, prop *pb.Proposal, spec interface{}) (*pb.Response, *pb.ChaincodeEvent, error) {
	if spec != nil {
		if _, istype := spec.(*pb.ChaincodeDeploymentSpec); istype {
			return s.ExecuteCDSResp, s.ExecuteCDSEvent, s.ExecuteCDSError
		}
	}

	return s.ExecuteResp, s.ExecuteEvent, s.ExecuteError
}

func (s *MockSupport) GetChaincodeDefinition(ctx context.Context, chainID string, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, chaincodeID string, txsim ledger.TxSimulator) (resourcesconfig.ChaincodeDefinition, error) {
	return s.ChaincodeDefinitionRv, s.ChaincodeDefinitionError
}

func (s *MockSupport) GetChaincodeDeploymentSpecFS(cds *pb.ChaincodeDeploymentSpec) (*pb.ChaincodeDeploymentSpec, error) {
	return cds, nil
}

func (s *MockSupport) CheckACL(signedProp *pb.SignedProposal, chdr *common.ChannelHeader, shdr *common.SignatureHeader, hdrext *pb.ChaincodeHeaderExtension) error {
	return s.CheckACLErr
}

func (s *MockSupport) IsJavaCC(buf []byte) (bool, error) {
	return s.IsJavaRV, s.IsJavaErr
}

func (s *MockSupport) CheckInstantiationPolicy(name, version string, cd resourcesconfig.ChaincodeDefinition) error {
	return s.CheckInstantiationPolicyError
}

func (s *MockSupport) GetApplicationConfig(cid string) (channelconfig.Application, bool) {
	return s.GetApplicationConfigRv, s.GetApplicationConfigBoolRv
}
