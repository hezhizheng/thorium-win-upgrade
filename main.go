package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"strconv"
	"thorium-win-upgrade/language"
	"thorium-win-upgrade/service"
	"thorium-win-upgrade/service/helper"
)

// 1、爬取 github 的页面 获取最新版的chrome版本
//2、与本地 thorium 当前版本比较，大于当前版本则下载到本地、解压(询问提示)
//3、覆盖旧版数据，老版本BIN重命名为BIN2
//4、删除下载的文件

type Config struct {
	LocalChromePath string `json:"local_chrome_path"`
	Lang            string `json:"lang"`
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
		log.Fatal("The configuration file was loaded incorrectly, please check!", e)
	}
	json.Unmarshal(jsonStr, &_config)

	fmt.Printf(language.LanguageMap[_config.Lang]["welcome"] + "\n\n")
	fmt.Printf("github：https://github.com/hezhizheng/thorium-win-upgrade\n\n")
	fmt.Printf(language.LanguageMap[_config.Lang]["local_path"] + _config.LocalChromePath + "\n\n")
	fmt.Printf(language.LanguageMap[_config.Lang]["tips"] + "\n\n")
	fmt.Printf(language.LanguageMap[_config.Lang]["check_update"] + "......\n\n")

	// 获取本地chrome版本
	f := &service.FileInfo{
		FileDir: _config.LocalChromePath + "\\BIN\\",
	}
	localVersionName := service.GetLocalVersionName(f)

	if localVersionName == "" {
		fmt.Scanln(&exit)
		return
	}
	//获取 thorium 最新版本
	chromeFileName, latestVersionName := service.GetLatestVersionName()
	// 比较版本号
	if helper.CompareVersion(latestVersionName, localVersionName) == 1 {
		fmt.Printf(language.LanguageMap[_config.Lang]["local_thorium_version"] + localVersionName + "，" + language.LanguageMap[_config.Lang]["last_thorium_version"] + latestVersionName + " " + language.LanguageMap[_config.Lang]["is_update"] + " " + language.LanguageMap[_config.Lang]["y_or_n"] + "\n")
		fmt.Printf(language.LanguageMap[_config.Lang]["upgrade_tip"] + "！！！\n")
	} else {
		fmt.Printf(language.LanguageMap[_config.Lang]["local_thorium_version"] + localVersionName + "，" + language.LanguageMap[_config.Lang]["last_thorium_version"] + latestVersionName + " " + language.LanguageMap[_config.Lang]["no_need_update"] + "\n")
		fmt.Printf(language.LanguageMap[_config.Lang]["input_any_exit"] + "\n")
		fmt.Scanln(&exit)
		return
	}

	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		line := input.Text()
		fmt.Printf(language.LanguageMap[_config.Lang]["input"] + "：" + line + "\n")
		if line != "1" {
			break
		}

		// 关闭 thorium
		_, e1 := exec.Command("taskkill", "/F", "/IM", "thorium.exe").Output()
		if e1 != nil {
			//fmt.Println("exit thorium fail")
		}

		fmt.Printf(language.LanguageMap[_config.Lang]["select_version"] + "\n")
		//c := make(map[int]string)
		for _, m := range service.DownloadableVersion {
			for key, value := range m {
				fmt.Println(key, ":", value)
			}
		}

		var (
			selectVersion string
		)

		fmt.Scanln(&selectVersion)
		fmt.Printf(language.LanguageMap[_config.Lang]["input"] + "：" + selectVersion + "\n")
		for _, m := range service.DownloadableVersion {
			intSelectVersion, _ := strconv.Atoi(selectVersion)
			value, exists := m[intSelectVersion]
			if exists {
				chromeFileName = value
				break
			}
		}

		if chromeFileName == "" {
			fmt.Printf(language.LanguageMap[_config.Lang]["not_found_version"] + "\n")
			return
		}

		fmt.Printf(language.LanguageMap[_config.Lang]["update_ing"] + "\n")
		service.DownloadChrome(latestVersionName, localVersionName, chromeFileName)
		fmt.Printf(language.LanguageMap[_config.Lang]["update_success"] + " " + language.LanguageMap[_config.Lang]["y_or_n"] + "\n")

		var (
			isDelete string
		)
		fmt.Scanln(&isDelete)
		fmt.Printf(language.LanguageMap[_config.Lang]["input"] + "：" + isDelete + "\n")
		if isDelete != "1" {
			break
		}
		fmt.Printf(language.LanguageMap[_config.Lang]["file_delete_ing"] + "......\n")
		service.DeleteDownloadFile(latestVersionName)
		fmt.Printf(language.LanguageMap[_config.Lang]["file_deleted"] + "......\n")
		break
	}

	fmt.Printf(language.LanguageMap[_config.Lang]["input_any_exit"] + "\n")
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
