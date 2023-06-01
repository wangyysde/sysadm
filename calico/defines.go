// Package calico is for agent installing calico automatic
// @Author:  Wayne Wang <yuying.wang@sincerecloud.com>
// @ Last Modified At: 2023.06.01
// @Copyright (c) 2023 Sincerecloud
// @HomePage: https://www.sincerecloud.com/
//
//	定义calico使用到的默认值
package calico

// default calico mode is to run
var defaultCalicoBackend = "vxlan"

// 默认镜像仓库地址
var defaultImageRepository = "docker.io/calico"

// calico版本
var defaultImageVersion = "cni:v3.24.1"
