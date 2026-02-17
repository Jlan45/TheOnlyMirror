package web

import (
	"TheOnlyMirror/config"
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed templates/index.html
var templateFS embed.FS

var indexTmpl *template.Template

func init() {
	indexTmpl = template.Must(template.ParseFS(templateFS, "templates/index.html"))
}

type CodeSection struct {
	Title string
	Code  string
}

type MirrorInfo struct {
	ID          string
	Name        string
	Icon        string
	Description string
	DetectClass string
	DetectLabel string
	Sections    []CodeSection
	Note        string
}

type IndexData struct {
	Domain  string
	Mirrors []MirrorInfo
}

func getMirrorInfos(domain, scheme string) []MirrorInfo {
	base := fmt.Sprintf("%s://%s", scheme, domain)

	return []MirrorInfo{
		{
			ID:          "dockerhub",
			Name:        "Docker Hub",
			Icon:        "\U0001F40B",
			Description: "Docker 容器镜像加速",
			DetectClass: "tag-ua",
			DetectLabel: "UA 自动检测",
			Sections: []CodeSection{
				{
					Title: "配置 Docker 守护进程",
					Code: fmt.Sprintf(`# 编辑 /etc/docker/daemon.json
{
  "registry-mirrors": ["%s"]
}

# 重启 Docker
sudo systemctl daemon-reload
sudo systemctl restart docker`, base),
				},
				{
					Title: "验证配置",
					Code:  `docker info | grep "Registry Mirrors" -A 1`,
				},
			},
			Note: "配置后所有 docker pull 命令将自动通过镜像加速。",
		},
		{
			ID:          "github",
			Name:        "GitHub",
			Icon:        "\U0001F4E6",
			Description: "GitHub 仓库克隆加速",
			DetectClass: "tag-ua",
			DetectLabel: "UA 自动检测",
			Sections: []CodeSection{
				{
					Title: "克隆仓库",
					Code:  fmt.Sprintf(`# 将 github.com 替换为镜像域名即可
git clone %s/<owner>/<repo>.git`, base),
				},
				{
					Title: "已有仓库切换 remote",
					Code: fmt.Sprintf(`git remote set-url origin %s/<owner>/<repo>.git
git pull`, base),
				},
			},
		},
		{
			ID:          "pypi",
			Name:        "PyPI",
			Icon:        "\U0001F40D",
			Description: "Python 包镜像加速",
			DetectClass: "tag-ua",
			DetectLabel: "UA 自动检测",
			Sections: []CodeSection{
				{
					Title: "临时使用",
					Code:  fmt.Sprintf(`pip install -i %s/pypi/simple/ <package-name>`, base),
				},
				{
					Title: "永久配置",
					Code: fmt.Sprintf(`pip config set global.index-url %s/pypi/simple/
pip config set global.trusted-host %s`, base, domain),
				},
				{
					Title: "或编辑 ~/.pip/pip.conf",
					Code: fmt.Sprintf(`[global]
index-url = %s/pypi/simple/
trusted-host = %s`, base, domain),
				},
			},
		},
		{
			ID:          "npm",
			Name:        "npm",
			Icon:        "\U0001F4E6",
			Description: "Node.js 包镜像加速",
			DetectClass: "tag-ua",
			DetectLabel: "UA 自动检测",
			Sections: []CodeSection{
				{
					Title: "临时使用",
					Code:  fmt.Sprintf(`npm install <package-name> --registry=%s`, base),
				},
				{
					Title: "永久配置",
					Code:  fmt.Sprintf(`npm config set registry %s`, base),
				},
				{
					Title: "使用 yarn",
					Code:  fmt.Sprintf(`yarn config set registry %s`, base),
				},
				{
					Title: "验证配置",
					Code:  `npm config get registry`,
				},
			},
		},
		{
			ID:          "go",
			Name:        "Go Modules",
			Icon:        "\U0001F439",
			Description: "Go 模块代理加速 (GitHub 源)",
			DetectClass: "tag-ua",
			DetectLabel: "UA 自动检测",
			Sections: []CodeSection{
				{
					Title: "设置 GOPROXY",
					Code: fmt.Sprintf(`export GOPROXY=%s,direct
export GONOSUMCHECK=*`, base),
				},
				{
					Title: "永久配置 (添加到 ~/.bashrc 或 ~/.zshrc)",
					Code: fmt.Sprintf(`echo 'export GOPROXY=%s,direct' >> ~/.bashrc
echo 'export GONOSUMCHECK=*' >> ~/.bashrc
source ~/.bashrc`, base),
				},
			},
			Note: "当前 Go 模块代理通过 GitHub 镜像实现，仅支持托管在 GitHub 上的模块。",
		},
		{
			ID:          "ubuntu",
			Name:        "Ubuntu",
			Icon:        "\U0001F427",
			Description: "Ubuntu APT 软件源镜像",
			DetectClass: "tag-path",
			DetectLabel: "路径匹配 /ubuntu",
			Sections: []CodeSection{
				{
					Title: "配置 sources.list (x86_64)",
					Code: fmt.Sprintf(`# 备份原文件
sudo cp /etc/apt/sources.list /etc/apt/sources.list.bak

# 替换源地址
sudo sed -i 's|http://archive.ubuntu.com|%s|g' /etc/apt/sources.list

# 更新索引
sudo apt update`, base),
				},
				{
					Title: "配置 sources.list (ARM / Ports)",
					Code: fmt.Sprintf(`sudo sed -i 's|http://ports.ubuntu.com|%s|g' /etc/apt/sources.list

sudo apt update`, base),
				},
				{
					Title: "手动编辑 /etc/apt/sources.list 示例 (Noble 24.04)",
					Code: fmt.Sprintf(`deb %s/ubuntu noble main restricted universe multiverse
deb %s/ubuntu noble-updates main restricted universe multiverse
deb %s/ubuntu noble-security main restricted universe multiverse`, base, base, base),
				},
			},
		},
		{
			ID:          "debian",
			Name:        "Debian",
			Icon:        "\U0001F427",
			Description: "Debian APT 软件源镜像",
			DetectClass: "tag-path",
			DetectLabel: "路径匹配 /debian",
			Sections: []CodeSection{
				{
					Title: "配置 sources.list",
					Code: fmt.Sprintf(`# 备份原文件
sudo cp /etc/apt/sources.list /etc/apt/sources.list.bak

# 替换源地址
sudo sed -i 's|http://deb.debian.org|%s|g' /etc/apt/sources.list

# 更新索引
sudo apt update`, base),
				},
				{
					Title: "手动编辑示例 (Bookworm)",
					Code: fmt.Sprintf(`deb %s/debian bookworm main contrib non-free non-free-firmware
deb %s/debian bookworm-updates main contrib non-free non-free-firmware
deb %s/debian-security bookworm-security main contrib non-free non-free-firmware`, base, base, base),
				},
			},
		},
		{
			ID:          "kali",
			Name:        "Kali Linux",
			Icon:        "\U0001F409",
			Description: "Kali Linux APT 软件源镜像",
			DetectClass: "tag-path",
			DetectLabel: "路径匹配 /kali",
			Sections: []CodeSection{
				{
					Title: "配置 sources.list",
					Code: fmt.Sprintf(`# 编辑 /etc/apt/sources.list
deb %s/kali kali-rolling main contrib non-free non-free-firmware

# 更新索引
sudo apt update`, base),
				},
			},
		},
		{
			ID:          "alpine",
			Name:        "Alpine Linux",
			Icon:        "\U0001F3D4\uFE0F",
			Description: "Alpine APK 软件源镜像",
			DetectClass: "tag-path",
			DetectLabel: "路径匹配 /alpine",
			Sections: []CodeSection{
				{
					Title: "配置 repositories",
					Code: fmt.Sprintf(`# 编辑 /etc/apk/repositories
%s/alpine/v3.21/main
%s/alpine/v3.21/community

# 更新索引
apk update`, base, base),
				},
				{
					Title: "一键替换",
					Code: fmt.Sprintf(`sed -i 's|https://dl-cdn.alpinelinux.org|%s|g' /etc/apk/repositories
apk update`, base),
				},
			},
		},
		{
			ID:          "aosp",
			Name:        "AOSP",
			Icon:        "\U0001F4F1",
			Description: "Android 开源项目镜像",
			DetectClass: "tag-ua",
			DetectLabel: "UA 自动检测",
			Sections: []CodeSection{
				{
					Title: "初始化 repo",
					Code: fmt.Sprintf(`repo init -u %s/platform/manifest -b main`, base),
				},
				{
					Title: "同步代码",
					Code:  `repo sync -c -j8`,
				},
			},
		},
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	scheme := "https"
	if !config.ServerConfig.Tls {
		scheme = "http"
	}

	allMirrors := getMirrorInfos(config.ServerConfig.Domain, scheme)

	// Filter to only show mirrors in the mirrorList config
	enabledSet := make(map[string]bool)
	for _, m := range config.ServerConfig.MirrorList {
		enabledSet[m] = true
	}

	var filtered []MirrorInfo
	for _, m := range allMirrors {
		if enabledSet[m.ID] {
			filtered = append(filtered, m)
		}
	}
	// If mirrorList is empty, show all mirrors
	if len(config.ServerConfig.MirrorList) == 0 {
		filtered = allMirrors
	}

	data := IndexData{
		Domain:  config.ServerConfig.Domain,
		Mirrors: filtered,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	indexTmpl.Execute(w, data)
}