package main

import (
	"net/http"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"bufio"
	"bytes"
	"io"
	"strings"
	"io/ioutil"
	"regexp"
	"github.com/fatih/color"

)
var (
	tr     = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	list = flag.String("list", "list.txt", "Default: list.txt")
	s3path = "s3-bucket.txt"
	cfpath = "cf-url.txt"

)


func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}



// Read a whole file into the memory and store it as array of lines
func readLines(path string) (lines []string, err error) {
	var (
		file *os.File
		part []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}





// write file of s3bucket
func writes3bucket(s3bucket,url string){

	//check if file exists
	if _, err := os.Stat("s3-bucket.txt"); os.IsNotExist(err) {
		// create file if not exists
		if os.IsNotExist(err) {
			var file, err = os.Create("s3-bucket.txt")
			if isError(err) { return }
			defer file.Close()
		}
	}


	fileHandle, _ := os.OpenFile(s3path, os.O_APPEND, 0666)
	writer := bufio.NewWriter(fileHandle)
	defer fileHandle.Close()
	stringtolog := fmt.Sprintf("%s,%s",url,s3bucket)
	fmt.Fprintln(writer,stringtolog)
	writer.Flush()
}

// write file of s3bucket
func writecf(cf,url string){

	//check if file exists
	if _, err := os.Stat("cf-dist.txt"); os.IsNotExist(err) {
		// create file if not exists
		if os.IsNotExist(err) {
			var file, err = os.Create("cf-dist.txt")
			if isError(err) { return }
			defer file.Close()
		}
	}


	fileHandle, _ := os.OpenFile(s3path, os.O_APPEND, 0666)
	writer := bufio.NewWriter(fileHandle)
	defer fileHandle.Close()
	stringtolog := fmt.Sprintf("%s,%s",url,cf)
	fmt.Fprintln(writer,stringtolog)
	writer.Flush()
}

func cfregex (responseString,url string) {
	re, _ := regexp.Compile(`https://(.+?).cloudfront.net/`)
	cf := re.FindString(responseString)

	if cf != "" {
		s3bucket := fmt.Sprintf("%s",cf)
		color.HiGreen("[*] Logging...Cloudfront Endpoint %s[*]\n",s3bucket)
		writecf(cf,url)
	}


}


func s3regex (responseString,url string) {
	re, _ := regexp.Compile(`https://s3.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}


}


func s3regex2 (responseString,url string) {
	re, _ := regexp.Compile(`http://s3.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}


func s3regex3 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-us-east-2.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}


func s3regex4 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-us-west-1.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}

func s3regex5 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-us-west-2.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}

func s3regex6 (responseString,url string) {
	re, _ := regexp.Compile(`//s3.ca-central-1.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}

func s3regex7 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-ap-south-1.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}

func s3regex8 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-ap-northeast-2.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}

func s3regex9 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-ap-southeast-1.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}

func s3regex10 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-ap-northeast-1.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}

func s3regex11 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-eu-central-1.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}
func s3regex12 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-eu-west-1.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}
func s3regex13 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-eu-west-2.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}
func s3regex14 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-eu-west-3.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}

func s3regex15 (responseString,url string) {
	re, _ := regexp.Compile(`//s3.sa-east-1.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}


func s3regex16 (responseString,url string) {
	re, _ := regexp.Compile(`https://(.+?).s3.amazonaws.com`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}

func s3regex17 (responseString,url string) {
	re, _ := regexp.Compile(`//s3.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}

func s3regex18 (responseString,url string) {
	re, _ := regexp.Compile(`//s3-ap-southeast-2.amazonaws.com/(.+?)/`)
	s3 := re.FindString(responseString)

	if s3 != "" {
		s3bucket := fmt.Sprintf("https:%s",s3)
		color.HiGreen("[*] Logging... %s[*]\n",s3bucket)
		writes3bucket(s3bucket,url)
	}

}





func main () {
	// Parse the flags provided
	flag.Parse()
	color.HiGreen("[*]                                 [*]\n")
	color.HiGreen("[*] AWS Scanner - By @random_robbie [*]\n")
	color.HiGreen("[*]                                 [*]\n")
	color.HiGreen("[*]                                 [*]\n")
	color.HiGreen("[*]    S3 Buckets for the Win       [*]\n")
	color.HiGreen("[*]                                 [*]\n")
	color.HiGreen("[*]                                 [*]\n\n\n\n\n\n\n")

	//Check the server is not blank or empty
	if *list == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	//read list line by line
	lines, err := readLines(*list)
	if err != nil {
		fmt.Println("Error: %s\n", err)
		return
	}
	for _, line := range lines {

		//Print Scanning Website
		attempt := fmt.Sprintf("[*] Scanning Website https://%s [*]\n", line)
		color.HiWhite(attempt)

		url := fmt.Sprintf("https://%s", line)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		// Do not verify certificates
		client := &http.Client{
			Transport: tr,
		}

		// Create HEAD request with random user agent.
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0")

		if resp, err := client.Do(req); err == nil {

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			responseString := string(body)
			if err != nil {
				panic(err)
			}
			//If amazonaws.com is found signal bucket found
			if strings.Contains(responseString, "amazonaws.com") == true {
				color.HiGreen("[*] S3 bucket found! [*]\n")
				s3regex(responseString, url)
				s3regex2(responseString, url)
				s3regex3(responseString, url)
				s3regex4(responseString, url)
				s3regex5(responseString, url)
				s3regex6(responseString, url)
				s3regex7(responseString, url)
				s3regex8(responseString, url)
				s3regex9(responseString, url)
				s3regex10(responseString, url)
				s3regex11(responseString, url)
				s3regex12(responseString, url)
				s3regex13(responseString, url)
				s3regex14(responseString, url)
				s3regex15(responseString, url)
				s3regex16(responseString, url)
				s3regex17(responseString, url)
				s3regex18(responseString, url)

			}
			//find Cloudfront distribution end point.
			if strings.Contains(responseString, "cloudfront.net") == true {
				color.HiGreen("[*] Cloudfront found! [*]\n")
				cfregex(responseString, url)

			} else {
				color.HiRed("[*] Moving to Next Site [*]\n")
			}
		}
	}

}


