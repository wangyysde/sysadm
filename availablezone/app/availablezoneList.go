package app

import (
	"github.com/wangyysde/sysadmServer"
	"net/http"
	"sysadm/sysadmLog"
	"sysadm/user"
)

func listHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(700050001, "debug", "now handling k8s cluster list"))

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700050003, "error", "user should login %s", e))
		runData.logEntity.LogErrors(errs)

		tplData := map[string]interface{}{
			"errormessage": "user should login",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, tplData)
		return
	}

}
