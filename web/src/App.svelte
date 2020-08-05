<svelte:head>
	<link rel="stylesheet" href="font-awesome.min.css">
	<link rel="stylesheet" href="wireframe.css">
</svelte:head>

<style>
.chart {
	width: 100%;
	height: calc(100% - 4em);
	min-height: 280px;
	margin: 0 auto;
}
</style>

<script>
	import Graph from './Graph.svelte';
	import Sidepanel from './Sidepanel.svelte';
</script>

<div>
	<div class="container-fluid h-100">
	  <div class="row h-100">
		<div class="col-md-10">
			<div class="chart">
				<Graph />
			</div>
		</div>
		<div class="col-md-2">
		  <div class="col-md-12" >
				<Sidepanel />
		  </div>
		</div>
	  </div>
	</div>
  </div>