<script>
  import NavBar from "$lib/components/NavBar.svelte";
</script>

<svelte:head>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@1/css/pico.min.css">
</svelte:head>

<NavBar/>
<main class="container">
    <slot></slot>
</main>

<style>
    main {
        padding-top: 80px;
        padding-bottom: 80px;
    }
</style>
