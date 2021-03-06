#!/bin/bash

#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#
# Test Java SDK e2e tests
#

WD="${GOPATH}/src/github.com/hyperledger/udo-sdk-java"
#WD="${WORKSPACE}/gopath/src/github.com/hyperledger/udo-sdk-java"
SDK_REPO_NAME=udo-sdk-java
git clone https://github.com/hyperledger/udo-sdk-java $WD
cd $WD
git checkout tags/v1.0.1
export GOPATH=$WD/src/test/fixture

cd $WD/src/test
chmod +x cirun.sh
source cirun.sh
