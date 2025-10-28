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






Install
------

**Requirements:** Go 1.16 or later

```bash
# Clone the repository
git clone https://github.com/random-robbie/AWS-Scanner.git
cd AWS-Scanner

# Download dependencies (automatically handled by Go modules)
go mod tidy

# Build the binary
go build -o aws-scanner main.go

# Run the scanner
./aws-scanner --list list.txt
```

Or build and run directly:
```bash
go run main.go --list list.txt
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


Thanks
----
[Glove](https://github.com/Glove)

Use a VPS from DO

[![DigitalOcean Referral Badge](https://web-platforms.sfo2.cdn.digitaloceanspaces.com/WWW/Badge%201.svg)](https://www.digitalocean.com/?refcode=e22bbff5f6f1&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)
