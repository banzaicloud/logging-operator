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

package render

import (
	"bytes"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/banzaicloud/logging-operator/pkg/model/input"
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
	"github.com/banzaicloud/logging-operator/pkg/plugins"
)

func toDirective(t *testing.T, converter plugins.DirectiveConverter) types.Directive {
	directive, err := converter.ToDirective(secret.NewSecretLoader(nil, "", "", nil), "test")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	return directive
}

type PluginUnitTest struct {
	System       *types.Builder
	FluentConfig *types.System
	Test         *testing.T
	Prefix       string
}

func (p PluginUnitTest) Render() string {
	b := &bytes.Buffer{}
	renderer := FluentRender{
		Out:    b,
		Indent: 2,
	}
	err := renderer.Render(p.FluentConfig)
	if err != nil {
		p.Test.Fatal(err)
	}
	return b.String()
}

func (p PluginUnitTest) DiffResult(expected string) {
	prepared := p.Prefix + expected + "</label>"
	if a, e := diff.TrimLinesInString(p.Render()), diff.TrimLinesInString(prepared); a != e {
		p.Test.Errorf("Result does not match (-actual vs +expected):\n%v\nActual: %s", diff.LineDiff(a, e), p.Render())
	}
}

func NewOutputPluginTest(t *testing.T, plugin plugins.DirectiveConverter) *PluginUnitTest {
	suite := &PluginUnitTest{
		Test: t,
		Prefix: `
<source>
  @type tail
  @id test_tail
  path input.log
</source>
<match **>
  @type label_router
  @id test_label_router
  <route>
    @label @a42fd8d29c181fcf9887280c4a51bd1e
    namespace ns-test
  </route>
</match>
<label @a42fd8d29c181fcf9887280c4a51bd1e>`,
	}
	suite.System = types.NewSystem(toDirective(t, input.NewTailInputConfig("input.log")), types.NewRouter("test"))

	flowObj, err := types.NewFlow(
		"ns-test",
		map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	flowObj.WithOutputs(toDirective(t, plugin))

	err = suite.System.RegisterFlow(flowObj)
	if err != nil {
		t.Fatal(err)
	}

	fluentConfig, err := suite.System.Build()
	if err != nil {
		t.Fatal(err)
	}
	suite.FluentConfig = fluentConfig
	return suite
}
