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
*
* NOTE:
* 本文件定义了将来准备独立出来的sysadmApiServer组件所使用的导出的数据结构、常量和全局变量
 */

package app

import (
	"github.com/wangyysde/sysadmLog"
)

// 说明：
// 1. 不管是服务端还是客户端，关于命令、命令状态和日志的侦听接口仅支持POST方法请求
// 2. 主动工作模式时，服务端以CommandData格式向客户端发送要执行的命令：此时不管命令是同步的还是异步的，客户端均以CommandStatus数据结构响应服务端，
//    只是：a.当命令是异步时，客户端对所接收的命令进行合法性检查通过后，以ComandStatusReceived状态响应，否则以ComandStatusSendError状态响应;b. 当
//    命令是同步的时候，则以命令的执行结果响应; c.当命令是异步时，服务端应用通过客户端的/getCommandStatus接口查询命令的执行状态，服务端请求数据结构
//    采用CommandStatusReq数据结构，此时，客户端应以CommandStatus格式响应此请求。
// 3. 被动工作模式时，客户端以CommandReq格式向服务端获取要执行的命令时，此时服务端不管是否有命令需要客户执行，均以CommandData数据结构形式响应客户
//    端。
// 4. 主动工作模式时，如果命令是同步的，则命令的执行状态在上述第1步时，客户端以响应服务端发送执行命令请求时，以CommandStatus格式反馈给服务端；当命
//    令是异步时，服务端需要通过请求客户端的/getCommandStatus接口以CommandStatusReq数据结构形式向客户端请求命令的执行状态，此时客户端以
//    CommandStatus数据结构响应此请求
// 5. 被动工作模式时，不管命令是否是异步，客户端均以CommandStatus数据结构形式向服务端/receiveCommandStatus接口发送命令的执行状态，服务端RepStatus
//    以数据结构的形式响应客户端的此请求
// 6. 主动工作模式时，服务端以LogReq数据结构形式向客户端/getLogs接口请求命令的执行日志，客户端以LogData数据结构的形式响应服务端的此请求。
// 7. 被动工作模式时，客户端以LogData数据结构形式向服务端/receiveLogs接口发送命令的执行日志，此时，服务端以RepStatus数据结构响应客户端

// 下面定义的是命令生命周期的各个状态常量. 后续如果需要添加新的状态，则应该在相应的状态之间添加，且对应的取值等分定义
type CommandStatusCode uint32

const (

	// 表示apiServer 已经创建好命令，等待下发给客户端执行
	CommandStatusCreated CommandStatusCode = 300

	// 表示服务端或客户端已经成功接收了命令或命令状态信息
	ComandStatusReceived CommandStatusCode = 500

	// 表示命令下发出错
	ComandStatusSendError CommandStatusCode = 600

	// 表示命令已经成功下发，但是apiServer尚未收到任何关于本命令的状态信息
	CommandStatusSent CommandStatusCode = 700

	// 表示命令已经成功下发，且apiServer已经收到关于本命令的至少一条状态信息
	CommandStatusRunning CommandStatusCode = 800

	// 表示命令已经成功下发，但是指定时间未能收到命令执行成功或错误的状态报告，通常此时本状态是由定时任务设置的
	CommandStatusTimeout CommandStatusCode = 900

	// 表示命令的子任务执行成功
	CommandStatusTaskOk CommandStatusCode = 950

	// 表示命令的子任务执行成功
	CommandStatusTaskError CommandStatusCode = 960

	// 表示apiServer已经接收到命令已经执行完成，但是命令执行错误
	CommandStatusError CommandStatusCode = 1000

	// 表示apiServer已经接收到命令已经执行完成，且命令已经正常成功执行
	CommandStatusOK CommandStatusCode = 1100

	// 不认识的命令。这通常表示agent收到了一个不认识的命名，或者apiserver收到了一个自己无法识别命令序列号的命令状态信息或命令日志信息
	CommandStatusUnrecognized CommandStatusCode = 1200

	// 表示状态是一个未知状态,这通常表示发生了一个未知错误
	CommandStatusUnkown CommandStatusCode = 9000
)

// 当添加或减少了上面定义的命令状态码的值，则下面这个切片的内容也需要相应的调整
var AllCommandStatusCode = []CommandStatusCode{
	CommandStatusCreated,
	ComandStatusReceived,
	ComandStatusSendError,
	CommandStatusSent,
	CommandStatusRunning,
	CommandStatusTimeout,
	CommandStatusError,
	CommandStatusOK,
	CommandStatusUnrecognized,
	CommandStatusUnkown,
}

// 记录命令的数据
type Command struct {
	// 命令序列号，字符串型，组成为发送给客户端的8字节的当天日期+对应命令在数据库的11位ID值，即YYYYMMDDCOMMANDID.如果对应命令
	// 在数据库中的COMMANDID不足11位，是以前导0补足。
	// 本字段的值是apiServer和其客户端在通信过程中标识一个命令的标识符，包括apiServer查询命令执行状态或客户端返回命令执行状态
	// 与执行结果，均使用本字段的值来标识一个命令。
	// 如果此字段的值为0000 0000 0000 0000 000， 表示当前没有命令要客户端执行
	CommandSeq string `form:"commandSeq" json:"commandSeq" yaml:"commandSeq" xml:"commandSeq"`

	// 具体需要客户端执行的命令,字符串类型,不区分大小写。客户端将依据此值进行事务的路由
	// 如果当前没有命令要执行，则此字段值为"""
	Command string `form:"command" json:"command" yaml:"command" xml:"command"`

	// 指示命令是否是属地同步命令。所谓同步命令是指，命令能够快速执行完成，即能在一个HTTP会话请求超时之前（一般超时时间为几秒内）执行完成的命令。
	Synchronized bool `form:"synchronized" json:"synchronized" yaml:"synchronized" xml:"synchronized"`

	// 运行命令所需要的参数，其中map中的key表示参数名，忽略大小写，且每个参数的长度不得大于64个字符。map中的value是参数的值，可以为空。
	// 不支持多层级的参数格式， 当需要多层级格式时，可以通过不同的key名展开为一个层级。客户端需要判断参数的合法性。
	Parameters map[string]string `form:"parameter" json:"parameter" yaml:"parameter" xml:"parameter"`
}

// 用于记录客户端根据配置或者apiServer的要求埴写的，用于标识自身的标识符。
type NodeIdentifier struct {
	// 记录客户端获取的本节点的IP地址列表，本列表的元素不得重复。如果配置或者apiServer没有要求用IP地址作为标识符之一，则本子段的值可以空
	Ips []string `form:"ips" json:"ips" yaml:"ips" xml:"ips"`

	// 记录客户端获取的本节点的MAC地址列表，本列表的元素不得重复。如果配置或者apiServer没有要求用MAC地址作为标识符之一，则本子段的值可以空
	Macs []string `form:"macs" json:"macs" yaml:"macs" xml:"macs"`

	// 记录客户端获取的本节点的主机名信息，如果配置或者apiServer没有要求用主机名作为标识符之一，则本子段的值可以空
	Hostname string `form:"hostname" json:"hostname" yaml:"hostname" xml:"hostname"`

	// 记录配置或apiServer提供的定义的节点标识符。本字段的值中不得包含IP,HOSTNAME或MAC中任何一个字符串，且本字符串是大小写敏感的。
	Customize string `form:"customize" json:"customize" yaml:"customize" xml:"customize"`
}

// 记录服务端和客户发送具体的命令数据的数据结构
type CommandData struct {
	// 本字段的值可以是IP,HOSTNAME和MAC三个字符串的任意组合(注意：必须是大写），或者是长度不超过256个字符且不包含前面三个字符串的任意字符串(注意：该字符串区分大小写），以标识一个节点。当本子段的值为空时，则表示使用客户端定义的节点标识
	NodeIdentiferStr string `form:"nodeIdentifierStr" json:"nodeIdentifierStr" yaml:"nodeIdentifierStr" xml:"nodeIdentifierStr"`

	// 需要发送的命令数据
	Command `form:"command" json:"command" yaml:"command" xml:"command"`
}

// 当运行于被动模式时，客户端向服务端发送的用于获取要执行的命令时，所发请求的数据结构
type CommandReq struct {
	// 用于记录客户端根据配置或者apiServer的要求埴写的，用于标识自身的标识符。
	NodeIdentifier `form:"nodeIdentifier" json:"nodeIdentifier" yaml:"nodeIdentifier" xml:"nodeIdentifier"`
}

// 当apiServer以主动模式运行时，apiServer向客户发送的获取命令的执行状态的数据结构。
type CommandStatusReq struct {
	// 命令序列号，字符串型，组成为发送给客户端的8字节的当天日期+对应命令在数据库的11位ID值，即YYYYMMDDCOMMANDID.如果对应命令
	// 在数据库中的COMMANDID不足11位，是以前导0补足。
	// 本字段的值是apiServer和其客户端在通信过程中标识一个命令的标识符，包括apiServer查询命令执行状态或客户端返回命令执行状态
	// 与执行结果，均使用本字段的值来标识一个命令。
	// 如果此字段的值为0000000000000000000， 表示当前没有命令要客户端执行
	CommandSeq string `form:"commandSeq" json:"commandSeq" yaml:"commandSeq" xml:"commandSeq"`

	// 本字段的值可以是IP,HOSTNAME和MAC三个字符串的任意组合(注意：必须是大写），或者是长度不超过256个字符且不包含前面三个字符串的任意字符串(注意：该字符串区分大小写），以标识一个节点。当本子段的值为空时，则表示使用客户端定义的节点标识
	NodeIdentiferStr string `form:"nodeIdentifierStr" json:"nodeIdentifierStr" yaml:"nodeIdentifierStr" xml:"nodeIdentifierStr"`
}

// 用于记录客户端主动上报或者接收查询时，客户端返回给服务端关于命令执行状态的数据结构。
type CommandStatus struct {
	// 命令序列号，字符串型，组成为发送给客户端的8字节的当天日期+对应命令在数据库的11位ID值，即YYYYMMDDCOMMANDID.如果对应命令
	// 在数据库中的COMMANDID不足11位，是以前导0补足。
	// 本字段的值是apiServer和其客户端在通信过程中标识一个命令的标识符，包括apiServer查询命令执行状态或客户端返回命令执行状态
	// 与执行结果，均使用本字段的值来标识一个命令。
	// 如果此字段的值为0000000000000000000， 表示当前没有命令要客户端执行
	CommandSeq string `form:"commandSeq" json:"commandSeq" yaml:"commandSeq" xml:"commandSeq"`

	// 用于记录客户端根据配置或者apiServer的要求埴写的，用于标识自身的标识符。
	NodeIdentifier `form:"nodeIdentifier" json:"nodeIdentifier" yaml:"nodeIdentifier" xml:"nodeIdentifier"`

	// 命令的执行状态，对应前面定义的常量中之一
	StatusCode CommandStatusCode `form:"statusCode" json:"statusCode" yaml:"statusCode" xml:"statusCode"`

	// 命令执行状态的信息，最长不超过256个字符。对应执行成功的命令，本字段的值可以为空
	StatusMessage string `form:"statusMessage" json:"statusMessage" yaml:"statusMessage" xml:"statusMessage"`

	// 当命令执行成功后，本字段用于记录命令执行后的结果集，可以嵌套，嵌套内的数据类型为map[string]interface{} 或map[string]string
	// 且嵌套层级不超过3层。不执行结果没有数据集时，本字段的值可以为空。
	Data map[string]interface{} `form:"parameter" json:"parameter" yaml:"parameter" xml:"parameter"`

	// 当客户端没有发现所请求的命令的时候，应设置本字段的值为true，否则为false.当本字段的值为true时，通常意味以下几种情况：
	// 1. 请求中所提供的CommandSeq不正确；2. 客户端之前已经完成对所有日志的发送清除了对应命令的日志；
	NotCommand bool `form:"notCommand" json:"notCommand" yaml:"notCommand" xml:"notCommand"`
}

// 用户记录发送日志信息的数据结构
type Log struct {
	// 日志序列号，字符串型，组成为日志产的8字节当天日期+由6位组的000001开始的按步长为1增加的序号。apiServer根据这个序列号判断所接
	// 收的命令日志是否是已经接的重复日志
	// 当本字段的值为YYYYMMDD111111时，表示本条日志是对应命令的结束行日志。结束行日志不是命令的实际日志，只是类似文件EOF标识的标识
	LogSeq string `form:"logSeq" json:"logSeq" yaml:"logSeq" xml:"logSeq"`

	// 日志级别，对应于logrus包中定义的日志级别。对于结束行日志，本字段的值应设置为logrus.InfoLevel
	Level sysadmLog.Level `form:"level" json:"level" yaml:"level" xml:"level"`

	// 日志信息，除结束行的日志外，不得为空。
	Message string `form:"message" json:"message" yaml:"message" xml:"message"`
}

// 当apiServer以主动模式运行时，apiServer向客户发送的获取客户执行某条命令日志信息的请求数据结构。
type LogReq struct {
	// 命令序列号，字符串型，组成为发送给客户端的8字节的当天日期+对应命令在数据库的11位ID值，即YYYYMMDDCOMMANDID.如果对应命令
	// 在数据库中的COMMANDID不足11位，是以前导0补足。
	// 本字段的值是apiServer和其客户端在通信过程中标识一个命令的标识符，包括apiServer查询命令执行状态或客户端返回命令执行状态
	// 与执行结果，均使用本字段的值来标识一个命令。
	// 如果此字段的值为0000000000000000000， 表示当前没有命令要客户端执行
	CommandSeq string `form:"commandSeq" json:"commandSeq" yaml:"commandSeq" xml:"commandSeq"`

	// 本字段的值可以是IP,HOSTNAME和MAC三个字符串的任意组合(注意：必须是大写），或者是长度不超过256个字符且不包含前面三个字符串的任意字符串(注意：该字符串区分大小写），以标识一个节点。当本子段的值为空时，则表示使用客户端定义的节点标识
	NodeIdentiferStr string `form:"nodeIdentifierStr" json:"nodeIdentifierStr" yaml:"nodeIdentifierStr" xml:"nodeIdentifierStr"`

	// 日志序列号，字符串型，组成为日志产的8字节当天日期+由6位组的000001开始的按步长为1增加的序号。apiServer根据这个序列号判断所接
	// 收的命令日志是否是已经接的重复日志
	// 请求日志的起始序列号，即客户端需将序列号为本字段值的日志开始的日志发送给服务端
	// 当本字段的值为YYYYMMDD000000时，表示请求客户端发送从第一行开始的日志
	StartSeq string `form:"startSeq" json:"startSeq" yaml:"startSeq" xml:"startSeq"`

	// 表示本次apiServer可以接收的最大日志条数
	Num int `form:"num" json:"num" yaml:"num" xml:"num"`
}

// 记录客户端向服务端发送日志的数据结构，这包括当apiServer运行于主动模式时，客户端响应服务端的请求向服务端发送的日志数据；
//
//	这包括当apiServer运行于被动模式时，客户端主动向服务端发送的日志数据
type LogData struct {
	// 命令序列号，字符串型，组成为发送给客户端的8字节的当天日期+对应命令在数据库的11位ID值，即YYYYMMDDCOMMANDID.如果对应命令
	// 在数据库中的COMMANDID不足11位，是以前导0补足。
	// 本字段的值是apiServer和其客户端在通信过程中标识一个命令的标识符，包括apiServer查询命令执行状态或客户端返回命令执行状态
	// 与执行结果，均使用本字段的值来标识一个命令。
	// 如果此字段的值为0000000000000000000， 表示当前没有命令要客户端执行
	CommandSeq string `form:"commandSeq" json:"commandSeq" yaml:"commandSeq" xml:"commandSeq"`

	// 用于记录客户端根据配置或者apiServer的要求埴写的，用于标识自身的标识符。
	NodeIdentifier `form:"nodeIdentifier" json:"nodeIdentifier" yaml:"nodeIdentifier" xml:"nodeIdentifier"`

	// 记录客户端发送给服务端的具体的日志数据。注意：当本字段的值为空时，仅表示指定命令当前没有所需要的日志，并不表示日志已经发送结束。
	// 当命令已经成功执行,且命令执行过程中未产生任何日志时，仍需要发送一条标识日志已经发送结束的结束日志。
	Logs []Log `form:"logs" json:"logs" yaml:"logs" xml:"logs"`

	// 本次所发送的日志的总条数
	Total int `form:"total" json:"total" yaml:"total" xml:"total"`

	// 日志结束标志.当本字段值为true时，则表示当前数据包含对应命令所有的日志，否则表示对应命令还有日志未发送，
	// 如果是主动模式下，则apiserver需要再次请求余下日志; 被动模式下表示客户端下一次再发剩余日志
	EndFlag bool `form:"endflag" json:"endflag" yaml:"endflag" xml:"endflag"`

	// 当客户端没有发现所请求的命令的时候，应设置本字段的值为true，否则为false.当本字段的值为true时，通常意味以下几种情况：
	// 1. 请求中所提供的CommandSeq不正确；2. 客户端之前已经完成对所有日志的发送清除了对应命令的日志；
	NotCommand bool `form:"notCommand" json:"notCommand" yaml:"notCommand" xml:"notCommand"`
}

// 通用的，用户响应一般请求的响应数据的数据结构
type RepStatus struct {
	// 命令序列号，字符串型，组成为发送给客户端的8字节的当天日期+对应命令在数据库的11位ID值，即YYYYMMDDCOMMANDID.如果对应命令
	// 在数据库中的COMMANDID不足11位，是以前导0补足。
	// 本字段的值是apiServer和其客户端在通信过程中标识一个命令的标识符，包括apiServer查询命令执行状态或客户端返回命令执行状态
	// 与执行结果，均使用本字段的值来标识一个命令。
	// 如果此字段的值为0000000000000000000， 表示当前没有命令要客户端执行
	CommandSeq string `form:"commandSeq" json:"commandSeq" yaml:"commandSeq" xml:"commandSeq"`

	// 标识请求是否被正确接收，如果因为各种原因导到请求没有被正确接受，则将本字段设置成ComandStatusSendError状态，否则本字段设置成
	// ComandStatusReceived状态
	StatusCode CommandStatusCode `form:"statusCode" json:"statusCode" yaml:"statusCode" xml:"statusCode"`

	// 响应的信息，可以为空
	Message string `form:"message" json:"message" yaml:"message" xml:"message"`

	// 当客户端没有发现所请求的命令的时候，应设置本字段的值为true，否则为false.当本字段的值为true时，通常意味以下几种情况：
	// 1. 请求中所提供的CommandSeq不正确；2. 客户端之前已经完成对所有日志的发送清除了对应命令的日志；
	NotCommand bool `form:"notCommand" json:"notCommand" yaml:"notCommand" xml:"notCommand"`
}

// 记录待发送的command数据
type commandDataBeSent struct {
	// hostid identified a host
	hostID int32

	// 当apiserver是以主动模式运行时，连接客户端agent的地址或域名
	agentAddress string

	// 当apiserver以主动模式运行时，apiserver向agent发送命令时，请求的发送目 标路径。如果本字段为空，则apiserver默认向/receiveCommand请求命令的状态
	commandUri string

	// 当apiserver以主动模式运行时，apiserver连接agent是否使用TLS.0表示否，否则表示是
	agentIsTls bool

	// 当apiserver以主动模式运行时，apiserver连接agent时主动采用TLS,本子段是CA证书内容
	agentCa string

	// 当apiserver以主动模式运行时，apiserver连接agent时主动采用TLS,本子段是证书内容
	agentCert string

	// 当apiserver以主动模式运行时，apiserver连接agent时主动采用TLS,本子段是密钥内容
	agentKey string

	// 当apiserver以主动模式运行时，apiserver连接agent时主动采用TLS,指定是否跳过检查不合法证书。1表示是，否则为否
	insecureSkipVerify bool

	// agent listen port when it runing in active mode
	agentPort int

	// 待发送的command 数据
	CommandData
}

// 记录当apiserver以主动模式运行时，apiserver主动向客户端发送请求时，客户端的参数数据，其中uri为请求的uri路径。
type clientRequestData struct {
	// hostid identified a host
	hostID int32

	// 当apiserver是以主动模式运行时，连接客户端agent的地址或域名
	agentAddress string

	// 当apiserver以主动模式运行时，apiserver连接agent是否使用TLS.0表示否，否则表示是
	agentIsTls bool

	// 当apiserver以主动模式运行时，apiserver连接agent时主动采用TLS,本子段是CA证书内容
	agentCa string

	// 当apiserver以主动模式运行时，apiserver连接agent时主动采用TLS,本子段是证书内容
	agentCert string

	// 当apiserver以主动模式运行时，apiserver连接agent时主动采用TLS,本子段是密钥内容
	agentKey string

	// 当apiserver以主动模式运行时，apiserver连接agent时主动采用TLS,指定是否跳过检查不合法证书。1表示是，否则为否
	insecureSkipVerify bool

	// agent listen port when it runing in active mode
	agentPort int

	// to save the uri address where the apiserver send command data to
	commandUri string

	// to save the uri address where the apiserver get command status from
	commandStatusUri string

	// to save the uri address where the apiserver get command logs from
	commandLogsUri string

	Command
}

// hold the running data related to command. this items can be configurable in the feature
type runDataForCommand struct {
	// 命令的最大重试次数
	maxTryTimes int

	// concurrency number of apiserver sending command data to agent when apiserver is running in active mode
	concurrencySendCommand int

	// concurrency number of apiserver get command status from agent when apiserver is running in active mode
	concurrencyGetCommandStatus int

	// concurrency number of apiserver get command log from agent when apiserver is running in active mode
	concurrencyGetCommandLog int

	// 日志信息在redis里存储的路径
	logRootPathInRedis string

	// 每次获取命令日志的最大条数
	maxGetLogNumPerTime int
}
