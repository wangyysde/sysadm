/* =============================================================
* @Author:  Wayne Wang <yuying.wang@sincerecloud.com>
* @ Last Modified At: 2023.05.10
* @Copyright (c) 2023 Sincerecloud
* @HomePage: https://www.sincerecloud.com/
*
*  定义keepalived使用到的默认值
 */
package keepalived

// default value of SMTP smtp_connect_timeout what are smtp server connection timeout in seconds.
var defaultSmtpConnectTimeout int = 30

// default router id. this value will be used when route id not specified and can not got hostname
var defaultRouterID string = "controlplane"

// default name of vrrp instance
var defaultVrrpInstanceName string = "kube-apiserver"

// default auth type
var defaultAuthType string = "PASS"

// default auth pass
var defaultAuthPass string = "sincerecloudPassword"

// default master priority value
var defaultMasterPriority int = 100

// default advert_int
var defaultAdvertInt int = 1

// default delay_loop
var defaultDelayLoop int = 60

// default lb_algo
// lb_algo is one of Round-Robin(rr),Weighted Round-Robin(wrr),Least-Connection(lc),Weighted Least-Connection(wlc),
// Locality-Based Least-Connection(lblc),Locality-Based Least-Connection Scheduling with Replication(lblcr),
// Destination Hash(dh),Source Hash(sh),Source Expected Delay(sed),Never Queue(nq)
var defaultLbAlgo string = "rr"

// default lb_kind
var defaultLbKind string = "DR"

// default persistence_timeout
var defaultPersistenceTimeout int = 60

// default max retry times for check
var defaultCheckRetry int = 3

// default check timeout
var defaultCheckTimeout int = 10
