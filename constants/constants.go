package constants

import "regexp"

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

// Loadbalancer
const (
	HttpProtocol            = "HTTP"
	TerminatedHttpsProtocol = "TERMINATED_HTTPS"
	TcpProtocol             = "TCP"
	UdpProtocol             = "UDP"
	ProxyProtocol           = "PROXY"
	HttpsProtocol           = "HTTPS"
	UdpConnectProtocol      = "UDP-CONNECT"
	PingProtocol            = "PING"
	SctpProtocol            = "SCTP"
	TlsHelloProtocol        = "TLS-HELLO"

	GetMethod     = "GET"
	PostMethod    = "POST"
	HeadMethod    = "HEAD"
	PutMethod     = "PUT"
	DeleteMethod  = "DELETE"
	TraceMethod   = "TRACE"
	OptionsMethod = "OPTIONS"
	PatchMethod   = "PATCH"
	ConnectMethod = "CONNECT"

	RoundRobin       = "ROUND_ROBIN"
	LeastConnections = "LEAST_CONNECTIONS"
	SourceIp         = "SOURCE_IP"

	HttpCookie = "HTTP_COOKIE"
	AppCookie  = "APP_COOKIE"

	OkStatus               = "200"
	CreatedStatus          = "201"
	AcceptedStatus         = "202"
	NonAuthoritativeStatus = "203"
	NoContentStatus        = "204"

	InternalNetworkType = "internal"
	ExternalNetworkType = "external"

	LbSmallType  = "small"
	LbMediumType = "medium"
	LbLargeType  = "large"
)

var (
	ValidListenerProtocols = []string{
		HttpProtocol,
		TerminatedHttpsProtocol,
		TcpProtocol,
		UdpProtocol,
	}
	ValidAlgorithms = []string{
		RoundRobin,
		LeastConnections,
		SourceIp,
	}
	ValidPoolProtocols = []string{
		HttpProtocol,
		TcpProtocol,
		ProxyProtocol,
		UdpProtocol,
	}
	ValidHealthMonitorProtocols = []string{
		HttpProtocol,
		HttpsProtocol,
		PingProtocol,
		SctpProtocol,
		TcpProtocol,
		TlsHelloProtocol,
		UdpConnectProtocol,
	}
	ValidHealthMonitorMethods = []string{
		GetMethod,
		PostMethod,
		HeadMethod,
		PutMethod,
		DeleteMethod,
		TraceMethod,
		OptionsMethod,
		PatchMethod,
		ConnectMethod,
	}
	ValidHealthMonitorExceptedCodes = []string{
		OkStatus,
		CreatedStatus,
		AcceptedStatus,
		NonAuthoritativeStatus,
		NoContentStatus,
	}
	ValidStickySessions = []string{
		AppCookie,
		HttpCookie,
		SourceIp,
	}
	ValidUrlPathRegex, _ = regexp.Compile(`^/.*`)
	ValidLbNetworkTypes  = []string{ExternalNetworkType, InternalNetworkType}
	ValidLbTypes         = []string{LbSmallType, LbMediumType, LbLargeType}
)
