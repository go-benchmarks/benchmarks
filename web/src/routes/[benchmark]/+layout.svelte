<script>
  import {pageTitle} from "$lib/store";
  let slug = $pageTitle.toLowerCase().replaceAll(" ", "-");

  pageTitle.subscribe(value => {
    slug = value.toLowerCase().replaceAll(" ", "-");
  });
</script>

<svelte:head>
    <title>Go {$pageTitle} Benchmarks</title>
</svelte:head>


<hgroup>
    <h1>Go {$pageTitle} Benchmark</h1>
    <p><a target="_blank" class="link text-blue-400" href="//github.com/go-benchmarks/benchmarks/tree/main/benchmarks/{slug}">github.com/go-benchmarks/benchmarks/tree/main/benchmarks/{slug}</a>
    </p>
</hgroup>
<br/>
<slot></slot>

