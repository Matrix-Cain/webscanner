package utility

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/panjf2000/ants"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/transform"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var mutex = sync.Mutex{}
var UserAgents = []string{
	"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; AcooBrowser; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Acoo Browser; SLCC1; .NET CLR 2.0.50727; Media Center PC 5.0; .NET CLR 3.0.04506)",
	"Mozilla/4.0 (compatible; MSIE 7.0; AOL 9.5; AOLBuild 4337.35; Windows NT 5.1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
	"Mozilla/5.0 (Windows; U; MSIE 9.0; Windows NT 9.0; en-US)",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 2.0.50727; Media Center PC 6.0)",
	"Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 1.0.3705; .NET CLR 1.1.4322)",
	"Mozilla/4.0 (compatible; MSIE 7.0b; Windows NT 5.2; .NET CLR 1.1.4322; .NET CLR 2.0.50727; InfoPath.2; .NET CLR 3.0.04506.30)",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN) AppleWebKit/523.15 (KHTML, like Gecko, Safari/419.3) Arora/0.3 (Change: 287 c9dfb30)",
	"Mozilla/5.0 (X11; U; Linux; en-US) AppleWebKit/527+ (KHTML, like Gecko, Safari/419.3) Arora/0.6",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.2pre) Gecko/20070215 K-Ninja/2.1.1",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.9) Gecko/20080705 Firefox/3.0 Kapiko/3.0",
	"Mozilla/5.0 (X11; Linux i686; U;) Gecko/20070322 Kazehakase/0.4.5",
	"Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.8) Gecko Fedora/1.9.0.8-1.fc10 Kazehakase/0.5.6",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/535.20 (KHTML, like Gecko) Chrome/19.0.1036.7 Safari/535.20",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; fr) Presto/2.9.168 Version/11.52",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.11 TaoBrowser/2.0 Safari/536.11",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.71 Safari/537.1 LBBROWSER",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; LBBROWSER)",
	"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E; LBBROWSER)",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.84 Safari/535.11 LBBROWSER",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E)",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; QQBrowser/7.0.3698.400)",
	"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E)",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SV1; QQDownload 732; .NET4.0C; .NET4.0E; 360SE)",
	"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E)",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E)",
	"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.89 Safari/537.1",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.89 Safari/537.1",
	"Mozilla/5.0 (iPad; U; CPU OS 4_2_1 like Mac OS X; zh-cn) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8C148 Safari/6533.18.5",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:2.0b13pre) Gecko/20110307 Firefox/4.0b13pre",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:16.0) Gecko/20100101 Firefox/16.0",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11",
	"Mozilla/5.0 (X11; U; Linux x86_64; zh-CN; rv:1.9.2.10) Gecko/20100922 Ubuntu/10.10 (maverick) Firefox/3.6.10",
}
var Headers RequestHeaders
var result = make(map[string]map[string]string) // Protected by mutex since this program is network i/o bound
var once sync.Once
var config *Config

type Config struct {
	Headers      RequestHeaders
	FingerPrints []configTemplate
}
type RequestHeaders struct {
	userAgent []string
}
type ProcessedResponse struct {
	CertOption *tls.ConnectionState
	Header     http.Header
	Body       string
}
type PostProcessingResponse struct { // Since Banner and Title are static so no need to check every time
	ProcessedResponse
	Banners []string
	Titles  []string
}

func (response ProcessedResponse) proceedPostProcessing(resp ProcessedResponse) PostProcessingResponse {
	post := PostProcessingResponse{resp, GetBanners(resp), GetTitle(resp)}
	return post
}

func ReadConfig() {
	once.Do(func() {

		Headers = RequestHeaders{}
		Headers.userAgent = UserAgents
		config = &Config{Headers: Headers}
		OpenConfigFile("fofa.json")
	})
}

func IsValidUrl(inputString string) bool {
	_, err := url.ParseRequestURI(inputString)
	if err != nil {
		return false
	}

	u, err := url.Parse(inputString)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return false
	}

	return true
}

func GetTitle(resp ProcessedResponse) []string {
	r := regexp.MustCompile(`<title>(?P<title>(.*))</title>`)
	titles := r.FindStringSubmatch(resp.Body)
	return titles
}

func GetBanners(resp ProcessedResponse) []string {
	r := regexp.MustCompile(`<\s*banner.*>(?P<banner>(.*?))<\s*\s*banner>`)
	banners := r.FindStringSubmatch(resp.Body)
	return banners
}

func BodyContains(resp PostProcessingResponse, content string) bool {
	if strings.Contains(resp.Body, content) {
		return true
	} else {
		return false
	}

}

func ServerContains(resp PostProcessingResponse, content string) bool {
	for _, headerSlice := range resp.Header {
		for _, header := range headerSlice {
			if strings.Contains(strings.ToLower(header), strings.ToLower(content)) {
				return true
			}
		}
	}
	return false
}

func CertContains(resp PostProcessingResponse, content string) bool {
	if resp.CertOption != nil {
		certificates := resp.CertOption.PeerCertificates
		for _, certificate := range certificates {
			if bytes.Contains(certificate.Raw, []byte(content)) { //Not sure if certificate.Raw works
				return true
			}
		}
		return false
	}
	return false // omit the certificate match
}

func LoadSingleUrl(urls interface{}) {
	targetUrl := urls.(string)
	if !IsValidUrl(targetUrl) {
		log.Fatal("Invalid URL Input")
	}
	preResp, err := checkSingleUrl(targetUrl)
	if err != nil {
		log.Fatal(err)
	}
	resp := ProcessedResponse{}.proceedPostProcessing(preResp)
	rules := config.FingerPrints
	var data = map[string]string{}
	matchAND := false
	for _, rule := range rules {
		//if rule.RuleID == "392" {
		//	log.Println(rule)
		//}
		for _, check := range rule.Rules { //check for or relation
			for _, checkRule := range check { //check for and relation
				// omit certain option for we only check http fingerprints
				if strings.HasPrefix(checkRule.Match, "header") {
					matchAND = ServerContains(resp, strings.ToLower(checkRule.Content))
				}
				if strings.HasPrefix(checkRule.Match, "body") {
					matchAND = BodyContains(resp, strings.ToLower(checkRule.Content))
				}
				if strings.HasPrefix(checkRule.Match, "banner") {
					if resp.Banners == nil {
						matchAND = false
					} else {
						for _, banner := range resp.Banners {
							matchAND = strings.Contains(strings.ToLower(banner), strings.ToLower(checkRule.Content))
						}
					}
				}
				if strings.HasPrefix(checkRule.Match, "cert") {
					matchAND = CertContains(resp, checkRule.Content)
				}
				if strings.HasPrefix(checkRule.Match, "title") {
					if resp.Titles == nil {
						matchAND = false
					} else {
						for _, title := range resp.Titles {
							matchAND = strings.Contains(strings.ToLower(title), strings.ToLower(checkRule.Content))
						}
					}
				}
				if !matchAND {
					break
				}
			}
			if matchAND {

				//data = " ❤ " + url + " ❤ \n"
				//data += "Product: " + rule.Product + "\n"
				//data += "Company: " + rule.Company + "\n"

				//log.Println("Level: " + rule.RuleID)
				//log.Println("Softhard: " + rule.SoftHardID)
				data[rule.Product] = rule.Company

				log.Println(" ❤ " + targetUrl + " ❤ ")
				log.Println("Product: " + rule.Product)
				log.Println("Company: " + rule.Company)
				//log.Println("Category: " + rule.Category)
				//log.Println("Parent_Category: " + rule.ParentCategory)
				//log.Println()
				fmt.Println()
				matchAND = false
			}
		}
	}
	mutex.Lock()
	for product, company := range data {
		if result[targetUrl] == nil {
			result[targetUrl] = make(map[string]string)
		}
		result[targetUrl][product] = company
	}

	mutex.Unlock()
}

func LoadMultipleUrl(name string, poolSize int, save bool) {
	start := time.Now()
	var urlList []string
	file, err := os.Open(name)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	// read line by line
	for fileScanner.Scan() {
		if fileScanner.Text() != "" {

			urlList = append(urlList, strings.Trim(fileScanner.Text(), "\t"))
		}
	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	if len(urlList) == 0 {
		log.Fatal("Empty File Input")
	}
	routineWork(urlList, poolSize)
	if save {
		putResultsToFiles()
	}

	cost := time.Since(start)
	log.Printf("cost=[%s]", cost)
}

func routineWork(urlList []string, poolSize int) {
	defer ants.Release()
	var wg sync.WaitGroup
	runTimes := len(urlList)

	// Use the pool with a function,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	p, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
		LoadSingleUrl(i)
		wg.Done()
	})
	defer p.Release()
	// Submit tasks one by one.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(urlList[i])
	}

	wg.Wait()

	fmt.Printf("running goroutines: %d\n", p.Running())

}

func checkSingleUrl(url string) (ProcessedResponse, error) {
	rand.Seed(time.Now().Unix())
	ReadConfig() // Only do once
	client := http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", config.Headers.userAgent[rand.Intn(len(config.Headers.userAgent))])

	resp, err := client.Do(req)
	if err != nil {
		//log.Println(err)
		return ProcessedResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	return ProcessedResponse{CertOption: resp.TLS, Header: resp.Header, Body: decodeHTMLBody(resp.Body, "")}, nil
}

func putResultsToFiles() {

	file, err := os.OpenFile(generateFileName(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Unable to create save file handle", err)
	}
	//及时关闭file句柄
	defer file.Close()

	for tagretUrl, info := range result {
		for product, company := range info {
			file.WriteString(" ❤ " + tagretUrl + " ❤ \n")
			file.WriteString("Product: " + product + "\n")
			file.WriteString("Company: " + company + "\n")
			_, err := file.WriteString(" \n")
			if err != nil {
				log.Fatal("Unexpected error occurred while saving results")
			}
		}

	}
}

func detectContentCharset(r *bufio.Reader) string {

	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

// DecodeHTMLBody returns an decoding reader of the html Body for the specified `charset`
// If `charset` is empty, DecodeHTMLBody tries to guess the encoding from the content
func decodeHTMLBody(body io.Reader, charset string) string {
	r := bufio.NewReader(body)
	if charset == "" {
		charset = detectContentCharset(r)
	}
	e, err := htmlindex.Get(charset)
	if err != nil {
		log.Println(err)
		return ""
	}
	if name, _ := htmlindex.Name(e); name != "utf-8" {
		body = transform.NewReader(r, e.NewDecoder())
	}
	bodyContent, err := io.ReadAll(r)
	if err != nil {
		log.Println(err)
		return ""
	}
	return strings.ToLower(string(bodyContent))
}

func generateFileName() string {
	timeStamp := time.Now().Format("20060102150405")
	filename := timeStamp + ".txt"
	return filename

}
