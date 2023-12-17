<script>
  import {pageTitle} from "$lib/store";
  var slug = $pageTitle.toLowerCase().replaceAll(" ", "-");
</script>


<hgroup>
    <h1>Go {$pageTitle} Benchmark</h1>
    <p><a target="_blank" class="link text-blue-400" href="//github.com/go-benchmarks/benchmarks/tree/main/benchmarks/{slug}">github.com/go-benchmarks/benchmarks/tree/main/benchmarks/{slug}</a>
    </p>
</hgroup>
<br/>
<slot></slot>

