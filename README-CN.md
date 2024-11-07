## thorium-win-upgrade

> 一个可以升级 [Thorium Browser ](https://thorium.rocks/)  的工具

> [github](https://github.com/hezhizheng/thorium-win-upgrade)

### 功能
- 简单交互式操作
- 自动检测最新的 `Thorium` 版本
- 用户决定是否进行升级操作(自动下载、解压、重命名文件等)
- windows下 配合.bat文件 实现开机自动检测更新功能

### 流程
![free-pic](./images/1.png)


### 使用
自定义config.json配置文件(thorium 的安装目录)

例：假如我的 thorium 安装解压目录为

![free-pic](./images/2.png)

那么 local_chrome_path 就定义为 `D:\test1107\thorium`。如下：
```
# 参数说明
{
  "app": {
    "local_chrome_path": "D:\\test1107\\thorium"
    ,"proxy_url": ""
  }
}
- proxy_url： 下载代理，检查最新版本跟下载可能需要翻墙，如程序有错误抛出，请尝试使用代理(http://127.0.0.1:7890)解决(不使用代理则无需配置该项，或者设置为空字符串)。
- local_chrome_path：本地 thorium 的安装路径
```

编译 (提供编译好的文件 thorium-win-upgrade.7z
下载 [releases](https://github.com/hezhizheng/thorium-win-upgrade/releases) )

手动编译
```
go build -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"
```

运行
- 请不要随意更改`thorium`原本的目录结构
- 保证编译的文件与 config.json、7z.dll、7z.exe 文件 在同级目录
- 执行 ./thorium-win-upgrade.exe 或者双击启动，根据提示输入指令完成升级

升级

![free-pic](./images/4.png)
![free-pic](./images/6.png)


无需升级

![free-pic](./images/3.png)

windows 开机自动检测(创建.bat文件)

[./thorium-win-upgrade.bat](./thorium-win-upgrade.bat)

创建快捷方式，设定开机自启即可