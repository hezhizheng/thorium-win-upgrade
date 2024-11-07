package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"thorium-win-upgrade/service"
	"thorium-win-upgrade/service/helper"
)

// 1、爬取 github 的页面 获取最新版的chrome版本
//2、与本地 thorium 当前版本比较，大于当前版本则下载到本地、解压(询问提示)
//3、覆盖旧版数据，老版本BIN重命名为BIN2
//4、删除下载的文件

type Config struct {
	LocalChromePath string `json:"local_chrome_path"`
}

func init() {
	initConfig()
}

func main() {

	var _config Config
	var (
		exit string
	)

	configStr := viper.Get(`app`)
	jsonStr, e := json.Marshal(configStr)
	if e != nil {
		log.Fatal("配置文件加载有误，请检查！", e)
	}
	json.Unmarshal(jsonStr, &_config)

	fmt.Printf("欢迎使用 thorium-win-upgrade 工具\n\n")
	fmt.Printf("github：https://github.com/hezhizheng/thorium-win-upgrade\n\n")
	fmt.Printf("当前定义的本地chrome的安装路径为：" + _config.LocalChromePath + "\n\n")
	fmt.Printf("请根据提示输入相关指令进行操作\n\n")
	fmt.Printf("检查更新中......\n\n")

	// 获取本地chrome版本
	f := &service.FileInfo{
		FileDir: _config.LocalChromePath + "\\BIN\\",
	}
	localVersionName := service.GetLocalVersionName(f)

	//获取 thorium 最新版本
	chromeFileName, latestVersionName := service.GetLatestVersionName()
	// 比较版本号
	if helper.CompareVersion(latestVersionName, localVersionName) == 1 {
		fmt.Printf("当前本地 Thorium 的版本为：" + localVersionName + "，" + "最新 Thorium 版本为：" + latestVersionName + " 是否进行升级？1：是 2：否\n")
		fmt.Printf("提示：升级前请确保浏览器已处于退出状态！！！\n")
	} else {
		fmt.Printf("当前本地 Thorium 的版本为：" + localVersionName + "，" + "最新 Thorium 版本为：" + latestVersionName + " 无需升级\n")
		fmt.Printf("输入任意键退出\n")
		fmt.Scanln(&exit)
		return
	}

	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		line := input.Text()
		fmt.Printf("输入了：" + line + "\n")
		if line != "1" {
			break
		}

		// 关闭 thorium
		_, e1 := exec.Command("taskkill", "/F", "/IM", "thorium.exe").Output()
		if e1 != nil {
			fmt.Println("关闭 thorium 失败")
			panic(e1)
		}

		fmt.Printf("升级中，请等待，此过程中请不要做任何输入。\n")
		service.DownloadChrome(latestVersionName, localVersionName, chromeFileName)
		fmt.Printf("升级成功，是否删除下载/解压的文件？（建议先检查是否升级成功在执行此操作！！！）1：是 2：否\n")

		var (
			isDelete string
		)
		fmt.Scanln(&isDelete)
		fmt.Printf("输入了：" + isDelete + "\n")
		if isDelete != "1" {
			break
		}
		fmt.Printf("文件删除中......\n")
		service.DeleteDownloadFile(latestVersionName)
		fmt.Printf("删除完成......\n")
		break
	}

	fmt.Printf("输入任意键退出\n")
	fmt.Scanln(&exit)
	return
}

func initConfig() {
	viper.SetConfigType("json") // 设置配置文件的类型
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			log.Println("read config error")
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}
}
