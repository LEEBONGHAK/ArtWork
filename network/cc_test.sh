#!/bin/bash

set -x

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["addUser", "myCompany"]}' --peerAddresses peer0.org1.artwork.com:7051 peer0.org1.artwork.com:8051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["addUser", "u1"]}' --peerAddresses peer0.org1.artwork.com:7051 peer0.org1.artwork.com:8051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["addUser", "u2"]}' --peerAddresses peer0.org1.artwork.com:7051 peer0.org1.artwork.com:8051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["addWork", "a1", "title", "m1", "100", "100"]}' --peerAddresses peer0.org1.artwork.com:7051 peer0.org1.artwork.com:8051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["tradeProps", "myCompany", "u1", "a1", "30"]}' --peerAddresses peer0.org1.artwork.com:7051 peer0.org1.artwork.com:8051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["tradeProps", "u1", "u2", "a1", "30"]}' --peerAddresses peer0.org1.artwork.com:7051 peer0.org1.artwork.com:8051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["endTradeProps", "a1", "1000"]}' --peerAddresses peer0.org1.artwork.com:7051 peer0.org1.artwork.com:8051
sleep 3

docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getHistory", "u1"]}'
docker exec cli peer chaincode invoke -n ArtWork -C myart -c '{"Args":["getHistory", "a1"]}'