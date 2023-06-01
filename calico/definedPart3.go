/* =============================================================
* @Author:  Wayne Wang <yuying.wang@sincerecloud.com>
* @ Last Modified At: 2023.06.01
* @Copyright (c) 2023 Sincerecloud
* @HomePage: https://www.sincerecloud.com/
*
*  定义部署calico的yaml内容，本内容使用模板变量，这是第三部分内容。
*  程序应当按顺序先后apply各部分的yaml内容。
 */

package calico

// 部署calico的yaml内容第三部分内容
var tplPart3 string = `
---
# Source: calico/templates/kdd-crds.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ipamblocks.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: IPAMBlock
    listKind: IPAMBlockList
    plural: ipamblocks
    singular: ipamblock
  preserveUnknownFields: false
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: IPAMBlockSpec contains the specification for an IPAMBlock
              resource.
            properties:
              affinity:
                description: Affinity of the block, if this block has one. If set,
                  it will be of the form "host:<hostname>". If not set, this block
                  is not affine to a host.
                type: string
              allocations:
                description: Array of allocations in-use within this block. nil entries
                  mean the allocation is free. For non-nil entries at index i, the
                  index is the ordinal of the allocation within this block and the
                  value is the index of the associated attributes in the Attributes
                  array.
                items:
                  type: integer
                  # TODO: This nullable is manually added in. We should update controller-gen
                  # to handle []*int properly itself.
                  nullable: true
                type: array
              attributes:
                description: Attributes is an array of arbitrary metadata associated
                  with allocations in the block. To find attributes for a given allocation,
                  use the value of the allocation's entry in the Allocations array
                  as the index of the element in this array.
                items:
                  properties:
                    handle_id:
                      type: string
                    secondary:
                      additionalProperties:
                        type: string
                      type: object
                  type: object
                type: array
              cidr:
                description: The block's CIDR.
                type: string
              deleted:
                description: Deleted is an internal boolean used to workaround a limitation
                  in the Kubernetes API whereby deletion will not return a conflict
                  error if the block has been updated. It should not be set manually.
                type: boolean
              sequenceNumber:
                default: 0
                description: We store a sequence number that is updated each time
                  the block is written. Each allocation will also store the sequence
                  number of the block at the time of its creation. When releasing
                  an IP, passing the sequence number associated with the allocation
                  allows us to protect against a race condition and ensure the IP
                  hasn't been released and re-allocated since the release request.
                format: int64
                type: integer
              sequenceNumberForAllocation:
                additionalProperties:
                  format: int64
                  type: integer
                description: Map of allocated ordinal within the block to sequence
                  number of the block at the time of allocation. Kubernetes does not
                  allow numerical keys for maps, so the key is cast to a string.
                type: object
              strictAffinity:
                description: StrictAffinity on the IPAMBlock is deprecated and no
                  longer used by the code. Use IPAMConfig StrictAffinity instead.
                type: boolean
              unallocated:
                description: Unallocated is an ordered list of allocations which are
                  free in the block.
                items:
                  type: integer
                type: array
            required:
            - allocations
            - attributes
            - cidr
            - strictAffinity
            - unallocated
            type: object
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []

---
# Source: calico/templates/kdd-crds.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ipamconfigs.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: IPAMConfig
    listKind: IPAMConfigList
    plural: ipamconfigs
    singular: ipamconfig
  preserveUnknownFields: false
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: IPAMConfigSpec contains the specification for an IPAMConfig
              resource.
            properties:
              autoAllocateBlocks:
                type: boolean
              maxBlocksPerHost:
                description: MaxBlocksPerHost, if non-zero, is the max number of blocks
                  that can be affine to each host.
                maximum: 2147483647
                minimum: 0
                type: integer
              strictAffinity:
                type: boolean
            required:
            - autoAllocateBlocks
            - strictAffinity
            type: object
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []

---
# Source: calico/templates/kdd-crds.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ipamhandles.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: IPAMHandle
    listKind: IPAMHandleList
    plural: ipamhandles
    singular: ipamhandle
  preserveUnknownFields: false
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: IPAMHandleSpec contains the specification for an IPAMHandle
              resource.
            properties:
              block:
                additionalProperties:
                  type: integer
                type: object
              deleted:
                type: boolean
              handleID:
                type: string
            required:
            - block
            - handleID
            type: object
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []

---
# Source: calico/templates/kdd-crds.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ippools.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: IPPool
    listKind: IPPoolList
    plural: ippools
    singular: ippool
  preserveUnknownFields: false
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: IPPoolSpec contains the specification for an IPPool resource.
            properties:
              allowedUses:
                description: AllowedUse controls what the IP pool will be used for.  If
                  not specified or empty, defaults to ["Tunnel", "Workload"] for back-compatibility
                items:
                  type: string
                type: array
              blockSize:
                description: The block size to use for IP address assignments from
                  this pool. Defaults to 26 for IPv4 and 122 for IPv6.
                type: integer
              cidr:
                description: The pool CIDR.
                type: string
              disableBGPExport:
                description: 'Disable exporting routes from this IP Pool''s CIDR over
                  BGP. [Default: false]'
                type: boolean
              disabled:
                description: When disabled is true, Calico IPAM will not assign addresses
                  from this pool.
                type: boolean
              ipip:
                description: 'Deprecated: this field is only used for APIv1 backwards
                  compatibility. Setting this field is not allowed, this field is
                  for internal use only.'
                properties:
                  enabled:
                    description: When enabled is true, ipip tunneling will be used
                      to deliver packets to destinations within this pool.
                    type: boolean
                  mode:
                    description: The IPIP mode.  This can be one of "always" or "cross-subnet".  A
                      mode of "always" will also use IPIP tunneling for routing to
                      destination IP addresses within this pool.  A mode of "cross-subnet"
                      will only use IPIP tunneling when the destination node is on
                      a different subnet to the originating node.  The default value
                      (if not specified) is "always".
                    type: string
                type: object
              ipipMode:
                description: Contains configuration for IPIP tunneling for this pool.
                  If not specified, then this is defaulted to "Never" (i.e. IPIP tunneling
                  is disabled).
                type: string
              nat-outgoing:
                description: 'Deprecated: this field is only used for APIv1 backwards
                  compatibility. Setting this field is not allowed, this field is
                  for internal use only.'
                type: boolean
              natOutgoing:
                description: When nat-outgoing is true, packets sent from Calico networked
                  containers in this pool to destinations outside of this pool will
                  be masqueraded.
                type: boolean
              nodeSelector:
                description: Allows IPPool to allocate for a specific node by label
                  selector.
                type: string
              vxlanMode:
                description: Contains configuration for VXLAN tunneling for this pool.
                  If not specified, then this is defaulted to "Never" (i.e. VXLAN
                  tunneling is disabled).
                type: string
            required:
            - cidr
            type: object
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []

---
# Source: calico/templates/kdd-crds.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: (devel)
  creationTimestamp: null
  name: ipreservations.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: IPReservation
    listKind: IPReservationList
    plural: ipreservations
    singular: ipreservation
  preserveUnknownFields: false
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: IPReservationSpec contains the specification for an IPReservation
              resource.
            properties:
              reservedCIDRs:
                description: ReservedCIDRs is a list of CIDRs and/or IP addresses
                  that Calico IPAM will exclude from new allocations.
                items:
                  type: string
                type: array
            type: object
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []

---
# Source: calico/templates/kdd-crds.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: kubecontrollersconfigurations.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: KubeControllersConfiguration
    listKind: KubeControllersConfigurationList
    plural: kubecontrollersconfigurations
    singular: kubecontrollersconfiguration
  preserveUnknownFields: false
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: KubeControllersConfigurationSpec contains the values of the
              Kubernetes controllers configuration.
            properties:
              controllers:
                description: Controllers enables and configures individual Kubernetes
                  controllers
                properties:
                  namespace:
                    description: Namespace enables and configures the namespace controller.
                      Enabled by default, set to nil to disable.
                    properties:
                      reconcilerPeriod:
                        description: 'ReconcilerPeriod is the period to perform reconciliation
                          with the Calico datastore. [Default: 5m]'
                        type: string
                    type: object
                  node:
                    description: Node enables and configures the node controller.
                      Enabled by default, set to nil to disable.
                    properties:
                      hostEndpoint:
                        description: HostEndpoint controls syncing nodes to host endpoints.
                          Disabled by default, set to nil to disable.
                        properties:
                          autoCreate:
                            description: 'AutoCreate enables automatic creation of
                              host endpoints for every node. [Default: Disabled]'
                            type: string
                        type: object
                      leakGracePeriod:
                        description: 'LeakGracePeriod is the period used by the controller
                          to determine if an IP address has been leaked. Set to 0
                          to disable IP garbage collection. [Default: 15m]'
                        type: string
                      reconcilerPeriod:
                        description: 'ReconcilerPeriod is the period to perform reconciliation
                          with the Calico datastore. [Default: 5m]'
                        type: string
                      syncLabels:
                        description: 'SyncLabels controls whether to copy Kubernetes
                          node labels to Calico nodes. [Default: Enabled]'
                        type: string
                    type: object
                  policy:
                    description: Policy enables and configures the policy controller.
                      Enabled by default, set to nil to disable.
                    properties:
                      reconcilerPeriod:
                        description: 'ReconcilerPeriod is the period to perform reconciliation
                          with the Calico datastore. [Default: 5m]'
                        type: string
                    type: object
                  serviceAccount:
                    description: ServiceAccount enables and configures the service
                      account controller. Enabled by default, set to nil to disable.
                    properties:
                      reconcilerPeriod:
                        description: 'ReconcilerPeriod is the period to perform reconciliation
                          with the Calico datastore. [Default: 5m]'
                        type: string
                    type: object
                  workloadEndpoint:
                    description: WorkloadEndpoint enables and configures the workload
                      endpoint controller. Enabled by default, set to nil to disable.
                    properties:
                      reconcilerPeriod:
                        description: 'ReconcilerPeriod is the period to perform reconciliation
                          with the Calico datastore. [Default: 5m]'
                        type: string
                    type: object
                type: object
              debugProfilePort:
                description: DebugProfilePort configures the port to serve memory
                  and cpu profiles on. If not specified, profiling is disabled.
                format: int32
                type: integer
              etcdV3CompactionPeriod:
                description: 'EtcdV3CompactionPeriod is the period between etcdv3
                  compaction requests. Set to 0 to disable. [Default: 10m]'
                type: string
              healthChecks:
                description: 'HealthChecks enables or disables support for health
                  checks [Default: Enabled]'
                type: string
              logSeverityScreen:
                description: 'LogSeverityScreen is the log severity above which logs
                  are sent to the stdout. [Default: Info]'
                type: string
              prometheusMetricsPort:
                description: 'PrometheusMetricsPort is the TCP port that the Prometheus
                  metrics server should bind to. Set to 0 to disable. [Default: 9094]'
                type: integer
            required:
            - controllers
            type: object
          status:
            description: KubeControllersConfigurationStatus represents the status
              of the configuration. It's useful for admins to be able to see the actual
              config that was applied, which can be modified by environment variables
              on the kube-controllers process.
            properties:
              environmentVars:
                additionalProperties:
                  type: string
                description: EnvironmentVars contains the environment variables on
                  the kube-controllers that influenced the RunningConfig.
                type: object
              runningConfig:
                description: RunningConfig contains the effective config that is running
                  in the kube-controllers pod, after merging the API resource with
                  any environment variables.
                properties:
                  controllers:
                    description: Controllers enables and configures individual Kubernetes
                      controllers
                    properties:
                      namespace:
                        description: Namespace enables and configures the namespace
                          controller. Enabled by default, set to nil to disable.
                        properties:
                          reconcilerPeriod:
                            description: 'ReconcilerPeriod is the period to perform
                              reconciliation with the Calico datastore. [Default:
                              5m]'
                            type: string
                        type: object
                      node:
                        description: Node enables and configures the node controller.
                          Enabled by default, set to nil to disable.
                        properties:
                          hostEndpoint:
                            description: HostEndpoint controls syncing nodes to host
                              endpoints. Disabled by default, set to nil to disable.
                            properties:
                              autoCreate:
                                description: 'AutoCreate enables automatic creation
                                  of host endpoints for every node. [Default: Disabled]'
                                type: string
                            type: object
                          leakGracePeriod:
                            description: 'LeakGracePeriod is the period used by the
                              controller to determine if an IP address has been leaked.
                              Set to 0 to disable IP garbage collection. [Default:
                              15m]'
                            type: string
                          reconcilerPeriod:
                            description: 'ReconcilerPeriod is the period to perform
                              reconciliation with the Calico datastore. [Default:
                              5m]'
                            type: string
                          syncLabels:
                            description: 'SyncLabels controls whether to copy Kubernetes
                              node labels to Calico nodes. [Default: Enabled]'
                            type: string
                        type: object
                      policy:
                        description: Policy enables and configures the policy controller.
                          Enabled by default, set to nil to disable.
                        properties:
                          reconcilerPeriod:
                            description: 'ReconcilerPeriod is the period to perform
                              reconciliation with the Calico datastore. [Default:
                              5m]'
                            type: string
                        type: object
                      serviceAccount:
                        description: ServiceAccount enables and configures the service
                          account controller. Enabled by default, set to nil to disable.
                        properties:
                          reconcilerPeriod:
                            description: 'ReconcilerPeriod is the period to perform
                              reconciliation with the Calico datastore. [Default:
                              5m]'
                            type: string
                        type: object
                      workloadEndpoint:
                        description: WorkloadEndpoint enables and configures the workload
                          endpoint controller. Enabled by default, set to nil to disable.
                        properties:
                          reconcilerPeriod:
                            description: 'ReconcilerPeriod is the period to perform
                              reconciliation with the Calico datastore. [Default:
                              5m]'
                            type: string
                        type: object
                    type: object
                  debugProfilePort:
                    description: DebugProfilePort configures the port to serve memory
                      and cpu profiles on. If not specified, profiling is disabled.
                    format: int32
                    type: integer
                  etcdV3CompactionPeriod:
                    description: 'EtcdV3CompactionPeriod is the period between etcdv3
                      compaction requests. Set to 0 to disable. [Default: 10m]'
                    type: string
                  healthChecks:
                    description: 'HealthChecks enables or disables support for health
                      checks [Default: Enabled]'
                    type: string
                  logSeverityScreen:
                    description: 'LogSeverityScreen is the log severity above which
                      logs are sent to the stdout. [Default: Info]'
                    type: string
                  prometheusMetricsPort:
                    description: 'PrometheusMetricsPort is the TCP port that the Prometheus
                      metrics server should bind to. Set to 0 to disable. [Default:
                      9094]'
                    type: integer
                required:
                - controllers
                type: object
            type: object
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []

---
# Source: calico/templates/kdd-crds.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: networkpolicies.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: NetworkPolicy
    listKind: NetworkPolicyList
    plural: networkpolicies
    singular: networkpolicy
  preserveUnknownFields: false
  scope: Namespaced
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            properties:
              egress:
                description: The ordered set of egress rules.  Each rule contains
                  a set of packet match criteria and a corresponding action to apply.
                items:
                  description: "A Rule encapsulates a set of match criteria and an
                    action.  Both selector-based security Policy and security Profiles
                    reference rules - separated out as a list of rules for both ingress
                    and egress packet matching. \n Each positive match criteria has
                    a negated version, prefixed with \"Not\". All the match criteria
                    within a rule must be satisfied for a packet to match. A single
                    rule can contain the positive and negative version of a match
                    and both must be satisfied for the rule to match."
                  properties:
                    action:
                      type: string
                    destination:
                      description: Destination contains the match criteria that apply
                        to destination entity.
                      properties:
                        namespaceSelector:
                          description: "NamespaceSelector is an optional field that
                            contains a selector expression. Only traffic that originates
                            from (or terminates at) endpoints within the selected
                            namespaces will be matched. When both NamespaceSelector
                            and another selector are defined on the same rule, then
                            only workload endpoints that are matched by both selectors
                            will be selected by the rule. \n For NetworkPolicy, an
                            empty NamespaceSelector implies that the Selector is limited
                            to selecting only workload endpoints in the same namespace
                            as the NetworkPolicy. \n For NetworkPolicy, ` + "`global()`" + `
                            NamespaceSelector implies that the Selector is limited
                            to selecting only GlobalNetworkSet or HostEndpoint. \n
                            For GlobalNetworkPolicy, an empty NamespaceSelector implies
                            the Selector applies to workload endpoints across all
                            namespaces."
                          type: string
                        nets:
                          description: Nets is an optional field that restricts the
                            rule to only apply to traffic that originates from (or
                            terminates at) IP addresses in any of the given subnets.
                          items:
                            type: string
                          type: array
                        notNets:
                          description: NotNets is the negated version of the Nets
                            field.
                          items:
                            type: string
                          type: array
                        notPorts:
                          description: NotPorts is the negated version of the Ports
                            field. Since only some protocols have ports, if any ports
                            are specified it requires the Protocol match in the Rule
                            to be set to "TCP" or "UDP".
                          items:
                            anyOf:
                            - type: integer
                            - type: string
                            pattern: ^.*
                            x-kubernetes-int-or-string: true
                          type: array
                        notSelector:
                          description: NotSelector is the negated version of the Selector
                            field.  See Selector field for subtleties with negated
                            selectors.
                          type: string
                        ports:
                          description: "Ports is an optional field that restricts
                            the rule to only apply to traffic that has a source (destination)
                            port that matches one of these ranges/values. This value
                            is a list of integers or strings that represent ranges
                            of ports. \n Since only some protocols have ports, if
                            any ports are specified it requires the Protocol match
                            in the Rule to be set to \"TCP\" or \"UDP\"."
                          items:
                            anyOf:
                            - type: integer
                            - type: string
                            pattern: ^.*
                            x-kubernetes-int-or-string: true
                          type: array
                        selector:
                          description: "Selector is an optional field that contains
                            a selector expression (see Policy for sample syntax).
                            \ Only traffic that originates from (terminates at) endpoints
                            matching the selector will be matched. \n Note that: in
                            addition to the negated version of the Selector (see NotSelector
                            below), the selector expression syntax itself supports
                            negation.  The two types of negation are subtly different.
                            One negates the set of matched endpoints, the other negates
                            the whole match: \n \tSelector = \"!has(my_label)\" matches
                            packets that are from other Calico-controlled \tendpoints
                            that do not have the label \"my_label\". \n \tNotSelector
                            = \"has(my_label)\" matches packets that are not from
                            Calico-controlled \tendpoints that do have the label \"my_label\".
                            \n The effect is that the latter will accept packets from
                            non-Calico sources whereas the former is limited to packets
                            from Calico-controlled endpoints."
                          type: string
                        serviceAccounts:
                          description: ServiceAccounts is an optional field that restricts
                            the rule to only apply to traffic that originates from
                            (or terminates at) a pod running as a matching service
                            account.
                          properties:
                            names:
                              description: Names is an optional field that restricts
                                the rule to only apply to traffic that originates
                                from (or terminates at) a pod running as a service
                                account whose name is in the list.
                              items:
                                type: string
                              type: array
                            selector:
                              description: Selector is an optional field that restricts
                                the rule to only apply to traffic that originates
                                from (or terminates at) a pod running as a service
                                account that matches the given label selector. If
                                both Names and Selector are specified then they are
                                AND'ed.
                              type: string
                          type: object
                        services:
                          description: "Services is an optional field that contains
                            options for matching Kubernetes Services. If specified,
                            only traffic that originates from or terminates at endpoints
                            within the selected service(s) will be matched, and only
                            to/from each endpoint's port. \n Services cannot be specified
                            on the same rule as Selector, NotSelector, NamespaceSelector,
                            Nets, NotNets or ServiceAccounts. \n Ports and NotPorts
                            can only be specified with Services on ingress rules."
                          properties:
                            name:
                              description: Name specifies the name of a Kubernetes
                                Service to match.
                              type: string
                            namespace:
                              description: Namespace specifies the namespace of the
                                given Service. If left empty, the rule will match
                                within this policy's namespace.
                              type: string
                          type: object
                      type: object
                    http:
                      description: HTTP contains match criteria that apply to HTTP
                        requests.
                      properties:
                        methods:
                          description: Methods is an optional field that restricts
                            the rule to apply only to HTTP requests that use one of
                            the listed HTTP Methods (e.g. GET, PUT, etc.) Multiple
                            methods are OR'd together.
                          items:
                            type: string
                          type: array
                        paths:
                          description: 'Paths is an optional field that restricts
                            the rule to apply to HTTP requests that use one of the
                            listed HTTP Paths. Multiple paths are OR''d together.
                            e.g: - exact: /foo - prefix: /bar NOTE: Each entry may
                            ONLY specify either a ` + "`exact` or a `prefix`" + ` match. The
                            validator will check for it.'
                          items:
                            description: 'HTTPPath specifies an HTTP path to match.
                              It may be either of the form: exact: <path>: which matches
                              the path exactly or prefix: <path-prefix>: which matches
                              the path prefix'
                            properties:
                              exact:
                                type: string
                              prefix:
                                type: string
                            type: object
                          type: array
                      type: object
                    icmp:
                      description: ICMP is an optional field that restricts the rule
                        to apply to a specific type and code of ICMP traffic.  This
                        should only be specified if the Protocol field is set to "ICMP"
                        or "ICMPv6".
                      properties:
                        code:
                          description: Match on a specific ICMP code.  If specified,
                            the Type value must also be specified. This is a technical
                            limitation imposed by the kernel's iptables firewall,
                            which Calico uses to enforce the rule.
                          type: integer
                        type:
                          description: Match on a specific ICMP type.  For example
                            a value of 8 refers to ICMP Echo Request (i.e. pings).
                          type: integer
                      type: object
                    ipVersion:
                      description: IPVersion is an optional field that restricts the
                        rule to only match a specific IP version.
                      type: integer
                    metadata:
                      description: Metadata contains additional information for this
                        rule
                      properties:
                        annotations:
                          additionalProperties:
                            type: string
                          description: Annotations is a set of key value pairs that
                            give extra information about the rule
                          type: object
                      type: object
                    notICMP:
                      description: NotICMP is the negated version of the ICMP field.
                      properties:
                        code:
                          description: Match on a specific ICMP code.  If specified,
                            the Type value must also be specified. This is a technical
                            limitation imposed by the kernel's iptables firewall,
                            which Calico uses to enforce the rule.
                          type: integer
                        type:
                          description: Match on a specific ICMP type.  For example
                            a value of 8 refers to ICMP Echo Request (i.e. pings).
                          type: integer
                      type: object
                    notProtocol:
                      anyOf:
                      - type: integer
                      - type: string
                      description: NotProtocol is the negated version of the Protocol
                        field.
                      pattern: ^.*
                      x-kubernetes-int-or-string: true
                    protocol:
                      anyOf:
                      - type: integer
                      - type: string
                      description: "Protocol is an optional field that restricts the
                        rule to only apply to traffic of a specific IP protocol. Required
                        if any of the EntityRules contain Ports (because ports only
                        apply to certain protocols). \n Must be one of these string
                        values: \"TCP\", \"UDP\", \"ICMP\", \"ICMPv6\", \"SCTP\",
                        \"UDPLite\" or an integer in the range 1-255."
                      pattern: ^.*
                      x-kubernetes-int-or-string: true
                    source:
                      description: Source contains the match criteria that apply to
                        source entity.
                      properties:
                        namespaceSelector:
                          description: "NamespaceSelector is an optional field that
                            contains a selector expression. Only traffic that originates
                            from (or terminates at) endpoints within the selected
                            namespaces will be matched. When both NamespaceSelector
                            and another selector are defined on the same rule, then
                            only workload endpoints that are matched by both selectors
                            will be selected by the rule. \n For NetworkPolicy, an
                            empty NamespaceSelector implies that the Selector is limited
                            to selecting only workload endpoints in the same namespace
                            as the NetworkPolicy. \n For NetworkPolicy, ` + "`global()`" + `
                            NamespaceSelector implies that the Selector is limited
                            to selecting only GlobalNetworkSet or HostEndpoint. \n
                            For GlobalNetworkPolicy, an empty NamespaceSelector implies
                            the Selector applies to workload endpoints across all
                            namespaces."
                          type: string
                        nets:
                          description: Nets is an optional field that restricts the
                            rule to only apply to traffic that originates from (or
                            terminates at) IP addresses in any of the given subnets.
                          items:
                            type: string
                          type: array
                        notNets:
                          description: NotNets is the negated version of the Nets
                            field.
                          items:
                            type: string
                          type: array
                        notPorts:
                          description: NotPorts is the negated version of the Ports
                            field. Since only some protocols have ports, if any ports
                            are specified it requires the Protocol match in the Rule
                            to be set to "TCP" or "UDP".
                          items:
                            anyOf:
                            - type: integer
                            - type: string
                            pattern: ^.*
                            x-kubernetes-int-or-string: true
                          type: array
                        notSelector:
                          description: NotSelector is the negated version of the Selector
                            field.  See Selector field for subtleties with negated
                            selectors.
                          type: string
                        ports:
                          description: "Ports is an optional field that restricts
                            the rule to only apply to traffic that has a source (destination)
                            port that matches one of these ranges/values. This value
                            is a list of integers or strings that represent ranges
                            of ports. \n Since only some protocols have ports, if
                            any ports are specified it requires the Protocol match
                            in the Rule to be set to \"TCP\" or \"UDP\"."
                          items:
                            anyOf:
                            - type: integer
                            - type: string
                            pattern: ^.*
                            x-kubernetes-int-or-string: true
                          type: array
                        selector:
                          description: "Selector is an optional field that contains
                            a selector expression (see Policy for sample syntax).
                            \ Only traffic that originates from (terminates at) endpoints
                            matching the selector will be matched. \n Note that: in
                            addition to the negated version of the Selector (see NotSelector
                            below), the selector expression syntax itself supports
                            negation.  The two types of negation are subtly different.
                            One negates the set of matched endpoints, the other negates
                            the whole match: \n \tSelector = \"!has(my_label)\" matches
                            packets that are from other Calico-controlled \tendpoints
                            that do not have the label \"my_label\". \n \tNotSelector
                            = \"has(my_label)\" matches packets that are not from
                            Calico-controlled \tendpoints that do have the label \"my_label\".
                            \n The effect is that the latter will accept packets from
                            non-Calico sources whereas the former is limited to packets
                            from Calico-controlled endpoints."
                          type: string
                        serviceAccounts:
                          description: ServiceAccounts is an optional field that restricts
                            the rule to only apply to traffic that originates from
                            (or terminates at) a pod running as a matching service
                            account.
                          properties:
                            names:
                              description: Names is an optional field that restricts
                                the rule to only apply to traffic that originates
                                from (or terminates at) a pod running as a service
                                account whose name is in the list.
                              items:
                                type: string
                              type: array
                            selector:
                              description: Selector is an optional field that restricts
                                the rule to only apply to traffic that originates
                                from (or terminates at) a pod running as a service
                                account that matches the given label selector. If
                                both Names and Selector are specified then they are
                                AND'ed.
                              type: string
                          type: object
                        services:
                          description: "Services is an optional field that contains
                            options for matching Kubernetes Services. If specified,
                            only traffic that originates from or terminates at endpoints
                            within the selected service(s) will be matched, and only
                            to/from each endpoint's port. \n Services cannot be specified
                            on the same rule as Selector, NotSelector, NamespaceSelector,
                            Nets, NotNets or ServiceAccounts. \n Ports and NotPorts
                            can only be specified with Services on ingress rules."
                          properties:
                            name:
                              description: Name specifies the name of a Kubernetes
                                Service to match.
                              type: string
                            namespace:
                              description: Namespace specifies the namespace of the
                                given Service. If left empty, the rule will match
                                within this policy's namespace.
                              type: string
                          type: object
                      type: object
                  required:
                  - action
                  type: object
                type: array
              ingress:
                description: The ordered set of ingress rules.  Each rule contains
                  a set of packet match criteria and a corresponding action to apply.
                items:
                  description: "A Rule encapsulates a set of match criteria and an
                    action.  Both selector-based security Policy and security Profiles
                    reference rules - separated out as a list of rules for both ingress
                    and egress packet matching. \n Each positive match criteria has
                    a negated version, prefixed with \"Not\". All the match criteria
                    within a rule must be satisfied for a packet to match. A single
                    rule can contain the positive and negative version of a match
                    and both must be satisfied for the rule to match."
                  properties:
                    action:
                      type: string
                    destination:
                      description: Destination contains the match criteria that apply
                        to destination entity.
                      properties:
                        namespaceSelector:
                          description: "NamespaceSelector is an optional field that
                            contains a selector expression. Only traffic that originates
                            from (or terminates at) endpoints within the selected
                            namespaces will be matched. When both NamespaceSelector
                            and another selector are defined on the same rule, then
                            only workload endpoints that are matched by both selectors
                            will be selected by the rule. \n For NetworkPolicy, an
                            empty NamespaceSelector implies that the Selector is limited
                            to selecting only workload endpoints in the same namespace
                            as the NetworkPolicy. \n For NetworkPolicy, ` + "`global()`" + `
                            NamespaceSelector implies that the Selector is limited
                            to selecting only GlobalNetworkSet or HostEndpoint. \n
                            For GlobalNetworkPolicy, an empty NamespaceSelector implies
                            the Selector applies to workload endpoints across all
                            namespaces."
                          type: string
                        nets:
                          description: Nets is an optional field that restricts the
                            rule to only apply to traffic that originates from (or
                            terminates at) IP addresses in any of the given subnets.
                          items:
                            type: string
                          type: array
                        notNets:
                          description: NotNets is the negated version of the Nets
                            field.
                          items:
                            type: string
                          type: array
                        notPorts:
                          description: NotPorts is the negated version of the Ports
                            field. Since only some protocols have ports, if any ports
                            are specified it requires the Protocol match in the Rule
                            to be set to "TCP" or "UDP".
                          items:
                            anyOf:
                            - type: integer
                            - type: string
                            pattern: ^.*
                            x-kubernetes-int-or-string: true
                          type: array
                        notSelector:
                          description: NotSelector is the negated version of the Selector
                            field.  See Selector field for subtleties with negated
                            selectors.
                          type: string
                        ports:
                          description: "Ports is an optional field that restricts
                            the rule to only apply to traffic that has a source (destination)
                            port that matches one of these ranges/values. This value
                            is a list of integers or strings that represent ranges
                            of ports. \n Since only some protocols have ports, if
                            any ports are specified it requires the Protocol match
                            in the Rule to be set to \"TCP\" or \"UDP\"."
                          items:
                            anyOf:
                            - type: integer
                            - type: string
                            pattern: ^.*
                            x-kubernetes-int-or-string: true
                          type: array
                        selector:
                          description: "Selector is an optional field that contains
                            a selector expression (see Policy for sample syntax).
                            \ Only traffic that originates from (terminates at) endpoints
                            matching the selector will be matched. \n Note that: in
                            addition to the negated version of the Selector (see NotSelector
                            below), the selector expression syntax itself supports
                            negation.  The two types of negation are subtly different.
                            One negates the set of matched endpoints, the other negates
                            the whole match: \n \tSelector = \"!has(my_label)\" matches
                            packets that are from other Calico-controlled \tendpoints
                            that do not have the label \"my_label\". \n \tNotSelector
                            = \"has(my_label)\" matches packets that are not from
                            Calico-controlled \tendpoints that do have the label \"my_label\".
                            \n The effect is that the latter will accept packets from
                            non-Calico sources whereas the former is limited to packets
                            from Calico-controlled endpoints."
                          type: string
                        serviceAccounts:
                          description: ServiceAccounts is an optional field that restricts
                            the rule to only apply to traffic that originates from
                            (or terminates at) a pod running as a matching service
                            account.
                          properties:
                            names:
                              description: Names is an optional field that restricts
                                the rule to only apply to traffic that originates
                                from (or terminates at) a pod running as a service
                                account whose name is in the list.
                              items:
                                type: string
                              type: array
                            selector:
                              description: Selector is an optional field that restricts
                                the rule to only apply to traffic that originates
                                from (or terminates at) a pod running as a service
                                account that matches the given label selector. If
                                both Names and Selector are specified then they are
                                AND'ed.
                              type: string
                          type: object
                        services:
                          description: "Services is an optional field that contains
                            options for matching Kubernetes Services. If specified,
                            only traffic that originates from or terminates at endpoints
                            within the selected service(s) will be matched, and only
                            to/from each endpoint's port. \n Services cannot be specified
                            on the same rule as Selector, NotSelector, NamespaceSelector,
                            Nets, NotNets or ServiceAccounts. \n Ports and NotPorts
                            can only be specified with Services on ingress rules."
                          properties:
                            name:
                              description: Name specifies the name of a Kubernetes
                                Service to match.
                              type: string
                            namespace:
                              description: Namespace specifies the namespace of the
                                given Service. If left empty, the rule will match
                                within this policy's namespace.
                              type: string
                          type: object
                      type: object
                    http:
                      description: HTTP contains match criteria that apply to HTTP
                        requests.
                      properties:
                        methods:
                          description: Methods is an optional field that restricts
                            the rule to apply only to HTTP requests that use one of
                            the listed HTTP Methods (e.g. GET, PUT, etc.) Multiple
                            methods are OR'd together.
                          items:
                            type: string
                          type: array
                        paths:
                          description: 'Paths is an optional field that restricts
                            the rule to apply to HTTP requests that use one of the
                            listed HTTP Paths. Multiple paths are OR''d together.
                            e.g: - exact: /foo - prefix: /bar NOTE: Each entry may
                            ONLY specify either a ` + "`exact` or a `prefix`" + ` match. The
                            validator will check for it.'
                          items:
                            description: 'HTTPPath specifies an HTTP path to match.
                              It may be either of the form: exact: <path>: which matches
                              the path exactly or prefix: <path-prefix>: which matches
                              the path prefix'
                            properties:
                              exact:
                                type: string
                              prefix:
                                type: string
                            type: object
                          type: array
                      type: object
                    icmp:
                      description: ICMP is an optional field that restricts the rule
                        to apply to a specific type and code of ICMP traffic.  This
                        should only be specified if the Protocol field is set to "ICMP"
                        or "ICMPv6".
                      properties:
                        code:
                          description: Match on a specific ICMP code.  If specified,
                            the Type value must also be specified. This is a technical
                            limitation imposed by the kernel's iptables firewall,
                            which Calico uses to enforce the rule.
                          type: integer
                        type:
                          description: Match on a specific ICMP type.  For example
                            a value of 8 refers to ICMP Echo Request (i.e. pings).
                          type: integer
                      type: object
                    ipVersion:
                      description: IPVersion is an optional field that restricts the
                        rule to only match a specific IP version.
                      type: integer
                    metadata:
                      description: Metadata contains additional information for this
                        rule
                      properties:
                        annotations:
                          additionalProperties:
                            type: string
                          description: Annotations is a set of key value pairs that
                            give extra information about the rule
                          type: object
                      type: object
                    notICMP:
                      description: NotICMP is the negated version of the ICMP field.
                      properties:
                        code:
                          description: Match on a specific ICMP code.  If specified,
                            the Type value must also be specified. This is a technical
                            limitation imposed by the kernel's iptables firewall,
                            which Calico uses to enforce the rule.
                          type: integer
                        type:
                          description: Match on a specific ICMP type.  For example
                            a value of 8 refers to ICMP Echo Request (i.e. pings).
                          type: integer
                      type: object
                    notProtocol:
                      anyOf:
                      - type: integer
                      - type: string
                      description: NotProtocol is the negated version of the Protocol
                        field.
                      pattern: ^.*
                      x-kubernetes-int-or-string: true
                    protocol:
                      anyOf:
                      - type: integer
                      - type: string
                      description: "Protocol is an optional field that restricts the
                        rule to only apply to traffic of a specific IP protocol. Required
                        if any of the EntityRules contain Ports (because ports only
                        apply to certain protocols). \n Must be one of these string
                        values: \"TCP\", \"UDP\", \"ICMP\", \"ICMPv6\", \"SCTP\",
                        \"UDPLite\" or an integer in the range 1-255."
                      pattern: ^.*
                      x-kubernetes-int-or-string: true
                    source:
                      description: Source contains the match criteria that apply to
                        source entity.
                      properties:
                        namespaceSelector:
                          description: "NamespaceSelector is an optional field that
                            contains a selector expression. Only traffic that originates
                            from (or terminates at) endpoints within the selected
                            namespaces will be matched. When both NamespaceSelector
                            and another selector are defined on the same rule, then
                            only workload endpoints that are matched by both selectors
                            will be selected by the rule. \n For NetworkPolicy, an
                            empty NamespaceSelector implies that the Selector is limited
                            to selecting only workload endpoints in the same namespace
                            as the NetworkPolicy. \n For NetworkPolicy, ` + "`global()`" + `
                            NamespaceSelector implies that the Selector is limited
                            to selecting only GlobalNetworkSet or HostEndpoint. \n
                            For GlobalNetworkPolicy, an empty NamespaceSelector implies
                            the Selector applies to workload endpoints across all
                            namespaces."
                          type: string
                        nets:
                          description: Nets is an optional field that restricts the
                            rule to only apply to traffic that originates from (or
                            terminates at) IP addresses in any of the given subnets.
                          items:
                            type: string
                          type: array
                        notNets:
                          description: NotNets is the negated version of the Nets
                            field.
                          items:
                            type: string
                          type: array
                        notPorts:
                          description: NotPorts is the negated version of the Ports
                            field. Since only some protocols have ports, if any ports
                            are specified it requires the Protocol match in the Rule
                            to be set to "TCP" or "UDP".
                          items:
                            anyOf:
                            - type: integer
                            - type: string
                            pattern: ^.*
                            x-kubernetes-int-or-string: true
                          type: array
                        notSelector:
                          description: NotSelector is the negated version of the Selector
                            field.  See Selector field for subtleties with negated
                            selectors.
                          type: string
                        ports:
                          description: "Ports is an optional field that restricts
                            the rule to only apply to traffic that has a source (destination)
                            port that matches one of these ranges/values. This value
                            is a list of integers or strings that represent ranges
                            of ports. \n Since only some protocols have ports, if
                            any ports are specified it requires the Protocol match
                            in the Rule to be set to \"TCP\" or \"UDP\"."
                          items:
                            anyOf:
                            - type: integer
                            - type: string
                            pattern: ^.*
                            x-kubernetes-int-or-string: true
                          type: array
                        selector:
                          description: "Selector is an optional field that contains
                            a selector expression (see Policy for sample syntax).
                            \ Only traffic that originates from (terminates at) endpoints
                            matching the selector will be matched. \n Note that: in
                            addition to the negated version of the Selector (see NotSelector
                            below), the selector expression syntax itself supports
                            negation.  The two types of negation are subtly different.
                            One negates the set of matched endpoints, the other negates
                            the whole match: \n \tSelector = \"!has(my_label)\" matches
                            packets that are from other Calico-controlled \tendpoints
                            that do not have the label \"my_label\". \n \tNotSelector
                            = \"has(my_label)\" matches packets that are not from
                            Calico-controlled \tendpoints that do have the label \"my_label\".
                            \n The effect is that the latter will accept packets from
                            non-Calico sources whereas the former is limited to packets
                            from Calico-controlled endpoints."
                          type: string
                        serviceAccounts:
                          description: ServiceAccounts is an optional field that restricts
                            the rule to only apply to traffic that originates from
                            (or terminates at) a pod running as a matching service
                            account.
                          properties:
                            names:
                              description: Names is an optional field that restricts
                                the rule to only apply to traffic that originates
                                from (or terminates at) a pod running as a service
                                account whose name is in the list.
                              items:
                                type: string
                              type: array
                            selector:
                              description: Selector is an optional field that restricts
                                the rule to only apply to traffic that originates
                                from (or terminates at) a pod running as a service
                                account that matches the given label selector. If
                                both Names and Selector are specified then they are
                                AND'ed.
                              type: string
                          type: object
                        services:
                          description: "Services is an optional field that contains
                            options for matching Kubernetes Services. If specified,
                            only traffic that originates from or terminates at endpoints
                            within the selected service(s) will be matched, and only
                            to/from each endpoint's port. \n Services cannot be specified
                            on the same rule as Selector, NotSelector, NamespaceSelector,
                            Nets, NotNets or ServiceAccounts. \n Ports and NotPorts
                            can only be specified with Services on ingress rules."
                          properties:
                            name:
                              description: Name specifies the name of a Kubernetes
                                Service to match.
                              type: string
                            namespace:
                              description: Namespace specifies the namespace of the
                                given Service. If left empty, the rule will match
                                within this policy's namespace.
                              type: string
                          type: object
                      type: object
                  required:
                  - action
                  type: object
                type: array
              order:
                description: Order is an optional field that specifies the order in
                  which the policy is applied. Policies with higher "order" are applied
                  after those with lower order.  If the order is omitted, it may be
                  considered to be "infinite" - i.e. the policy will be applied last.  Policies
                  with identical order will be applied in alphanumerical order based
                  on the Policy "Name".
                type: number
              selector:
                description: "The selector is an expression used to pick pick out
                  the endpoints that the policy should be applied to. \n Selector
                  expressions follow this syntax: \n \tlabel == \"string_literal\"
                  \ ->  comparison, e.g. my_label == \"foo bar\" \tlabel != \"string_literal\"
                  \  ->  not equal; also matches if label is not present \tlabel in
                  { \"a\", \"b\", \"c\", ... }  ->  true if the value of label X is
                  one of \"a\", \"b\", \"c\" \tlabel not in { \"a\", \"b\", \"c\",
                  ... }  ->  true if the value of label X is not one of \"a\", \"b\",
                  \"c\" \thas(label_name)  -> True if that label is present \t! expr
                  -> negation of expr \texpr && expr  -> Short-circuit and \texpr
                  || expr  -> Short-circuit or \t( expr ) -> parens for grouping \tall()
                  or the empty selector -> matches all endpoints. \n Label names are
                  allowed to contain alphanumerics, -, _ and /. String literals are
                  more permissive but they do not support escape characters. \n Examples
                  (with made-up labels): \n \ttype == \"webserver\" && deployment
                  == \"prod\" \ttype in {\"frontend\", \"backend\"} \tdeployment !=
                  \"dev\" \t! has(label_name)"
                type: string
              serviceAccountSelector:
                description: ServiceAccountSelector is an optional field for an expression
                  used to select a pod based on service accounts.
                type: string
              types:
                description: "Types indicates whether this policy applies to ingress,
                  or to egress, or to both.  When not explicitly specified (and so
                  the value on creation is empty or nil), Calico defaults Types according
                  to what Ingress and Egress are present in the policy.  The default
                  is: \n - [ PolicyTypeIngress ], if there are no Egress rules (including
                  the case where there are   also no Ingress rules) \n - [ PolicyTypeEgress
                  ], if there are Egress rules but no Ingress rules \n - [ PolicyTypeIngress,
                  PolicyTypeEgress ], if there are both Ingress and Egress rules.
                  \n When the policy is read back again, Types will always be one
                  of these values, never empty or nil."
                items:
                  description: PolicyType enumerates the possible values of the PolicySpec
                    Types field.
                  type: string
                type: array
            type: object
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`
