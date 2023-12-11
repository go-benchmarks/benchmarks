<script>
  import NavBar from "$lib/components/NavBar.svelte";
</script>

<svelte:head>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@1/css/pico.min.css">
</svelte:head>

<NavBar />
<main class="container">
    <header>
        <hgroup>
            <h1>Go Benchmarks</h1>
            <p>subtitle</p>
        </hgroup>
    </header>
        <slot></slot>
</main>

<style>
    main {
        padding-top: 80px;
        padding-bottom: 80px;
    }

    article {
        padding-top: calc(1.5 * 1rem);
        padding-bottom: calc(1.5 * 1rem);
    }
</style>
