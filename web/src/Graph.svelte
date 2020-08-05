<script>
  import { onMount } from "svelte";
  import { scaleLinear, scaleOrdinal } from "d3-scale";
  import { zoom, zoomIdentity } from "d3-zoom";
  import { schemeCategory10 } from "d3-scale-chromatic";
  import { select, selectAll } from "d3-selection";
  import { drag } from "d3-drag";
  import { store, current }from './store.js';
  import {
    forceSimulation,
    forceLink,
    forceManyBody,
    forceCenter,
  } from "d3-force";

  import { event as currentEvent } from "d3-selection"; // Needed to get drag working, see: https://github.com/d3/d3/issues/2733
  let d3 = {
    zoom,
    zoomIdentity,
    scaleLinear,
    scaleOrdinal,
    schemeCategory10,
    select,
    selectAll,
    drag,
    forceSimulation,
    forceLink,
    forceManyBody,
    forceCenter,
  };

  const colors = {
    "skip": "grey",
    "fail": "red",
    "pass": "green"
  }

  let graph;

  let simulation, svg;
  let width, height;

  let source = new EventSource("http://localhost:3000/events/");

  source.onmessage = function (e) {
    graph = JSON.parse(e.data);
		store.update(graphs => [ graph, ...graphs])
    
    d3.selectAll("svg > *").remove();
    let links = graph.links.map((d) => Object.create(d));
    let nodes = graph.nodes.map((d) => Object.create(d));
    renderGraph(nodes, links);
  };

  current.subscribe(graph => {
    if (graph.links === undefined) {
      return
    }
    refreshGraph(graph)
  });

  function refreshGraph(graph){
    d3.selectAll("svg > *").remove();
    let links = graph.links.map((d) => Object.create(d));
    let nodes = graph.nodes.map((d) => Object.create(d));
    renderGraph(nodes, links);
  }

  onMount(() => {
    height = width = window.innerHeight;
    // width = document.querySelector(".chart").clientWidth;

    svg = d3
      .select(".chart")
      .append("svg")
      .attr("width", width)
      .attr("height", height);
  });

  function renderGraph(nodes, links) {
    simulation = d3
      .forceSimulation(nodes)
      .force(
        "link",
        d3
          .forceLink(links)
          .id((d) => d.id)
          .distance(function () {
            return height / Math.sqrt(nodes.length);
          })
          .strength(1)
      )
      .force("charge", d3.forceManyBody())
      .force("center", d3.forceCenter(width / 2, height / 2))
      .on("tick", simulationUpdate);

    const g = svg.append("g");

    svg
      .append("defs")
      .selectAll("marker")
      .data(["suit", "licensing", "resolved"])
      .enter()
      .append("marker")
      .attr("id", function (d) {
        return d;
      })
      .attr("viewBox", "0 -5 10 10")
      .attr("refX", 25)
      .attr("refY", 0)
      .attr("markerWidth", 6)
      .attr("markerHeight", 6)
      .attr("orient", "auto")
      .append("path")
      .attr("d", "M0,-5L10,0L0,5 L10,0 L0, -5")
      .style("stroke", "#4679BD")
      .style("opacity", "0.6");

    const link = g
      .append("g")
      .attr("stroke", "#999")
      .attr("stroke-opacity", 0.6)
      .selectAll("line")
      .data(links)
      .join("line")
      .attr("stroke-width", (d) => Math.sqrt(d.value))
      .style("marker-end", "url(#suit)");

    const node = g
      .append("g")
      .attr("stroke", "#fff")
      .attr("stroke-width", 1.5)
      .selectAll("circle")
      .data(nodes)
      .join("circle")
      .attr("r", 15)
      .attr("fill", (d) => colors[d.status])
      .call(
        d3
          .drag()
          .on("start", dragstarted)
          .on("drag", dragged)
          .on("end", dragended)
      );

    var gnode = svg
      .append("g")
      .attr("class", "nodes")
      .selectAll("g")
      .data(nodes)
      .enter()
      .append("g");

    gnode.append("text").text(function (d) {
      return d.id;
    });

    function simulationUpdate() {
      link
        .attr("x1", (d) => d.source.x)
        .attr("y1", (d) => d.source.y)
        .attr("x2", (d) => d.target.x)
        .attr("y2", (d) => d.target.y);

      node.attr("cx", (d) => d.x).attr("cy", (d) => d.y);

      d3.selectAll("text")
        .attr("x", function (d) {
          return d.x + 15;
        })
        .attr("y", function (d) {
          return d.y;
        });
    }
  }

  function dragstarted() {
    if (!currentEvent.active) simulation.alphaTarget(0.3).restart();
    currentEvent.subject.fx = currentEvent.x;
    currentEvent.subject.fy = currentEvent.y;
  }

  function dragged() {
    currentEvent.subject.fx = currentEvent.x;
    currentEvent.subject.fy = currentEvent.y;
  }

  function dragended() {
    if (!currentEvent.active) simulation.alphaTarget(0);
    currentEvent.subject.fx = null;
    currentEvent.subject.fy = null;
  }
</script>

<div class="chartdiv" />
