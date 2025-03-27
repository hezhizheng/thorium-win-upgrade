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

var DownloadableVersion []map[int]string

// 获取版本
func (f *FileInfo) GetLocalVersion() (err error) {
	rd, e := ioutil.ReadDir(f.FileDir)

	if e != nil {
		log.Println("Load Dir Fail", err, f.FileDir)
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
		log.Println("Request", r.URL, "...")
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
		log.Println("Request" + CheckoutVersionApi + "Fail")
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
		log.Println("Request", r.URL, "...")
	})

	cVersion := ""

	c.OnResponse(func(resp *colly.Response) {
		//log.Println(string(resp.Body))
	})

	versionType := viper.Get(`app.version_type`)
	aIndex := 1
	log.Println("Downloadable version...")
	c.OnHTML("a", func(e *colly.HTMLElement) {
		// 获取 <a> 标签的 href 属性
		href := e.Attr("href")
		if strings.Contains(href, ".zip") && !strings.Contains(href, "tag") && !strings.Contains(href, "templates") {
			log.Println(href)
			myMap := map[int]string{
				aIndex: href,
			}
			DownloadableVersion = append(DownloadableVersion, myMap)
			aIndex++
		}
		if strings.Contains(href, versionType.(string)) && strings.Contains(href, ".zip") {
			fileName = href
		}
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

	proxyUrl := viper.GetString(`app.proxy_url`)

	if proxyUrl != "" {
		c.SetProxy(proxyUrl)
	}

	vURL := CheckoutVersionApiNext + "/" + tag
	visitError := c.Visit(vURL)

	if visitError != nil {
		log.Println("Request" + vURL + "Fail")
		panic(visitError)
	}
	c.Wait()

	if fileName == "" {
		// 没有找到对应版本，请确认 version_type 是否设置正确
		panic("No corresponding version found, please confirm whether version_type is set correctly.")
	}

	cVersion = strings.Split(tag, "M")[1]
	return fileName, cVersion
}
