package main

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/mirrors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var DebianDistribution = map[string][]string{
	"ubuntu": []string{"bionic", "devel", "focal", "jammy", "lunar", "mantic", "noble", "oracular", "trusty", "xenial"},
	"debian": []string{"buster", "bullseye", "jessie", "sid", "stretch", "wheezy"},
}
var funcMap = map[string]func(w http.ResponseWriter, r *http.Request){
	"pypi":         mirrors.Pypi,
	"dockerhub":    mirrors.Docker,
	"ubuntu":       mirrors.Ubuntu,
	"ubuntu_ports": mirrors.UbuntuPorts,
	"debian":       mirrors.Debian,
	"kali":         mirrors.Kali,
	"alpine":       mirrors.Alpine,
	"npm":          mirrors.Npm,
	"github":       mirrors.Github,
}
var uaMap = map[string]string{
	//存放不同UA的特征，比如docker的特征就是ua中包含docker，key特征，value是上面funcMap中定义的镜像类型
	"docker":         "dockerhub",
	"git":            "github",
	"pip":            "pypi",
	"npm":            "npm",
	"node":           "npm",
	"Go-http-client": "go",
	"APT-HTTP":       "debian-linux",
}

func whichDebianDistribution(request *http.Request) string {
	if strings.HasPrefix(request.URL.Path, "/ubuntu-ports") {
		return "ubuntu_ports"
	}
	if strings.HasPrefix(request.URL.Path, "/ubuntu") {
		return "ubuntu"
	}

	if strings.HasPrefix(request.URL.Path, "/debian") {
		return "debian"
	}
	if strings.HasPrefix(request.URL.Path, "/kali") {
		return "kali"
	}
	return ""
}
func whichMirror(request *http.Request) string {
	typeFromUA := ""
	//初步判断
	for key, value := range uaMap {
		uas := strings.Split(request.UserAgent(), " ")
		if strings.Contains(uas[0], key) {
			log.Println("Mirror:", value)
			log.Println(request.URL.Path)
			typeFromUA = value
		}
	}
	//针对特殊情况的判断
	//debian系
	if typeFromUA == "debian-linux" {
		return whichDebianDistribution(request)
	}
	//alpine
	if strings.HasPrefix(request.URL.Path, "/alpine") {
		return "alpine"
	}
	//go 包获取，暂时只支持github
	if request.URL.Query()["go-get"] != nil {
		return "github"
	}
	mirrorType := typeFromUA
	return mirrorType
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
		if config.ServerConfig.HostControll {
			if !config.ServerConfig.CheckHost(r.Host) {
				http.NotFound(w, r)
				return
			}
		}
		mirrorType := whichMirror(r)
		if mirrorType == "" {
			http.NotFound(w, r)
			return
		}
		if f, ok := funcMap[mirrorType]; ok {
			f(w, r)
			return
		}
	})
	if config.ServerConfig.Tls {
		tlssrv := &http.Server{
			Addr:    fmt.Sprintf(":%d", config.ServerConfig.TlsPort),
			Handler: server,
		}
		go tlssrv.ListenAndServeTLS(config.ServerConfig.CertFile, config.ServerConfig.KeyFile)
		log.Println("TLS Server started")
	}
	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.ServerConfig.Port),
		Handler: server,
		//TLSConfig:    cfg,
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Println("Server started")
	srv.ListenAndServe()
}
