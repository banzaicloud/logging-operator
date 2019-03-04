/*
 * Copyright © 2019 Banzai Cloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fluentbit

import (
	loggingv1alpha1 "github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// TODO in case of rbac add created serviceAccount name
func (r *Reconciler) daemonSet() runtime.Object {
	return &appsv1.DaemonSet{
		ObjectMeta: templates.FluentbitObjectMeta(fluentbitDeaemonSetName, util.MergeLabels(r.Fluentbit.Labels, labelSelector), r.Fluentbit),
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: util.MergeLabels(r.Fluentbit.Labels, labelSelector)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: util.MergeLabels(r.Fluentbit.Labels, labelSelector),
					// TODO Move annotations to configuration
					Annotations: map[string]string{
						"prometheus.io/scrape": "true",
						"prometheus.io/path":   "/api/v1/metrics/prometheus",
						"prometheus.io/port":   "2020",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "logging",
					Volumes:            generateVolume(r.Fluentbit),
					Containers: []corev1.Container{
						{
							// TODO move to configuration
							Name:  "fluent-bit",
							Image: "fluent/fluent-bit:latest",
							// TODO get from config translate to const
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "monitor",
									ContainerPort: 2020,
									Protocol:      "TCP",
								},
							},
							// TODO Get this from config
							Resources: corev1.ResourceRequirements{
								Limits:   nil,
								Requests: nil,
							},
							VolumeMounts: generateVolumeMounts(r.Fluentbit),
						},
					},
				},
			},
		},
	}
}

func generateVolumeMounts(fluentbit *loggingv1alpha1.Fluentbit) (v []corev1.VolumeMount) {
	v = []corev1.VolumeMount{
		{
			Name:      "varlibcontainers",
			ReadOnly:  true,
			MountPath: "/var/lib/docker/containers",
		},
		{
			Name:      "config",
			MountPath: "/fluent-bit/etc/fluent-bit.conf",
			SubPath:   "fluent-bit.conf",
		},
		{
			Name:      "positions",
			MountPath: "/tail-db",
		},
		{
			Name:      "varlogs",
			ReadOnly:  true,
			MountPath: "/var/log/",
		},
	}
	if fluentbit.Spec.TLS.Enabled {
		tlsRelatedVolume := []corev1.VolumeMount{
			{
				Name:      "fluent-tls",
				MountPath: "/fluent-bit/tls/caCert",
				SubPath:   "caCert",
			},
			{
				Name:      "fluent-tls",
				MountPath: "/fluent-bit/tls/clientCert",
				SubPath:   "clientCert",
			},
			{
				Name:      "fluent-tls",
				MountPath: "/fluent-bit/tls/clientKey",
				SubPath:   "clientKey",
			},
		}
		v = append(v, tlsRelatedVolume...)
	}
	return
}

func generateVolume(fluentbit *loggingv1alpha1.Fluentbit) (v []corev1.Volume) {
	v = []corev1.Volume{
		{
			Name: "varlibcontainers",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/lib/docker/containers",
				},
			},
		},
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "fluent-bit-config",
					},
				},
			},
		},
		{
			Name: "varlogs",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/log",
				},
			},
		},
		{
			Name: "positions",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}
	if fluentbit.Spec.TLS.Enabled {
		tlsRelatedVolume := corev1.Volume{
			Name: "fluent-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: fluentbit.Spec.TLS.SecretName,
				},
			},
		}
		v = append(v, tlsRelatedVolume)
	}
	return
}
