/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
* @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html
 */

package app

import (
	sessions "github.com/wangyysde/sysadmSessions"
	sysadmDB "sysadm/db"
	sysadmObjects "sysadm/objects/app"
	"sysadm/sysadmLog"
	sysadmSetting "sysadm/syssetting/app"
)

// 存储运行期数据
type runingData struct {
	dbConf        *sysadmDB.DbConfig
	logEntity     *sysadmLog.LoggerConfig
	workingRoot   string
	sessionName   string
	sessionOption sessions.Options
	pageInfo      sysadmSetting.PageInfo
	objectEntiy   sysadmObjects.ObjectEntity
}

type VolumeMount struct {
	// 基本信息
	BasicInfo VolumeMountBasic `xml:"basicInfo,omitempty" json:"basicInfo,omitempty" yaml:"basicInfo,omitempty"`

	// PVC命名
	PvcData PvcData `json:"pvcData,omitempty" yaml:"pvcData,omitempty" xml:"pvcData,omitempty"`

	// HostPath数据
	HostPathData HostPathData `json:"hostPathData,omitempty" yaml:"hostPathData,omitempty" xml:"hostPathData,omitempty"`

	// ConfigMap名称
	CmName string `json:"cmData,omitempty" yaml:"cmData,omitempty" xml:"cmData,omitempty"`

	// Secret名称
	SecretName string `json:"secretData,omitempty" yaml:"secretData,omitempty" xml:"secretData,omitempty"`

	// 容器挂载点数据
	ContainerData ContainerVolumeMountData `json:"containerData,omitempty" yaml:"containerData,omitempty" xml:"containerData,omitempty"`
}

type VolumeMountBasic struct {
	// 卷名称
	VolumeName string `json:"volumeName,omitempty" yaml:"volumeName,omitempty" xml:"volumeName,omitempty"`

	// 卷类型
	VolumeType string `json:"volumeType,omitempty" yaml:"volumeType,omitempty" xml:"volumeType,omitempty"`
}

type PvcData struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
}

type HostPathData struct {
	// 主机路径
	HostPath string `json:"hostPath,omitempty" yaml:"hostPath,omitempty" xml:"hostPath,omitempty"`

	// 类型
	HostPathType string `json:"hostPathType,omitempty" yaml:"hostPathType,omitempty" xml:"hostPathType,omitempty"`
}

type ContainerVolumeMountData struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`

	MountPath string `json:"mountPath,omitempty" yaml:"mountPath,omitempty" xml:"mountPath,omitempty"`

	SubPath string `json:"subPath,omitempty" yaml:"subPath,omitempty" xml:"subPath,omitempty"`
}

type ContainerData struct {
	// 容器的基本信息
	BasicInfo ContainerBasicInfo `json:"basicInfo,omitempty" yaml:"basicInfo,omitempty" xml:"basicInfo,omitempty"`

	// 容器的Quota数据
	Quota ContainerQuota `json:"quota,omitempty" yaml:"quota,omitempty" xml:"quota,omitempty"`

	// 容器启动探测
	Startup ContainerProbe `json:"startup,omitempty" yaml:"startup,omitempty" xml:"startup,omitempty"`

	// 容器就绪探测
	Readiness ContainerProbe `json:"readiness,omitempty" yaml:"readiness,omitempty" xml:"readiness,omitempty"`

	// 容器存活探测
	Liveness ContainerProbe `json:"liveness,omitempty" yaml:"liveness,omitempty" xml:"liveness,omitempty"`
}
type ContainerBasicInfo struct {
	// 容器名
	ContainerName string `json:"containerName,omitempty" yaml:"containerName,omitempty" xml:"containerName,omitempty"`

	// 容器类型，普通容器其值为0,特权容器为1,初始化容器为2. 1,2表示特权的初始化容器
	ContainerType string `json:"containerType,omitempty" yaml:"containerType,omitempty" xml:"containerType,omitempty"`

	// 容器镜像地址
	ContainerImage string `json:"containerImage,omitempty" yaml:"containerImage,omitempty" xml:"containerImage,omitempty"`

	// 镜像摘取策略
	ImagePullPolicy string `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty" xml:"imagePullPolicy,omitempty"`

	// 容器启动命令
	StartCommand string `json:"startCommand,omitempty" yaml:"startCommand,omitempty" xml:"startCommand,omitempty"`

	// 环境变量
	Environment string `json:"environment,omitempty" yaml:"environment,omitempty" xml:"environment,omitempty"`
}

type ContainerQuota struct {
	// 容器Quota的 CPU Request值
	CpuRequest string `json:"cpuRequest,omitempty" yaml:"cpuRequest,omitempty" xml:"cpuRequest,omitempty"`

	// 容器Quota的 CPU Limit值
	CpuLimit string `json:"cpuLimit,omitempty" yaml:"cpuLimit,omitempty" xml:"cpuLimit,omitempty"`

	MemRequest string `json:"memRequest,omitempty" yaml:"memRequest,omitempty" xml:"memRequest,omitempty"`

	MemLimit string `json:"memLimit,omitempty" yaml:"memLimit,omitempty" xml:"memLimit,omitempty"`
}

type ContainerProbe struct {
	ProbeType           string   `json:"probeType,omitempty" yaml:"probeType,omitempty" xml:"probeType,omitempty"`
	HttpProtocol        string   `json:"httpProtocol,omitempty" yaml:"httpProtocol,omitempty" xml:"httpProtocol,omitempty"`
	HttpPort            string   `json:"httpPort,omitempty" yaml:"httpPort,omitempty" xml:"httpPort,omitempty"`
	HttpPath            string   `json:"httpPath,omitempty" yaml:"httpPath,omitempty" xml:"httpPath,omitempty"`
	HttpHeaderKey       []string `json:"httpHeaderKey,omitempty" yaml:"httpHeaderKey,omitempty" xml:"httpHeaderKey,omitempty"`
	HttpHeaderValue     []string `json:"httpHeaderValue,omitempty" yaml:"httpHeaderValue,omitempty" xml:"httpHeaderValue,omitempty"`
	TcpPort             string   `json:"tcpPort,omitempty" yaml:"tcpPort,omitempty" xml:"tcpPort,omitempty"`
	Command             string   `json:"command,omitempty" yaml:"command,omitempty" xml:"command,omitempty"`
	InitialDelaySeconds string   `json:"initialDelaySeconds,omitempty" yaml:"initialDelaySeconds,omitempty" xml:"initialDelaySeconds,omitempty"`
	PeriodSeconds       string   `json:"periodSeconds,omitempty" yaml:"periodSeconds,omitempty" xml:"periodSeconds,omitempty"`
	TimeoutSeconds      string   `json:"timeoutSeconds,omitempty" yaml:"timeoutSeconds,omitempty" xml:"timeoutSeconds,omitempty"`
	FailureThreshold    string   `json:"failureThreshold,omitempty" yaml:"failureThreshold,omitempty" xml:"failureThreshold,omitempty"`
	SuccessThreshold    string   `json:"successThreshold,omitempty" yaml:"successThreshold,omitempty" xml:"successThreshold,omitempty"`
}
