package main

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/mirrors"
	"log"
	"net/http"
	"strings"
)

func github(w http.ResponseWriter, r *http.Request) {

}

var funcMap = map[string]func(w http.ResponseWriter, r *http.Request){
	"pypi":      mirrors.Pypi,
	"dockerhub": mirrors.Dockerhub,
}
var uaMap = map[string]string{
	//存放不同UA的特征，比如docker的特征就是ua中包含docker，key特征，value是上面funcMap中定义的镜像类型
	"docker":         "dockerhub",
	"git":            "github",
	"pip":            "pypi",
	"npm":            "npm",
	"node":           "npm",
	"Go-http-client": "go",
}

func whichMirror(UA string) string {
	for key, value := range uaMap {
		if strings.Contains(UA, key) {
			log.Println("UA:", UA, "Mirror:", value)
			return value
		}
	}
	return ""
}
func main() {
	//要做的镜像： github dockerhub pypi npm golang
	// 定义多个目标服务器
	if config.LoadConfig() != nil {
		log.Fatal("load config error")
		return
	}
	server := http.NewServeMux()
	// 创建反向代理处理函数
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mirrorType := whichMirror(r.UserAgent())
		if mirrorType == "" {
			http.NotFound(w, r)
			return
		}
		if f, ok := funcMap[mirrorType]; ok {
			f(w, r)
			return
		}
	})

	// 启动服务器
	srv := &http.Server{
		Addr:    ":8090",
		Handler: server,
		//TLSConfig:    cfg,
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	srv.ListenAndServe()
}
