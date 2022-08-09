# paasta-monitoring-agent

## 1. 개요
- Bosh VM, PaaS-TA VM에 각각 배포되어 Metric 데이터를 수집하는 Agent 모듈.
- metrics_agent 라는 이름으로 monit 툴에 의해 관리 및 동작함.

## 2. 수집하는 항목
- CPU : Load Average 1M / 5M / 15M
- Memory : Used / Cached / Free / Buffer / Total Memory
- Disk : Read Bytes / Write Bytes / Read time / Write time
- Network : Byte received, sended / Drop in, out / Error in, out
- Process : Process name / Process id / Process's memory usage 

## 3. 수집하는 항목
- Metric 원천데이터는 gopsutil 라이브러리를 활용하여 수집
- 수집한 metric 데이터는 InfluxDB에 적재
- InfluxDB는 version 1.x 기준

## 4. 유의사항
- 차후 InfluxDB를 version 2.x로 업그레이드 하게되면 influx-go-client 라이브러리도 교체해야 함.
