/*
Copyright 2017 Hitachi America

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

package metadata_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/hyperledger/udo/common/tools/configtxgen/metadata"
	"github.com/stretchr/testify/assert"
)

func TestGetVersionInfo(t *testing.T) {
	testVersions := []string{"", "TestVersion"}

	for _, version := range testVersions {
		metadata.Version = version
		if version == "" {
			version = "development build"
		}

		expected := fmt.Sprintf("%s:\n Version: %s\n Go version: %s\n OS/Arch: %s",
			metadata.ProgramName, version, runtime.Version(),
			fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
		assert.Equal(t, expected, metadata.GetVersionInfo())
	}
}
