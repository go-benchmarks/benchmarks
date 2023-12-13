<script lang="ts">
    import {getBenchmarkGroups} from "$lib/model";

    const benchmarkGroups = getBenchmarkGroups()
</script>

<svelte:head>
    <title>Go Benchmarks</title>
</svelte:head>

<section>
    <div class="grid grid-cols-2">
        {#each benchmarkGroups as group}
            <div class="card">
                <h2><a href="/{group.Name}">{group.Name}</a></h2>
                <h3>Functions:</h3>
                <ul>
                    {#each Array.from(new Set(group.Benchmarks[0].Variations.map(v => v.Name))) as v}
                        <li>{v}</li>
                    {/each}
                </ul>
                <h3>Implementations:</h3>
                <ul>
                    {#each group.Benchmarks as benchmark}
                        <li>{benchmark.Name}</li>
                    {/each}
                </ul>
            </div>
        {/each}
    </div>
</section>

<style lang="scss">
  .card {
    @apply bg-gray-800 rounded-lg p-4;

    h2 {
      @apply text-xl;
    }

    h3 {
      @apply font-bold text-sm pb-0 mb-2;
    }
  }
</style>
