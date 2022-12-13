package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var start_time = int(time.Now().Unix() - 1670949827)

func main() {
	println(strconv.Itoa(start_time))
	osname, arch := determineEnvironment()
	println("OS:" + osname + " Arch:" + arch)

	proxyString := flag.String("p", "", "proxy string")
	flag.Parse()

	python_achive_file := downloadPython(osname, arch, *proxyString)
	installPython(osname, arch, python_achive_file, *proxyString)

	mcl_file := downloadMCLInstaller(osname, arch, *proxyString)
	installMCL(osname, arch, mcl_file, *proxyString)

	cloneSource()
	makeConfig(osname)

	writeLaunchScript(osname, arch)
	println("安装完成!")
	println("请先运行run-mirai.bat登录qq号成功之后，保持运行状态，运行run-bot.bat")
}

// 确定OS和架构
func determineEnvironment() (osname, arch string) {
	return runtime.GOOS, runtime.GOARCH
}

func downloadPython(osname, arch, proxy string) string {
	python_url := ""
	if osname == "windows" {
		if arch == "386" {
			python_url = "https://www.python.org/ftp/python/3.10.9/python-3.10.9-embed-win32.zip"
		} else if arch == "amd64" {
			python_url = "https://www.python.org/ftp/python/3.10.9/python-3.10.9-embed-amd64.zip"
		} else {
			panic("不支持的操作系统、架构组合")
		}
	} else if osname == "linux" {
		python_url = "https://www.python.org/ftp/python/3.10.9/Python-3.10.9.tgz"
	}

	println("下载Python环境:" + python_url)
	return DownloadFile(python_url, "./python", proxy)
}

func installPython(osname, arch, achive_file, proxy string) {
	println("安装Python环境")
	if osname == "windows" {
		//解压归档文件
		DeCompress(achive_file, "./python/")
		//下载pip
		println("下载pip")
		pip_url := "https://bootstrap.pypa.io/get-pip.py"
		pip_file := DownloadFile(pip_url, "./python/", proxy)
		//安装pip
		println("安装pip")
		RunCMDPipe("安装pip", ".", "./python/python.exe ", pip_file)
		ReplaceStringInFile("./python/python310._pth", "#import site", "import site")

		//安装依赖
		println("安装依赖")
		RunCMDPipe("安装依赖", ".", "./python/Scripts/pip.exe ", "install", "pymysql", "yiri-mirai", "openai", "colorlog", "func_timeout")
		RunCMDPipe("安装依赖", ".", "./python/Scripts/pip.exe ", "install", "websockets", "--upgrade")

	} else if osname == "linux" {
		// DeCompress(achive_file,"./python/")
		RunCMDPipe("解压Python源码", ".", "tar", "zxvf", achive_file, "-C", "./python")
		linux_installerCompiler()
		pwd, _ := RunCMDPipe("检查pwd", "./python/", "pwd")
		pwd = strings.Trim(pwd, "\n")
		RunCMDPipe("配置编译环境", "./python/Python-3.10.9", "./configure", "--prefix="+pwd)
		RunCMDPipe("编译Python", "./python/Python-3.10.9", "make")
		RunCMDPipe("安装Python", "./python/Python-3.10.9", "make", "install")

		println("安装依赖")
		RunCMDPipe("安装依赖", ".", "python/bin/pip3", "install", "pymysql", "yiri-mirai", "openai", "colorlog", "func_timeout")
		RunCMDPipe("安装依赖", ".", "python/bin/pip3", "install", "websockets", "--upgrade")
	}
}

func linux_installerCompiler() {

	result, _ := RunCMDPipe("检查包管理器", ".", "apt")
	print(result)
	if result == "" {
		result, _ := RunCMDPipe("检查包管理器", ".", "yum")
		if result == "" {
			fmt.Println("不支持的Linux操作系统")
			os.Exit(-1)
		} else {
			RunCMDPipe("安装编译环境", ".", "yum", "install", "zlib-devel", "bzip2-devel", "openssl", "openssl-devel", "ncurses-devel", "sqlite-devel",
				"readline-devel", "tk-devel", "gcc", "make", "readline", "libffi-devel", "-y") //zlib-devel bzip2-devel openssl openssl-devel ncurses-devel sqlite-devel readline-devel tk-devel gcc make readline libffi-devel -y
		}
	} else {
		RunCMDPipe("安装编译环境", ".", "apt", "update")
		RunCMDPipe("安装编译环境", ".", "apt", "install", "build-essential", "zlib1g-dev", "libncurses5-dev", "libgdbm-dev", "libnss3-dev", "libssl-dev", "libreadline-dev", "libffi-dev", "libsqlite3-dev", "wget", "libbz2-dev")
	}
}

func downloadMCLInstaller(osname, arch, proxy string) string {
	mcl_url := ""
	if osname == "windows" {
		if arch == "386" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-windows-x86.exe"
		} else if arch == "amd64" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-windows-amd64.exe"
		} else {
			panic("不支持的操作系统、架构组合")
		}
	} else if osname == "linux" {
		if arch == "386" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-linux-amd64-musl"
		} else if arch == "amd64" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-linux-amd64-musl"
		} else if arch == "arm" {
			mcl_url = "https://github.com/iTXTech/mcl-installer/releases/download/a02f711/mcl-installer-a02f711-linux-arm-musl"
		} else {
			panic("不支持的操作系统、架构组合")
		}
	}

	println("下载MCL安装器:" + mcl_url)
	return DownloadFile(mcl_url, "./mirai", proxy)
}

func installMCL(osname, arch, installer_file, proxy string) {
	println("安装mirai")
	installer_file = strings.ReplaceAll(installer_file, "mirai/", "")
	println(installer_file)
	if osname == "windows" {
		RunCMDPipe("安装mirai", "./mirai", installer_file)
	} else if osname == "linux" {
		RunCMDPipe("安装mirai", "./mirai", "chmod", "+x", installer_file)
		RunCMDPipe("安装mirai", "./mirai", installer_file)
	}

	RunCMDTillStringOutput("配置mirai", "./mirai", "I/main: mirai-console started successfully.", "./java/bin/java", "-jar", "mcl.jar")
	RunCMDPipe("配置mirai", "./mirai", "./java/bin/java", "-jar", "mcl.jar", "--update-package", "net.mamoe:mirai-api-http", "--channel", "stable-v2", "--type", "plugin")
	RunCMDTillStringOutput("配置mirai", "./mirai", "I/main: mirai-console started successfully.", "./java/bin/java", "-jar", "mcl.jar")

	//更改协议
	ReplaceStringInFile("./mirai/config/Console/AutoLogin.yml", "protocol: ANDROID_PHONE", "protocol: ANDROID_PAD")
}

func cloneSource() {
	println("下载源代码")
	RunCMDPipe("下载源代码", ".", "git", "clone", "https://gitee.com/RockChin/QChatGPT")
}

func makeConfig(osname string) {
	println("生成配置文件")
	RunCMDPipe("生成配置文件", "./QChatGPT", "../python/bin/python3", "main.py")
	// RunCMDPipe("./QChatGPT", "../python/python", "main.py", "init_db")
	mirai_api_http_config := `adapters:
  - ws
debug: true
enableVerify: true
verifyKey: yirimirai
singleMode: false
cacheSize: 4096
adapterSettings:
  ws:
    host: localhost
    port: 8080
    reservedSyncId: -1`
	ioutil.WriteFile("./mirai/config/net.mamoe.mirai-api-http/setting.yml", []byte(mirai_api_http_config), 0644)

	println("=============================================")

	api_key := ""
	print("请输入OpenAI账号的api_key: ")
	fmt.Scanf("%s", &api_key)
	ReplaceStringInFile("./QChatGPT/config.py", "openai_api_key", api_key)

	qqn := 0
	print("请输入QQ号: ")
	if osname == "windows" {
		fmt.Scanf("%d", &qqn)
	}
	fmt.Scanf("%d", &qqn)
	ReplaceStringInFile("./QChatGPT/config.py", "1234567890", strconv.Itoa(qqn))
}

func writeLaunchScript(osname, arch string) {
	println("生成启动脚本")
	if osname == "windows" {
		ioutil.WriteFile("./run-mirai.bat", []byte(`cd mirai/
java\bin\java -jar mcl.jar`), 0644)
		ioutil.WriteFile("./run-bot.bat", []byte(`cd QChatGPT
..\python\python.exe main.py`), 0644)
	} else if osname == "linux" {
		ioutil.WriteFile("./run-mirai.sh", []byte(`cd mirai/
java/bin/java -jar mcl.jar`), 0644)
		RunCMDPipe("修改脚本权限", ".", "chmod", "+x", "run-mirai.sh")
		ioutil.WriteFile("./run-bot.sh", []byte(`cd QChatGPT
../python/bin/python3 main.py`), 0644)
		RunCMDPipe("修改脚本权限", ".", "chmod", "+x", "run-bot.sh")
	}
}
