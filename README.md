PaaS_TA_Monitoring-v2.0
=======================


1. [개요](#1)
    * [문서 개요](#2)  
        * [목적](#3)
        * [범위](#4)
        * [참고자료](#5)

2. [IaaS Monitoring 애플리케이션 환결설정](#6)
    * [개요](#7)
	* [개발환경 구성](#8)
	    * [Git 설치](#9)
        * [IntelliJ IDEA설치](#10)
        * [Go Lang설치](#11)
        * [Go 애플리케이션 환경설정](#12)
    * [소스 다운로드](#13)
    * [애플리케이션 프로젝트 Open](#14)
        * [PaaS-TA-Monitoring Management 프로젝트 Open](#15)
    
  
    
<br /><br /><br />

#   1. 개요 <div id='1' />

##  1.1. 문서개요 <div id='2' />

<br />

### 1.1.1. 목적 <div id='3' />

> 본 문서(개방형 PaaS 플랫폼 고도화 및 개발자 지원환경 개발 가이드)는 Paas-TA 프로젝트의 Monitoring 애플리케이션을 개발 및 배포하는 방법에 대해 제시하는 문서이다.

<br />

###  1.1.2. 범위 <div id='4' />

> 본 문서의 범위는 PaaS-Ta 서비스 들의 시스템 상태를 조회하고, 임계치 정보와의 비교를 통해 관리자에게 관련 정보를 전달하는 방법에 대한 내용으로 한정되어 있다.

<br />

### 1.1.3. 참고자료 <div id='5' />
- https://golang.org/
- https://git-scm.com
- https://github.com/tedsuo/ifrit
- https://github.com/tedsuo/rata
- https://github.com/go-sql-driver/mysql
- https://github.com/jinzhu/gorm
- https://github.com/influxdata/influxdb/client/v2
- https://github.com/gorilla/handlers
- https://github.com/gorilla/mux
- https://github.com/stretchr/testify/assert
- https://github.com/onsi/ginkgo
- https://github.com/onsi/gomega
- https://github.com/tools/godep
- https://github.com/davecgh/go-spew/spew
- https://github.com/pmezard/go-difflib/difflib
- https://github.com/cloudfoundry-community/gogobosh
- https://github.com/onsi/ginkgo
- https://github.com/onsi/gomega

<br /><br /><br />

#   2. IaaS Monitoring 애플리케이션 환경 설정 <div id='6' />

##  2.1. 개요 <div id='7' />

> 개방형 플랫폼 프로젝트의 모니터링 시스템에서 IaaS(Openstack)시스템의 상태를 조회하여, 사전에 설정한 임계치 값과 비교 후, 초과된 시스템 자원을 사용중인 서비스들의 목록을 관리자에게 통보하기 위한 애플리케이션 개발하고, 배포하는 방법을 설명한다.

<br />

##  2.2. 개발환경 구성 <div id='8' />

> 애플리케이션 개발을 위해 다음과 같은 환경으로 개발환경을 구성 한다.

    - OS : Window/Ubuntu
    - Golang : 1.8.3
    - Dependencies :  github.com/tedsuo/ifrit
                      github.com/tedsuo/rata
                      github.com/influxdata/influxdb/client/v2
                      github.com/rackspace/gophercloud
                      github.com/go-sql-driver/mysql
                      github.com/jinzhu/gorm
                      github.com/cihub/seelog
                      github.com/monasca/golang-monascaclient/monascaclient
                      github.com/gophercloud/gophercloud/
                      github.com/alexedwards/scs
                      gopkg.in/olivere/elastic.v3
                      github.com/onsi/ginkgo
                      github.com/onsi/gomega
                      github.com/stretchr/testify
    - IDE : Intellij IDEA 2017.
    - 형상관리: Git

> ※	Intellij IDEA 는 Commnuity와 Ultimate 버전이 있는데, Community 버전은 Free이고, Ultimate 버전은 은 30-day trial버전이다

<br />

###  2.2.1.	Git 설치 <div id='9' />

> 아래 URL에서 자신에 OS에 맞는 Git client를 다운로드 받아 설치 한다.
    
    https://git-scm.com/downloads

<br />

###  2.2.2.	Go Lang 설치 <div id='10' />

> 아래 URL에서 자신에 OS에 맞는 go SDK를 다운로드 받아 설치 한다. (1.8 이상)

    https://golang.org/dl

> GOROOT, 및 PATH를 설정한다.

<br />

###  2.2.3.	IntelliJ IDEA설치 <div id='11' />

> **IDEA 다운로드**

    https://www.jetbrains.com/idea/?fromMenu

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.2.2_1.png)
&nbsp;&nbsp;&nbsp; ![intellj_install_2](images/2.2.2_2.png)
&nbsp;&nbsp;&nbsp; ![intellj_install_3](images/2.2.2_3.png)

<br />

> **IntelliJ IDEA 설치**

&nbsp;&nbsp;&nbsp; idealC-2017.2.5.exe 더블클릭하여 설치를 실행한다.

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_3](images/2.2.3_1.png)

- 'Next' 버튼 클릭

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.2.3_2.png)

- 설치위치 지정 후 'Next' 버튼 클릭

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.2.3_3.png)

- 'Next' 버튼 클릭

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.2.3_4.png)

- 'Install' 버튼 클릭

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.2.3_5.png)

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.2.3_6.png)

- Run IntelliJ IDEA Community Edition” 체크 선택
- 'Finish' 버튼 클릭

<br />

###  2.2.4.	Go 애플리케이션 환경설정 <div id='12' />

> 만약, Go SDK 설정이 되어 있지 않을 경우, 아래 절차를 통해 SDK를 등록한다.

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.2.4_2.png)

- 화면상단 메뉴에서 File > Setting 을 클릭한다.

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.2.4_2.png)

- 왼쪽 메뉴에서 "Plugin"을 선택 후, "Browse repositories" 버튼을 클릭한다.

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.2.4_3.png)

- 검색어 입력란에 "Go"  입력 후, 조회된 결과에서 "Go"를 선택한 뒤, "Install" 버튼을 클릭한다.

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.2.4_4.png)

- "Restart IntelliJ IDEA" 선택
- IntellJ를 재시작 한다.

<br />

##   2.3. 소스 다운로드 <div id='13' />

> PaaS-TA 소스를 다운로드 받는다.
    
    $ cd workspace
    $ git clone https://github.com/PaaS-TA/PaaS-TA-Monitoring

<br />

##   2.4. 애플리케이션 프로젝트 Open <div id='14' />

> Paasta-monitoring 애플리케이션을 개발하기 위한 애플리케이션의 생성과 환경설정, Dependencies Module에 대한 Import 방법에 대하여 설명한다.

> Paasta-monitoring은  pasta-monitoring-management와  pasta-monitoring-batch 프로젝트로 구성되어 있다.

<table>
    <tr>
        <th>프로젝트</th>
        <th>목적</th>
    </tr>
    <tr>
        <td>pasta-monitoring-management</td>
        <td>paasta 모니터링 관련하여 Container 배포 상태, 알람설정, Alarm 정보등을 조회 가능하다. 이 프로젝트를 실행하기 전에 pasta-monitoring-batch가 먼저 실행되어야 한다.</td>
    </tr>
    <tr>
        <td>pasta-monitoring-batch</td>
        <td>Bosh, pasta vm, container 를 모니터링 및 알람 설정값이 벗어났을때 paasta portal에 Auto-scaing을 요청한다.</td>
    </tr>
</table>

<br />

###  2.4.1. PaaS-TA-Monitoring Management 프로젝트 Open <div id='15' />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.4.1_1.png)

- IntellJ 실행 후 "Open" 을 선택한다.

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.4.1_2.png)

- 화면상단 메뉴에서 File > Open 을 클릭한다.

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.4.1_3.png)

- "Languages & Frameworks" 를 클릭한다.
- "Go"를 클릭한다

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.4.1_4.png)

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.4.1_5.png)

- "GOROOT" 를 클릭한다. SDK를 아래와 같이 선택한다. 

- (※ GO Lang의 버전은 1.8 이상 설치)

<br />

&nbsp;&nbsp;&nbsp; ![intellj_install_1](images/2.4.1_6.png)

- Global GOPATH 우측 + 버튼을 클릭하여 "C:\Go\bin" 설정한다. 
- Project GOPATH 우측 + 버튼을 클릭하여 "\…\PaaS-Monitoring\src\paasta-monitoring-batch" 로 설정한다.
- Project GOPATH 우측 + 버튼을 클릭하여 "\…\PaaS-Monitoring\src\paasta-monitoring-management" 로 설정한다.
- IntellJ 를 재시작한다.


<br /><br /><br />



## Developer Guide
#### [Developer Guide - Batch](./paasta-monitoring-batch/doc/PaaS-TA-Monitoring-2.0_Batch.md)

#### [Developer Guide - Management](./paasta-monitoring-management/doc/PaaS-TA-Monitoring-2.0_Management.md)
