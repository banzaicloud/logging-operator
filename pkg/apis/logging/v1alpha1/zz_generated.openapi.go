// +build !ignore_autogenerated

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.Fluentbit":           schema_pkg_apis_logging_v1alpha1_Fluentbit(ref),
		"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentbitSpec":       schema_pkg_apis_logging_v1alpha1_FluentbitSpec(ref),
		"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentbitStatus":     schema_pkg_apis_logging_v1alpha1_FluentbitStatus(ref),
		"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.Fluentd":             schema_pkg_apis_logging_v1alpha1_Fluentd(ref),
		"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentdSpec":         schema_pkg_apis_logging_v1alpha1_FluentdSpec(ref),
		"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentdStatus":       schema_pkg_apis_logging_v1alpha1_FluentdStatus(ref),
		"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.LoggingPlugin":       schema_pkg_apis_logging_v1alpha1_LoggingPlugin(ref),
		"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.LoggingPluginSpec":   schema_pkg_apis_logging_v1alpha1_LoggingPluginSpec(ref),
		"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.LoggingPluginStatus": schema_pkg_apis_logging_v1alpha1_LoggingPluginStatus(ref),
	}
}

func schema_pkg_apis_logging_v1alpha1_Fluentbit(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Fluentbit is the Schema for the fluentbits API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentbitSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentbitStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentbitSpec", "github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentbitStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_logging_v1alpha1_FluentbitSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "FluentbitSpec defines the desired state of Fluentbit",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_logging_v1alpha1_FluentbitStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "FluentbitStatus defines the observed state of Fluentbit",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_logging_v1alpha1_Fluentd(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Fluentd is the Schema for the fluentds API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentdSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentdStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentdSpec", "github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.FluentdStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_logging_v1alpha1_FluentdSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "FluentdSpec defines the desired state of Fluentd",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_logging_v1alpha1_FluentdStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "FluentdStatus defines the observed state of Fluentd",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_logging_v1alpha1_LoggingPlugin(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "LoggingPlugin is the Schema for the loggingplugins API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.LoggingPluginSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.LoggingPluginStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.LoggingPluginSpec", "github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1.LoggingPluginStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_logging_v1alpha1_LoggingPluginSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "LoggingPluginSpec defines the desired state of LoggingPlugin",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_logging_v1alpha1_LoggingPluginStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "LoggingPluginStatus defines the observed state of LoggingPlugin",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
