/* =============================================================
* @Author:  Wayne Wang <yuying.wang@sincerecloud.com>
* @ Last Modified At: 2023.05.10
* @Copyright (c) 2023 Sincerecloud
* @HomePage: https://www.sincerecloud.com/
*
*  定义与keepalived相关的公共数据结构
 */
package keepalived

// 全局配置
// REF: https://www.keepalived.org/manpage.html
type GlobalDefs struct {
	// 当主从切换时通知用户的邮件列表
	NotificationEmail []string
	// 发送通知邮件的邮箱帐号
	// email from address that will be in the header
	// (default: keepalived@<local host name>)
	NotificationEmailFrom string
	// 发送邮件的smtp服务器地址
	// Remote SMTP server used to send notification email.
	SmtpServer string
	// Remote SMTP server port
	SmtpPort string
	// SMTP服务器连接超时时间
	SmtpConnectTimeout int
	// 当前节点的服务器标识，一般是主机名。这个值主备不一定非要一样
	RouterId string
}

// 认证配置
type Authentication struct {
	// 要有PASS和AH两种
	AuthType string
	// 验证密码，同一个vrrp_instance下MASTER和BACKUP密码必须相同
	AuthPass string
}

// vrrp实例配置
type VrrpInstance struct {
	// 实例名
	Name string
	// 置lvs的状态，MASTER和BACKUP两种，必须大写
	State string
	// VIP地址绑定的接口名称，lo或eth0等
	Interface string
	// 设置虚拟路由标识，这个标识是一个数字，同一个vrrp实例使用唯一标识。同一个虚拟路由不同节点的这个值应该相同
	VirtualRouterId int
	// 义优先级，数字越大优先级越高，在一个vrrp——instance下，master的优先级必须大于backup。数字应该在0~255之间
	Priority int
	// 设定master与backup负载均衡器之间同步检查的时间间隔，单位是秒
	AdvertInt int
	// 验证类型和密码
	Authentication
	// 设置虚拟ip地址，可以设置多个，每行一个
	VirtualIpaddress []string
}

// realserver状态监测设置部分单位秒
type TcpCheck struct {
	// 连接超时时间
	ConnectTimeout int
	// 重试次数
	Retry int
}

type RealServer struct {
	// 真实服务器的IP地址
	RealIp string
	// 真实服务器服务的端口号
	Port int
	// 权重，数字越大权重越高
	Weight int
	// realserver状态监测设置部分单位秒
	TcpCheck
}

// 虚拟服务器设置
type VirtualServer struct {
	// 虚拟ip地址,不带端口
	VirtualIp string
	// 对外服务端口号
	Port int
	// 健康检查时间间隔
	DelayLoop int
	// 负载均衡调度算法
	LbAlgo string
	// 负载均衡转发规则
	LbKind string
	// 会话保持时间
	PersistenceTimeout int
	// 转发协议类型，有TCP和UDP两种
	Protocol string
	// 虚拟服务器设置
	RealServers []RealServer
}

// keepalived 配置信息数据结构
type KeepalivedConf struct {
	// 全局配置
	GlobalDefs
	// VRRP实例配置
	VrrpInstances []VrrpInstance
	// 虚拟服务器配置
	VirtualServers []VirtualServer
}
