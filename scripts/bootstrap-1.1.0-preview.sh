#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

export VERSION=1.1.0-preview
export ARCH=$(echo "$(uname -s|tr '[:upper:]' '[:lower:]'|sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')
#Set MARCH variable i.e ppc64le,s390x,x86_64,i386
MARCH=`uname -m`

dockerUdoPull() {
  local UDO_TAG=$1
  for IMAGES in peer orderer couchdb ccenv javaenv kafka zookeeper tools; do
      echo "==> UDO IMAGE: $IMAGES"
      echo
      docker pull hyperledger/udo-$IMAGES:$UDO_TAG
      docker tag hyperledger/udo-$IMAGES:$UDO_TAG hyperledger/udo-$IMAGES
  done
}

dockerCaPull() {
      local CA_TAG=$1
      echo "==> UDO CA IMAGE"
      echo
      docker pull hyperledger/fabric-ca:$CA_TAG
      docker tag hyperledger/fabric-ca:$CA_TAG hyperledger/fabric-ca
}

: ${CA_TAG:="$MARCH-$VERSION"}
: ${UDO_TAG:="$MARCH-$VERSION"}

echo "===> Downloading platform binaries"
curl https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/udo/hyperledger-udo/${ARCH}-${VERSION}/hyperledger-udo-${ARCH}-${VERSION}.tar.gz | tar xz

echo "===> Pulling udo Images"
dockerUdoPull ${UDO_TAG}

echo "===> Pulling udo ca Image"
dockerCaPull ${CA_TAG}
echo
echo "===> List out hyperledger docker images"
docker images | grep hyperledger*
