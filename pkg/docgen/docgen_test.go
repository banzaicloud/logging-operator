// Copyright © 2020 Banzai Cloud
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

package docgen_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/banzaicloud/logging-operator/pkg/docgen"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var logger logr.Logger

func init() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = zapr.NewLogger(log)
}

func TestGenParse(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	docItem := docgen.DocItem{
		Name:       "sample-name",
		SourcePath: filepath.Join(currentDir, "testdata", "sample.go"),
		DestPath:   filepath.Join(currentDir, "../../build/test/docgen"),
	}

	parser := docgen.GetDocumentParser(docItem, logger)
	parser.Generate()
}
