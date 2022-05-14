package resource

import (
	appv1 "github.com/daviderli614/vm-controller-operator/api/v1"
	"github.com/daviderli614/vm-controller-operator/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewDeployment(vmc *appv1.VMClaim) *appsv1.Deployment {
	podSpec := VMControllerPodSpec(vmc)
	labels := map[string]string{
		utils.VMControllerLabels: utils.VMControllerName,
	}
	var replicas int32 = 1
	return &appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{
			APIVersion: utils.VMControllerAPIVersion,
			Kind:       utils.KindDeployment,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            utils.VMControllerName,
			Namespace:       utils.KubeSystemNamespace,
			Labels:          labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: *podSpec,
			},
		},
	}
}

func VMControllerPodSpec(vmc *appv1.VMClaim) *corev1.PodSpec {
	args := utils.VMControllerArgs(vmc)

	spec := &corev1.PodSpec{
		ServiceAccountName: utils.VMControllerServiceAccountName,
		PriorityClassName:  utils.VMControllerPriorityClassName,
		Affinity:           &corev1.Affinity{
			NodeAffinity:   &corev1.NodeAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
					{
						Weight:     100,
						Preference: corev1.NodeSelectorTerm{
							MatchExpressions: []corev1.NodeSelectorRequirement{
								{
									Key:      "node.uk8s.ucloud.cn/asg_id",
									Operator: "DoesNotExist",
								},
							},
						},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:            utils.VMControllerName,
				Image:           vmc.Spec.Image,
				Command:         args,
				Env:             vmc.Spec.Envs,
				Resources:       vmc.Spec.Resources,
				ImagePullPolicy: corev1.PullAlways,
			},
		},
	}

	return spec
}