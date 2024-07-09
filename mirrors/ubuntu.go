package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
)

func Ubuntu(w http.ResponseWriter, r *http.Request) {
	ubuntuProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("ubuntu"))
	ubuntuProxy.ModifyResponse = func(resp *http.Response) error {
		println("ubuntu")
		return nil
	}
	ubuntuProxy.ServeHTTP(w, r)
}
