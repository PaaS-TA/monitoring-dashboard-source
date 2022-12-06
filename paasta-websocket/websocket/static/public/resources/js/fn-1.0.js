
const fnComm = {
	url : '/v2/',
	timeRange : '1440',
	timeRefresh : '1d',
	defaultTimeRange : '15m',
	groupBy : '1m',

	setTimer(){
		return `?defaultTimeRange=${fnComm.defaultTimeRange}&groupBy=${fnComm.groupBy}`;
	},

	init(){
		document.querySelector('.timeSetting').addEventListener('click', () => {
			document.querySelector('.timePop').classList.toggle('on', true);
		});

		document.querySelector('.timePop > a').addEventListener('click', () => {
			document.querySelector('.timePop').classList.toggle('on', false);
		});

		sessionStorage.setItem('login',false);

		fnComm.loginCheck();
		fnComm.alarmCount();
	},

	// 로그인 체크 ////////////////////////////////////////////////////////////////
	loginCheck(user, pw){
		var request = new XMLHttpRequest();

		request.open('POST', `${fnComm.url}login`, false);
		request.setRequestHeader('X-XSRF-TOKEN', sessionStorage.getItem('token'));

		request.onreadystatechange = () => {
			if (request.readyState === XMLHttpRequest.DONE){
				if(request.status === 200){
					var userInfo = JSON.parse(request.responseText);
					sessionStorage.setItem('login', true);
					sessionStorage.setItem('user', user);
					sessionStorage.setItem('mail', userInfo.userEmail);
					sessionStorage.setItem('sysType', userInfo.sysType);

					var localInfo = {"name":user,"email":userInfo.userEmail,"sysType":"IaaS","i1":"I","i2":"S","p1":"","p2":""};

					localStorage.setItem('ls.user', JSON.stringify(localInfo));

					document.location.href = 'index.html';
				} else {
					fnComm.alertPopup('ERROR', JSON.parse(request.responseText).message);
				};
			};
		};

		request.send(`{"username":"${user}","password":"${pw}"}`);
	},

	// 로그인 TOKEN 할당 ////////////////////////////////////////////////////////////////
	getToken(user, pw){
		var request = new XMLHttpRequest();

		try {
			request.open('GET', `${fnComm.url}ping`, false);
			request.send();
			console.log(request);

			var tokenArray = request.getAllResponseHeaders().toLowerCase().split('\n');

			for(let value of tokenArray){
				if(value.indexOf('x-xsrf-token') != -1){
					console.log(value.split(': ')[1]);
					sessionStorage.setItem('token', value.split(': ')[1]);
					localStorage.setItem('ls.token', value.split(': ')[1]);
					fnComm.loginCheck(user, pw);
				};
			};
		}
		catch {
			fnComm.alertPopup('ERROR', 'Authentication failed\nContact your manager')
		}
	},

	// 알람 카운트 //////////////////////////////////////////////////////////////////////
	alarmCount(type){
		document.querySelector('.outBtn strong').innerHTML = sessionStorage.getItem('user');
		document.querySelector('.outBtn span').innerHTML = sessionStorage.getItem('mail');

		// 오늘 날짜
		var today = new Date();
		var dd = today.getDate();
		var mm = today.getMonth()+1;
		var nn = today.getMonth();
		var yyyy = today.getFullYear();

		if(dd < 10) dd = '0'+dd;
		if(mm < 10) mm = '0'+mm;
		if(nn < 10) nn = '0'+nn;

		// 검색 Default set
		var form = `${yyyy}-${nn}-${dd}`;
		var to = `${yyyy}-${mm}-${dd}`;

		var url = '';
		switch(type){
			case 'paas':
				url = `${fnComm.url}paas/alarm/status/count?resolveStatus=1&searchDateFrom=${form}&searchDateTo=${to}`;
				break;
			case 'caas':
				url = `${fnComm.url}caas/monitoring/alarmCount?searchDateFrom=${form}&searchDateTo=${to}`;
				break;
			case 'saas':
				url = `${fnComm.url}saas/app/application/alarmCount?searchDateFrom=${form}&searchDateTo=${to}`;
				break;
		};

		var request = new XMLHttpRequest();

		try {
			request.open('GET', url, false);
			request.setRequestHeader('X-XSRF-TOKEN', sessionStorage.getItem('token'));

			request.onreadystatechange = () => {
				if (request.readyState === XMLHttpRequest.DONE){
					if(request.status === 200){
						if(JSON.parse(request.responseText).totalCnt > 0){
							document.querySelector('.alarmView').classList.toggle('on', true);
							document.querySelector('.alarmView span').innerHTML = JSON.parse(request.responseText).totalCnt;

							document.querySelector('.alarmView span').addEventListener('click', () => {
								location.href = 'alarm_status.html';
							});

							//fnComm.timeSetting();
						} else {
							document.querySelector('.alarmView').classList.toggle('on', false);
						}
					};
				};
			};

			request.send();

			var type = sessionStorage.getItem('sysType').split(',');

			// config.ini에 설정된 타입에 따라 헤더 메뉴의 show/hide 여부를 제어함
			for(var i=0 ; i<type.length ; i++){
				if (type[i] !== 'ALL') {
					document.querySelector(`.global .${type[i]}`).style.display = 'inline-block';
				} else {
					document.querySelector(`.global .IaaS`).style.display = 'inline-block';
					document.querySelector(`.global .PaaS`).style.display = 'inline-block';
					document.querySelector(`.global .SaaS`).style.display = 'inline-block';
					document.querySelector(`.global .CaaS`).style.display = 'inline-block';
				}
			};

			// logout 이벤트
			document.querySelector('.logout').addEventListener('click', (e) => {
				sessionStorage.clear();
				document.location.href = '../login.html';
			}, false);
		}
		catch {
			fnComm.alertPopup('ERROR', 'Authentication failed\nContact your manager', fnComm.winReCall);
		}
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// 검색 타이머 설정 - timeSetting()
	/////////////////////////////////////////////////////////////////////////////////////
	timeSetting(){
		// 검색 타이머 설정여부
		if(sessionStorage.getItem('defaultTimeRange')){
			fnComm.defaultTimeRange = sessionStorage.getItem('defaultTimeRange');
			fnComm.groupBy = sessionStorage.getItem('groupBy');
		} else {
			sessionStorage.setItem('defaultTimeRange', fnComm.defaultTimeRange);
			sessionStorage.setItem('groupBy', fnComm.groupBy);
		};

		document.querySelector('.timeSetting').addEventListener('click', (e) => {
			document.querySelector('.timePop').classList.toggle('on');
		}, false);

		document.querySelector('.timePop .close').addEventListener('click', (e) => {
			document.querySelector('.timePop').classList.toggle('on');
		}, false);

		// 검색 타이머 라디오 이벤트
		var radio = document.querySelectorAll('.timePop [name=timeRange]');

		for(var i=0 ; i<radio.length ; i++){
			radio[i].addEventListener('click', (e) => {
				sessionStorage.setItem('defaultTimeRange', e.target.value);
				sessionStorage.setItem('groupBy', e.target.getAttribute('data-group'));
				fnComm.defaultTimeRange = e.target.value;
				fnComm.groupBy = e.target.getAttribute('data-group');
			}, false);
		};
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// 데이터 로드 - loadData(method, url, callbackFunction)
	// (전송타입, url, 콜백함수)
	/////////////////////////////////////////////////////////////////////////////////////
	loadData(method, url, callbackFunction, list){
		if(sessionStorage.getItem('token') == null){
			let href = document.location.href;
			if(href.includes('/public/index')){
				document.location.href = '../public/login.html';
			}else{
				document.location.href = '../login.html';
			}
			console.log("token expired..");
			return;
		}
		var request = new XMLHttpRequest();
		request.open(method, url);
		request.setRequestHeader('X-XSRF-TOKEN', sessionStorage.getItem('token'));

		request.onreadystatechange = () => {
			if (request.readyState === XMLHttpRequest.DONE){
				if(request.status === 200 && request.responseText != ''){
					callbackFunction(JSON.parse(request.responseText), list);
				} else if(request.status === 401){
					sessionStorage.clear();
					document.location.href = '../login.html';
				} else if (request.status === 500) {
					fnComm.alertPopup('ERROR', JSON.parse(request.responseText).message);
				};
			};
		};

		request.send();
	},


	requestAjax(method, url, callbackFunction, errCallback, list){
		if(sessionStorage.getItem('token') == null){
			document.location.href = '../login.html';
		}
		var request = new XMLHttpRequest();
		request.open(method, url);
		request.setRequestHeader('X-XSRF-TOKEN', sessionStorage.getItem('token'));

		request.onreadystatechange = () => {
			if (request.readyState === XMLHttpRequest.DONE){
				if(request.status === 200 && request.responseText != ''){
					callbackFunction(JSON.parse(request.responseText), list);
				} else if(request.status === 401){
					sessionStorage.clear();
					document.location.href = '../login.html';
				} else if (request.status === 500) {
					errCallback(JSON.parse(request.responseText).message);
					//fnComm.alertPopup('ERROR', JSON.parse(request.responseText).message);
				};
			};
		};

		request.send();
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// 데이터 저장 - saveData(method, url, data)
	// (전송타입, url, 데이터)
	/////////////////////////////////////////////////////////////////////////////////////
	saveData(method, url, data, bull){
		var request = new XMLHttpRequest();
		request.open(method, url);
		request.setRequestHeader('X-XSRF-TOKEN', sessionStorage.getItem('token'));

		console.log(sessionStorage.getItem('token'));

		request.onreadystatechange = () => {
			if (request.readyState === XMLHttpRequest.DONE){
				if(request.status === 201 || request.status === 200){
					if(method == 'POST'){
						fnComm.alertPopup('SAVE', 'COMPLETE', fnComm.winReload);
					} else if(method == 'PATCH'){
						fnComm.alertPopup('PATCH', 'COMPLETE', fnComm.winReload);
					} else if(method == 'PUT'){
						if(bull != 'false'){
							fnComm.alertPopup('PUT', 'COMPLETE', fnComm.winReload);
						};
					} else if(method == 'DELETE'){
						fnComm.alertPopup('DELETE', 'COMPLETE', fnComm.winReload);
					};
				} else {
					if(method == 'DELETE'){
						fnComm.alertPopup('DELETE', 'DELETE FAILED', fnComm.winReload);
					} else {
						fnComm.alertPopup('SAVE', 'SAVE FAILED', fnComm.winReload);
					};
					//sessionStorage.setItem('login', false);
					//sessionStorage.setItem('token', '');
					//document.location.href = 'login.html';
				};
			};
		};

		request.send(data);
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// 소수점 자릿수 제어 - numberComma(digit, number)
	// (소수점 자릿수, 데이터 숫자)
	/////////////////////////////////////////////////////////////////////////////////////
	numberComma(digit, number){
		return number.toFixed(digit).toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
	},

	convertUnits(size, unit){
		var convertedSize, convertedUnit;
		if (unit == "MB") {
			if (size >= 1048576) {
				convertedSize = size / 1048576;
				convertedUnit = "TB";
			} else if (size >= 1024) {
				convertedSize = size / 1024;
				convertedUnit = "GB";
			} else {
				convertedSize = size;
				convertedUnit = "MB";
			};
		} else if (unit == "GB") {
			if (size >= 1048576) {
				convertedSize = size / 1048576;
				convertedUnit = "PB";
			} else if (size >= 1024) {
				convertedSize = size / 1024;
				convertedUnit = "TB";
			} else {
				convertedSize = size;
				convertedUnit = "GB";
			};
		};
		return {convertedSize, convertedUnit};
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// 날짜 컨버팅 - unixTime(time)
	// (20190601)
	/////////////////////////////////////////////////////////////////////////////////////
	unixTime : function(time){
		var timestamp = new Date((time+32400)*1000);

		return ('0' + timestamp.getUTCHours()).slice(-2) + ':' + ('0' + timestamp.getUTCMinutes()).slice(-2);
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// html 생성 - appendHtml(target, html)
	// (삽입 타겟, html)
	/////////////////////////////////////////////////////////////////////////////////////
	appendHtml(target, html, type){
		var div = document.createElement(type);
		div.innerHTML = html;
		while (div.children.length > 0){
			target.appendChild(div.children[0]);
		};
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// html 삭제 - removeHtml(target)
	// (타겟 : 타겟의 자식요소 전부 삭제)
	/////////////////////////////////////////////////////////////////////////////////////
	removeHtml(target){
		while(target.hasChildNodes()){
			target.removeChild(target.firstChild);
		};
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// 공통링크 삽입 - linkHtml(target, url)
	// (링크타겟, 이동 URL)
	/////////////////////////////////////////////////////////////////////////////////////
	linkHtml(target, url){
		for(var i=0 ; i<target.length ; i++){
			target[i].addEventListener('click', () => {
				document.location.href  = url;
			}, false);
		};
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// 공통 경고 팝업 - alertPopup(title, text, fn)
	// (제목, 경고문구, 콜백함수)
	/////////////////////////////////////////////////////////////////////////////////////
	alertPopup(title, text, callback){
		var html = `<div id="alertPop"><div><h3>${title}</h3><p>${text}</p><div><button class="closed">Confirm</button></div></div></div>`;

		fnComm.appendHtml(document.body, html, 'body');

		document.getElementById('alertPop').querySelector('.closed').addEventListener('click', (e) => {
			if(callback) callback();

			document.body.removeChild(document.getElementById('alertPop'));
		}, false);
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// Window Reload
	/////////////////////////////////////////////////////////////////////////////////////
	winReload(){
		window.location.reload();
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// Window ReCall
	/////////////////////////////////////////////////////////////////////////////////////
	winReCall(){
		console.log(window.location.href);
		if(window.location.href.indexOf('paas') != -1 || window.location.href.indexOf('caas') != -1 || window.location.href.indexOf('saas') != -1){
			window.location.href = '../login.html';
		} else {
			window.location.href = 'login.html';
		};
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// Count UP
	/////////////////////////////////////////////////////////////////////////////////////
	countUp(el, num) {
		var cnt = 0;
		var dif = 0;

		var thisID = setInterval(function(){
			if(cnt < num){
				dif = num - cnt;

				if(dif > 0) {
					cnt += Math.ceil(dif / 5);
				};

				el.innerHTML = cnt;
			} else {
				clearInterval(thisID);
			};
		}, 20);
	},

	// BOSH 컨디션 차트 ////////////////////////////////////////////////////////////////
	boshConditionChart(data){
		var chart = c3.generate({
			bindto: '#boshChart',
			data: {
				columns: [
					['Running', 0, data.running, 0, 0, 0, 0, 0],
					['Warning', 0, 0, data.warning, 0, 0, 0, 0],
					['Critical', 0, 0, 0, data.critical, 0, 0, 0],
					['Failed', 0, 0, 0, 0, data.failed, 0, 0],
					['Total', 0, 0, 0, 0, 0, data.total, 0]
				],
				labels: true,
				type: 'spline'
			},
			color: {
				pattern: ['#43be42', '#f6a200', '#f4256c', '#c048c8', '#ccc']
			},
			axis: {
				x: {
					show: false
				},
				y: {
					show: false
				}
			},
			tooltip: {
				//show: false
			}
		});
	},

	// PaaS 컨디션 차트 ////////////////////////////////////////////////////////////////
	paasConditionChart(data){
		var chart = c3.generate({
			bindto: '#paasChart',
			data: {
				columns: [
					['Running', 0, data.Running, 0, 0, 0, 0, 0],
					['Warning', 0, 0, data.Warning, 0, 0, 0, 0],
					['Critical', 0, 0, 0, data.Critical, 0, 0, 0],
					['Failed', 0, 0, 0, 0, data.Failed, 0, 0],
					['Total', 0, 0, 0, 0, 0, data.Total, 0]
				],
				labels: true,
				type: 'spline'
			},
			color: {
				pattern: ['#43be42', '#f6a200', '#f4256c', '#c048c8', '#ccc']
			},
			axis: {
				x: {
					show: false
				},
				y: {
					show: false
				}
			},
		});
	},

	// CONTAINER 컨디션 차트 ////////////////////////////////////////////////////////////////
	contConditionChart(data){
		var chart = c3.generate({
			bindto: '#contChart',
			data: {
				columns: [
					['Running', 0, data.Running, 0, 0, 0, 0, 0],
					['Warning', 0, 0, data.Warning, 0, 0, 0, 0],
					['Critical', 0, 0, 0, data.Critical, 0, 0, 0],
					['Failed', 0, 0, 0, 0, data.Failed, 0, 0],
					['Total', 0, 0, 0, 0, 0, data.Total, 0]
				],
				labels: true,
				type: 'spline'
			},
			color: {
				pattern: ['#43be42', '#f6a200', '#f4256c', '#c048c8', '#ccc']
			},
			axis: {
				x: {
					show: false
				},
				y: {
					show: false
				}
			},
			tooltip: {
				//show: false
			}
		});
	},

	// Detail 차트 //////////////////////////////////////////////////////////////////////////
	detailChart(data, target){
		let cnt = 0;
		let yPos = 9;
		let xPosArr = [];
		let xPos = 0;
		let detailTime = [];
		let detailData = [{time: 1669892640, usage: 0.22},{time: 1669892700, usage: 0.22999999999999998},{time: 1669892760, usage: 0.08},{time: 1669892820, usage: 0.185},{time: 1669892880, usage: 0.16999999999999998},
			{time: 1669892940, usage: 0.06}];
		let chartColor;
		let xPosDiv;
		let sizeHeightDiv;
		let sizeWidthDiv;

		console.log("ㅎㅇ")
		console.log("detailData : ", detailData);

		data.forEach((load, idx) => {
			if(idx === 0){
				detailTime.push('time');
			};

			detailData[idx] = [load.name];

			load.metric.forEach(value => {

				if(idx === 0){

					var timeStamp = new Date((value.time)*1000);

					timeStamp.setHours(timeStamp.getHours()+9);

					let utcFullYear = timeStamp.getUTCFullYear();
					let utcMonth = timeStamp.getUTCMonth()+1;
					let utcDate = timeStamp.getUTCDate();

					detailTime.push(utcFullYear + '-' + utcMonth + '-' + utcDate + ' ' + ('0' + timeStamp.getUTCHours()).slice(-2) + ':' + ('0' + timeStamp.getUTCMinutes()).slice(-2) + ':' + ('0' + timeStamp.getUTCSeconds()).slice(-2));
					xPosArr.push(Number(value.usage));
				};

				detailData[idx].push(Number(Math.ceil(value.usage)));

				if(yPos < Number(value.usage))
					yPos = Math.ceil(value.usage);

			});

			cnt = idx;
		});

		yPos = Math.ceil((yPos / 10)) * 10;
		xPos = Math.floor((Math.min.apply(null, xPosArr) / 10))* 10;

		var dataType = [];
		dataType.push(detailTime);
		for (var i = 0; i <= cnt; i++) {
			dataType.push(detailData[i]);
		};

		switch (target){
			case '#cpuUsageChart':
			case '#cpuLoadChart':
				chartColor = ['#ff015a', '#fc6604', '#fcce34', '#cc86cc'];
				break;
			case '#memoryUsageChart':
				chartColor = ['#9cce34'];
				break;
			case '#diskUsageChart':
			case '#diskIoChart':
				chartColor = ['#3d003d', '#c701c7', '#cc86cc', '#fccefc'];
				break;
			case '#networkByteChart':
			case '#networkPacktesChart':
			case '#networkDropChart':
			case '#networkErrorChart':
				chartColor = ['#649afc', '#b7c3d8'];
				break;
			case '#podChart':
				chartColor = ['#649afc'];
				break;
			case '#cpuChart' :
				chartColor = ['#b7c3d8'];
				break;
			case '#memoryChart':
				chartColor = ['#9cce34'];
				break;
			case '#diskChart' :
				chartColor = ['#ff015a'];
				break;
		};

		switch (target){
			case '#cpuUsageChart':
			case '#cpuLoadChart':
			case '#memoryUsageChart':
			case '#diskUsageChart':
			case '#diskIoChart':
			case '#networkByteChart':
			case '#networkPacktesChart':
			case '#networkDropChart':
			case '#networkErrorChart':
				xPosDiv = 0;
				sizeHeightDiv = 0;
				sizeWidthDiv = 0;
				break;
			case '#podChart':
			case '#cpuChart' :
			case '#memoryChart':
			case '#diskChart' :
				xPosDiv = xPos;
				sizeHeightDiv = 140;
				sizeWidthDiv = 400;
				break;
		};

		console.log(dataType);

		var xGridline = [];
		var chart = c3.generate({
			bindto: target,
			data: {
				x: 'time',
				xFormat: '%Y-%m-%d %H:%M:%S',
				columns: dataType,
				labels: false,
				type: 'spline',
			},
			color: {
				pattern: chartColor
			},
			axis: {
				x: {
					//type: 'category',
					type: 'timeseries',
					localtime: true,
					tick: {
						count: 3,
						format: '%H:%M'
					},
				},
				y: {
					max: yPos,
					min: xPosDiv,
					tick: {
						count: 5,
					},
					padding: {top: 0, bottom: 0}
				}
			},
			point: {
				show: true
			},
			grid: {
				x: {
					lines: function() {
						for(var i =1; detailTime.length > i; i++) {
							if(i % 4 === 0) {
								xGridline.push({value: detailTime[i], class:'dash-line'})
							}
						}
						return xGridline;
					},
				},
				y: {
					//show: true
				}
			},
			padding: {
				right:20
			},
			size: {
				height: sizeHeightDiv,
				width: sizeWidthDiv
			}
		});
	},
	// Gauge 차트 //////////////////////////////////////////////////////////////////////////
	gaugeChart(data, category, target){
		var chart = c3.generate({
			bindto: target,
			data: {
				columns: [
					[category, 0]
				],
				type: 'gauge',
			},
			gauge: {
				label: {
					format: function(value, ratio) {
						return value + ' %';
					},
					show: false,
				},
				units: ' %',
				width: 3,
			},
			color: {
				pattern: ['#55c554', '#fba602', '#F97600', '#e91a61'],
				threshold: {
					unit: 'value',
					max: 200,
					values: [30, 60, 90, 100]
				}
			},
			size: {
				//height: 150,
			}
		});

		setTimeout(function () {
			chart.load({
				columns: [[category, data]]
			});
		}, 500);
	},

	detailLog(data){
		console.log(data);
	},

	calElapsedTime(dateTimeStr, timeGap) {
		var now = new Date();
		now.setHours(now.getHours() - timeGap);   // 시차가 있다면 시차를 반영한다.
		var targetDate = new Date(dateTimeStr);

		var elapsed;
		if (now.getFullYear() > targetDate.getFullYear()) {
			elapsed = now.getFullYear() - targetDate.getFullYear();
			elapsed = elapsed + ' year';
		} else if (now.getMonth() > targetDate.getMonth()) {
			elapsed = now.getMonth() - targetDate.getMonth();
			elapsed = elapsed + ' month';
		} else if (now.getDate() > targetDate.getDate()) {
			elapsed = now.getDate() - targetDate.getDate();
			elapsed = elapsed + ' day';
		} else if (now.getDate() == targetDate.getDate()) {
			var nowTime = now.getTime();
			var targetTime = targetDate.getTime();

			if (nowTime > targetTime) {
				var elapsedTime = nowTime - targetTime;
				const hour = String(Math.floor((elapsedTime/ (1000 * 60 *60 )) % 24 )).padStart(2, "0"); // 시
				const minutes = String(Math.floor((elapsedTime  / (1000 * 60 )) % 60 )).padStart(2, "0"); // 분
				const second = String(Math.floor((elapsedTime / 1000 ) % 60)).padStart(2, "0"); // 초
				elapsed = hour + ':' + minutes + ':' + second;
			}
		}
		return elapsed;

	}
};