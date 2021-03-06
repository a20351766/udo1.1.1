#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# if version not passed in, default to latest released version
export VERSION=${1:-1.1.0}
# if ca version not passed in, default to latest released version
export CA_VERSION=${2:-$VERSION}
# current version of thirdparty images (couchdb, kafka and zookeeper) released
export THIRDPARTY_IMAGE_VERSION=0.4.6
export ARCH=$(echo "$(uname -s|tr '[:upper:]' '[:lower:]'|sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')
#Set MARCH variable i.e ppc64le,s390x,x86_64,i386
MARCH=`uname -m`

dockerUdoPull() {
  local UDO_TAG=$1
  for IMAGES in peer orderer ccenv javaenv tools; do
      echo "==> UDO IMAGE: $IMAGES"
      echo
      docker pull hyperledger/udo-$IMAGES:$UDO_TAG
      docker tag hyperledger/udo-$IMAGES:$UDO_TAG hyperledger/udo-$IMAGES
  done
}

dockerThirdPartyImagesPull() {
  local THIRDPARTY_TAG=$1
  for IMAGES in couchdb kafka zookeeper; do
      echo "==> THIRDPARTY DOCKER IMAGE: $IMAGES"
      echo
      docker pull hyperledger/udo-$IMAGES:$THIRDPARTY_TAG
      docker tag hyperledger/udo-$IMAGES:$THIRDPARTY_TAG hyperledger/udo-$IMAGES
  done
}

dockerCaPull() {
      local CA_TAG=$1
      echo "==> UDO CA IMAGE"
      echo
      docker pull hyperledger/fabric-ca:$CA_TAG
      docker tag hyperledger/fabric-ca:$CA_TAG hyperledger/fabric-ca
}

: ${CA_TAG:="$MARCH-$CA_VERSION"}
: ${UDO_TAG:="$MARCH-$VERSION"}
: ${THIRDPARTY_TAG:="$MARCH-$THIRDPARTY_IMAGE_VERSION"}

echo "===> Downloading platform specific udo binaries"
curl https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/udo/hyperledger-udo/${ARCH}-${VERSION}/hyperledger-udo-${ARCH}-${VERSION}.tar.gz | tar xz

echo "===> Downloading platform specific fabric-ca-client binary"
curl https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric-ca/hyperledger-fabric-ca/${ARCH}-${VERSION}/hyperledger-fabric-ca-${ARCH}-${VERSION}.tar.gz | tar xz
if [ $? != 0 ]; then
     echo
     echo "------> $VERSION fabric-ca-client binary is not available to download  (Avaialble from 1.1.0-rc1) <----"
     echo
fi
echo "===> Pulling udo Images"
dockerUdoPull ${UDO_TAG}

echo "===> Pulling udo ca Image"
dockerCaPull ${CA_TAG}

echo "===> Pulling thirdparty docker images"
dockerThirdPartyImagesPull ${THIRDPARTY_TAG}

echo
echo "===> List out hyperledger docker images"
docker images | grep hyperledger*
