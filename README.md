# Installer for QChatGPT

为[QChatGPT项目](https://github.com/RockChinQ/QChatGPT)使用Go语言编写的一键部署脚本，自动化部署所需依赖  

- 注意：下载的Python和mirai均为免安装版，不影响系统其他环境

## 使用方法

- 提前准备好需要使用的QQ号
- 注册并获取OpenAI账号，参考以下文章，注册完成之后到账户设置获取`api-key`
    - [只需 1 元搞定 ChatGPT 注册](https://zhuanlan.zhihu.com/p/589470082)
    - [手把手教你如何注册ChatGPT，超级详细](https://guxiaobei.com/51461)

从Release中下载发行可执行文件，直接运行，等待环境配置完毕后根据提示输入`api-key`和`QQ号`  
运行完毕后即可运行`run-mirai.bat`启动mirai并输入`login <QQ号> <QQ密码>`根据提示登录账号，登录完成后运行`run-bot.bat`启动机器人  

- 若下载速度较慢需要使用网络代理，请在启动安装器时提供参数`-p <代理地址>`,如：
```
qcg-installer-0.1-windows-x64.exe -p http://localhost:7890
```

## 目前支持的平台和架构

- Windows x64
