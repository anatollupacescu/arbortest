<script>
import { store, current }from './store.js';
let records = []

let len = 0

store.subscribe(value => {
    records = value;
    len = records.length
});

function selectGraph(index) {
    current.update(() => records[index])
}
</script>

{#each records as { commit, message }, i}
<div class="card mb-4 box-shadow" style="cursor: pointer;" on:click={() => selectGraph(i)}>
  <div class="card-header">
    <h6 class="my-0 font-weight-normal">{len - i}: {commit}</h6>
  </div>
  <div class="card-body">{message}</div>
</div>
{/each}
