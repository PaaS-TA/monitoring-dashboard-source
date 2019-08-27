
const psRoot = {
	url : 'http://13.125.236.133:8090/v2/',
	timeRange : '1440',
	timeRefresh : '1d',
	defaultTimeRange : '15m',
	groupBy : '1m',
}

const ssRoot = {
	url : 'http://13.125.236.133:8090/v2/',
	timeRange : '1440',
	timeRefresh : '1d',
	defaultTimeRange : '15m',
	groupBy : '1m',
}

const csRoot = {
	url : 'http://13.125.236.133:8090/v2/',
	timeRange : '1440',
	timeRefresh : '1d',
	defaultTimeRange : '15m',
	groupBy : '1m',
}

const isRoot = {
	url : 'http://13.125.236.133:8090/v2/',
	timeRange : '1440',
	timeRefresh : '1d',
	defaultTimeRange : '15m',
	groupBy : '1m',
}

const fnComm = {
	url : 'http://13.125.236.133:8090/v2/',
	timeRange : '1440',
	timeRefresh : '1d',
	defaultTimeRange : '15m',
	groupBy : '1m',

	setTimer(){
		return `?defaultTimeRange=${psRoot.defaultTimeRange}&groupBy=${psRoot.groupBy}`;
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
		//if(!sessionStorage.getItem('login')){
			//fnComm.getToken();

			var request = new XMLHttpRequest();
			request.open('POST', `${psRoot.url}login`, false);

			request.onreadystatechange = () => {
				if (request.readyState === XMLHttpRequest.DONE){
					if(request.status === 200){
						console.log(JSON.parse(request.responseText));
						sessionStorage.setItem('login', true);
						sessionStorage.setItem('user', user);
						sessionStorage.setItem('mail', JSON.parse(request.responseText).userEmail);
						document.location.href = 'index.html';
					} else {
						alert(JSON.parse(request.responseText).message);
					};
				};
			};
			
			request.setRequestHeader('X-XSRF-TOKEN', sessionStorage.getItem('token'));
			request.send(`{"username":"${user}","password":"${pw}"}`);
			
		//};
	},

	// 로그인 TOKEN 할당 ////////////////////////////////////////////////////////////////
    getToken(user, pw){
		var request = new XMLHttpRequest();

		try {
			request.open('GET', `${psRoot.url}ping`, false);
			request.send();
			console.log(request);

			//alert('Authentication failed\nContact your manager');

			var tokenArray = request.getAllResponseHeaders().toLowerCase().split('\n');
			
			for(let value of tokenArray){
				if(value.indexOf('x-xsrf-token') != -1){
					//console.log(value.split(': ')[1]);
					sessionStorage.setItem('token', value.split(': ')[1]);
					fnComm.loginCheck(user, pw);
				};
			};
		}
		catch {
			fnComm.alertPopup('ERROR', 'Authentication failed\nContact your manager')
		}
	},

	// 알람 카운트 //////////////////////////////////////////////////////////////////////
	alarmCount(){
		document.querySelector('.outBtn strong').innerHTML = sessionStorage.getItem('user');
		document.querySelector('.outBtn span').innerHTML = sessionStorage.getItem('mail');
		
		var request = new XMLHttpRequest();
		request.open('GET', `${psRoot.url}paas/alarm/status/count?resolveStatus=1&state=ALARM`, false);
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

						fnComm.timeSetting();
					} else {
						document.querySelector('.alarmView').classList.toggle('on', false);
					}
				} else {
					console.log(JSON.parse(request.responseText).HttpStatus+' Error!\n'+JSON.parse(request.responseText).message);
				};
			};
		};

		request.send();
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// 검색 타이머 설정 - timeSetting()
	/////////////////////////////////////////////////////////////////////////////////////
	timeSetting(){
		// 검색 타이머 설정여부
		if(sessionStorage.getItem('defaultTimeRange')){
			psRoot.defaultTimeRange = sessionStorage.getItem('defaultTimeRange');
			psRoot.groupBy = sessionStorage.getItem('groupBy');
		} else {
			sessionStorage.setItem('defaultTimeRange', psRoot.defaultTimeRange);
			sessionStorage.setItem('groupBy', psRoot.groupBy);
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
				psRoot.defaultTimeRange = e.target.value;
				psRoot.groupBy = e.target.getAttribute('data-group');
			}, false);
		};
		
		// logout 이벤트
		document.querySelector('.logout').addEventListener('click', (e) => {
			sessionStorage.setItem('login', false);
			sessionStorage.setItem('token', '');
			document.location.href = 'login.html';
		}, false);
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// 데이터 로드 - loadData(method, url, callbackFunction)
	// (전송타입, url, 콜백함수)
	/////////////////////////////////////////////////////////////////////////////////////
	loadData(method, url, callbackFunction, list){
		var request = new XMLHttpRequest();
		request.open(method, url);
		//request.setRequestHeader('X-XSRF-TOKEN', sessionStorage.getItem('token'));

		request.onreadystatechange = () => {
			if (request.readyState === XMLHttpRequest.DONE){
				if(request.status === 200 && request.responseText != ''){
					callbackFunction(JSON.parse(request.responseText), list);
				} else {
					//sessionStorage.setItem('login', false);
					//sessionStorage.setItem('token', '');
					//document.location.href = 'login.html';
				};
			};
		};

		request.send();
	},

	/////////////////////////////////////////////////////////////////////////////////////
	// 데이터 로드 - saveData(method, url, callbackFunction)
	// (전송타입, url, 콜백함수)
	/////////////////////////////////////////////////////////////////////////////////////
	saveData(method, url, data){
		console.log(data);
		var request = new XMLHttpRequest();
		request.open(method, url);
		request.setRequestHeader('X-XSRF-TOKEN', sessionStorage.getItem('token'));

		request.onreadystatechange = () => {
			if (request.readyState === XMLHttpRequest.DONE){
				if(request.status === 200 && request.responseText != ''){
					console.log('success');
					callbackFunction(JSON.parse(request.responseText), list);
				} else {
					console.log('error');
					sessionStorage.setItem('login', false);
					sessionStorage.setItem('token', '');
					document.location.href = 'login.html';
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
	alertPopup(title, text, fn){
		var html = `<div id="alertPop"><div><h3>${title}</h3><p>${text}</p><div><button class="closed">Confirm</button></div></div></div>`;

		fnComm.appendHtml(document.body, html, 'body');

		document.getElementById('alertPop').querySelector('.closed').addEventListener('click', (e) => {
			if(fn) fn();
			
			document.body.removeChild(document.getElementById('alertPop'));
		}, false);
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
		});
	},

	// CONTAINER 컨디션 차트 ////////////////////////////////////////////////////////////////
	contConditionChart(data){
		console.log(data);
		var chart = c3.generate({
    		bindto: '#contChart',
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

	// Detail 차트 //////////////////////////////////////////////////////////////////////////
	detailChart(data, target){
		let cnt = 0;
		let yPos = 9;
		let detailTime = [];
		let detailData = [];
		let chartColor;
		
		data.forEach((load, idx) => {
			if(idx === 0){
				detailTime.push('time');
			};

			detailData[idx] = [load.name];

			load.metric.forEach(value => {
				if(idx === 0){
					var timeStamp = new Date((value.time+32400)*1000);

					detailTime.push(('0' + timeStamp.getUTCHours()).slice(-2) + ':' + ('0' + timeStamp.getUTCMinutes()).slice(-2) + ':' + ('0' + timeStamp.getUTCSeconds()).slice(-2));
				};

				detailData[idx].push(Number(Math.ceil(value.usage)));

				if(yPos < Number(value.usage)) yPos = Math.ceil(value.usage);
			});

			cnt = idx;
		});

		yPos = Math.ceil((yPos / 10)) * 10;
		
		var dataType = [];
		dataType.push(detailTime);
		for (var i = 0; i <= cnt; i++) {
			dataType.push(detailData[i]);
		};
		console.log(dataType);

		switch (target){
			case '#cpuUsageChart':
			case '#cpuLoadChart':
				chartColor = ['#ff015a', '#fc6604', '#fcce34'];
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
		};

		var chart = c3.generate({
    		bindto: target,
			data: {
				x: 'time',
				xFormat: '%H:%M:%S',
				columns: dataType,
				labels: false,
				type: 'spline',
			},
			color: {
				pattern: chartColor
			},
			axis: {
				x: {
					type: 'timeseries',
					localtime: false,
					tick: {
						count: 2,
					}
				},
				y: {
					max: yPos,
					min: 0,
					tick: {
						count: 5,
					}
				}
			},
			point: {
				show: false
			}
		});
	},

	detailLog(data){
		console.log(data);
	},
};