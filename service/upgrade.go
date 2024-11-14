package service

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"path"
	"thorium-win-upgrade/language"
	"thorium-win-upgrade/service/helper"
)

func DownloadChrome(latestVersionName, localVersionName, chromeFileName string) {

	lang := viper.GetString(`app.lang`)
	needDownload := false

	if latestVersionName != "" && localVersionName != "" {
		if helper.CompareVersion(latestVersionName, localVersionName) == 1 {
			needDownload = true
		}
	} else {
		panic(language.LanguageMap[lang]["version_err"])
	}

	url := DownloadHost + chromeFileName
	path := viper.Get(`app.local_chrome_path`)

	filename := path.(string) + "\\" + latestVersionName + ".zip"

	if needDownload && !fileExists(filename) {
		fmt.Println(language.LanguageMap[lang]["download_start"] + " " + url + " ........")
		err := DownloadFile(filename, url)

		if err != nil {
			fmt.Println(language.LanguageMap[lang]["download_file"] + url + language.LanguageMap[lang]["fail"])
			panic(err)
		}

		fmt.Println(language.LanguageMap[lang]["download_done"])
	}

	if fileExists(filename) {

		// 先删除旧版本升级遗留的文件夹
		os.RemoveAll(path.(string) + "\\" + "BIN2")
		os.RemoveAll(path.(string) + "\\" + "thorium_tmp")

		fmt.Println(language.LanguageMap[lang]["unzip_file"] + "........")

		//_, e1 := exec.Command("./7z.exe", "x", filename, "-o"+path.(string)+"\\"+"thorium_tmp").Output()
		//if e1 != nil {
		//	fmt.Println(language.LanguageMap[lang]["unzip_file_fail"])
		//	panic(e1)
		//}

		unzipErr := helper.Unzip(filename, path.(string)+"\\"+"thorium_tmp")
		if unzipErr != nil {
			fmt.Println(language.LanguageMap[lang]["unzip_file_fail"])
			panic(unzipErr)
		}

		fmt.Println(language.LanguageMap[lang]["unzip_done"])

		renameErr := os.Rename(path.(string)+"\\"+"BIN", path.(string)+"\\"+"BIN2")

		if renameErr != nil {
			fmt.Println(language.LanguageMap[lang]["rename_file_name_fail"])
			panic(renameErr)
		}

		e2 := copyDir(path.(string)+"\\"+"thorium_tmp\\BIN", path.(string)+"\\"+"BIN")

		if e2 != nil {
			fmt.Println(language.LanguageMap[lang]["copy_dir_fail"])
			panic(e2)
		}

		return
	}

	panic(language.LanguageMap[lang]["update_fail"])

}

func DeleteDownloadFile(latestVersionName string) {
	localChromePath := viper.Get(`app.local_chrome_path`).(string)
	// os.RemoveAll(localChromePath + "\\" + "App2")
	os.RemoveAll(localChromePath + "\\" + "thorium_tmp")
	os.RemoveAll(localChromePath + "\\" + latestVersionName + ".zip")
}

func fileForCopyDir(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func copyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = fileForCopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
