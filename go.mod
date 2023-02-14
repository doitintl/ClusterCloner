module clustercloner

go 1.14

require (
	cloud.google.com/go v0.43.0
	github.com/Azure/azure-sdk-for-go v39.2.0+incompatible
	github.com/Azure/go-autorest/autorest v0.10.0
	github.com/Azure/go-autorest/autorest/adal v0.8.2
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2
	github.com/Azure/go-autorest/autorest/to v0.3.0
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.5.1
	github.com/urfave/cli/v2 v2.2.0
	github.com/weaveworks/eksctl v0.0.0-20200521211543-1564272ac86f
	google.golang.org/genproto v0.0.0-20190716160619-c506a9f90610
)

require (
	github.com/aws/aws-sdk-go v1.34.0
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0
	github.com/kris-nova/logger v0.0.0-20181127235838-fd0d87064b06
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/tjarratt/babble v0.0.0-20191209142150-eecdf8c2339d
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.0.0+incompatible
	github.com/awslabs/goformation => github.com/errordeveloper/goformation v0.0.0-20190507151947-a31eae35e596
	// Override version since auto-detected one fails with GOPROXY
	github.com/census-instrumentation/opencensus-proto => github.com/census-instrumentation/opencensus-proto v0.2.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1
	// k8s.io/kops is still using old version of component-base
	// which uses an older version of the following package
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.4
	// Used to pin the k8s library versions regardless of what other dependencies enforce
	k8s.io/api => k8s.io/api v0.16.8
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.8
	k8s.io/apiserver => k8s.io/apiserver v0.16.8
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.8
	k8s.io/client-go => k8s.io/client-go v0.16.8
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.16.8
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.8
	k8s.io/code-generator => k8s.io/code-generator v0.16.8
	k8s.io/component-base => k8s.io/component-base v0.16.8
	k8s.io/cri-api => k8s.io/cri-api v0.16.8
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.16.8
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.8
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.16.8
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.16.8
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.16.8
	k8s.io/kubectl => k8s.io/kubectl v0.16.8
	k8s.io/kubelet => k8s.io/kubelet v0.16.8
	k8s.io/kubernetes => k8s.io/kubernetes v1.16.8
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.16.8
	k8s.io/metrics => k8s.io/metrics v0.16.8
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.16.8
)
