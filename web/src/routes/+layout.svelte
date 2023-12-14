<script>
  import "../app.scss";
  import NavBar from "$lib/components/NavBar.svelte";
  import Footer from "$lib/components/Footer.svelte";
</script>

<NavBar/>
<main class="min-h-screen">
    <slot></slot>
</main>
<Footer/>

<style lang="scss">
    main {
        @apply mt-28 mx-auto sm:w-3/4 max-w-7xl;
    }
</style>
