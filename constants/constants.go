package constants

const (
	KubernetesKubeRouter = "kube-router"
	KubernetesCilium     = "cilium"
	FreeDatatransfer     = "free_datatransfer"
	FreeBandwidth        = "free_bandwidth"
	SavingPlan           = "saving_plan"
	OnDemand             = "on_demand"
	NoSchedule           = "NoSchedule"
	PreferNoSchedule     = "PreferNoSchedule"
	NoExecute            = "NoExecute"
)

var (
	ValidCNIPlugins   = []string{KubernetesKubeRouter, KubernetesCilium}
	ValidBillingPlans = []string{SavingPlan, OnDemand}
	ValidNetworkPlans = []string{FreeDatatransfer, FreeBandwidth}
	ValidEffects      = []string{NoSchedule, PreferNoSchedule, NoExecute}
)
