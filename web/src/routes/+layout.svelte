<script>
  import NavBar from "$lib/components/NavBar.svelte";
  import {pageTitle} from "$lib/store";
</script>

<svelte:head>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@1/css/pico.min.css">
</svelte:head>

<NavBar/>
<main class="container">
    <header>
        <hgroup>
            <h1>Go {$pageTitle} Benchmark</h1>
            <p><a href="//github.com/go-benchmarks/benchmarks/tree/main/benchmarks/{$pageTitle}">github.com/go-benchmarks/benchmarks/tree/main/benchmarks/{$pageTitle}</a></p>
        </hgroup>
    </header>
    <br/>
        <slot></slot>
</main>

<style>
    main {
        padding-top: 80px;
        padding-bottom: 80px;
    }
</style>
