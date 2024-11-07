package service

import (
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const DownloadHost = "https://github.com"
const CheckoutVersionApi = "https://github.com/Alex313031/Thorium-Win/tags"
const CheckoutVersionApiNext = "https://github.com/Alex313031/Thorium-Win/releases/expanded_assets"

type FileInfo struct {
	FileDir string
	Version string
}

// 获取版本
func (f *FileInfo) GetLocalVersion() (err error) {
	rd, e := ioutil.ReadDir(f.FileDir)

	if e != nil {
		log.Println("目录读取失败", err, f.FileDir)
		return nil
	}

	// 第一个文件夹名字即版本号
	f.Version = rd[0].Name()

	return nil
}

func GetLocalVersionName(f *FileInfo) string {
	f.GetLocalVersion()
	return f.Version
}

func GetLatestVersionName() (string, string) {
	//return "fileName", "cVersion"
	fileName := ""

	c := colly.NewCollector(
		colly.Async(true),
	)

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("请求", r.URL, "...")
	})

	c.OnResponse(func(resp *colly.Response) {
		//log.Println(string(resp.Body))
	})

	retryCount := 0
	c.OnError(func(res *colly.Response, err error) {
		log.Println("Something went wrong:", err)
		if retryCount < 3 {
			retryCount += 1
			_retryErr := res.Request.Retry()
			log.Println("retry wrong:", _retryErr)
		}
	})

	c.OnXML(`//*[@id="repo-content-pjax-container"]/div/div[2]/div/div[1]/div[2]/div[1]/div/div/div[1]/div[1]/h2/a`, func(element *colly.XMLElement) {
		fileName = element.Attr("href")
	})

	proxyUrl := viper.GetString(`app.proxy_url`)

	if proxyUrl != "" {
		c.SetProxy(proxyUrl)
	}

	visitError := c.Visit(CheckoutVersionApi)

	if visitError != nil {
		log.Println("访问" + CheckoutVersionApi + "失败")
		panic(visitError)
	}
	c.Wait()

	if fileName != "" {
		FStrSplit := strings.Split(fileName, "/Alex313031/Thorium-Win/releases/tag/")[1]
		fileName = FStrSplit
	}

	return GetLatestVersionNameNext(fileName)
}

func GetLatestVersionNameNext(tag string) (string, string) {
	//return "fileName", "cVersion"
	fileName := ""

	c := colly.NewCollector(
		colly.Async(true),
	)

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("请求", r.URL, "...")
	})

	cVersion := ""

	c.OnResponse(func(resp *colly.Response) {
		//log.Println(string(resp.Body))
	})

	retryCount := 0
	c.OnError(func(res *colly.Response, err error) {
		log.Println("Something went wrong:", err)
		if retryCount < 3 {
			retryCount += 1
			_retryErr := res.Request.Retry()
			log.Println("retry wrong:", _retryErr)
		}
	})

	c.OnXML(`/html/body/div/ul/li[1]/div[1]/a`, func(element *colly.XMLElement) {
		fileName = element.Attr("href")
	})

	proxyUrl := viper.GetString(`app.proxy_url`)

	if proxyUrl != "" {
		c.SetProxy(proxyUrl)
	}

	vURL := CheckoutVersionApiNext + "/" + tag
	visitError := c.Visit(vURL)

	if visitError != nil {
		log.Println("访问" + vURL + "失败")
		panic(visitError)
	}
	c.Wait()

	cVersion = strings.Split(tag, "M")[1]
	return fileName, cVersion
}
