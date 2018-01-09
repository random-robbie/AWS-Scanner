# AWS-Scanner


```
  /$$$$$$  /$$      /$$  /$$$$$$           /$$$$$$                                                             
 /$$__  $$| $$  /$ | $$ /$$__  $$         /$$__  $$                                                            
| $$  \ $$| $$ /$$$| $$| $$  \__/        | $$  \__/  /$$$$$$$  /$$$$$$  /$$$$$$$  /$$$$$$$   /$$$$$$   /$$$$$$ 
| $$$$$$$$| $$/$$ $$ $$|  $$$$$$  /$$$$$$|  $$$$$$  /$$_____/ |____  $$| $$__  $$| $$__  $$ /$$__  $$ /$$__  $$
| $$__  $$| $$$$_  $$$$ \____  $$|______/ \____  $$| $$        /$$$$$$$| $$  \ $$| $$  \ $$| $$$$$$$$| $$  \__/
| $$  | $$| $$$/ \  $$$ /$$  \ $$         /$$  \ $$| $$       /$$__  $$| $$  | $$| $$  | $$| $$_____/| $$      
| $$  | $$| $$/   \  $$|  $$$$$$/        |  $$$$$$/|  $$$$$$$|  $$$$$$$| $$  | $$| $$  | $$|  $$$$$$$| $$      
|__/  |__/|__/     \__/ \______/          \______/  \_______/ \_______/|__/  |__/|__/  |__/ \_______/|__/   
```





Scans a list of websites for Cloudfront or S3 Buckets.

This will output a text file with URL,CF or URL,s3bucketurl






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


Notes
-----

Input must be with out https:// as this is hardcoded at the moment.



To Do
-----

[  ] Fix bugs like the one on CNN.com where they dont close the dir and it grabs a load of messy js
[  ] Auto detect if prefix is needed or not
[  ] get rid of the mass regex functions and try do it as one.
