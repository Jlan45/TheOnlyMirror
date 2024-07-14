package mirrors

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/utils"
	"net/http"
	"net/url"
	"strings"
)

func Docker(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/registry") {
		Dockerblob(w, r)
	} else {
		Dockerhub(w, r)
	}
}
func Dockerhub(w http.ResponseWriter, r *http.Request) {
	dockerhubProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("dockerhub"))
	dockerhubProxy.ModifyResponse = func(resp *http.Response) error {
		if resp.StatusCode == 307 {
			//处理跳转blob站点
			orgLocation, _ := url.Parse(resp.Header.Get("location"))
			orgLocation.Host = config.ServerConfig.Domain
			orgLocation.Scheme = config.ServerConfig.GetMyUrl().Scheme
			resp.Header.Set("location", orgLocation.String())

		}
		return nil
	}
	dockerhubProxy.ServeHTTP(w, r)
}
func Dockerblob(w http.ResponseWriter, r *http.Request) {
	dockerblobProxy := utils.GetSimpleReverseProxy(config.ServerConfig.GetSourceUrl("dockerblob"))
	dockerblobProxy.ServeHTTP(w, r)
}
