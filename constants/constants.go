package constants

// Kubernetes
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

// Loadbalancer L7 policy
const (
	RedirectToPoolAction = "REDIRECT_TO_POOL"
	RejectAction         = "REJECT"
	RedirectToUrlAction  = "REDIRECT_TO_URL"
	RedirectPrefixAction = "REDIRECT_PREFIX"

	ACLsTypeHostName = "HOST_NAME"
	ACLsTypePath     = "PATH"
	ACLsTypeHeader   = "HEADER"
	ACLsTypeFileType = "FILE_TYPE"

	ACLsCompareTypeEqualTo    = "EQUAL_TO"
	ACLsCompareTypeRegex      = "REGEX"
	ACLsCompareTypeContains   = "CONTAINS"
	ACLsCompareTypeEndsWith   = "ENDS_WITH"
	ACLsCompareTypeStartsWith = "STARTS_WITH"
)

var (
	ValidL7PolicyActions = []string{
		RedirectToPoolAction,
		RejectAction,
		RedirectToUrlAction,
		RedirectPrefixAction,
	}
	ValidACLsTypes = []string{
		ACLsTypeHostName,
		ACLsTypePath,
		ACLsTypeHeader,
		ACLsTypeFileType,
	}
	ValidACLsCompareType = []string{
		ACLsCompareTypeEqualTo,
		ACLsCompareTypeRegex,
		ACLsCompareTypeContains,
		ACLsCompareTypeEndsWith,
		ACLsCompareTypeStartsWith,
	}
)
