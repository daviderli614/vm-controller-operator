package utils

import (
	"fmt"
	appv1 "github.com/daviderli614/vm-controller-operator/api/v1"
)

const (
	VMControllerName               = "vm-controller"
	VMControllerLabels             = "app"
	VMControllerAPIVersion         = "apps/v1"
	KindDeployment                 = "Deployment"
	VMControllerPriorityClassName  = "system-cluster-critical"
	VMControllerServiceAccountName = "vm-controller"
)

const (
	KubeSystemNamespace = "kube-system"
)

const (
	StdErrThresholdArg               VMControllerArg = "--stderrthreshold"
	CloudProviderArg                 VMControllerArg = "--cloud-provider"
	ScaleDownEnabledArg              VMControllerArg = "--scale-down-enabled"
	ScaleDownDelayAfterAddArg        VMControllerArg = "--scale-down-delay-after-add"
	ScaleDownUnneededTimeArg         VMControllerArg = "--scale-down-unneeded-time"
	ScaleDownUtilizationThresholdArg VMControllerArg = "--scale-down-utilization-threshold"
	ScaleUpUtilizationThresholdArg   VMControllerArg = "--scale-up-utilization-threshold"
	VerbosityArg                     VMControllerArg = "--v"
	StartArg                         VMControllerArg = "./vm-controller"
	BalanceSimilarNodeGroupsArg      VMControllerArg = "--balance-similar-node-groups"
	SkipNodesWithLocalStorageArg     VMControllerArg = "--skip-nodes-with-local-storage"
	ScanIntervalArg                  VMControllerArg = "--scan-interval"
    Expander                         VMControllerArg = "--expander"
)

type VMControllerArg string

func (a VMControllerArg) String() string {
	return string(a)
}

func (a VMControllerArg) Value(v interface{}) string {
	return fmt.Sprintf("%s=%v", a.String(), v)
}

func VMControllerArgs(c *appv1.VMClaim) []string {
	s := &c.Spec

	args := []string{
		StartArg.String(),
		VerbosityArg.Value(4),
		StdErrThresholdArg.Value("info"),
		CloudProviderArg.Value("ucloud"),
	}

	if c.Spec.ScaleDown != nil {
		args = append(args, ScaleDownArgs(s.ScaleDown)...)
	}

	if c.Spec.ScaleUp != nil {
		args = append(args, ScaleUpArgs(s.ScaleUp)...)
	}

	if c.Spec.ScanInterval != nil {
		args = append(args, ScanIntervalArg.Value(*c.Spec.ScanInterval))
	}

	if c.Spec.BalanceSimilarNodeGroups != nil {
		args = append(args, BalanceSimilarNodeGroupsArg.Value(*c.Spec.BalanceSimilarNodeGroups))
	}

	if c.Spec.SkipNodesWithLocalStorage != nil {
		args = append(args, SkipNodesWithLocalStorageArg.Value(*c.Spec.SkipNodesWithLocalStorage))
	}

	if c.Spec.Expander != nil {
		args = append(args, Expander.Value(*c.Spec.Expander))
	}

	return args
}

func ScaleDownArgs(sd *appv1.ScaleDownConfig) []string {
	if !sd.Enabled {
		return []string{ScaleDownEnabledArg.Value(false)}
	}

	args := []string{
		ScaleDownEnabledArg.Value(true),
	}

	if sd.DelayAfterAdd != nil {
		args = append(args, ScaleDownDelayAfterAddArg.Value(*sd.DelayAfterAdd))
	}

	if sd.UnneededTime != nil {
		args = append(args, ScaleDownUnneededTimeArg.Value(*sd.UnneededTime))
	}

	if sd.ScaleDownUtilizationThreshold != nil {
		args = append(args, ScaleDownUtilizationThresholdArg.Value(*sd.ScaleDownUtilizationThreshold))
	}

	return args
}

func ScaleUpArgs(su *appv1.ScaleUpConfig) []string {
	args := []string{}
	if su.ScaleUpUtilizationThreshold != nil {
		args = append(args, ScaleUpUtilizationThresholdArg.Value(*su.ScaleUpUtilizationThreshold))
	}

	return args
}