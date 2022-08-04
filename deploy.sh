GOARCH=amd64 GOOS=linux go build -o main
zip main.zip main
rm main
aws lambda update-function-code  --function-name go-crawler --zip-file fileb://main.zip
rm main.zip