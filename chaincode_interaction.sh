#!/usr/bin/env sh
#
# SPDX-License-Identifier: Apache-2.0
#

# . peer1admin.sh

{
  peer chaincode query -C mychannel -n chaincode1 -c '{"Args":["InitLedger"]}' --waitForEvent --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
  # query all proposol
  peer chaincode query -C mychannel -n chaincode1 -c '{"Args":["QueryAllProposal"]}'
    
  # create proposal
  peer chaincode invoke -o 127.0.0.1:6050 -C mychannel -n chaincode1 -c '{"Args":["CreateProposal","proposal1","Akshay","999"]}' --waitForEvent --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
  
  # approve proposal
  peer chaincode invoke -o 127.0.0.1:6050 -C mychannel -n chaincode1 -c '{"Args":["ApproveProposal","proposal1"]}' --waitForEvent --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
  
  # setCIBILTrack
  peer chaincode invoke -o 127.0.0.1:6050 -C mychannel -n chaincode1 -c '{"Args":["SetCIBILTrack", "proposal1", "true", "true"]}' --waitForEvent --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
  
  # approve proposal
  peer chaincode invoke -o 127.0.0.1:6050 -C mychannel -n chaincode1 -c '{"Args":["ApproveProposal","proposal1"]}' --waitForEvent --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
  
  # query proposal
  peer chaincode query -C mychannel -n chaincode1 -c '{"Args":["QueryProposal","proposal1"]}'
  
  # getHistoryOf proposal
  peer chaincode query -C mychannel -n chaincode1 -c '{"Args":["GetHistoryOfProposal","proposal1"]}'
}
