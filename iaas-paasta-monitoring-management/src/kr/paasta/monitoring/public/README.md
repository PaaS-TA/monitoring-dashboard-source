# PaaS-TA-IaaS


## Getting Started

To get you started you can simply clone the `PaaS-TA-IaaS` repository and install the dependencies:


### Clone `PaaS-TA-IaaS`

Clone the `angular-seed` repository using git:

```
git clone https://github.com/CrossentCloud/PaaS-TA-IaaS
cd PaaS-TA-IaaS/src/openstack-monitoring-portal/src/kr/paasta/monitoring/openstack/public
```


### Install

```
npm install

bower install
```

Behind the scenes this will also call `bower install`. After that, you should find out that you have
two new folders in your project.

* `node_modules` - contains the npm packages for the tools we need
* `bower_components` - contains the Angular framework files


### Run

build your application in folder dist
```
gulp package
```

start BrowserSync server on your source files with live reload
```
gulp serve
```
