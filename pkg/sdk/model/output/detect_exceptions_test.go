// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output_test

import (
	"testing"

	"github.com/banzaicloud/logging-operator/pkg/sdk/model/output"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/render"
	"github.com/ghodss/yaml"
)

func TestDetectExceptions(t *testing.T) {
	CONFIG := []byte(`
remove_tag_prefix: foo
multiline_flush_interval: 0.1
languages: 
  - java
  - python
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
@type detect_exceptions
@id test_detect_exceptions
    <buffer tag,time>
      @type file
      path /buffers/test_detect_exceptions.*.buffer
      retry_forever true
      timekey 1m
      timekey_use_utc true
      timekey_wait 30s
    </buffer>
  </match>
`
	ed := &output.ExceptionDetectorOutputConfig{}
	yaml.Unmarshal(CONFIG, ed)
	test := render.NewOutputPluginTest(t, ed)
	test.DiffResult(expected)
}
