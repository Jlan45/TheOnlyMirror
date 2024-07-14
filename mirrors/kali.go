package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
)

func Kali(w http.ResponseWriter, r *http.Request) {
	kaliProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("kali"))
	kaliProxy.ServeHTTP(w, r)
}
