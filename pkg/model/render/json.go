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
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/banzaicloud/logging-operator/pkg/model/types"
	"github.com/goph/emperror"
)

type JsonRender struct {
	out    io.Writer
	indent int
}

func (t *JsonRender) Render(config types.FluentConfig) error {
	var out []byte
	var err error
	if t.indent > 0 {
		out, err = json.MarshalIndent(config, "", strings.Repeat(" ", t.indent))
	} else {
		out, err = json.Marshal(config)
	}

	if err != nil {
		return emperror.Wrap(err, "Failed to marshal model into yaml")
	}
	fmt.Fprintf(t.out, "%s", out)
	return nil
}
