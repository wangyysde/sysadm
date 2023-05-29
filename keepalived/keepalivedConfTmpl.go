/* =============================================================
* @Author:  Wayne Wang <yuying.wang@sincerecloud.com>
* @ Last Modified At: 2023.05.10
* @Copyright (c) 2023 Sincerecloud
* @HomePage: https://www.sincerecloud.com/
*
*  定义keepalived配置文件模板
 */

package keepalived

// keepalived 配置文件的模板
var keepConfigTmpl string = `
! Configuration File for keepalived

global_defs {
   {{ $emailLen := len .NotificationEmail -}}
   {{ if gt $emailLen 0 -}}
   notification_email {
     {{- range $i, $v := .NotificationEmail }}
     {{ $v }} 
	 {{- end }}
   }
   {{- end }}
   {{- if .NotificationEmailFrom }}
   notification_email_from {{ .NotificationEmailFrom }}
   {{- end }}
   {{- if .SmtpServer }} 
   smtp_server {{ .SmtpServer }} {{ .SmtpPort }}
   {{- end }}
   {{- if .SmtpConnectTimeout }}
   smtp_connect_timeout {{ .SmtpConnectTimeout }}
   {{- end }}
   {{- if .RouterId }}
   router_id {{ .RouterId }}
   {{- end }}
}
{{ $instanceLen := len .VrrpInstances -}}
{{ if gt $instanceLen 0 -}}
{{- range $i, $v := .VrrpInstances }}
vrrp_instance {{ $v.Name }} {
    state {{ $v.State }}
    interface {{ $v.Interface }}
    virtual_router_id {{ $v.VirtualRouterId }}
    priority {{ $v.Priority }}
    advert_int {{ $v.AdvertInt }}
    {{- if $v.Authentication }}
    authentication {
        auth_type {{ $v.Authentication.AuthType }}
        auth_pass {{ $v.Authentication.AuthPass }}
    }
	{{- end }}
	{{ $vipLen := len $v.VirtualIpaddress -}}
	{{ if gt $vipLen 0 }}
    virtual_ipaddress {
		{{- range $ii, $vv := $v.VirtualIpaddress }}
        {{ $vv }}
		{{- end }}
    }
	{{- end }}
}

{{ end }}
{{- end }}
{{ $virturalLen := len .VirtualServers -}}
{{ if gt $virturalLen 0 -}}
{{- range $i, $v := .VirtualServers }}
virtual_server {{ $v.VirtualIp }} {{ $v.Port }} {
    delay_loop {{ $v.DelayLoop }}
    lb_algo {{ $v.LbAlgo }}
    lb_kind {{ $v.LbKind }}
    persistence_timeout {{ $v.PersistenceTimeout }} 
    protocol {{ $v.Protocol }} 
	{{ $rsLen := len $v.RealServers -}}
	{{ if gt $rsLen 0 }}
	{{- range $ii, $vv := $v.RealServers }}
    real_server {{ $vv.RealIp }} {{ $vv.Port }} 
        weight {{ $vv.Weight }}
        TCP_CHECK  {
            connect_timeout {{ $vv.TcpCheck.ConnectTimeout }}
            retry: {{ $vv.TcpCheck.Retry }}
        }
    }
	{{- end }}
	{{- end }}
}
{{ end }}
{{- end }}
`
