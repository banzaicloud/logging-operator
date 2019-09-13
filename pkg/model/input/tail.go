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

package input

import (
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
)

// +kubebuilder:object:generate=true

type TailInputConfig struct {
	Path string `json:"path"`
}

func NewTailInputConfig(path string) *TailInputConfig {
	return &TailInputConfig{Path: path}
}

func (c *TailInputConfig) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	return types.NewFlatDirective(types.PluginMeta{
		Type:      "tail",
		Directive: "source",
	}, c, secretLoader)
}
