/* =============================================================
* @Author:  Wayne Wang <yuying.wang@sincerecloud.com>
* @ Last Modified At: 2023.06.01
* @Copyright (c) 2023 Sincerecloud
* @HomePage: https://www.sincerecloud.com/
*
*  定义公共函数
 */

package calico

import (
	"bytes"
	"fmt"
	"goclient/k8sclient"
	"strings"
	"text/template"
)

func ValidConfig(conf *GlobalConf) error {
	if strings.TrimSpace(conf.CalicoBackend) == "" {
		conf.CalicoBackend = defaultCalicoBackend
	}

	if strings.TrimSpace(conf.ImageRepository) == "" {
		conf.ImageRepository = "defaultImageRepository"
	}

	if strings.TrimSpace(conf.ImageVersion) == "" {
		conf.ImageVersion = defaultImageVersion
	}

	return nil
}

func NewConfig(backend, repository, version string) (*GlobalConf, error) {
	return &GlobalConf{CalicoBackend: backend, ImageRepository: repository, ImageVersion: version}, nil
}

// convert template of calico to the content in yaml then apply them to k8s cluster
func ApplyCalico(kubeconf string, conf *GlobalConf) error {
	config, e := k8sclient.BuildKubeConf("", kubeconf)
	if e != nil {
		return e
	}

	dyClient, e := k8sclient.BuildDynamicClient(config)
	if e != nil {
		return e
	}

	apiRL, e := k8sclient.GetApiResourcesList(config)
	if e != nil {
		return e
	}

	partContent, e := buildYamlContent(tplPart1, conf)
	if e != nil {
		return e
	}

	e = k8sclient.ApplyFromYaml(string(partContent), config, dyClient, apiRL)
	if e != nil {
		return e
	}

	partContent, e = buildYamlContent(tplPart2, conf)
	if e != nil {
		return e
	}

	partContent, e = buildYamlContent(tplPart3, conf)
	if e != nil {
		return e
	}

	partContent, e = buildYamlContent(tplPart4, conf)
	if e != nil {
		return e
	}

	return nil
}

// build the content in yaml format for applying calico
func buildYamlContent(tmpContent string, conf *GlobalConf) (string, error) {
	tmpContent = strings.TrimSpace(tmpContent)
	if tmpContent == "" {
		return "", fmt.Errorf("can not build yaml content with empty string")
	}

	tmpl, e := template.New("calico").Parse(tmpContent)
	if e != nil {
		return "", e
	}

	var tpl bytes.Buffer
	e = tmpl.Execute(&tpl, *conf)
	if e != nil {
		return "", e
	}

	return tpl.String(), nil
}
