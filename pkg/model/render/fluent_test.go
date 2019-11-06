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

package render_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/andreyvit/diff"
	"github.com/banzaicloud/logging-operator/pkg/model/filter"
	"github.com/banzaicloud/logging-operator/pkg/model/input"
	"github.com/banzaicloud/logging-operator/pkg/model/output"
	"github.com/banzaicloud/logging-operator/pkg/model/render"
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
	"github.com/banzaicloud/logging-operator/pkg/plugins"
)

func TestRenderDirective(t *testing.T) {
	var tests = []struct {
		name      string
		directive types.Directive
		expected  string
		reproduce int
	}{
		{
			name: "single level just name",
			directive: &types.GenericDirective{
				PluginMeta: types.PluginMeta{
					Directive: "match",
				},
			},
			expected: heredoc.Doc(`
            <match>
            </match>`,
			),
		},
		{
			name: "single level with tag and attributes",
			directive: &types.GenericDirective{
				PluginMeta: types.PluginMeta{
					Directive: "match",
					Tag:       "tag",
				},
				Params: map[string]string{
					"path": "file",
				},
			},
			expected: heredoc.Doc(`
            <match tag>
              path file
            </match>`,
			),
		},
		{
			name: "single level with just tag",
			directive: &types.GenericDirective{
				PluginMeta: types.PluginMeta{
					Directive: "match",
					Tag:       "tag",
				},
			},
			expected: heredoc.Doc(`
            <match tag>
            </match>`,
			),
		},
		{
			name: "single level with just attributes",
			directive: &types.GenericDirective{
				PluginMeta: types.PluginMeta{
					Directive: "match",
				},
				Params: map[string]string{
					"path": "file",
				},
			},
			expected: heredoc.Doc(`
            <match>
              path file
            </match>`,
			),
		},
		{
			name: "two levels",
			directive: &types.GenericDirective{
				PluginMeta: types.PluginMeta{
					Directive: "match",
				},
				Params: map[string]string{
					"path": "file",
				},
				SubDirectives: []types.Directive{
					&types.GenericDirective{
						PluginMeta: types.PluginMeta{
							Directive: "router1",
						},
						Params: map[string]string{
							"namespace": "asd",
							"labels":    "{\"a\":\"b\"}",
						},
					},
					&types.GenericDirective{
						PluginMeta: types.PluginMeta{
							Directive: "router2",
						},
						Params: map[string]string{
							"namespace": "asd2",
						},
					},
				},
			},
			expected: heredoc.Doc(`
            <match>
              path file
              <router1>
                labels {"a":"b"}
                namespace asd
              </router1>
              <router2>
                namespace asd2
              </router2>
            </match>`,
			),
		},
		{
			name:      "tail input",
			directive: toDirective(t, input.NewTailInputConfig("/path/to/input")),
			expected: heredoc.Doc(`
            <source>
              @type tail
              @id test_tail
              path /path/to/input
            </source>`,
			),
		},
		{
			name:      "stdout filter",
			directive: toDirective(t, filter.NewStdOutFilterConfig()),
			expected: heredoc.Doc(`
            <filter **>
              @type stdout
              @id test_stdout
            </filter>`,
			),
		},
		{
			name:      "stdout filter",
			directive: toDirective(t, output.NewNullOutputConfig()),
			expected: heredoc.Doc(`
            <match **>
              @type null
              @id test_null
            </match>`,
			),
		},
		{
			name:      "empty flow",
			directive: newFlowOrPanic("", nil),
			expected: heredoc.Doc(`
            <label @d41d8cd98f00b204e9800998ecf8427e>
            </label>`,
			),
		},
		{
			name:      "namespace flow",
			directive: newFlowOrPanic("test", nil),
			expected: heredoc.Doc(`
            <label @098f6bcd4621d373cade4e832627b4f6>
            </label>`,
			),
		},
		{
			name: "namespace and labels flow",
			directive: newFlowOrPanic("test", map[string]string{
				"key": "value",
				"a":   "b",
			}),
			expected: heredoc.Doc(`
            <label @e02a5a13f3f75484debfe1f11fecb65f>
            </label>`,
			),
			// run multiple times to make sure the label is stable
			reproduce: 10,
		},
		{
			name: "global router",
			directive: types.NewRouter("test").
				AddRoute(
					newFlowOrPanic("", nil),
				),
			expected: heredoc.Doc(`
            <match **>
              @type label_router
              @id test_label_router
              <route>
                @label @d41d8cd98f00b204e9800998ecf8427e
              </route>
            </match>`,
			),
		},
		{
			name: "namespaced router",
			directive: types.NewRouter("test").
				AddRoute(
					newFlowOrPanic("test", nil),
				),
			expected: heredoc.Doc(`
            <match **>
              @type label_router
              @id test_label_router
              <route>
                @label @098f6bcd4621d373cade4e832627b4f6
                namespace test
              </route>
            </match>`,
			),
		},
		{
			name: "namespaced router with labels",
			directive: types.NewRouter("test").
				AddRoute(
					newFlowOrPanic("test", map[string]string{"a": "b", "c": "d"}),
				),
			expected: heredoc.Doc(`
            <match **>
              @type label_router
              @id test_label_router
              <route>
                @label @092f5fa58e4f619d739f5b65f2ed38bc
                labels a:b,c:d
                namespace test
              </route>
            </match>`,
			),
			// run multiple times to make sure the label is stable
			reproduce: 10,
		},
	}
	for _, test := range tests {
		for i := 0; i <= test.reproduce; i++ {
			b := bytes.Buffer{}
			renderer := render.FluentRender{
				Out:    &b,
				Indent: 2,
			}
			_ = renderer.RenderDirectives([]types.Directive{
				test.directive,
			}, 0)
			if a, e := diff.TrimLinesInString(b.String()), diff.TrimLinesInString(test.expected); a != e {
				t.Errorf("[%s] Result does not match (-actual vs +expected):\n%v", test.name, diff.LineDiff(a, e))
			}
		}
	}
}

func TestMultipleOutput(t *testing.T) {
	system := types.NewSystem(toDirective(t, input.NewTailInputConfig("input.log")), types.NewRouter("test"))

	flowObj, err := types.NewFlow(
		"ns-test",
		map[string]string{
			"key1": "val1",
			"key2": "val2",
		})
	if err != nil {
		t.Fatal(err)
	}
	flowObj.
		WithFilters(toDirective(t, filter.NewStdOutFilterConfig())).
		WithOutputs(toDirective(t, output.NewNullOutputConfig())).
		WithOutputs(toDirective(t, output.NewNullOutputConfig()))

	err = system.RegisterFlow(flowObj)
	if err != nil {
		t.Fatal(err)
	}

	fluentConfig, err := system.Build()
	if err != nil {
		t.Fatal(err)
	}

	b := &bytes.Buffer{}
	renderer := render.FluentRender{
		Out:    b,
		Indent: 2,
	}
	err = renderer.Render(fluentConfig)
	if err != nil {
		t.Fatal(err)
	}

	expected := `
		<source>
          @type tail
          @id test_tail
          path input.log
        </source>
        <match **>
          @type label_router
          @id test_label_router
          <route>
            @label @901f778f9602a78e8fd702c1973d8d8d
            labels key1:val1,key2:val2
            namespace ns-test
          </route>
        </match>
        <label @901f778f9602a78e8fd702c1973d8d8d>
          <filter **>
            @type stdout
            @id test_stdout
          </filter>
          <match **>
            @type copy
            <store>
              @type null
              @id test_null
            </store>
            <store>
              @type null
              @id test_null
            </store>
          </match>
        </label>`

	if a, e := diff.TrimLinesInString(b.String()), diff.TrimLinesInString(expected); a != e {
		t.Errorf("Result does not match (-actual vs +expected):\n%v\nActual: %s", diff.LineDiff(a, e), b.String())
	}
}

func TestRenderFullFluentConfig(t *testing.T) {
	system := types.NewSystem(toDirective(t, input.NewTailInputConfig("input.log")), types.NewRouter("test"))

	flowObj, err := types.NewFlow(
		"ns-test",
		map[string]string{
			"key1": "val1",
			"key2": "val2",
		})
	if err != nil {
		t.Fatal(err)
	}
	flowObj.
		WithFilters(toDirective(t, filter.NewStdOutFilterConfig())).
		WithOutputs(toDirective(t, output.NewNullOutputConfig()))

	err = system.RegisterFlow(flowObj)
	if err != nil {
		t.Fatal(err)
	}

	fluentConfig, err := system.Build()
	if err != nil {
		t.Fatal(err)
	}

	b := &bytes.Buffer{}
	renderer := render.FluentRender{
		Out:    b,
		Indent: 2,
	}
	err = renderer.Render(fluentConfig)
	if err != nil {
		t.Fatal(err)
	}

	expected := `
		<source>
          @type tail
          @id test_tail
          path input.log
        </source>
        <match **>
          @type label_router
          @id test_label_router
          <route>
            @label @901f778f9602a78e8fd702c1973d8d8d
            labels key1:val1,key2:val2
            namespace ns-test
          </route>
        </match>
        <label @901f778f9602a78e8fd702c1973d8d8d>
          <filter **>
            @type stdout
            @id test_stdout
          </filter>
          <match **>
            @type null
            @id test_null
          </match>
        </label>`

	if a, e := diff.TrimLinesInString(b.String()), diff.TrimLinesInString(expected); a != e {
		t.Errorf("Result does not match (-actual vs +expected):\n%v\nActual: %s", diff.LineDiff(a, e), b.String())
	}
}

func TestRenderS3(t *testing.T) {
	table := []struct {
		name     string
		s3Config output.S3OutputConfig
		expected string
		err      string
	}{
		{
			name: "assumerole",
			s3Config: output.S3OutputConfig{
				Path:     "/var/buffer",
				S3Bucket: "test_bucket",
				Buffer: &output.Buffer{
					RetryForever: true,
					Path:         "asd",
				},
				AssumeRoleCredentials: &output.S3AssumeRoleCredentials{
					RoleArn:         "asd",
					RoleSessionName: "lkj",
				},
			},
			expected: ` @type s3
                        @id test_s3
						path /var/buffer
						s3_bucket test_bucket
						<buffer tag,time>
						  @type file
                          path asd
						  retry_forever true
                          timekey 10m
						</buffer>
						<assume_role_credentials>
							role_arn asd
							role_session_name lkj
						</assume_role_credentials>`,
		},
		{
			name: "instanceprofile",
			s3Config: output.S3OutputConfig{
				Path:                       "/var/buffer",
				S3Bucket:                   "test_bucket",
				InstanceProfileCredentials: &output.S3InstanceProfileCredentials{},
			},
			expected: ` @type s3
                        @id test_s3
						path /var/buffer
						s3_bucket test_bucket
						<instance_profile_credentials>
						</instance_profile_credentials>`,
		},
		{
			name: "shared",
			s3Config: output.S3OutputConfig{
				Path:     "/var/buffer",
				S3Bucket: "test_bucket",
				SharedCredentials: &output.S3SharedCredentials{
					Path:        "e",
					ProfileName: "f",
				},
			},
			expected: ` @type s3
                        @id test_s3
						path /var/buffer
						s3_bucket test_bucket
						<shared_credentials>
							path e
							profile_name f
						</shared_credentials>`,
		},
		{
			name: "missing auth",
			s3Config: output.S3OutputConfig{
				Path:     "/var/buffer",
				S3Bucket: "test_bucket",
			},
			err: "One of AssumeRoleCredentials or SharedCredentials or InstanceProfileCredentials must be configured",
		},
	}
	for _, item := range table {
		t.Logf("> %s\n", item.name)
		err := ValidateRenderS3(t, &item.s3Config, item.expected)
		if item.err != "" {
			if err == nil {
				t.Errorf("expected error: %s", item.err)
				continue
			}
			if err.Error() != item.err {
				t.Errorf("expected error: %s got %s", item.err, err)
				continue
			}
			continue
		}
		if err != nil {
			t.Error(err)
		}
	}
}

func ValidateRenderS3(t *testing.T, s3Config plugins.DirectiveConverter, expected string) error {
	system := types.NewSystem(toDirective(t, input.NewTailInputConfig("input.log")), types.NewRouter("test"))

	s3Plugin, err := s3Config.ToDirective(secret.NewSecretLoader(nil, "", "", nil), "test")
	if err != nil {
		return err
	}
	flowObj, err := types.NewFlow(
		"ns-test",
		map[string]string{
			"key1": "val1",
			"key2": "val2",
		})
	if err != nil {
		return err
	}
	flowObj.WithOutputs(s3Plugin)

	err = system.RegisterFlow(flowObj)
	if err != nil {
		return err
	}

	fluentConfig, err := system.Build()
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	renderer := render.FluentRender{
		Out:    b,
		Indent: 2,
	}
	err = renderer.Render(fluentConfig)
	if err != nil {
		return err
	}

	expected = fmt.Sprintf(`
		<source>
          @type tail
          @id test_tail
          path input.log
        </source>
        <match **>
          @type label_router
          @id test_label_router
          <route>
            @label @901f778f9602a78e8fd702c1973d8d8d
            labels key1:val1,key2:val2
            namespace ns-test
          </route>
        </match>
        <label @901f778f9602a78e8fd702c1973d8d8d>
          <match **>
            %s
          </match>
        </label>`, expected)
	if a, e := diff.TrimLinesInString(b.String()), diff.TrimLinesInString(expected); a != e {
		t.Errorf("Result does not match (-actual vs +expected):\n%v\nActual: %s", diff.LineDiff(a, e), b.String())
	}
	return nil
}

func newFlowOrPanic(namespace string, labels map[string]string) *types.Flow {
	flowObj, err := types.NewFlow(namespace, labels)
	if err != nil {
		panic(err)
	}
	return flowObj
}

func toDirective(t *testing.T, converter plugins.DirectiveConverter) types.Directive {
	directive, err := converter.ToDirective(secret.NewSecretLoader(nil, "", "", nil), "test")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	return directive
}
