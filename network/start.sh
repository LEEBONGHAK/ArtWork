#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -ev

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

docker-compose -f docker-compose.yml down

# docker-compose -> 컨테이터수행 및 net_basic 네트워크 생성
docker-compose -f docker-compose.yml up -d ca.org1.artwork.com ca.org2.artwork.com orderer.artwork.com peer0.org1.artwork.com  peer0.org2.artwork.com cli
docker ps -a
docker network ls
# wait for Hyperledger Fabric to start
# incase of errors when running later commands, issue export FABRIC_START_TIMEOUT=<larger number>
export FABRIC_START_TIMEOUT=5
#echo ${FABRIC_START_TIMEOUT}
sleep ${FABRIC_START_TIMEOUT}

# Create the channel -> myart.block cli working dir 복사
docker exec cli peer channel create -o orderer.artwork.com:7050 -c myart -f /etc/hyperledger/configtx/channel.tx
# clie workding dir (/etc/hyperledger/configtx/) myart.block

# Join peer0.org1.artwork.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.artwork.com/msp" peer0.org1.artwork.com peer channel join -b /etc/hyperledger/configtx/myart.block
sleep 3

# Join peer0.org2.artwork.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org2.artwork.com/msp" peer0.org2.artwork.com peer channel join -b /etc/hyperledger/configtx/myart.block
sleep 3

# anchor ORG1 myart update
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.artwork.com/msp" peer0.org1.artwork.com peer channel update -f /etc/hyperledger/configtx/Org1MSPanchors.tx -c myart -o orderer.artwork.com:7050
# anchor ORG2 myart update
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org2.artwork.com/msp" peer0.org2.artwork.com peer channel update -f /etc/hyperledger/configtx/Org2MSPanchors.tx -c myart -o orderer.artwork.com:7050
