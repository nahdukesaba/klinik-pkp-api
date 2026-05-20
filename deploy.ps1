cd E:\Server\Projects\klinik-pkp-api

git pull

go mod tidy

go build -o app.exe

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed"
    exit
}

net stop backend-service

Start-Sleep -Seconds 2

net start backend-service

Write-Host "Deployment complete"