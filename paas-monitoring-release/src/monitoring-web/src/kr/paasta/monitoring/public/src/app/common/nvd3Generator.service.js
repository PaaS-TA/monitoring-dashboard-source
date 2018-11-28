'use strict';

angular.module('monitoring')
    .service('nvd3Generator', function () {
        return {
            lineChart: {
                options: function () {
                    return {
                        chart: {
                            type: 'lineChart',
                            margin: {
                                top: 10,
                                right: 30,
                                bottom: 20,
                                left: 20
                            },
                            interpolate: 'linear',
                            isArea: true,
                            legendPosition: 'top',
                            x: function (d) {
                                return d.time;
                            },
                            y: function (d) {
                                return d.usage;
                            },
                            useInteractiveGuideline: true,
                            xAxis: {
                                tickFormat: function(d) {
                                    return d3.time.format('%H:%M:%S')(new Date(d*1000))
                                },
                                ticks: 0
                            },
                            yAxis: {
                                tickFormat: function (d) {
                                    return d3.format('.02f')(d);
                                },
                                axisLabelDistance: -10
                            },
                            color: ['#b9e866','#5cc6d9','#f5da25','#9041d2','#2172e4','#22a581','#ed3131','#9c78df','#ff9936','#2b6e90']
                        }
                    }
                }
            },
            lineWithFocusChart: {
                options: function () {
                    return {
                        chart: {
                            type: 'lineWithFocusChart',
                            height: 185,
                            margin : {
                                top: 10,
                                right: -20,
                                bottom: 20,
                                left: 40
                            },
                            duration: 500,
                            useInteractiveGuideline: true,
                            x: function (d) {
                                return d.time;
                            },
                            y: function (d) {
                                return d.usage;
                            },
                            xAxis: {
                                axisLabel: 'X Axis',
                                tickFormat: function(d){
                                    return d3.format(',f')(d);
                                },
                                ticks: 0
                            },
                            x2Axis: {
                                tickFormat: function(d){
                                    return d3.format(',f')(d);
                                }
                            },
                            yAxis: {
                                axisLabel: 'Y Axis',
                                tickFormat: function(d){
                                    return d3.format(',.2f')(d);
                                },
                                rotateYLabel: false
                            },
                            y2Axis: {
                                tickFormat: function(d){
                                    return d3.format(',.2f')(d);
                                }
                            }

                        }
                    }
                }
            },
            stackedAreaChart: {
                options: function(){
                    return {
                        chart: {
                            type: 'stackedAreaChart',
                            margin : {
                                top: 10,
                                right: -20,
                                bottom: 20,
                                left: 40
                            },
                            interpolate: 'linear',
                            showControls: false,
                            showLegend: true,
                            legendPosition: 'right',
                            x: function (d) {
                                return d.time;
                            },
                            y: function (d) {
                                return d.usage;
                            },
                            useVoronoi: false,
                            clipEdge: true,
                            transitionDuration: 500,
                            useInteractiveGuideline: true,
                            xAxis: {
                                tickFormat: function(d) {
                                    return d3.time.format('%H:%M:%S')(new Date(d*1000))
                                },
                                ticks: 0
                            },
                            yAxis: {
                                tickFormat: function (d) {
                                    return d3.format('.02f')(d);
                                },
                                axisLabelDistance: -10
                            },
                            color: ['#b9e866','#5cc6d9','#f5da25','#9041d2','#2172e4','#22a581','#ed3131','#9c78df','#ff9936','#2b6e90']
                        }
                    }
                }
            },
            multiBarChart: {
                options: function(){
                    return {
                        chart: {
                            type: 'multiBarChart',
                            margin : {
                                top: 10,
                                right: -20,
                                bottom: 20,
                                left: 40
                            },
                            x: function (d) {
                                return d.time;
                            },
                            y: function (d) {
                                return d.usage;
                            },
                            legendPosition: 'right',
                            clipEdge: true,
                            //staggerLabels: true,
                            duration: 500,
                            stacked: false,
                            showControls: false,
                            xAxis: {
                                showMaxMin: true,
                                tickFormat: function(d) {
                                    return d3.time.format('%H:%M:%S')(new Date(d*1000))
                                },
                                ticks: 0
                            },
                            yAxis: {
                                tickFormat: function (d) {
                                    return d3.format('.02f')(d);
                                },
                                axisLabelDistance: -10
                            },
                            color: ['#b9e866','#5cc6d9','#f5da25','#9041d2','#2172e4','#22a581','#ed3131','#9c78df','#ff9936','#2b6e90']
                        }
                    }
                }
            },
            donutChart: {
                options: function(){
                    return {
                        chart: {
                            type: 'pieChart',
                            margin : {
                                top: 0,
                                right: 0,
                                bottom: 0,
                                left: 0
                            },
                            x: function (d) {
                                return d.key;
                            },
                            y: function (d) {
                                return d.value;
                            },
                            legend: {
                                updateState: false
                            },
                            legendPosition: 'right',
                            showLabels: true,
                            showLegend: true,
                            duration: 500,
                            labelThreshold: 0.01,
                            labelSunbeamLayout: false,
                            title: true,
                            valueFormat: function(d) {
                                return d3.format('.0')(d);
                            },
                            cornerRadius: 0,
                            donutRatio: 0.4,
                            donut: true,
                            labelsOutside: false,
                            color: ['#b9e866','#5cc6d9','#f5da25','#9041d2','#2172e4','#22a581','#ed3131','#9c78df','#ff9936','#2b6e90']
                        }
                    }
                }
            },
            discreteBarChart: {
                options: function () {
                    return {
                        chart: {
                            type: 'discreteBarChart',
                            // height: 450,
                            margin : {
                                top: 20,
                                right: 20,
                                bottom: 50,
                                left: 55
                            },
                            x: function(d){
                                return d.key;
                            },
                            y: function(d){
                                return d.value;
                            },
                            showValues: true,
                            valueFormat: function(d){
                              return d3.format('.0')(d);
                            },
                            duration: 500,
                            xAxis: {
                                axisLabel: 'X Axis'
                            },
                            yAxis: {
                                axisLabel: 'Y Axis',
                                axisLabelDistance: -10
                            }
                        }
                    }
                }
            },
            list: {
                options: function () {
                    return {}
                }
            },
            log: {
                options: function () {
                    return {}
                }
            }
        };
    })
;
