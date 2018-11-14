#!/bin/bash -eu
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#


##################################################
# This script pulls docker images from hyperledger
# docker hub repository and Tag it as
# hyperledger/udo-<image> latest tag
##################################################

#Set ARCH variable i.e ppc64le,s390x,x86_64,i386
ARCH=`uname -m`

dockerUdoPull() {
  local UDO_TAG=$1
  for IMAGES in peer orderer couchdb ccenv javaenv kafka tools zookeeper; do
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
usage() {
      echo "Description "
      echo
      echo "Pulls docker images from hyperledger dockerhub repository"
      echo "tag as hyperledger/udo-<image>:latest"
      echo
      echo "USAGE: "
      echo
      echo "./download-dockerimages.sh [-c <fabric-ca tag>] [-f <udo tag>]"
      echo "      -c fabric-ca docker image tag"
      echo "      -f udo docker image tag"
      echo
      echo
      echo "EXAMPLE:"
      echo "./download-dockerimages.sh -c x86_64-1.0.0-beta -f x86_64-1.0.0-beta"
      echo
      echo "By default, pulls fabric-ca and udo 1.0.0-beta docker images"
      echo "from hyperledger dockerhub"
      exit 0
}

while getopts "\?hc:f:" opt; do
  case "$opt" in
     c) CA_TAG="$OPTARG"
        echo "Pull CA IMAGES"
        ;;

     f) UDO_TAG="$OPTARG"
        echo "Pull UDO TAG"
        ;;
     \?|h) usage
        echo "Print Usage"
        ;;
  esac
done

: ${CA_TAG:="$ARCH-1.0.0-beta"}
: ${UDO_TAG:="$ARCH-1.0.0-beta"}

echo "===> Pulling udo Images"
dockerUdoPull ${UDO_TAG}

echo "===> Pulling udo ca Image"
dockerCaPull ${CA_TAG}
echo
echo "===> List out hyperledger docker images"
docker images | grep hyperledger*
