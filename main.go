
package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
)




var (
        tr     = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
        list = flag.String("list", "list.txt", "Default: list.txt")
        s3path = "s3-bucket.txt"
        cfpath = "cf-url.txt"
        awsurls = []string {
			"http://s3.amazonaws.com/([-A-z0-9.]+)",
			"https://s3.amazonaws.com/([-A-z0-9.]+)",
			"http://([-A-z0-9.]+).amazonaws.com",
			"https://([-A-z0-9.]+).amazonaws.com",
			"http://(s3.|s3-)[a-zA-Z0-9-]*.amazonaws.com/([-A-z0-9.]+)",
			"https://(s3.|s3-)[a-zA-Z0-9-]*.amazonaws.com/([-A-z0-9.]+)",

        }

        cloudfronturls = []string {
                "https://(.+?).cloudfront.net/",
				"http://(.+?).cloudfront.net/",
        }

)


func isError(err error) bool {
        if err != nil {
                fmt.Println(err.Error())
        }

        return err != nil
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

func runAWSCheck (responseString,url string) {
        for _, each := range awsurls {
                re, _ := regexp.Compile(each)
                cf := re.FindString(responseString)
        
                if cf != "" {
                        s3bucket := fmt.Sprintf("%s",cf)
                        color.HiGreen("[*] Logging... %s\n",s3bucket)
                        writes3bucket(cf,url)
                }
        }
}

//noinspection GoUnresolvedReference
func runCloudfrontCheck (responseString,url string) {
        for _, each := range cloudfronturls {
                re, _ := regexp.Compile(each)
                cf := re.FindString(responseString)
        
                if cf != "" {
                        s3bucket := fmt.Sprintf("%s",cf)
                        color.HiGreen("[*] Logging... %s\n",s3bucket)
                        writecf(cf,url)
                }
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
        color.HiGreen("[*]       Version 1.0               [*]\n\n\n\n")

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
                attempt := fmt.Sprintf("[*] Scanning Website https://%s\n", line)
                color.HiWhite(attempt)

                url := fmt.Sprintf("https://%s", line)
                tr := &http.Transport{
                        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
                        Dial: (&net.Dialer{
                                Timeout:   30 * time.Second,
                                KeepAlive: 30 * time.Second,
                        }).Dial,
                        TLSHandshakeTimeout:   10 * time.Second,
                        ResponseHeaderTimeout: 10 * time.Second,
                        ExpectContinueTimeout: 1 * time.Second,
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
                                color.HiGreen("[*] S3 bucket found!\n")
                                runAWSCheck(responseString, url)
                        }
                        //find Cloudfront distribution end point.
                        if strings.Contains(responseString, "cloudfront.net") == true {
                                color.HiGreen("[*] Cloudfront found!\n")
                                runCloudfrontCheck(responseString, url)

                        } else {
                                color.HiRed("[*] Moving to Next Site\n")
                        }
                }
        }

}
