
// Cookie Script /////////////////////////////////////////////////////////
if($.cookie('defaultTimeRange') == null){
    $.cookie('defaultTimeRange', 15);
    $.cookie('groupBy', 1);
};

var webAdd = location.href;
webAdd = webAdd.replace('http://','').split('/')[0];
console.log(webAdd);


// Public Script /////////////////////////////////////////////////////////
var pass = {

    url : 'http://13.115.122.103:8080/v2/',
    //url : webAdd+'/v2/',
    // pass.chartColor(호출 URL)

    chartColor : ['#6f9654','#6f9654','#1c91c0','#43459d','#43459d','#e7711b','#e7711b','#e7711b','#e7711b'],

    tokenCompare : function(bull){
        var request = new XMLHttpRequest();

        request.open('GET', pass.url+'ping', false);
        request.send(null);

        var token;
        var tokenArray = request.getAllResponseHeaders().toLowerCase().split('\n');

        console.log(tokenArray);

        $.each(tokenArray, function(){
            if(String(this).indexOf('x-xsrf-token') != -1){
                token = String(this).split(': ')[1];
            }
        });

        return token;
    },

    // pass.ajaxLoad(호출 URL)
    ajaxLoad : function(url){
        var result = '';
        $.ajax({
            url:url,
            type:'GET',
            async:false,
            dataType:'json',
            beforeSend : function(xhr){
                xhr.setRequestHeader('X-XSRF-TOKEN', sessionStorage.getItem('token'));
            },
            success:function(data){
                result = data;
                //console.log('success');
            },
            error:function(data){
                console.log('ERROR');
                //location.href = 'login.html';
                return false;
            }
        });

        return result;
    },
    
    ajaxPost : function(url, rule, updata){
        console.log(url);
        console.log(rule);
        console.log(updata);
        $.ajax({
            url:url,
            method:rule,
            data:updata,
            contentType:'application/json',
            async:false,
            beforeSend : function(xhr){
                console.log(sessionStorage.getItem('token'));
                xhr.setRequestHeader('X-XSRF-TOKEN', sessionStorage.getItem('token'));
            },
            success:function(){
                console.log('SUCESS');
            },
            error:function(data){
                console.log(data);
                console.log('ERROR');
            }
        });
    },

    // pass.numberComma(소수점 자릿수, 데이터 숫자)
    numberComma : function(digit, number){
        return number.toFixed(digit).toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
    },

    //32400

    unixTime : function(time){
        var timestamp = new Date((time+32400)*1000);

        return ('0' + timestamp.getUTCHours()).slice(-2) + ':' + ('0' + timestamp.getUTCMinutes()).slice(-2);
    },

    space : function(type){
        if(!type){
            return '\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0';
        } else {
            return '\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0';
        }
    },

    countUp : function(el, str, num){
        var self = this;

        if(num > 0){
            var thisID = setInterval(function(){
                if(num > str){
                    str++;
                    
                    if(el.parent('dd').length > 0){
                        el.parent('dd').addClass('on');
                    } else {
                        el.parents('dl').addClass('on');
                    };

                    el.text(str);
                } else if(num == str) {
                    clearInterval(thisID);
                };
            }, 20);
        } else {
            if(el.parent('dd').length > 0){
                el.text(0).parents('dd').removeClass('on');
            } else {
                el.text(0).parents('dl').removeClass('on');
            };
        }
    },
};

// Public Event /////////////////////////////////////////////////////////
$(function() {
    console.log(sessionStorage.getItem('login'));
    if(location.href.indexOf('login') == -1 && location.href.indexOf('join') == -1){

        if(sessionStorage.getItem('login') == 'logout'){
            sessionStorage.setItem('login', '');
            sessionStorage.setItem('token', '');
            alert('로그아웃 되었습니다.\n다시 로그인해 주십시요.');
            location.href = 'login.html';
        } else if(sessionStorage.getItem('login') == '' || sessionStorage.getItem('login') == 'error' || sessionStorage.getItem('token') == ''){
            alert('정상적인 접근이 아닙니다.\n다시 로그인해 주십시요.');
            sessionStorage.setItem('login', '');
            sessionStorage.setItem('token', '');
            location.href = 'login.html';
        };

        // Header Alarm
        var alarmNumber = 0;
        var headerAlarm = pass.url + 'paas/alarm/status/count?resolveStatus=1&state=ALARM';
        var headerAlarmData = pass.ajaxLoad(headerAlarm);

        console.log(headerAlarmData);

        $('.alarmView span').text(headerAlarmData.totalCnt);


        // Select Event /////////////////////////////////////////////////////////
        var selectTarget = $('.select_wrap select');

        selectTarget.change(function(){
            var select_name = $(this).children('option:selected').text();
            $(this).siblings('label').text(select_name);
        });

        // Timer Event /////////////////////////////////////////////////////////
        // Start Timer Setting
        $('.timePop .timer input[value=' + $.cookie('defaultTimeRange') + ']').trigger('click');
        $('.timePop #timeSelect').val($.cookie('groupBy')).change();

        $('header .timeSetting').off().on('click', function(){
            $('header .timePop').fadeToggle(300);
        });
        $('header .timePop .save').off().on('click', function(){
            $.cookie('defaultTimeRange',  $('.timePop .timer').find('input:checked').val());
            $.cookie('groupBy', $('header .timePop').find('select').val());

            $('header .timePop').fadeOut(300);
        });
        $('header .timePop .close').off().on('click', function(){
            $('header .timePop').fadeOut(300);
        });

        // Alarm Event /////////////////////////////////////////////////////////
        $('header .alarmView').off().on('click', function(){
            if($(this).find('span').text() != '0'){
                $.cookie('alarm', true);
                location.href = 'alarm_status.html';
            } else {
                alert('NO Alarm Message');
            };
        });

        $('header .logout').off().on('click', function(){
            sessionStorage.setItem('login', 'logout');
            sessionStorage.setItem('token', 'logout');
            location.href = 'login.html';
        });

        $('#container a').off().on('click', function(){
            $.cookie('alarm', false);
        });

        $('.popWrap .status > a').on('click', function(){
            $(this).parents('.popWrap').fadeOut(300);
        });
    };
});
