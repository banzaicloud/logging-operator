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

package fluentd

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) role() (runtime.Object, reconciler.DesiredState, error) {
	if *r.Logging.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.Role{
			ObjectMeta: r.FluentdObjectMeta(roleName),
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{""},
					Resources: []string{"configmaps", "secrets"},
					Verbs:     []string{"*"},
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.Role{
		ObjectMeta: r.FluentdObjectMeta(roleName),
		Rules:      []rbacv1.PolicyRule{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) roleBinding() (runtime.Object, reconciler.DesiredState, error) {
	if *r.Logging.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate {
		return &rbacv1.RoleBinding{
			ObjectMeta: r.FluentdObjectMeta(roleBindingName),
			RoleRef: rbacv1.RoleRef{
				Kind:     "Role",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     r.Logging.QualifiedName(roleName),
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      r.Logging.QualifiedName(defaultServiceAccountName),
					Namespace: r.Logging.Spec.ControlNamespace,
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &rbacv1.RoleBinding{
		ObjectMeta: r.FluentdObjectMeta(roleBindingName),
		RoleRef:    rbacv1.RoleRef{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) serviceAccount() (runtime.Object, reconciler.DesiredState, error) {
	if *r.Logging.Spec.FluentdSpec.Security.RoleBasedAccessControlCreate && r.Logging.Spec.FluentdSpec.Security.ServiceAccount == "" {
		return &corev1.ServiceAccount{
			ObjectMeta: r.FluentdObjectMeta(defaultServiceAccountName),
		}, reconciler.StatePresent, nil
	}
	return &corev1.ServiceAccount{
		ObjectMeta: r.FluentdObjectMeta(defaultServiceAccountName),
	}, reconciler.StateAbsent, nil
}
