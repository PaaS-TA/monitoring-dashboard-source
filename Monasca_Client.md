Monasca Client 설치 가이드
==========================

1. [개요](#1.)
    * [문서 목적](#1.1.)
    * [범위](#1.2.)
    * [확인사항](#1.3.)
    * [참고자료](#1.4.)
2. [Monasca Agent 설치 및 설정](#2.)
    * [dependencies 설치](#2.1.)
    * [monasca agent 설치](#2.2.)
    * [설치확인](#2.3.)
    * [monasca-setup 실행](#2.4.)
        * [Controller Node의 경우](#2.4.1.)
        * [Compute Node의 경우 (System 정보 수집과 VM 정보 수집 setup)](#2.4.2.)
    * [monasca agent configuration 파일 수정](#2.5.)
    * [Monasca Agent 시스템 정보 수집 설정 파일 수정](#2.6.)
        * [/etc/monasca/agent/conf.d/cpu.yaml](#2.6.1.)
        * [/etc/monasca/agent/conf.d/disk.yaml](#2.6.2.)
        * [/etc/monasca/agent/conf.d/libvirt.yaml (Compute Node에 한함)](#2.6.3.)
    * [monasca agent 서비스 재시작](#2.7.)
    * [확인](#2.8.)
    * [서비스 자동등록 되지 않을경우](#2.9.)
    * [Agent 사용자 mon-agent 사용자 자동 등록 되지 않는경우](#2.10.)
    * [Compute Node VM메트릭 정보가 수집 되지 않는경우](#2.11.)
3. [FileBeat 설치 및 설정](#3.)
    * [filebeat repository 등록](#3.1.)
    * [filebeat 설치](#3.2.)
    * [filebeat configuration 파일 수정](#3.3.)
    * [Elasticsearch-Logstash Certificate 파일을 Client 환경에 복사한다.](#3.4.)
    * [/etc/host 파일에 Elasticsearch Server 정보를 등록한다.](#3.5.)
    * [filebeat 서비스를 재가동한다.](#3.6.)
    * [확인](#3.7.)

    
# 1.	개요  <div id='1.'/>
## 1.1.	문서 목적   <div id='1.1.'/>
본 문서(설치가이드)는, IaaS(Infrastructure as a Service) 중 하나인 Openstack 기반의 Cloud 서비스 상태 및 자원 정보, 그리고 VM Instance의 시스템 정보와 로그정보를 수집하여, 실시간으로 서버로 정보를 전송하기 위한 Agent를 설치하는데 그 목적이 있다.

## 1.2.	범위   <div id='1.2.'/>
본 문서의 범위는 Openstack 모니터링을 위한 오픈소스인 Monasca 제품군의 설치를 위한 내용으로 한정되어 있다.

## 1.3.	확인사항   <div id='1.3.'/>
- Openstack 기반 환경 구성에 따라 Agent Setup 설정이 달라짐을 확인한다.
- Openstack newton 버전
- Node OS (Ubuntu 16.0.14)
- 크게 Controller Node .와 Compute Node로 구분된다.
- Controller Node에는 Nova, Neutron, Cinder, Glance, Keystone, Swift 서비스가 설치되는 환경이고, Compute Node에는 VM(Instance)이 생성되어 실행되는 환경을 의미한다.

## 1.4.	참고자료   <div id='1.4.'/>
- https://wiki.openstack.org/wiki/Monasca
- https://github.com/openstack/monasca-agent (version 2.7.0)
- https://www.elastic.co/kr/products/beats/filebeat

# 2.	Monasca Agent 설치 및 설정   <div id='2.'/>
## 2.1.	dependencies 설치   <div id='2.1.'/>
<pre>
    $ sudo apt-get install python-pip
</pre>    
    
## 2.2.	monasca agent 설치   <div id='2.2.'/>
<pre>
    $ sudo pip install monasca-agent==2.7.0
</pre>
    
## 2.3.	설치확인   <div id='2.3.'/>
<pre>
    $  sudo pip list |grep monasca-agent
</pre>
    
## 2.4.	monasca-setup 실행   <div id='2.4.'/>
### 2.4.1.	Controller Node의 경우   <div id='2.4.1.'/>
<pre>
    $ sudo monasca-setup \
      --username “cross-tenant user id” \
      --password “cross-tenant user password” \
      --project_name “admin project name” \
      --project_id “admin project id” \
      --user_domain_id “domain id” \
      --project_domain_id “domain id” \
      --keystone_url http://“keystone ip”:”keystone auth port”/v3 \
      --monasca_url http://”monasca ip”:”monasca server port”/v2.0 \
      --check_frequency '15'  \
      --log_level 'DEBUG'  \
    --insecure true \
    --system_only 
</pre>
    
### 2.4.2.	Compute Node의 경우 (System 정보 수집과 VM 정보 수집 setup)   <div id='2.4.2.'/>
<pre>
    $ sudo monasca-setup \
      --username “cross-tenant user id” \
      --password “cross-tenant user password” \
      --project_name “admin project name” \
      --project_id “admin project id” \
      --user_domain_id “domain id” \
      --project_domain_id “domain id” \
      --keystone_url http://“keystone ip”:”keystone auth port”/v3 \
      --monasca_url http://”monasca ip”:”monasca server port”/v2.0 \
      --check_frequency '15'  \
      --log_level 'DEBUG'  \
    --insecure true \
    --system_only
    
    $ sudo monasca-setup -d libvirt -a 'ping_check=false alive_only=false'
</pre>

## 2.5.	monasca agent configuration 파일 수정.   <div id='2.5.'/>
<pre>
    $ sudo vi /etc/monasca/agent/agent.yml
    Api:
      amplifier: 0
      backlog_send_rate: 1000
      ca_file: null
      endpoint_type: null
      insecure: true
      keystone_url: http://”keystone ip” :”keystone auth port”/v3
      max_buffer_size: 1000
      max_measurement_buffer_size: -1
      password: cfmonit
      project_domain_id: default
      project_domain_name: default
      project_id: “admin project id”
      project_name: “admin project name”
      region_name: null
      service_type: null
      url: http:// “monasca server ip”:”monasca server port”/v2.0
      user_domain_id: “domain id”
      user_domain_name: default
      username: admin
    Logging:
      collector_log_file: /var/log/monasca/agent/collector.log
      enable_logrotate: true
      forwarder_log_file: /var/log/monasca/agent/forwarder.log
      log_level: DEBUG                                           # Log 레벨 설정
      statsd_log_file: /var/log/monasca/agent/statsd.log
    Main:
      check_freq: 15                                              # 수집 주기(초)
      collector_restart_interval: 24
      dimensions: {}
      hostname: controller
      num_collector_threads: 1
      pool_full_max_retries: 4
      sub_collection_warn: 6
    Statsd:
      monasca_statsd_port: 8125
</pre>

## 2.6.	Monasca Agent 시스템 정보 수집 설정 파일 수정   <div id='2.6.'/>
### 2.6.1.	/etc/monasca/agent/conf.d/cpu.yaml   <div id='2.6.1.'/>
<pre>
    init_config: null
    instances:
    - built_by: System
      name: cpu_stats
      send_rollup_stats: True    # vcpu measurement Option 추가
</pre>
      
### 2.6.2.	/etc/monasca/agent/conf.d/disk.yaml   <div id='2.6.2.'/>
<pre>
    init_config: null
    instances:
    - built_by: System
      device_blacklist_re: .*freezer_backup_snap.*
      ignore_filesystem_types: iso9660,tmpfs
      name: disk_stats
      send_rollup_stats: True    # Node disk 사용량  Option 추가
</pre>

### 2.6.3.	/etc/monasca/agent/conf.d/libvirt.yaml (Compute Node에 한함)   <div id='2.6.3.'/>
<pre>
    init_config:
      alive_only: false
      auth_url: http://controller:35357
      cache_dir: /dev/shm
      customer_metadata:
      - scale_group
      disk_collection_period: 0
      max_ping_concurrency: 8
      metadata:
      - scale_group
      nova_refresh: 14400
      password: cfmonit
      ping_check: false
      project_name: admin
      username: admin
      vm_cpu_check_enable: true
      vm_disks_check_enable: true
      vm_extended_disks_check_enable: true   # vm disk 사용량 추가
      vm_network_check_enable: true
      vm_ping_check_enable: true
      vm_probation: 300
      vnic_collection_period: 0
    instances: []
</pre>
    
## 2.7.	monasca agent 서비스 재시작.   <div id='2.7.'/>
<pre>
    $ sudo service monasca-agent restart
</pre>
 
- 서비스 등록이 되지 않을경우<br>
/etc/systemd/system/monasca-agent.service
<pre>
    [Unit]
    Description=Monasca Agent
    [Service]
    Type=simple
    User=mon-agent
    Group=mon-agent
    Restart=on-failure
    ExecStart=/usr/local/bin/supervisord -c /etc/monasca/agent/supervisor.conf -n
        
    [Install]
    WantedBy=multi-user.target
</pre>
<pre>    
    $cd /etc/systemd/system/multi-user.target.wants
    sudo ln –s /etc/systemd/system/monasca-agent.service /etc/systemd/system/monasca-agent.service
</pre>        

- cf-mon os user 자동 등록되지 않을경우 사용자 수동 등록
<pre>
    $ sudo useradd mon-agent
</pre>
    
## 2.8.	확인   <div id='2.8.'/>
![](images/Monasca/2.8.png)

## 2.9. 서비스 자동등록 되지 않을경우    <div id='2.9.'/>
/etc/systemd/system/monasca-agent.service
<pre>    
    [Unit]
    Description=Monasca Agent
    [Service]
    Type=simple
    User=mon-agent
    Group=mon-agent
    Restart=on-failure
    ExecStart=/usr/local/bin/supervisord -c /etc/monasca/agent/supervisor.conf -n
    
    [Install]
    WantedBy=multi-user.target
</pre>
<pre>    
    $cd /etc/systemd/system/multi-user.target.wants
    sudo ln –s /etc/systemd/system/monasca-agent.service /etc/systemd/system/monasca-agent.service
</pre>

## 2.10. Agent 사용자 mon-agent 사용자 자동 등록 되지 않는경우    <div id='2.10.'/>
<pre>
    $ sudo useradd mon-agent
</pre>
    
## 2.11.  Compute Node VM메트릭 정보가 수집 되지 않는경우    <div id='2.11.'/>
<pre>
    $ cd /
    $ sudo chmod 757 /run
</pre>
    
# 3.	FileBeat 설치 및 설정   <div id='3.'/>
## 3.1.	filebeat repository 등록   <div id='3.1.'/>
<pre>
    $ echo "deb https://artifacts.elastic.co/packages/5.x/apt stable main" | sudo tee -a     /etc/apt/sources.list.d/elastic-5.x.list
    $ sudo apt-get update
</pre>
    
## 3.2.	filebeat 설치   <div id='3.2.'/>
<pre>
    $ sudo apt-get install -y filebeat
</pre>

## 3.3.	filebeat configuration 파일 수정   <div id='3.3.'/>
<pre>
    $ sudo vi /etc/filebeat/filebeat.yml
    ---
    ...
    # Add log files 
    paths:
        - /var/log/cinder/*.log                # 수집하고자 하는 로그파일을 지정한다.
        - /var/log/glance/*.log               # 여러개를 지정할 수 있다.
        - /var/log/neutron/*.log
    # Add document type
    document-type: syslog
    ...
    #Disable Elasticsearch output
    #output.elasticsearch:
     # Array of hosts to connect to.
     #hosts: ["localhost:9200"] 
    ...
    # output.logstash:
      # The Logstash hosts
      hosts: ["elasticsearch server ip:5443"]                   #elasticsearch server ip address
      bulk_max_size: 2048
      ssl.certificate_authorities: ["/etc/filebeat/logstash.crt"]    #logstash certificate file location
      template.name: "filebeat"
      template.path: "filebeat.template.json"
      template.overwrite: false
</pre>
      
## 3.4.	Elasticsearch-Logstash Certificate 파일을 Client 환경에 복사한다.   <div id='3.4.'/>
<pre>
    $ sudo scp ubuntu@”elasticsearch server ip”:/etc/logstash/logstash.crt /etc/filebeat/
</pre>
    
## 3.5.	/etc/host 파일에 Elasticsearch Server 정보를 등록한다.   <div id='3.5.'/>
<pre>
    $ sudo vi /etc/hosts
    ---
    “elasticsearch server ip”    “hostname”
    
    ex) 10.10.10.10  elasticsearch-server
</pre>
    
## 3.6.	filebeat 서비스를 재가동한다.   <div id='3.6.'/>
<pre>
    $ sudo service filebeat restart
</pre>
    
## 3.7.	확인.   <div id='3.7.'/>
<pre>
    $ ps -ef |grep filebeat
</pre>    
![](images/Monasca/3.7.png)