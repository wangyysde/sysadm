/* =============================================================
* @Author:  Wayne Wang <yuying.wang@sincerecloud.com>
* @ Last Modified At: 2023.05.10
* @Copyright (c) 2023 Sincerecloud
* @HomePage: https://www.sincerecloud.com/
*
 */

package keepalived

import (
	"testing"
)

func TestCreateConf(t *testing.T) {
	var (
		notificationEmail     = []string{"1@163.com", "2@163.com", "3@163.com"}
		notificationEmailFrom = "wangyuying<yuying.wang@sincerecloud.com>"
		smtpServer            = "smtp.163.com"
		smtpConnectTimeout    = 30
		smtpPort              = 625
		routerId              = "testhostname"
		instanceName          = "kube-apiserver"
		vipNicName            = "ens33"
		authPass              = "testAuthPassword"
		routerid              = 50
		priority              = 100
		vips                  = []string{"172.28.2.50", "172.28.2.51"}
		testVip1              = "172.28.2.50"
		lbAlgo                = "rr"
		vsPort                = 8081
		rs1                   = "172.28.2.10"
		rs2                   = "172.28.2.11"
		confFile              = "/etc/kubernetes/keepalived/keepalived.conf"
	)

	conf := setGlobalConf(nil, notificationEmailFrom, smtpServer, routerId, notificationEmail, smtpPort, smtpConnectTimeout)
	conf, e := addVrrpInstance(conf, instanceName, vipNicName, authPass, routerid, priority, 0, vips, true)
	if e != nil {
		t.Fatalf("add vrrp instance error %s", e)
	}

	vs, e := newVirtualServer(testVip1, lbAlgo, vsPort, 0, 0, true)
	if e != nil {
		t.Fatal("create a new virtual server error")
	}

	vs, e = addRealServer(vs, rs1, 80, 100, 3, 5)
	if e != nil {
		t.Fatalf("add a new realserver to virtual server error %s", e)
	}

	vs, e = addRealServer(vs, rs2, 81, 99, 3, 5)
	if e != nil {
		t.Fatalf("add a new realserver to virtual server error %s", e)
	}

	conf, e = addVirtualServer(conf, vs)
	if e != nil {
		t.Fatalf("add virtual server to keepalived configuration error %s", e)
	}

	err := createKeepalivedConf(conf, confFile)
	if err != nil {
		t.Fatalf("create keepalived configuration file error: %s\n ", err)
	}

	t.Log("keepalived configuration  file has be created\n")
}

func TestCreateStaticPod(t *testing.T) {
	var (
		confFile                   = "/etc/kubernetes/keepalived/keepalived.conf"
		registryUrl                = "hb.sincerecloud.com/k8s/v1.26.3"
		image               string = "kube-keepalived:v1.0"
		manifestDir         string = "/etc/kubernetes/manifests"
		keepalivedPodFile   string = "kube-keepalived.yaml"
		initialDelaySeconds int32  = int32(10)
	)

	err := createKeepalivedPod(registryUrl, image, confFile, manifestDir, "", keepalivedPodFile, initialDelaySeconds)
	if err != nil {
		t.Fatalf("can not create keepalived pod static file error %s\n", err)
	}
	t.Log("keepalived pod static file has be created\n")
}
