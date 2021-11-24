#!/bin/bash

# 설치 1 -> cli -> peer0.org1.artwork.com
docker exec cli peer chaincode install -n ArtWork -v 0.9 -p github.com/ArtWork 
# linux /home/bstudent/ArtWork/contract/ArtWork -> cli /opt/gopath/src/github.com/ArtWork

docker exec  cli peer chaincode list --installed # 설치된 체인코드 쿼리 -> ID부여된 설치 체인코드이름 버전

# 설치 2 -> cli -> peer0.org2.artwork.com
docker exec -e "CORE_PEER_ADDRESS=peer0.org2.artwork.com:8051" -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.artwork.com/users/Admin@org2.artwork.com/msp" cli peer chaincode install -n ArtWork -v 0.9 -p github.com/ArtWork

docker exec -e "CORE_PEER_ADDRESS=peer0.org2.artwork.com:8051" -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.artwork.com/users/Admin@org2.artwork.com/msp" cli peer chaincode list --installed

# 배포 peer0.org1.artwork.com  -> dev-ArtWork 인도서 피어 컨테이너가 생성, 커미터 피어 couchdb myart_papercontract 테이블이생성
docker exec cli peer chaincode instantiate -n ArtWork -v 0.9 -C myart -c '{"Args":[]}' -P 'AND ("Org1MSP.member","Org2MSP.member")' 
# 체인코드 같은 이름으로 배포 -> upgrade
sleep 3

# 배포 확인 
docker exec cli peer chaincode list --instantiated -C myart