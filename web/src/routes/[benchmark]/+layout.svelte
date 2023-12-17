<script>
  import {pageTitle} from "$lib/store";
</script>


<hgroup>
    <h1>Go {$pageTitle} Benchmark</h1>
    <p><a target="_blank" class="link text-blue-400" href="//github.com/go-benchmarks/benchmarks/tree/main/benchmarks/{$pageTitle.toLowerCase()}">github.com/go-benchmarks/benchmarks/tree/main/benchmarks/{$pageTitle.toLowerCase()}</a>
    </p>
</hgroup>
<br/>
<slot></slot>

