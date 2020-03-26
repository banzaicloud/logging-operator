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

func TestS3(t *testing.T) {
	CONFIG := []byte(`
assume_role_credentials: arn:aws:iam::123456789012:role/logs
s3_bucket: logging-amazon-s3
s3_region: eu-central-1
path: logs/${tag}/%Y/%m/%d/
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
    @type s3
    @id test
    path logs/${tag}/%Y/%m/%d/
    s3_bucket logging-amazon-s3
    s3_region eu-central-1
    <buffer tag,time>
      @type file
      path /buffers/test.*.buffer
      retry_forever true
      timekey 1m
      timekey_use_utc true
      timekey_wait 30s
    </buffer>
    <assume_role_credentials>
      role_arn
      role_session_name
    </assume_role_credentials>
  </match>
`
	s3 := &output.S3OutputConfig{}
	yaml.Unmarshal(CONFIG, s3)
	test := render.NewOutputPluginTest(t, s3)
	test.DiffResult(expected)
}
