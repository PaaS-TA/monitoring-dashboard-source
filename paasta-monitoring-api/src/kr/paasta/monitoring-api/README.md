
Environment 

1. IDE : Goland 
2. GOROOT : Go 1.17.8 
3. GOPATH : Use GOPATH that’s defined in system environment” 체크
4. GO Modules : “Enable Go modules integration”을 체크
5. ~\paasta-monitoring-api\src\kr\paasta\monitoring-api\에서 
    go get -u 실행

Run / Debug Configuration

1. Name : paasta-monitoring-api
2. Runkind : File 
3. Files : ~\paasta-monitoring-api\src\kr\paasta\monitoring-api\main.go
4. Output directory : ~\paasta-monitoring-api\src\kr\paasta\monitoring-api\
5. Working directory : ~\paasta-monitoring-api\src\kr\paasta\monitoring-api\
6. Go tool argument : -i