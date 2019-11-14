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

package filter

import (
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
)

// +kubebuilder:object:generate=true
// +docName:"Parser"
// https://docs.fluentd.org/filter/parser
type ParserConfig struct {
	// Specify field name in the record to parse. If you leave empty the Container Runtime default will be used.
	KeyName string `json:"key_name,omitempty"`
	// Keep original event time in parsed result.
	ReserveTime bool `json:"reserve_time,omitempty"`
	// Keep original key-value pair in parsed result.
	ReserveData bool `json:"reserve_data,omitempty"`
	// Remove key_name field when parsing is succeeded
	RemoveKeyNameField bool `json:"remove_key_name_field,omitempty"`
	// If true, invalid string is replaced with safe characters and re-parse it.
	ReplaceInvalidSequence bool `json:"replace_invalid_sequence,omitempty"`
	// Store parsed values with specified key name prefix.
	InjectKeyPrefix string `json:"inject_key_prefix,omitempty"`
	// Store parsed values as a hash value in a field.
	HashValueField string `json:"hash_value_field,omitempty"`
	// Emit invalid record to @ERROR label. Invalid cases are: key not exist, format is not matched, unexpected error
	EmitInvalidRecordToError bool `json:"emit_invalid_record_to_error,omitempty"`
	// +docLink:"Parse Section,#Parse-Section"
	Parsers []ParseSection `json:"parsers,omitempty"`
}

// +kubebuilder:object:generate=true
// +docName:"Parse Section"
type ParseSection struct {
	// Parse type: apache2, apache_error, nginx, syslog, csv, tsv, ltsv, json, multiline, none, logfmt
	Type string `json:"type,omitempty"`
	// Regexp expression to evaluate
	Expression string `json:"expression,omitempty"`
	// Specify time field for event time. If the event doesn't have this field, current time is used.
	TimeKey string `json:"time_key,omitempty"`
	//  Specify null value pattern.
	NullValuePattern string `json:"null_value_pattern,omitempty"`
	// If true, empty string field is replaced with nil
	NullEmptyString bool `json:"null_empty_string,omitempty"`
	// If true, use Fluent::EventTime.now(current time) as a timestamp when time_key is specified.
	EstimateCurrentEvent bool `json:"estimate_current_event,omitempty"`
	// If true, keep time field in the record.
	KeepTimeKey bool `json:"keep_time_key,omitempty"`
	// Types casting the fields to proper types example: field1:type, field2:type
	Types string `json:"types,omitempty"`
	// Process value using specified format. This is available only when time_type is string
	TimeFormat string `json:"time_format,omitempty"`
	// Parse/format value according to this type available values: float, unixtime, string (default: string)
	TimeType string `json:"time_type,omitempty"`
	// Ff true, use local time. Otherwise, UTC is used. This is exclusive with utc. (default: true)
	LocalTime bool `json:"local_time,omitempty"`
	// If true, use UTC. Otherwise, local time is used. This is exclusive with localtime (default: false)
	UTC bool `json:"utc,omitempty"`
	// Use specified timezone. one can parse/format the time value in the specified timezone. (default: nil)
	Timezone string `json:"timezone,omitempty"`
}

func (p *ParseSection) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	parseMeta := types.PluginMeta{
		Directive: "parse",
		Type:      p.Type,
	}
	p.Type = ""
	return types.NewFlatDirective(parseMeta, p, secretLoader)
}

func NewParserConfig() *ParserConfig {
	return &ParserConfig{}
}

func (p *ParserConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "parser"
	parser := &types.GenericDirective{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "filter",
			Tag:       "**",
			Id:        id + "_" + pluginType,
		},
	}
	parserConfig := p.DeepCopy()
	if parserConfig.KeyName == "" {
		parserConfig.KeyName = types.GetLogKey()
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(parserConfig); err != nil {
		return nil, err
	} else {
		parser.Params = params
	}
	if len(parserConfig.Parsers) > 0 {
		for _, parseRule := range parserConfig.Parsers {
			if meta, err := parseRule.ToDirective(secretLoader, ""); err != nil {
				return nil, err
			} else {
				parser.SubDirectives = append(parser.SubDirectives, meta)
			}
		}
	}
	return parser, nil
}
