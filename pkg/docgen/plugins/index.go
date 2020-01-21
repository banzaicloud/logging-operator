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

package plugins

import (
	"fmt"
	"path/filepath"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/docgen"
	"github.com/go-logr/logr"
)

func GenerateWithIndex(lister *PluginLister, log logr.Logger) error {
	index := docgen.NewDoc(docgen.DocItem{
		Name:     "Readme",
		DestPath: "docs/plugins",
	}, log.WithName("index"))

	index.Append("# Supported Plugins\n\n")
	index.Append("For more information please click on the plugin name")
	index.Append("<center>\n")
	index.Append("| Name | Type | Description | Status |Version |")
	index.Append("|:---|---|:---|:---:|---:|")

	plugins, err := lister.GetPlugins()
	if err != nil {
		return errors.WrapIf(err, "failed to get plugin list")
	}

	for _, plugin := range plugins {
		log.Info("plugin", "Name", plugin.Item.SourcePath)
		document := docgen.GetDocumentParser(plugin.Item, log.WithName("docgen"))
		if err := document.Generate(); err != nil {
			return err
		}

		relPath, err := filepath.Rel("docs/plugins", document.Item.DestPath)
		if err != nil {
			return errors.WrapIff(err, "failed to determine relpath for %s", document.Item.DestPath)
		}

		index.Append(fmt.Sprintf("| **[%s](%s)** | %s | %s | %s | [%s](%s) |",
			document.DisplayName,
			filepath.Join(relPath, document.Item.Name+".md"),
			plugin.Category,
			document.Desc,
			document.Status,
			document.Version,
			document.Url))
	}
	index.Append("</center>")

	if err := index.Generate(); err != nil {
		return err
	}

	return nil
}
