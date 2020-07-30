// set a width and height for our SVG
var width = window.screen.width;
var height = window.screen.height;

var links =[
  {source: "Create", target: "inventory"},
  {source: "CreateWithEmptyName", target: "inventory"},
  {source: "Duplicate", target: "inventory"},
  {source: "Disable", target: "inventory"},
  {source: "GetAll", target: "inventory"},
  {source: "GetOne", target: "inventory"},
  {source: "CreateRecipe", target: "recipe"},
  {source: "CreateRecipeNoName", target: "recipe"},
  {source: "CreateRecipeNoItems", target: "recipe"},
  {source: "GetRecipe", target: "recipe"},
  {source: "GetAllRecipes", target: "recipe"},
  {source: "DisableRecipe", target: "recipe"},
  {source: "Provision", target: "stock"},
  {source: "GetStockPos", target: "stock"},
  {source: "GetAllStockPos", target: "stock"},
  {source: "CreateOrderOK", target: "order"},
  {source: "CreateOrderWhenNotEnoughStock", target: "order"},
  {source: "recipe-ext", target: "recipe"},
  {source: "inventory-ext", target: "inventory"},
  {source: "recipe-ext", target: "inventory-ext"},
  {source: "stock-ext", target: "stock"},
  {source: "stock-ext", target: "inventory-ext"},
  {source: "order-ext", target: "order"},
  {source: "order-ext", target: "recipe-ext"}
]

// create empty nodes array
var nodes = {};

// compute nodes from links data
links.forEach(function (link) {
  link.source =
    nodes[link.source] || (nodes[link.source] = { name: link.source });
  link.target =
    nodes[link.target] || (nodes[link.target] = { name: link.target });
});

// add a SVG to the body for our viz
var svg = d3
  .select("body")
  .append("svg")
  .attr("width", width)
  .attr("height", height);

// use the force
var force = d3.layout
  .force()
  .charge(-1200)
  .size([width, height])
  .nodes(d3.values(nodes))
  .links(links)
  .linkDistance(function (n) {
    return 150;
  })
  .start();

  svg.append("defs").selectAll("marker")
    .data(["suit", "licensing", "resolved"])
  .enter().append("marker")
    .attr("id", function(d) { return d; })
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

// add links
var link = svg
  .selectAll(".link")
  .data(links)
  .enter()
  .append("line")
  .attr("class", "link") 
  .style("marker-end",  "url(#suit)");

  var node = svg.selectAll(".node")
    .data(force.nodes())
    .enter().append("g")
    .attr("class", "node")
    .call(force.drag);

  node.append("circle")
    .attr("r", function(e) { 
      return width*0.01 })
    .style("stroke", "orange")
    .style("fill", function (d) {
        return "green";
    })

  node.append("text")
      .attr("dx", 20)
      .attr("dy", ".35em")
      .text(function(d) { return d.name })
      .style("stroke", "gray");

    force.on("tick", function (e) {
      link.attr("x1", function (d) {
          return d.source.x;
      }).attr("y1", function (d) {
          return d.source.y;
      }).attr("x2", function (d) {
          return d.target.x;
      }).attr("y2", function (d) {
          return d.target.y;
      });

      d3.selectAll("circle").attr("cx", function (d) {
          return d.x;
      }).attr("cy", function (d) {
        return d.y;
      });
      
      d3.selectAll("text").attr("x", function (d) {
          return d.x;
      }).attr("y", function (d) {
          return d.y;
      });
});
