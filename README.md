# Official Taucoder Go Client

Taucoder Go CLI Client

## Quick go install &amp; run

```sh
go install github.com/Taucoder-com/taucoder-go-client@latest
taucoder-go-client -apikey "YOUR-API-KEY" -output output-directory -quality 50 image1.png image2.jpg
```

## Install &amp; run from source code

```sh
git clone https://github.com/Taucoder-com/taucoder-go-client.git
cd taucoder-go-client
go run main.go -apikey "YOUR-API-KEY" -output output-directory -quality 50 image1.png image2.jpg
```

## Command line options

- `-apikey` string. API key for authentication (required)
- `-output` string. Output directory where the optimized images will be saved (required)
- `-quality` int. Quality of the output image (optional, default: 50, range: 25-95)
