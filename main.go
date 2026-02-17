package main

import (
	"TheOnlyMirror/config"
	"TheOnlyMirror/mirrors"
	"TheOnlyMirror/web"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
	"aosp":         mirrors.Aosp,
	"go":           mirrors.Github, // Go HTTP client uses the same handler as Github
}
var uaMap = map[string]string{
	//存放不同UA的特征，比如docker的特征就是ua中包含docker，key特征，value是上面funcMap中定义的镜像类型
	"docker":         "dockerhub",
	"git":            "github",
	"git-repo":       "aosp",
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
		if strings.Contains(request.UserAgent(), key) {
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
	if err := config.LoadConfig(); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Config file 'config.json' not found. Please rename 'config.example.json' to 'config.json' first.")
		}
		log.Fatal("load config error:", err)
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
			if r.URL.Path == "/" {
				web.IndexHandler(w, r)
				return
			}
			http.NotFound(w, r)
			return
		}
		if f, ok := funcMap[mirrorType]; ok {
			f(w, r)
			return
		}
		// If mirror type is not recognized, return 404
		http.NotFound(w, r)
	})
	if config.ServerConfig.Tls {
		tlssrv := &http.Server{
			Addr:              fmt.Sprintf(":%d", config.ServerConfig.TlsPort),
			Handler:           server,
			ReadHeaderTimeout: 10 * time.Second,
			IdleTimeout:       120 * time.Second,
		}
		go func() {
			if err := tlssrv.ListenAndServeTLS(config.ServerConfig.CertFile, config.ServerConfig.KeyFile); err != nil {
				log.Fatal("TLS server error: ", err)
			}
		}()
		log.Println("TLS Server started")
	}
	// 启动服务器
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.ServerConfig.Port),
		Handler:           server,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	log.Println("Server started")
	srv.ListenAndServe()
}
