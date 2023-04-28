module sysadm

go 1.16

require (
	github.com/go-playground/validator/v10 v10.7.0 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sql-driver/mysql v1.6.0
	github.com/lib/pq v1.10.6
	github.com/lithammer/dedent v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.6.0
	github.com/spf13/viper v1.12.0
	github.com/wangyysde/sshclient v0.0.0-20220914100106-2be5f8e32ceb
	github.com/wangyysde/sysadmLog v0.0.0-20210915071829-f43fc1c68a76
	github.com/wangyysde/sysadmServer v0.0.0-20220719023015-af14b6af71e5
	github.com/wangyysde/sysadmSessions v0.0.0-20211222125714-def5d4b4f078
	github.com/wangyysde/yaml v1.5.0
	golang.org/x/tools v0.8.0 // indirect
	k8s.io/api v0.26.3 // indirect
	k8s.io/apimachinery v0.26.3
	k8s.io/client-go v0.26.3
	k8s.io/cluster-bootstrap v0.0.0 // indirect
	k8s.io/kube-proxy v0.0.0
	k8s.io/kubernetes v1.26.3
	k8s.io/utils v0.0.0-20221107191617-1a15be271d1d
)

replace k8s.io/api => k8s.io/api v0.26.3

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.26.3

replace k8s.io/apimachinery => k8s.io/apimachinery v0.26.4

replace k8s.io/apiserver => k8s.io/apiserver v0.26.3

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.26.3

replace k8s.io/client-go => k8s.io/client-go v0.26.3

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.26.3

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.26.3

replace k8s.io/code-generator => k8s.io/code-generator v0.26.4

replace k8s.io/component-base => k8s.io/component-base v0.26.3

replace k8s.io/component-helpers => k8s.io/component-helpers v0.26.3

replace k8s.io/controller-manager => k8s.io/controller-manager v0.26.3

replace k8s.io/cri-api => k8s.io/cri-api v0.26.4

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.26.3

replace k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.26.3

replace k8s.io/kms => k8s.io/kms v0.26.4

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.26.3

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.26.3

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.26.3

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.26.3

replace k8s.io/kubectl => k8s.io/kubectl v0.26.3

replace k8s.io/kubelet => k8s.io/kubelet v0.26.3

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.26.3

replace k8s.io/metrics => k8s.io/metrics v0.26.3

replace k8s.io/mount-utils => k8s.io/mount-utils v0.26.4

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.26.3

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.26.3

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.26.3

replace k8s.io/sample-controller => k8s.io/sample-controller v0.26.3
