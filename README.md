# AWS-Scanner

Scans a list of websites for Cloudfront or S3 Bucketsbut it will only check the first page and pull the first one.

it will not go through every page and list every one.

you can how ever feed it a list of urls to check.

note input must be with out https:// as this is prefixed 




THIS IS A WORK IN PROGRESS!


Install
------

```
go get github.com/fatih/color
go build main.go
./main --list list.txt
```

Release's
-----

Pre-built binarys are avaliable on the [releases](https://github.com/random-robbie/AWS-Scanner/releases/download/v0.1/Releases-Beta.zip) tab.





Screenshot
------

[![Capture.png](https://s9.postimg.org/a0a819pnj/Capture.png)](https://postimg.org/image/y40zpk84b/)

To Do
-----

Fix bugs like the one on CNN.com where they dont close the dir and it grabs a load of messy js
Auto detect if prefix is needed or not
