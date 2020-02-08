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

func TestLogZ(t *testing.T) {
	CONFIG := []byte(`
endpoint:
  port: 8107
http_idle_timeout: 10
output_include_tags: true
output_include_time: true
buffer:
  tags: tag
  type: memory
  chunk_limit_size: 187m
  flush_interval: 3s
  flush_thread_count: 4
  queue_limit_length: 4096
`)
	expected := `
  <match **>
	@type logzio_buffered
	@id test_logzio_buffered
	endpoint_url https://listener.logz.io:8071?token=ekXLvrbriwkuTmhZrELsYbSznxqnhyQj\u0026type=my_type
	http_idle_timeout 10
	output_include_tags true
	output_include_time true
    <buffer tag,time>
	  @type file
	  chunk_limit_size 16m
	  flush_interval 3s
	  flush_thread_count 4
	  path /buffers/test_logzio_buffered.*.buffer
	  queue_limit_length 4096
	  retry_forever true
      timekey 10m
      timekey_wait 10m
    </buffer>
  </match>
`
	es := &output.LogZOutput{}
	yaml.Unmarshal(CONFIG, es)
	test := render.NewOutputPluginTest(t, es)
	test.DiffResult(expected)
}
