/* =============================================================
* @Author:  Wayne Wang <yuying.wang@sincerecloud.com>
* @ Last Modified At: 2023.06.01
* @Copyright (c) 2023 Sincerecloud
* @HomePage: https://www.sincerecloud.com/
*
*  定义部署calico的yaml内容，本内容使用模板变量，这是第一部分。
*  程序应当按顺序先后apply各部分的yaml内容。
 */

package calico

// 部署calico的yaml内容第一部分
var tplPart1 string = `
---
# Source: calico/templates/calico-kube-controllers.yaml
# This manifest creates a Pod Disruption Budget for Controller to allow K8s Cluster Autoscaler to evict

apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: calico-kube-controllers
  namespace: kube-system
  labels:
    k8s-app: calico-kube-controllers
spec:
  maxUnavailable: 1
  selector:
    matchLabels:
      k8s-app: calico-kube-controllers
---
# Source: calico/templates/calico-typha.yaml
# This manifest creates a Pod Disruption Budget for Typha to allow K8s Cluster Autoscaler to evict

apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: calico-typha
  namespace: kube-system
  labels:
    k8s-app: calico-typha
spec:
  maxUnavailable: 1
  selector:
    matchLabels:
      k8s-app: calico-typha
---
# Source: calico/templates/calico-kube-controllers.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: calico-kube-controllers
  namespace: kube-system
---
# Source: calico/templates/calico-node.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: calico-node
  namespace: kube-system
---
# Source: calico/templates/calico-config.yaml
# This ConfigMap is used to configure a self-hosted Calico installation.
kind: ConfigMap
apiVersion: v1
metadata:
  name: calico-config
  namespace: kube-system
data:
  # You must set a non-zero value for Typha replicas below.
  typha_service_name: "calico-typha"
  # Configure the backend to use.
  # use vxlan mode
  calico_backend: "{{ .CalicoBackend }}"

  # Configure the MTU to use for workload interfaces and tunnels.
  # By default, MTU is auto-detected, and explicitly setting this field should not be required.
  # You can override auto-detection by providing a non-zero value.
  veth_mtu: "0"

  # The CNI network configuration to install on each node. The special
  # values in this config will be automatically populated.
  cni_network_config: |-
    {
      "name": "k8s-pod-network",
      "cniVersion": "0.3.1",
      "plugins": [
        {
          "type": "calico",
          "log_level": "info",
          "log_file_path": "/var/log/calico/cni/cni.log",
          "datastore_type": "kubernetes",
          "nodename": "__KUBERNETES_NODE_NAME__",
          "mtu": __CNI_MTU__,
          "ipam": {
              "type": "calico-ipam"
          },
          "policy": {
              "type": "k8s"
          },
          "kubernetes": {
              "kubeconfig": "__KUBECONFIG_FILEPATH__"
          }
        },
        {
          "type": "portmap",
          "snat": true,
          "capabilities": {"portMappings": true}
        },
        {
          "type": "bandwidth",
          "capabilities": {"bandwidth": true}
        }
      ]
    }
---
# Source: calico/templates/kdd-crds.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: bgpconfigurations.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: BGPConfiguration
    listKind: BGPConfigurationList
    plural: bgpconfigurations
    singular: bgpconfiguration
  preserveUnknownFields: false
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        description: BGPConfiguration contains the configuration for any BGP routing.
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
            description: BGPConfigurationSpec contains the values of the BGP configuration.
            properties:
              asNumber:
                description: 'ASNumber is the default AS number used by a node. [Default:
                  64512]'
                format: int32
                type: integer
              bindMode:
                description: BindMode indicates whether to listen for BGP connections
                  on all addresses (None) or only on the node's canonical IP address
                  Node.Spec.BGP.IPvXAddress (NodeIP). Default behaviour is to listen
                  for BGP connections on all addresses.
                type: string
              communities:
                description: Communities is a list of BGP community values and their
                  arbitrary names for tagging routes.
                items:
                  description: Community contains standard or large community value
                    and its name.
                  properties:
                    name:
                      description: Name given to community value.
                      type: string
                    value:
                      description: Value must be of format 'aa:nn' or 'aa:nn:mm'.
                        For standard community use 'aa:nn' format, where 'aa' and
                        'nn' are 16 bit number. For large community use 'aa:nn:mm'
                        format, where 'aa', 'nn' and 'mm' are 32 bit number. Where,
                        'aa' is an AS Number, 'nn' and mm are per-AS identifier.
                      pattern: ^(\d+):(\d+)$|^(\d+):(\d+):(\d+)$
                      type: string
                  type: object
                type: array
              listenPort:
                description: ListenPort is the port where BGP protocol should listen.
                  Defaults to 179
                maximum: 65535
                minimum: 1
                type: integer
              logSeverityScreen:
                description: 'LogSeverityScreen is the log severity above which logs
                  are sent to the stdout. [Default: INFO]'
                type: string
              nodeMeshMaxRestartTime:
                description: Time to allow for software restart for node-to-mesh peerings.  When
                  specified, this is configured as the graceful restart timeout.  When
                  not specified, the BIRD default of 120s is used. This field can
                  only be set on the default BGPConfiguration instance and requires
                  that NodeMesh is enabled
                type: string
              nodeMeshPassword:
                description: Optional BGP password for full node-to-mesh peerings.
                  This field can only be set on the default BGPConfiguration instance
                  and requires that NodeMesh is enabled
                properties:
                  secretKeyRef:
                    description: Selects a key of a secret in the node pod's namespace.
                    properties:
                      key:
                        description: The key of the secret to select from.  Must be
                          a valid secret key.
                        type: string
                      name:
                        description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                          TODO: Add other useful fields. apiVersion, kind, uid?'
                        type: string
                      optional:
                        description: Specify whether the Secret or its key must be
                          defined
                        type: boolean
                    required:
                    - key
                    type: object
                type: object
              nodeToNodeMeshEnabled:
                description: 'NodeToNodeMeshEnabled sets whether full node to node
                  BGP mesh is enabled. [Default: true]'
                type: boolean
              prefixAdvertisements:
                description: PrefixAdvertisements contains per-prefix advertisement
                  configuration.
                items:
                  description: PrefixAdvertisement configures advertisement properties
                    for the specified CIDR.
                  properties:
                    cidr:
                      description: CIDR for which properties should be advertised.
                      type: string
                    communities:
                      description: Communities can be list of either community names
                        already defined in "Specs.Communities" or community value
                        of format "aa:nn" or 'aa:nn:mm'. For standard community use
                        'aa:nn' format, where 'aa' and 'nn' are 16 bit number. For
                        large community use 'aa:nn:mm' format, where 'aa'', 'nn' and
                        'mm' are 32 bit number. Where,'aa' is an AS Number, 'nn' and
                        'mm' are per-AS identifier.
                      items:
                        type: string
                      type: array
                  type: object
                type: array
              serviceClusterIPs:
                description: ServiceClusterIPs are the CIDR blocks from which service
                  cluster IPs are allocated. If specified, Calico will advertise these
                  blocks, as well as any cluster IPs within them.
                items:
                  description: ServiceClusterIPBlock represents a single allowed ClusterIP
                    CIDR block.
                  properties:
                    cidr:
                      type: string
                  type: object
                type: array
              serviceExternalIPs:
                description: ServiceExternalIPs are the CIDR blocks for Kubernetes
                  Service External IPs. Kubernetes Service ExternalIPs will only be
                  advertised if they are within one of these blocks.
                items:
                  description: ServiceExternalIPBlock represents a single allowed
                    External IP CIDR block.
                  properties:
                    cidr:
                      type: string
                  type: object
                type: array
              serviceLoadBalancerIPs:
                description: ServiceLoadBalancerIPs are the CIDR blocks for Kubernetes
                  Service LoadBalancer IPs. Kubernetes Service status.LoadBalancer.Ingress
                  IPs will only be advertised if they are within one of these blocks.
                items:
                  description: ServiceLoadBalancerIPBlock represents a single allowed
                    LoadBalancer IP CIDR block.
                  properties:
                    cidr:
                      type: string
                  type: object
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
  name: bgppeers.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: BGPPeer
    listKind: BGPPeerList
    plural: bgppeers
    singular: bgppeer
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
            description: BGPPeerSpec contains the specification for a BGPPeer resource.
            properties:
              asNumber:
                description: The AS Number of the peer.
                format: int32
                type: integer
              keepOriginalNextHop:
                description: Option to keep the original nexthop field when routes
                  are sent to a BGP Peer. Setting "true" configures the selected BGP
                  Peers node to use the "next hop keep;" instead of "next hop self;"(default)
                  in the specific branch of the Node on "bird.cfg".
                type: boolean
              maxRestartTime:
                description: Time to allow for software restart.  When specified,
                  this is configured as the graceful restart timeout.  When not specified,
                  the BIRD default of 120s is used.
                type: string
              node:
                description: The node name identifying the Calico node instance that
                  is targeted by this peer. If this is not set, and no nodeSelector
                  is specified, then this BGP peer selects all nodes in the cluster.
                type: string
              nodeSelector:
                description: Selector for the nodes that should have this peering.  When
                  this is set, the Node field must be empty.
                type: string
              numAllowedLocalASNumbers:
                description: Maximum number of local AS numbers that are allowed in
                  the AS path for received routes. This removes BGP loop prevention
                  and should only be used if absolutely necesssary.
                format: int32
                type: integer
              password:
                description: Optional BGP password for the peerings generated by this
                  BGPPeer resource.
                properties:
                  secretKeyRef:
                    description: Selects a key of a secret in the node pod's namespace.
                    properties:
                      key:
                        description: The key of the secret to select from.  Must be
                          a valid secret key.
                        type: string
                      name:
                        description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                          TODO: Add other useful fields. apiVersion, kind, uid?'
                        type: string
                      optional:
                        description: Specify whether the Secret or its key must be
                          defined
                        type: boolean
                    required:
                    - key
                    type: object
                type: object
              peerIP:
                description: The IP address of the peer followed by an optional port
                  number to peer with. If port number is given, format should be ` + "`[<IPv6>]:port`" + `
                  or ` + "`<IPv4>:<port>` for IPv4. If optional port number is not set," + `
                  and this peer IP and ASNumber belongs to a calico/node with ListenPort
                  set in BGPConfiguration, then we use that port to peer.
                type: string
              peerSelector:
                description: Selector for the remote nodes to peer with.  When this
                  is set, the PeerIP and ASNumber fields must be empty.  For each
                  peering between the local node and selected remote nodes, we configure
                  an IPv4 peering if both ends have NodeBGPSpec.IPv4Address specified,
                  and an IPv6 peering if both ends have NodeBGPSpec.IPv6Address specified.  The
                  remote AS number comes from the remote node's NodeBGPSpec.ASNumber,
                  or the global default if that is not set.
                type: string
              sourceAddress:
                description: Specifies whether and how to configure a source address
                  for the peerings generated by this BGPPeer resource.  Default value
                  "UseNodeIP" means to configure the node IP as the source address.  "None"
                  means not to configure a source address.
                type: string
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
  name: blockaffinities.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: BlockAffinity
    listKind: BlockAffinityList
    plural: blockaffinities
    singular: blockaffinity
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
            description: BlockAffinitySpec contains the specification for a BlockAffinity
              resource.
            properties:
              cidr:
                type: string
              deleted:
                description: Deleted indicates that this block affinity is being deleted.
                  This field is a string for compatibility with older releases that
                  mistakenly treat this field as a string.
                type: string
              node:
                type: string
              state:
                type: string
            required:
            - cidr
            - deleted
            - node
            - state
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
  name: caliconodestatuses.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: CalicoNodeStatus
    listKind: CalicoNodeStatusList
    plural: caliconodestatuses
    singular: caliconodestatus
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
            description: CalicoNodeStatusSpec contains the specification for a CalicoNodeStatus
              resource.
            properties:
              classes:
                description: Classes declares the types of information to monitor
                  for this calico/node, and allows for selective status reporting
                  about certain subsets of information.
                items:
                  type: string
                type: array
              node:
                description: The node name identifies the Calico node instance for
                  node status.
                type: string
              updatePeriodSeconds:
                description: UpdatePeriodSeconds is the period at which CalicoNodeStatus
                  should be updated. Set to 0 to disable CalicoNodeStatus refresh.
                  Maximum update period is one day.
                format: int32
                type: integer
            type: object
          status:
            description: CalicoNodeStatusStatus defines the observed state of CalicoNodeStatus.
              No validation needed for status since it is updated by Calico.
            properties:
              agent:
                description: Agent holds agent status on the node.
                properties:
                  birdV4:
                    description: BIRDV4 represents the latest observed status of bird4.
                    properties:
                      lastBootTime:
                        description: LastBootTime holds the value of lastBootTime
                          from bird.ctl output.
                        type: string
                      lastReconfigurationTime:
                        description: LastReconfigurationTime holds the value of lastReconfigTime
                          from bird.ctl output.
                        type: string
                      routerID:
                        description: Router ID used by bird.
                        type: string
                      state:
                        description: The state of the BGP Daemon.
                        type: string
                      version:
                        description: Version of the BGP daemon
                        type: string
                    type: object
                  birdV6:
                    description: BIRDV6 represents the latest observed status of bird6.
                    properties:
                      lastBootTime:
                        description: LastBootTime holds the value of lastBootTime
                          from bird.ctl output.
                        type: string
                      lastReconfigurationTime:
                        description: LastReconfigurationTime holds the value of lastReconfigTime
                          from bird.ctl output.
                        type: string
                      routerID:
                        description: Router ID used by bird.
                        type: string
                      state:
                        description: The state of the BGP Daemon.
                        type: string
                      version:
                        description: Version of the BGP daemon
                        type: string
                    type: object
                type: object
              bgp:
                description: BGP holds node BGP status.
                properties:
                  numberEstablishedV4:
                    description: The total number of IPv4 established bgp sessions.
                    type: integer
                  numberEstablishedV6:
                    description: The total number of IPv6 established bgp sessions.
                    type: integer
                  numberNotEstablishedV4:
                    description: The total number of IPv4 non-established bgp sessions.
                    type: integer
                  numberNotEstablishedV6:
                    description: The total number of IPv6 non-established bgp sessions.
                    type: integer
                  peersV4:
                    description: PeersV4 represents IPv4 BGP peers status on the node.
                    items:
                      description: CalicoNodePeer contains the status of BGP peers
                        on the node.
                      properties:
                        peerIP:
                          description: IP address of the peer whose condition we are
                            reporting.
                          type: string
                        since:
                          description: Since the state or reason last changed.
                          type: string
                        state:
                          description: State is the BGP session state.
                          type: string
                        type:
                          description: Type indicates whether this peer is configured
                            via the node-to-node mesh, or via en explicit global or
                            per-node BGPPeer object.
                          type: string
                      type: object
                    type: array
                  peersV6:
                    description: PeersV6 represents IPv6 BGP peers status on the node.
                    items:
                      description: CalicoNodePeer contains the status of BGP peers
                        on the node.
                      properties:
                        peerIP:
                          description: IP address of the peer whose condition we are
                            reporting.
                          type: string
                        since:
                          description: Since the state or reason last changed.
                          type: string
                        state:
                          description: State is the BGP session state.
                          type: string
                        type:
                          description: Type indicates whether this peer is configured
                            via the node-to-node mesh, or via en explicit global or
                            per-node BGPPeer object.
                          type: string
                      type: object
                    type: array
                required:
                - numberEstablishedV4
                - numberEstablishedV6
                - numberNotEstablishedV4
                - numberNotEstablishedV6
                type: object
              lastUpdated:
                description: LastUpdated is a timestamp representing the server time
                  when CalicoNodeStatus object last updated. It is represented in
                  RFC3339 form and is in UTC.
                format: date-time
                nullable: true
                type: string
              routes:
                description: Routes reports routes known to the Calico BGP daemon
                  on the node.
                properties:
                  routesV4:
                    description: RoutesV4 represents IPv4 routes on the node.
                    items:
                      description: CalicoNodeRoute contains the status of BGP routes
                        on the node.
                      properties:
                        destination:
                          description: Destination of the route.
                          type: string
                        gateway:
                          description: Gateway for the destination.
                          type: string
                        interface:
                          description: Interface for the destination
                          type: string
                        learnedFrom:
                          description: LearnedFrom contains information regarding
                            where this route originated.
                          properties:
                            peerIP:
                              description: If sourceType is NodeMesh or BGPPeer, IP
                                address of the router that sent us this route.
                              type: string
                            sourceType:
                              description: Type of the source where a route is learned
                                from.
                              type: string
                          type: object
                        type:
                          description: Type indicates if the route is being used for
                            forwarding or not.
                          type: string
                      type: object
                    type: array
                  routesV6:
                    description: RoutesV6 represents IPv6 routes on the node.
                    items:
                      description: CalicoNodeRoute contains the status of BGP routes
                        on the node.
                      properties:
                        destination:
                          description: Destination of the route.
                          type: string
                        gateway:
                          description: Gateway for the destination.
                          type: string
                        interface:
                          description: Interface for the destination
                          type: string
                        learnedFrom:
                          description: LearnedFrom contains information regarding
                            where this route originated.
                          properties:
                            peerIP:
                              description: If sourceType is NodeMesh or BGPPeer, IP
                                address of the router that sent us this route.
                              type: string
                            sourceType:
                              description: Type of the source where a route is learned
                                from.
                              type: string
                          type: object
                        type:
                          description: Type indicates if the route is being used for
                            forwarding or not.
                          type: string
                      type: object
                    type: array
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
  name: clusterinformations.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: ClusterInformation
    listKind: ClusterInformationList
    plural: clusterinformations
    singular: clusterinformation
  preserveUnknownFields: false
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        description: ClusterInformation contains the cluster specific information.
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
            description: ClusterInformationSpec contains the values of describing
              the cluster.
            properties:
              calicoVersion:
                description: CalicoVersion is the version of Calico that the cluster
                  is running
                type: string
              clusterGUID:
                description: ClusterGUID is the GUID of the cluster
                type: string
              clusterType:
                description: ClusterType describes the type of the cluster
                type: string
              datastoreReady:
                description: DatastoreReady is used during significant datastore migrations
                  to signal to components such as Felix that it should wait before
                  accessing the datastore.
                type: boolean
              variant:
                description: Variant declares which variant of Calico should be active.
                type: string
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
  name: felixconfigurations.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: FelixConfiguration
    listKind: FelixConfigurationList
    plural: felixconfigurations
    singular: felixconfiguration
  preserveUnknownFields: false
  scope: Cluster
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        description: Felix Configuration contains the configuration for Felix.
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
            description: FelixConfigurationSpec contains the values of the Felix configuration.
            properties:
              allowIPIPPacketsFromWorkloads:
                description: 'AllowIPIPPacketsFromWorkloads controls whether Felix
                  will add a rule to drop IPIP encapsulated traffic from workloads
                  [Default: false]'
                type: boolean
              allowVXLANPacketsFromWorkloads:
                description: 'AllowVXLANPacketsFromWorkloads controls whether Felix
                  will add a rule to drop VXLAN encapsulated traffic from workloads
                  [Default: false]'
                type: boolean
              awsSrcDstCheck:
                description: 'Set source-destination-check on AWS EC2 instances. Accepted
                  value must be one of "DoNothing", "Enable" or "Disable". [Default:
                  DoNothing]'
                enum:
                - DoNothing
                - Enable
                - Disable
                type: string
              bpfConnectTimeLoadBalancingEnabled:
                description: 'BPFConnectTimeLoadBalancingEnabled when in BPF mode,
                  controls whether Felix installs the connection-time load balancer.  The
                  connect-time load balancer is required for the host to be able to
                  reach Kubernetes services and it improves the performance of pod-to-service
                  connections.  The only reason to disable it is for debugging purposes.  [Default:
                  true]'
                type: boolean
              bpfDataIfacePattern:
                description: 'BPFDataIfacePattern is a regular expression that controls
                  which interfaces Felix should attach BPF programs to in order to
                  catch traffic to/from the network.  This needs to match the interfaces
                  that Calico workload traffic flows over as well as any interfaces
                  that handle incoming traffic to nodeports and services from outside
                  the cluster.  It should not match the workload interfaces (usually
                  named cali...).'
                type: string
              bpfDisableUnprivileged:
                description: 'BPFDisableUnprivileged, if enabled, Felix sets the kernel.unprivileged_bpf_disabled
                  sysctl to disable unprivileged use of BPF.  This ensures that unprivileged
                  users cannot access Calico''s BPF maps and cannot insert their own
                  BPF programs to interfere with Calico''s. [Default: true]'
                type: boolean
              bpfEnabled:
                description: 'BPFEnabled, if enabled Felix will use the BPF dataplane.
                  [Default: false]'
                type: boolean
              bpfEnforceRPF:
                description: 'BPFEnforceRPF enforce strict RPF on all interfaces with
                  BPF programs regardless of what is the per-interfaces or global
                  setting. Possible values are Disabled or Strict. [Default: Strict]'
                type: string
              bpfExtToServiceConnmark:
                description: 'BPFExtToServiceConnmark in BPF mode, control a 32bit
                  mark that is set on connections from an external client to a local
                  service. This mark allows us to control how packets of that connection
                  are routed within the host and how is routing interpreted by RPF
                  check. [Default: 0]'
                type: integer
              bpfExternalServiceMode:
                description: 'BPFExternalServiceMode in BPF mode, controls how connections
                  from outside the cluster to services (node ports and cluster IPs)
                  are forwarded to remote workloads.  If set to "Tunnel" then both
                  request and response traffic is tunneled to the remote node.  If
                  set to "DSR", the request traffic is tunneled but the response traffic
                  is sent directly from the remote node.  In "DSR" mode, the remote
                  node appears to use the IP of the ingress node; this requires a
                  permissive L2 network.  [Default: Tunnel]'
                type: string
              bpfKubeProxyEndpointSlicesEnabled:
                description: 'BPFKubeProxyEndpointSlicesEnabled in BPF mode, controls
                  whether Felixs embedded kube-proxy accepts EndpointSlices or not.'
                type: boolean
              bpfKubeProxyIptablesCleanupEnabled:
                description: 'BPFKubeProxyIptablesCleanupEnabled, if enabled in BPF
                  mode, Felix will proactively clean up the upstream Kubernetes kube-proxy''s
                  iptables chains.  Should only be enabled if kube-proxy is not running.  [Default:
                  true]'
                type: boolean
              bpfKubeProxyMinSyncPeriod:
                description: 'BPFKubeProxyMinSyncPeriod, in BPF mode, controls the
                  minimum time between updates to the dataplane for Felix''s embedded
                  kube-proxy.  Lower values give reduced set-up latency.  Higher values
                  reduce Felix CPU usage by batching up more work.  [Default: 1s]'
                type: string
              bpfLogLevel:
                description: 'BPFLogLevel controls the log level of the BPF programs
                  when in BPF dataplane mode.  One of "Off", "Info", or "Debug".  The
                  logs are emitted to the BPF trace pipe, accessible with the command
                  ` + "`tc exec bpf debug`." + ` [Default: Off].'
                type: string
              bpfMapSizeConntrack:
                description: 'BPFMapSizeConntrack sets the size for the conntrack
                  map.  This map must be large enough to hold an entry for each active
                  connection.  Warning: changing the size of the conntrack map can
                  cause disruption.'
                type: integer
              bpfMapSizeIPSets:
                description: BPFMapSizeIPSets sets the size for ipsets map.  The IP
                  sets map must be large enough to hold an entry for each endpoint
                  matched by every selector in the source/destination matches in network
                  policy.  Selectors such as "all()" can result in large numbers of
                  entries (one entry per endpoint in that case).
                type: integer
              bpfMapSizeIfState:
                description: BPFMapSizeIfState sets the size for ifstate map.  The
                  ifstate map must be large enough to hold an entry for each device
                  (host + workloads) on a host.
                type: integer
              bpfMapSizeNATAffinity:
                type: integer
              bpfMapSizeNATBackend:
                description: BPFMapSizeNATBackend sets the size for nat back end map.
                  This is the total number of endpoints. This is mostly more than
                  the size of the number of services.
                type: integer
              bpfMapSizeNATFrontend:
                description: BPFMapSizeNATFrontend sets the size for nat front end
                  map. FrontendMap should be large enough to hold an entry for each
                  nodeport, external IP and each port in each service.
                type: integer
              bpfMapSizeRoute:
                description: BPFMapSizeRoute sets the size for the routes map.  The
                  routes map should be large enough to hold one entry per workload
                  and a handful of entries per host (enough to cover its own IPs and
                  tunnel IPs).
                type: integer
              bpfPSNATPorts:
                anyOf:
                - type: integer
                - type: string
                description: 'BPFPSNATPorts sets the range from which we randomly
                  pick a port if there is a source port collision. This should be
                  within the ephemeral range as defined by RFC 6056 (1024–65535) and
                  preferably outside the  ephemeral ranges used by common operating
                  systems. Linux uses 32768–60999, while others mostly use the IANA
                  defined range 49152–65535. It is not necessarily a problem if this
                  range overlaps with the operating systems. Both ends of the range
                  are inclusive. [Default: 20000:29999]'
                pattern: ^.*
                x-kubernetes-int-or-string: true
              bpfPolicyDebugEnabled:
                description: BPFPolicyDebugEnabled when true, Felix records detailed
                  information about the BPF policy programs, which can be examined
                  with the calico-bpf command-line tool.
                type: boolean
              chainInsertMode:
                description: 'ChainInsertMode controls whether Felix hooks the kernel''s
                  top-level iptables chains by inserting a rule at the top of the
                  chain or by appending a rule at the bottom. insert is the safe default
                  since it prevents Calico''s rules from being bypassed. If you switch
                  to append mode, be sure that the other rules in the chains signal
                  acceptance by falling through to the Calico rules, otherwise the
                  Calico policy will be bypassed. [Default: insert]'
                type: string
              dataplaneDriver:
                description: DataplaneDriver filename of the external dataplane driver
                  to use.  Only used if UseInternalDataplaneDriver is set to false.
                type: string
              dataplaneWatchdogTimeout:
                description: 'DataplaneWatchdogTimeout is the readiness/liveness timeout
                  used for Felix''s (internal) dataplane driver. Increase this value
                  if you experience spurious non-ready or non-live events when Felix
                  is under heavy load. Decrease the value to get felix to report non-live
                  or non-ready more quickly. [Default: 90s]'
                type: string
              debugDisableLogDropping:
                type: boolean
              debugMemoryProfilePath:
                type: string
              debugSimulateCalcGraphHangAfter:
                type: string
              debugSimulateDataplaneHangAfter:
                type: string
              defaultEndpointToHostAction:
                description: 'DefaultEndpointToHostAction controls what happens to
                  traffic that goes from a workload endpoint to the host itself (after
                  the traffic hits the endpoint egress policy). By default Calico
                  blocks traffic from workload endpoints to the host itself with an
                  iptables "DROP" action. If you want to allow some or all traffic
                  from endpoint to host, set this parameter to RETURN or ACCEPT. Use
                  RETURN if you have your own rules in the iptables "INPUT" chain;
                  Calico will insert its rules at the top of that chain, then "RETURN"
                  packets to the "INPUT" chain once it has completed processing workload
                  endpoint egress policy. Use ACCEPT to unconditionally accept packets
                  from workloads after processing workload endpoint egress policy.
                  [Default: Drop]'
                type: string
              deviceRouteProtocol:
                description: This defines the route protocol added to programmed device
                  routes, by default this will be RTPROT_BOOT when left blank.
                type: integer
              deviceRouteSourceAddress:
                description: This is the IPv4 source address to use on programmed
                  device routes. By default the source address is left blank, leaving
                  the kernel to choose the source address used.
                type: string
              deviceRouteSourceAddressIPv6:
                description: This is the IPv6 source address to use on programmed
                  device routes. By default the source address is left blank, leaving
                  the kernel to choose the source address used.
                type: string
              disableConntrackInvalidCheck:
                type: boolean
              endpointReportingDelay:
                type: string
              endpointReportingEnabled:
                type: boolean
              externalNodesList:
                description: ExternalNodesCIDRList is a list of CIDR's of external-non-calico-nodes
                  which may source tunnel traffic and have the tunneled traffic be
                  accepted at calico nodes.
                items:
                  type: string
                type: array
              failsafeInboundHostPorts:
                description: 'FailsafeInboundHostPorts is a list of UDP/TCP ports
                  and CIDRs that Felix will allow incoming traffic to host endpoints
                  on irrespective of the security policy. This is useful to avoid
                  accidentally cutting off a host with incorrect configuration. For
                  back-compatibility, if the protocol is not specified, it defaults
                  to "tcp". If a CIDR is not specified, it will allow traffic from
                  all addresses. To disable all inbound host ports, use the value
                  none. The default value allows ssh access and DHCP. [Default: tcp:22,
                  udp:68, tcp:179, tcp:2379, tcp:2380, tcp:6443, tcp:6666, tcp:6667]'
                items:
                  description: ProtoPort is combination of protocol, port, and CIDR.
                    Protocol and port must be specified.
                  properties:
                    net:
                      type: string
                    port:
                      type: integer
                    protocol:
                      type: string
                  required:
                  - port
                  - protocol
                  type: object
                type: array
              failsafeOutboundHostPorts:
                description: 'FailsafeOutboundHostPorts is a list of UDP/TCP ports
                  and CIDRs that Felix will allow outgoing traffic from host endpoints
                  to irrespective of the security policy. This is useful to avoid
                  accidentally cutting off a host with incorrect configuration. For
                  back-compatibility, if the protocol is not specified, it defaults
                  to "tcp". If a CIDR is not specified, it will allow traffic from
                  all addresses. To disable all outbound host ports, use the value
                  none. The default value opens etcd''s standard ports to ensure that
                  Felix does not get cut off from etcd as well as allowing DHCP and
                  DNS. [Default: tcp:179, tcp:2379, tcp:2380, tcp:6443, tcp:6666,
                  tcp:6667, udp:53, udp:67]'
                items:
                  description: ProtoPort is combination of protocol, port, and CIDR.
                    Protocol and port must be specified.
                  properties:
                    net:
                      type: string
                    port:
                      type: integer
                    protocol:
                      type: string
                  required:
                  - port
                  - protocol
                  type: object
                type: array
              featureDetectOverride:
                description: FeatureDetectOverride is used to override the feature
                  detection. Values are specified in a comma separated list with no
                  spaces, example; "SNATFullyRandom=true,MASQFullyRandom=false,RestoreSupportsLock=".
                  "true" or "false" will force the feature, empty or omitted values
                  are auto-detected.
                type: string
              floatingIPs:
                description: FloatingIPs configures whether or not Felix will program
                  floating IP addresses.
                enum:
                - Enabled
                - Disabled
                type: string
              genericXDPEnabled:
                description: 'GenericXDPEnabled enables Generic XDP so network cards
                  that don''t support XDP offload or driver modes can use XDP. This
                  is not recommended since it doesn''t provide better performance
                  than iptables. [Default: false]'
                type: boolean
              healthEnabled:
                type: boolean
              healthHost:
                type: string
              healthPort:
                type: integer
              interfaceExclude:
                description: 'InterfaceExclude is a comma-separated list of interfaces
                  that Felix should exclude when monitoring for host endpoints. The
                  default value ensures that Felix ignores Kubernetes'' IPVS dummy
                  interface, which is used internally by kube-proxy. If you want to
                  exclude multiple interface names using a single value, the list
                  supports regular expressions. For regular expressions you must wrap
                  the value with ''/''. For example having values ''/^kube/,veth1''
                  will exclude all interfaces that begin with ''kube'' and also the
                  interface ''veth1''. [Default: kube-ipvs0]'
                type: string
              interfacePrefix:
                description: 'InterfacePrefix is the interface name prefix that identifies
                  workload endpoints and so distinguishes them from host endpoint
                  interfaces. Note: in environments other than bare metal, the orchestrators
                  configure this appropriately. For example our Kubernetes and Docker
                  integrations set the ''cali'' value, and our OpenStack integration
                  sets the ''tap'' value. [Default: cali]'
                type: string
              interfaceRefreshInterval:
                description: InterfaceRefreshInterval is the period at which Felix
                  rescans local interfaces to verify their state. The rescan can be
                  disabled by setting the interval to 0.
                type: string
              ipipEnabled:
                description: 'IPIPEnabled overrides whether Felix should configure
                  an IPIP interface on the host. Optional as Felix determines this
                  based on the existing IP pools. [Default: nil (unset)]'
                type: boolean
              ipipMTU:
                description: 'IPIPMTU is the MTU to set on the tunnel device. See
                  Configuring MTU [Default: 1440]'
                type: integer
              ipsetsRefreshInterval:
                description: 'IpsetsRefreshInterval is the period at which Felix re-checks
                  all iptables state to ensure that no other process has accidentally
                  broken Calico''s rules. Set to 0 to disable iptables refresh. [Default:
                  90s]'
                type: string
              iptablesBackend:
                description: IptablesBackend specifies which backend of iptables will
                  be used. The default is legacy.
                type: string
              iptablesFilterAllowAction:
                type: string
              iptablesLockFilePath:
                description: 'IptablesLockFilePath is the location of the iptables
                  lock file. You may need to change this if the lock file is not in
                  its standard location (for example if you have mapped it into Felix''s
                  container at a different path). [Default: /run/xtables.lock]'
                type: string
              iptablesLockProbeInterval:
                description: 'IptablesLockProbeInterval is the time that Felix will
                  wait between attempts to acquire the iptables lock if it is not
                  available. Lower values make Felix more responsive when the lock
                  is contended, but use more CPU. [Default: 50ms]'
                type: string
              iptablesLockTimeout:
                description: 'IptablesLockTimeout is the time that Felix will wait
                  for the iptables lock, or 0, to disable. To use this feature, Felix
                  must share the iptables lock file with all other processes that
                  also take the lock. When running Felix inside a container, this
                  requires the /run directory of the host to be mounted into the calico/node
                  or calico/felix container. [Default: 0s disabled]'
                type: string
              iptablesMangleAllowAction:
                type: string
              iptablesMarkMask:
                description: 'IptablesMarkMask is the mask that Felix selects its
                  IPTables Mark bits from. Should be a 32 bit hexadecimal number with
                  at least 8 bits set, none of which clash with any other mark bits
                  in use on the system. [Default: 0xff000000]'
                format: int32
                type: integer
              iptablesNATOutgoingInterfaceFilter:
                type: string
              iptablesPostWriteCheckInterval:
                description: 'IptablesPostWriteCheckInterval is the period after Felix
                  has done a write to the dataplane that it schedules an extra read
                  back in order to check the write was not clobbered by another process.
                  This should only occur if another application on the system doesn''t
                  respect the iptables lock. [Default: 1s]'
                type: string
              iptablesRefreshInterval:
                description: 'IptablesRefreshInterval is the period at which Felix
                  re-checks the IP sets in the dataplane to ensure that no other process
                  has accidentally broken Calico''s rules. Set to 0 to disable IP
                  sets refresh. Note: the default for this value is lower than the
                  other refresh intervals as a workaround for a Linux kernel bug that
                  was fixed in kernel version 4.11. If you are using v4.11 or greater
                  you may want to set this to, a higher value to reduce Felix CPU
                  usage. [Default: 10s]'
                type: string
              ipv6Support:
                description: IPv6Support controls whether Felix enables support for
                  IPv6 (if supported by the in-use dataplane).
                type: boolean
              kubeNodePortRanges:
                description: 'KubeNodePortRanges holds list of port ranges used for
                  service node ports. Only used if felix detects kube-proxy running
                  in ipvs mode. Felix uses these ranges to separate host and workload
                  traffic. [Default: 30000:32767].'
                items:
                  anyOf:
                  - type: integer
                  - type: string
                  pattern: ^.*
                  x-kubernetes-int-or-string: true
                type: array
              logDebugFilenameRegex:
                description: LogDebugFilenameRegex controls which source code files
                  have their Debug log output included in the logs. Only logs from
                  files with names that match the given regular expression are included.  The
                  filter only applies to Debug level logs.
                type: string
              logFilePath:
                description: 'LogFilePath is the full path to the Felix log. Set to
                  none to disable file logging. [Default: /var/log/calico/felix.log]'
                type: string
              logPrefix:
                description: 'LogPrefix is the log prefix that Felix uses when rendering
                  LOG rules. [Default: calico-packet]'
                type: string
              logSeverityFile:
                description: 'LogSeverityFile is the log severity above which logs
                  are sent to the log file. [Default: Info]'
                type: string
              logSeverityScreen:
                description: 'LogSeverityScreen is the log severity above which logs
                  are sent to the stdout. [Default: Info]'
                type: string
              logSeveritySys:
                description: 'LogSeveritySys is the log severity above which logs
                  are sent to the syslog. Set to None for no logging to syslog. [Default:
                  Info]'
                type: string
              maxIpsetSize:
                type: integer
              metadataAddr:
                description: 'MetadataAddr is the IP address or domain name of the
                  server that can answer VM queries for cloud-init metadata. In OpenStack,
                  this corresponds to the machine running nova-api (or in Ubuntu,
                  nova-api-metadata). A value of none (case insensitive) means that
                  Felix should not set up any NAT rule for the metadata path. [Default:
                  127.0.0.1]'
                type: string
              metadataPort:
                description: 'MetadataPort is the port of the metadata server. This,
                  combined with global.MetadataAddr (if not ''None''), is used to
                  set up a NAT rule, from 169.254.169.254:80 to MetadataAddr:MetadataPort.
                  In most cases this should not need to be changed [Default: 8775].'
                type: integer
              mtuIfacePattern:
                description: MTUIfacePattern is a regular expression that controls
                  which interfaces Felix should scan in order to calculate the host's
                  MTU. This should not match workload interfaces (usually named cali...).
                type: string
              natOutgoingAddress:
                description: NATOutgoingAddress specifies an address to use when performing
                  source NAT for traffic in a natOutgoing pool that is leaving the
                  network. By default the address used is an address on the interface
                  the traffic is leaving on (ie it uses the iptables MASQUERADE target)
                type: string
              natPortRange:
                anyOf:
                - type: integer
                - type: string
                description: NATPortRange specifies the range of ports that is used
                  for port mapping when doing outgoing NAT. When unset the default
                  behavior of the network stack is used.
                pattern: ^.*
                x-kubernetes-int-or-string: true
              netlinkTimeout:
                type: string
              openstackRegion:
                description: 'OpenstackRegion is the name of the region that a particular
                  Felix belongs to. In a multi-region Calico/OpenStack deployment,
                  this must be configured somehow for each Felix (here in the datamodel,
                  or in felix.cfg or the environment on each compute node), and must
                  match the [calico] openstack_region value configured in neutron.conf
                  on each node. [Default: Empty]'
                type: string
              policySyncPathPrefix:
                description: 'PolicySyncPathPrefix is used to by Felix to communicate
                  policy changes to external services, like Application layer policy.
                  [Default: Empty]'
                type: string
              prometheusGoMetricsEnabled:
                description: 'PrometheusGoMetricsEnabled disables Go runtime metrics
                  collection, which the Prometheus client does by default, when set
                  to false. This reduces the number of metrics reported, reducing
                  Prometheus load. [Default: true]'
                type: boolean
              prometheusMetricsEnabled:
                description: 'PrometheusMetricsEnabled enables the Prometheus metrics
                  server in Felix if set to true. [Default: false]'
                type: boolean
              prometheusMetricsHost:
                description: 'PrometheusMetricsHost is the host that the Prometheus
                  metrics server should bind to. [Default: empty]'
                type: string
              prometheusMetricsPort:
                description: 'PrometheusMetricsPort is the TCP port that the Prometheus
                  metrics server should bind to. [Default: 9091]'
                type: integer
              prometheusProcessMetricsEnabled:
                description: 'PrometheusProcessMetricsEnabled disables process metrics
                  collection, which the Prometheus client does by default, when set
                  to false. This reduces the number of metrics reported, reducing
                  Prometheus load. [Default: true]'
                type: boolean
              prometheusWireGuardMetricsEnabled:
                description: 'PrometheusWireGuardMetricsEnabled disables wireguard
                  metrics collection, which the Prometheus client does by default,
                  when set to false. This reduces the number of metrics reported,
                  reducing Prometheus load. [Default: true]'
                type: boolean
              removeExternalRoutes:
                description: Whether or not to remove device routes that have not
                  been programmed by Felix. Disabling this will allow external applications
                  to also add device routes. This is enabled by default which means
                  we will remove externally added routes.
                type: boolean
              reportingInterval:
                description: 'ReportingInterval is the interval at which Felix reports
                  its status into the datastore or 0 to disable. Must be non-zero
                  in OpenStack deployments. [Default: 30s]'
                type: string
              reportingTTL:
                description: 'ReportingTTL is the time-to-live setting for process-wide
                  status reports. [Default: 90s]'
                type: string
              routeRefreshInterval:
                description: 'RouteRefreshInterval is the period at which Felix re-checks
                  the routes in the dataplane to ensure that no other process has
                  accidentally broken Calico''s rules. Set to 0 to disable route refresh.
                  [Default: 90s]'
                type: string
              routeSource:
                description: 'RouteSource configures where Felix gets its routing
                  information. - WorkloadIPs: use workload endpoints to construct
                  routes. - CalicoIPAM: the default - use IPAM data to construct routes.'
                type: string
              routeSyncDisabled:
                description: RouteSyncDisabled will disable all operations performed
                  on the route table. Set to true to run in network-policy mode only.
                type: boolean
              routeTableRange:
                description: Deprecated in favor of RouteTableRanges. Calico programs
                  additional Linux route tables for various purposes. RouteTableRange
                  specifies the indices of the route tables that Calico should use.
                properties:
                  max:
                    type: integer
                  min:
                    type: integer
                required:
                - max
                - min
                type: object
              routeTableRanges:
                description: Calico programs additional Linux route tables for various
                  purposes. RouteTableRanges specifies a set of table index ranges
                  that Calico should use. Deprecates` + "`RouteTableRange`, overrides `RouteTableRange`." + `
                items:
                  properties:
                    max:
                      type: integer
                    min:
                      type: integer
                  required:
                  - max
                  - min
                  type: object
                type: array
              serviceLoopPrevention:
                description: 'When service IP advertisement is enabled, prevent routing
                  loops to service IPs that are not in use, by dropping or rejecting
                  packets that do not get DNAT''d by kube-proxy. Unless set to "Disabled",
                  in which case such routing loops continue to be allowed. [Default:
                  Drop]'
                type: string
              sidecarAccelerationEnabled:
                description: 'SidecarAccelerationEnabled enables experimental sidecar
                  acceleration [Default: false]'
                type: boolean
              usageReportingEnabled:
                description: 'UsageReportingEnabled reports anonymous Calico version
                  number and cluster size to projectcalico.org. Logs warnings returned
                  by the usage server. For example, if a significant security vulnerability
                  has been discovered in the version of Calico being used. [Default:
                  true]'
                type: boolean
              usageReportingInitialDelay:
                description: 'UsageReportingInitialDelay controls the minimum delay
                  before Felix makes a report. [Default: 300s]'
                type: string
              usageReportingInterval:
                description: 'UsageReportingInterval controls the interval at which
                  Felix makes reports. [Default: 86400s]'
                type: string
              useInternalDataplaneDriver:
                description: UseInternalDataplaneDriver, if true, Felix will use its
                  internal dataplane programming logic.  If false, it will launch
                  an external dataplane driver and communicate with it over protobuf.
                type: boolean
              vxlanEnabled:
                description: 'VXLANEnabled overrides whether Felix should create the
                  VXLAN tunnel device for VXLAN networking. Optional as Felix determines
                  this based on the existing IP pools. [Default: nil (unset)]'
                type: boolean
              vxlanMTU:
                description: 'VXLANMTU is the MTU to set on the IPv4 VXLAN tunnel
                  device. See Configuring MTU [Default: 1410]'
                type: integer
              vxlanMTUV6:
                description: 'VXLANMTUV6 is the MTU to set on the IPv6 VXLAN tunnel
                  device. See Configuring MTU [Default: 1390]'
                type: integer
              vxlanPort:
                type: integer
              vxlanVNI:
                type: integer
              wireguardEnabled:
                description: 'WireguardEnabled controls whether Wireguard is enabled
                  for IPv4 (encapsulating IPv4 traffic over an IPv4 underlay network).
                  [Default: false]'
                type: boolean
              wireguardEnabledV6:
                description: 'WireguardEnabledV6 controls whether Wireguard is enabled
                  for IPv6 (encapsulating IPv6 traffic over an IPv6 underlay network).
                  [Default: false]'
                type: boolean
              wireguardHostEncryptionEnabled:
                description: 'WireguardHostEncryptionEnabled controls whether Wireguard
                  host-to-host encryption is enabled. [Default: false]'
                type: boolean
              wireguardInterfaceName:
                description: 'WireguardInterfaceName specifies the name to use for
                  the IPv4 Wireguard interface. [Default: wireguard.cali]'
                type: string
              wireguardInterfaceNameV6:
                description: 'WireguardInterfaceNameV6 specifies the name to use for
                  the IPv6 Wireguard interface. [Default: wg-v6.cali]'
                type: string
              wireguardKeepAlive:
                description: 'WireguardKeepAlive controls Wireguard PersistentKeepalive
                  option. Set 0 to disable. [Default: 0]'
                type: string
              wireguardListeningPort:
                description: 'WireguardListeningPort controls the listening port used
                  by IPv4 Wireguard. [Default: 51820]'
                type: integer
              wireguardListeningPortV6:
                description: 'WireguardListeningPortV6 controls the listening port
                  used by IPv6 Wireguard. [Default: 51821]'
                type: integer
              wireguardMTU:
                description: 'WireguardMTU controls the MTU on the IPv4 Wireguard
                  interface. See Configuring MTU [Default: 1440]'
                type: integer
              wireguardMTUV6:
                description: 'WireguardMTUV6 controls the MTU on the IPv6 Wireguard
                  interface. See Configuring MTU [Default: 1420]'
                type: integer
              wireguardRoutingRulePriority:
                description: 'WireguardRoutingRulePriority controls the priority value
                  to use for the Wireguard routing rule. [Default: 99]'
                type: integer
              workloadSourceSpoofing:
                description: WorkloadSourceSpoofing controls whether pods can use
                  the allowedSourcePrefixes annotation to send traffic with a source
                  IP address that is not theirs. This is disabled by default. When
                  set to "Any", pods can request any prefix.
                type: string
              xdpEnabled:
                description: 'XDPEnabled enables XDP acceleration for suitable untracked
                  incoming deny rules. [Default: true]'
                type: boolean
              xdpRefreshInterval:
                description: 'XDPRefreshInterval is the period at which Felix re-checks
                  all XDP state to ensure that no other process has accidentally broken
                  Calico''s BPF maps or attached programs. Set to 0 to disable XDP
                  refresh. [Default: 90s]'
                type: string
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
