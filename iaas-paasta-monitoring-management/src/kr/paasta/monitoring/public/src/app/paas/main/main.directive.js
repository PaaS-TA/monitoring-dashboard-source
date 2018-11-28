(function() {
  'use strict';

  angular
    .module('monitoring')
    .directive('collapsibleTree', collapsibleTree);

  /** @ngInject */
  function collapsibleTree($timeout) {
    var elems = [];
    return {
      restrict: 'E',
      scope: {
        width: '=',
        height: '=',
        radius: '=',
        data: '=',
        collapseDepth: '=',
        expandFn: '&',
        collapseFn: '&'
      },
      link: function (scope, elem) {
        elems.push(elem);
        var margin = {top: 5, right: 10, bottom: 5, left: 40},
          width = scope.width,
          height = scope.height;

        var i = 0,
          duration = 750,
          root;

        var tree = d3.layout.tree()
          .size([height, width]);

        var diagonal = d3.svg.diagonal()
          .projection(function (d) {
            return [d.y, d.x];
          });

        var svg = d3.select(elem[0]).append("svg");
        var vis = svg
          .attr("width", width + margin.right + margin.left)
          .attr("height", height + margin.top + margin.bottom)
          .append("g")
          .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

        var treeZoom = d3.behavior.zoom();
        treeZoom.on('zoom', zoomed);
        d3.select(elem[0]).select('svg').call(treeZoom);

        function zoomed() {
          var zoomTranslate = treeZoom.translate();
          d3.select(elem[0]).select('g').attr('transform', 'translate(' + zoomTranslate[0] + ',' + zoomTranslate[1] + ')');
        }

        var collapseDepth = 0;
        scope.$watch('data', function () {
          if (scope.data) {
            root = angular.copy(scope.data);
            root.x0 = height / 2;
            root.y0 = 0;

            // root.children.forEach(collapse);
            collapseDepth = scope.collapseDepth;
            if (scope.collapseDepth) {
              recursiveCollapse(root);
            } else {
              root.children.forEach(collapse);
            }

            update(root);
          }
        }, true);

        function recursiveCollapse(data) {
          data.children.forEach(function (child) {
            if (collapseDepth > 1) {
              collapseDepth--;
              recursiveCollapse(child);
            }
            if (child.children) {
              child.children.forEach(collapse);
            }
          });
        }

        d3.select(window).on('resize', resize);

        function resize(param) {
          if (param) tree = d3.layout.tree().size([param.height, param.width]);

          $timeout(function () {
            for (var i in elems) {
              var widgetWidth = (param != undefined && param.width != undefined) ? param.width : elems[i].parent().width();
              var svg = d3.select(elems[i][0]).select('svg');
              if (widgetWidth >= 0) {
                svg.attr('width', widgetWidth);
                scope.width = widgetWidth;
              }
            }
          }, 500);
        }

        scope.expandNodes = function (param) {
          if (scope.data) {
            resize(param);
            expand(root);
            update(root, param.depth);
          }
        };
        scope.collapseNodes = function (param) {
          if (scope.data) {
            param.width += (margin.right + margin.left);
            resize(param);
            collapseDepth = param.collapseDepth;
            if (param.collapseDepth) {
              recursiveCollapse(root);
            } else {
              root.children.forEach(collapse);
            }
            update(root, param.depth);
          }
        };
        scope.expandFn({theExpandFn: scope.expandNodes});
        scope.collapseFn({theCollapseFn: scope.collapseNodes});

        function update(source, depth) {
          // Compute the new tree layout.
          var nodes = tree.nodes(root).reverse(),
            links = tree.links(nodes);

          // Normalize for fixed-depth.
          var fixedDepth = depth == undefined ? 130 : depth;
          nodes.forEach(function (d) {
            d.y = d.depth * fixedDepth;
          });

          // Update the nodes
          var node = vis.selectAll("g.node")
            .data(nodes, function (d) {
              return d.id || (d.id = ++i);
            });

          // Enter any new nodes at the parent's previous position.
          var nodeEnter = node.enter().append("g")
            .attr("class", "node")
            .attr("transform", function () {
              return "translate(" + source.y0 + "," + source.x0 + ")";
            })
            .on("click", click);

          nodeEnter.append("circle")
            .attr("r", 1e-6)
            .style("fill", function (d) {
              return d._children ? "#00aacc" : "#fff";
            });

          nodeEnter.append("text")
            .attr("x", function (d) {
              return d.children ? -10 : 10;
            })
            .attr("dy", ".35em")
            .attr("text-anchor", function (d) {
              return d.children ? "end" : "start";
            })
            .text(function (d) {
              return d.name;
            })
            .style("fill", function (d) {
              var color = '#293133';
              if (d.status == 'warning') {
                color = '#f0a141';
              } else if (d.status == 'critical') {
                color = '#ad6de8';
              } else if (d.status == 'fail') {
                color = '#e66b6b';
              }
              return color;
            })
            .style("fill-opacity", 1e-6);

          // Transition nodes to their new position.
          var nodeUpdate = node.transition()
            .duration(duration)
            .attr("transform", function (d) {
              return "translate(" + d.y + "," + d.x + ")";
            });

          nodeUpdate.select("circle")
            .attr("r", 5)
            .style("fill", function (d) {
              var color = d._children ? "#00aacc" : "#fff";
              if (d.status == 'warning') {
                color = '#f0a141';
              } else if (d.status == 'critical') {
                color = '#ad6de8';
              } else if (d.status == 'fail') {
                color = '#e66b6b';
              }
              return color;
            })
            .style("stroke", function (d) {
              var color = '#00aacc';
              if (d.status == 'warning') {
                color = '#f0a141';
              } else if (d.status == 'critical') {
                color = '#ad6de8';
              } else if (d.status == 'fail') {
                color = '#e66b6b';
              }
              return color;
            });

          nodeUpdate.select("text")
            .attr("x", function (d) {
              return d.children ? -10 : 10;
            })
            .attr("text-anchor", function (d) {
              return d.children ? "end" : "start";
            })
            .style("fill-opacity", 0.8);

          // Transition exiting nodes to the parent's new position.
          var nodeExit = node.exit().transition()
            .duration(duration)
            .attr("transform", function () {
              return "translate(" + source.y + "," + source.x + ")";
            })
            .remove();

          nodeExit.select("circle")
            .attr("r", 1e-6);

          nodeExit.select("text")
            .style("fill-opacity", 1e-6);

          // Update the links
          var link = vis.selectAll("path.link")
            .data(links, function (d) {
              return d.target.id;
            });

          // Enter any new links at the parent's previous position.
          link.enter().insert("path", "g")
            .attr("class", "link")
            .attr("d", function () {
              var o = {x: source.x0, y: source.y0};
              return diagonal({source: o, target: o});
            });

          // Transition links to their new position.
          link.transition()
            .duration(duration)
            .attr("d", diagonal);

          // Transition exiting nodes to the parent's new position.
          link.exit().transition()
            .duration(duration)
            .attr("d", function () {
              var o = {x: source.x, y: source.y};
              return diagonal({source: o, target: o});
            })
            .remove();

          // Stash the old positions for transition.
          nodes.forEach(function (d) {
            d.x0 = d.x;
            d.y0 = d.y;
          });
        }

        var collapse = function(d) {
          if (d.children) {
            d._children = d.children;
            d._children.forEach(collapse);
            d.children = null;
          }
        };

        function expand(d) {
          var children = (d.children) ? d.children : d._children;
          if (d._children) {
            d.children = d._children;
            d._children = null;
          }
          if (children)
            children.forEach(expand);
        }

        // Toggle children on click.
        function click(d) {
          if (d.children) {
            d._children = d.children;
            d.children = null;
          } else {
            d.children = d._children;
            d._children = null;
          }
          update(d);
        }

      }
    };
  }

})();

