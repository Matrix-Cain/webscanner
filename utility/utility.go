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
var Headers RequestHeaders
var result = make(map[string]map[string]string) // Protected by mutex since this program is network i/o bound
var once sync.Once
var config *Config

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

func LoadUrl(isFile bool, name string, poolSize int, save bool, fileName string) {
	start := time.Now()
	var urlList []string
	if isFile {
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
	} else {
		urlList = append(urlList, name)
	}
	if len(urlList) == 0 {
		log.Fatal("Empty File Input")
	}
	routineWork(urlList, poolSize)
	if save {
		putResultsToFiles(fileName)
	}

	cost := time.Since(start)
	log.Printf("cost=[%s]", cost)
}

func ScanUrl(urls interface{}) {
	targetUrl := urls.(string)
	if !IsValidUrl(targetUrl) {
		log.Fatal("Invalid URL Input")
	}
	preResp, err := checkSingleUrl(targetUrl)
	if err != nil { // May be Network or target issue
		log.Println(err) // So basically this should not be fatal
		return
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
			if bytes.Contains(certificate.Raw, []byte(content)) { //Not sure if you certificate.Raw works
				return true
			}
		}
		return false
	}
	return false // omit the certificate match
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

func routineWork(urlList []string, poolSize int) {
	defer ants.Release()
	var wg sync.WaitGroup
	runTimes := len(urlList)

	// Use the pool with a function,
	// set 10 to the capacity of goroutine pool and 1 second for expired duration.
	p, _ := ants.NewPoolWithFunc(poolSize, func(i interface{}) {
		ScanUrl(i)
		wg.Done()
	})
	defer p.Release()
	// Submit tasks one by one.
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(urlList[i])
	}

	wg.Wait()

	//fmt.Printf("running goroutines: %d\n", p.Running())

}

func putResultsToFiles(fileName string) {
	var file *os.File
	var err error

	if len(result) == 0 { // Nothing to save avoid making empty output file
		return
	}

	if fileName == "" {
		file, err = os.OpenFile(generateFileName(), os.O_WRONLY|os.O_CREATE, 0666)
	} else {
		file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	}

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
