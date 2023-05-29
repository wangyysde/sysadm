/* =============================================================
* @Author:  Wayne Wang <yuying.wang@sincerecloud.com>
* @ Last Modified At: 2023.05.10
* @Copyright (c) 2023 Sincerecloud
* @HomePage: https://www.sincerecloud.com/
*
*  定义公共函数
 */

package keepalived

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	v1 "k8s.io/api/core/v1"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
	staticpodutil "k8s.io/kubernetes/cmd/kubeadm/app/util/staticpod"
)

// 根据keepalived的配置信息生成keepalived的配置文件，并将所生成的配置文件系统写到confFile指定的文件里
func createKeepalivedConf(conf *KeepalivedConf, confFile string) error {
	if conf == nil {
		return fmt.Errorf("can not create configuration file for keepalived without configurations")
	}

	confFile = strings.TrimSpace(confFile)
	if confFile == "" {
		return fmt.Errorf("the keepalived configuration file name must be not empty")
	}
	tmpl, e := template.New("keepalived").Parse(keepConfigTmpl)
	if e != nil {
		return e
	}

	var tpl bytes.Buffer
	e = tmpl.Execute(&tpl, *conf)
	if e != nil {
		return e
	}

	// creates target folder if not already exists
	confDir := filepath.Dir(confFile)
	if err := os.MkdirAll(confDir, 0700); err != nil {
		return fmt.Errorf("failed to create directory %s error: %s", confDir, err)
	}

	if err := os.WriteFile(confFile, tpl.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write keepalived configuration file(%s) for %s ", confFile, err)
	}

	return nil
}

// 根据参数信息生成keepalived的静态pod文件.所生成的静态pod文件的文件名和存储路径通过keepalivedPodFile和manifestDir指定
// 本函数调用了kubeadm工具的相关方法和函数.
func createKeepalivedPod(registryUrl, image, confFile, manifestDir, patchesDir, keepalivedPodFile string, initialDelaySeconds int32) error {
	registryUrl = strings.TrimSpace(registryUrl)
	image = strings.TrimSpace(image)
	volumeMounts := []v1.VolumeMount{
		{Name: "host-localtime", ReadOnly: true, MountPath: "/etc/localtime"},
		{Name: "config", ReadOnly: true, MountPath: "/etc/keepalived/keepalived.conf"},
	}

	execAction := v1.ExecAction{Command: []string{"pidof", "keepalived"}}
	livenessProbe := &v1.Probe{
		ProbeHandler: v1.ProbeHandler{
			Exec: &execAction,
		},
		InitialDelaySeconds: initialDelaySeconds,
	}
	hostPathFile := v1.HostPathFile
	volumes := make(map[string]v1.Volume, 0)
	vol := staticpodutil.NewVolume("host-localtime", "/etc/localtime", &hostPathFile)
	volumes[vol.Name] = vol
	vol = staticpodutil.NewVolume("config", confFile, &hostPathFile)
	volumes[vol.Name] = vol

	if registryUrl == "" || image == "" {
		return fmt.Errorf("address of registry and image must be not empty")
	}
	image = registryUrl + "/" + image

	keepalivedStaticPod := staticpodutil.ComponentPod(v1.Container{
		Name:            "kube-keepalived",
		Image:           image,
		ImagePullPolicy: v1.PullIfNotPresent,
		Command:         []string{"keepalived", "--vrrp", "--log-detail", "--dump-conf", "--use-file=/etc/keepalived/keepalived.conf"},
		VolumeMounts:    volumeMounts,
		LivenessProbe:   livenessProbe,
	}, volumes, map[string]string{})

	if patchesDir != "" {
		patchedSpec, err := staticpodutil.PatchStaticPod(&keepalivedStaticPod, patchesDir, os.Stdout)
		if err != nil {
			return fmt.Errorf("failed to patch static Pod manifest file for keepalived error %s", err)
		}
		keepalivedStaticPod = *patchedSpec
	}

	// creates target folder if not already exists
	if err := os.MkdirAll(manifestDir, 0700); err != nil {
		return fmt.Errorf("failed to create directory %s error: %s", manifestDir, err)
	}

	// writes the pod to disk
	serialized, err := kubeadmutil.MarshalToYaml(&keepalivedStaticPod, v1.SchemeGroupVersion)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest for keepalived to YAML error %s", err)
	}

	filename := filepath.Join(manifestDir, keepalivedPodFile)
	if err := os.WriteFile(filename, serialized, 0600); err != nil {
		return fmt.Errorf("failed to write static pod manifest file for keepalived error %s", err)
	}

	return nil
}

// 设置keepalived 配置文件的全局部分内容。如果keepalived配置未实例化，则先实例化，然后设置相应项目内容，并返回实例化的指针值
// 如果连接SMTP服务器超时时间，即smtp_connect_timeout值为0或者大于600(10分钟),则取默认值30
// 如果routerID值未指定，则尝试获取主动的主机名作为routerID，获取失败则使用默认值"controlplane"作为routerID
func setGlobalConf(conf *KeepalivedConf, emailFrom, smtpServer, routerID string, emails []string, smtpPort, smtpTimeout int) *KeepalivedConf {
	if conf == nil {
		conf = &KeepalivedConf{
			GlobalDefs:     GlobalDefs{},
			VrrpInstances:  make([]VrrpInstance, 0),
			VirtualServers: make([]VirtualServer, 0),
		}
	}

	smtpPortStr := ""
	if smtpPort != 0 && smtpPort != 25 {
		smtpPortStr = strconv.Itoa(smtpPort)
	}

	if smtpTimeout == 0 || smtpTimeout > 600 {
		smtpTimeout = defaultSmtpConnectTimeout
	}

	routerID = strings.TrimSpace(routerID)
	if routerID == "" {
		hostname, e := os.Hostname()
		if e != nil {
			hostname = defaultRouterID
		}
		routerID = hostname
	}

	conf.NotificationEmail = emails
	conf.NotificationEmailFrom = strings.TrimSpace(emailFrom)
	conf.SmtpServer = strings.TrimSpace(smtpServer)
	conf.SmtpPort = smtpPortStr
	conf.SmtpConnectTimeout = smtpTimeout
	conf.RouterId = routerID

	return conf
}

// addVrrpInstance 增加keepalived 配置文件中的vrrp_instance内容，同一个keepalived实例中可以有多个vrrp_instance
// 被增加的整个keepalived不能为空，否则出错。实例名name如果为空，则使用默认实例名，但是能一个keepalived里不能有重复的实例名
// vip地址不能为空，且其应与本地某块网卡上所绑定的IP地址属于同一个网段。nic如果指定则必须是本地某块网卡的完整的网卡名，否则会根据vip地址查询相应的网卡名称
// state的值如何是true表示当前VRRP instance为MASTER，此时如果routeID值为0，则设置成默认值
// 对于MASTER，默认priority为100，对于BACKUP priority值必须指定
func addVrrpInstance(conf *KeepalivedConf, name, nic, authPass string, routeID, priority, advertInt int, vips []string, state bool) (*KeepalivedConf, error) {
	if conf == nil {
		return nil, fmt.Errorf("can not add vrrp instance on nil")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = "defaultVrrpInstanceName"
	}

	if len(vips) < 1 {
		return nil, fmt.Errorf("at least one virtual ipaddress for a vrrp instance")
	}

	nic = strings.ToLower(strings.TrimSpace(nic))
	if nic == "" {
		localNics, e := net.Interfaces()
		if e != nil {
			return nil, fmt.Errorf("get localhost interfaces error %s", e)
		}
		for _, ln := range localNics {
			if nic != "" {
				break
			}
			addrs, e := ln.Addrs()
			if e != nil {
				return nil, fmt.Errorf("get ip address on nic error %s", e)
			}
			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if !ok {
					continue
				}
				if isVipsContained(ipnet, vips) {
					nic = ln.Name
					break
				}
			}
		}
	} else {
		localNics, e := net.Interfaces()
		if e != nil {
			return nil, fmt.Errorf("get localhost interfaces error %s", e)
		}

		foundNic := false
		for _, ln := range localNics {
			if nic == ln.Name {
				foundNic = true
				break
			}
		}
		if !foundNic {
			return nil, fmt.Errorf("interface %s has not found on the localhost", nic)
		}
	}

	if nic == "" {
		return nil, fmt.Errorf("interface name must be not empty")
	}

	authPass = strings.TrimSpace(authPass)
	if authPass == "" {
		authPass = defaultAuthPass
	}

	if routeID == 0 {
		return nil, fmt.Errorf("virtual_router_id must be not zero")
	}

	if priority == 0 {
		if state {
			priority = defaultMasterPriority
		} else {
			return nil, fmt.Errorf("priority for a vrrp instance be not zero")
		}
	}

	if advertInt == 0 {
		advertInt = defaultAdvertInt
	}

	stateStr := "BACKUP"
	if state {
		stateStr = "MASTER"
	}

	vrrpInstance := VrrpInstance{
		Name:            name,
		State:           stateStr,
		Interface:       nic,
		VirtualRouterId: routeID,
		Priority:        priority,
		AdvertInt:       advertInt,
		Authentication: Authentication{
			AuthType: defaultAuthType,
			AuthPass: authPass,
		},
		VirtualIpaddress: vips,
	}

	conf.VrrpInstances = append(conf.VrrpInstances, vrrpInstance)

	return conf, nil
}

// 检查是否所有vip都与ipnet属于同一网段，如果是返回true, 否则返回false
func isVipsContained(ipnet *net.IPNet, vips []string) bool {
	if ipnet == nil {
		return false
	}

	for _, v := range vips {
		ip := net.ParseIP(v)
		if !ipnet.Contains(ip) {
			return false
		}
	}

	return true
}

// newVirtualServer 根据参数实例化一个新的VirtualServer实例.
func newVirtualServer(vip, lbAlgo string, port, delayLoop, persistenceTimeout int, protocol bool) (*VirtualServer, error) {

	vip = strings.TrimSpace(vip)
	if vip == "" {
		return nil, fmt.Errorf("vip must be not empty")
	}

	if port < 1024 || port > 65535 {
		return nil, fmt.Errorf("virtual server port %d is not valid", port)
	}

	if delayLoop < 1 || delayLoop > 600 {
		delayLoop = defaultDelayLoop
	}

	lbAlgo = strings.TrimSpace(strings.ToLower(lbAlgo))
	if !isLbAlgo(lbAlgo) {
		lbAlgo = defaultLbAlgo
	}

	if persistenceTimeout < 1 || persistenceTimeout > 600 {
		persistenceTimeout = defaultPersistenceTimeout
	}

	protocolStr := "UDP"
	if protocol {
		protocolStr = "TCP"
	}

	return &VirtualServer{
		VirtualIp:          vip,
		Port:               port,
		DelayLoop:          delayLoop,
		LbAlgo:             lbAlgo,
		LbKind:             defaultLbKind,
		PersistenceTimeout: persistenceTimeout,
		Protocol:           protocolStr,
		RealServers:        make([]RealServer, 0),
	}, nil
}

// isLbAlgo 检查调度算法是否合法
func isLbAlgo(lbAlgo string) bool {
	algos := []string{"rr", "wrr", "lc", "wlc", "lblc", "lblcr", "dh", "sh", "sed", "nq"}
	for _, v := range algos {
		if lbAlgo == v {
			return true
		}
	}

	return false
}

// addVirtualServer 增加一个新的VirtualServer到全局配置中
func addVirtualServer(conf *KeepalivedConf, vs *VirtualServer) (*KeepalivedConf, error) {
	if conf == nil {
		return nil, fmt.Errorf("can not add virtual server on nil")
	}

	if vs == nil {
		return nil, fmt.Errorf("can not nil to configuration")
	}

	vip := vs.VirtualIp
	vrrpIns := conf.VrrpInstances
	foundVip := false
	for _, instance := range vrrpIns {
		if foundVip {
			break
		}
		for _, ip := range instance.VirtualIpaddress {
			if vip == ip {
				foundVip = true
				break
			}
		}
	}
	if !foundVip {
		return nil, fmt.Errorf("vip %s is not valid", vip)
	}

	conf.VirtualServers = append(conf.VirtualServers, *vs)

	return conf, nil
}

// addRealServer 根据参数增加一个realserver实例到VirtualServer实例中
func addRealServer(vs *VirtualServer, realServer string, port, weight, checkRetry, checkTimeout int) (*VirtualServer, error) {
	if vs == nil {
		return nil, fmt.Errorf("can not add realserver on nil")
	}

	realServer = strings.TrimSpace(realServer)
	if ip := net.ParseIP(realServer); ip == nil {
		return nil, fmt.Errorf("realserver address %s is not valid", realServer)
	}

	if checkRetry == 0 {
		checkRetry = defaultCheckRetry
	}

	if checkTimeout == 0 {
		checkTimeout = defaultCheckTimeout
	}

	rs := RealServer{
		RealIp: realServer,
		Port:   port,
		Weight: weight,
		TcpCheck: TcpCheck{
			ConnectTimeout: checkTimeout,
			Retry:          checkRetry,
		},
	}

	vs.RealServers = append(vs.RealServers, rs)

	return vs, nil
}
