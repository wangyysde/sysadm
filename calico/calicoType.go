/* =============================================================
* @Author:  Wayne Wang <yuying.wang@sincerecloud.com>
* @ Last Modified At: 2023.06.01
* @Copyright (c) 2023 Sincerecloud
* @HomePage: https://www.sincerecloud.com/
*
*  定义与calico相关的公共数据结构
 */

package calico

type GlobalConf struct {
	// CalicoBackend mode
	CalicoBackend string

	// 镜像仓库地址
	ImageRepository string

	// calico版本
	ImageVersion string
}
