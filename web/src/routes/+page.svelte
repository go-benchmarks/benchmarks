<script lang="ts">
    import {getBenchmarkGroups} from "$lib/model";

    const benchmarkGroups = getBenchmarkGroups()
</script>

<svelte:head>
    <title>Go Benchmarks</title>
</svelte:head>

<section class="text-center mb-32">
    <hgroup>
        <img class="inline" alt="Go Benchmarks Logo" width="200" src="/favicon.png"/>
        <h1>Go Benchmarks</h1>
        <p>Visual comparison of different Go benchmarks</p>
    </hgroup>
</section>
<section>
    <div class="grid grid-cols-3 gap-8">
        {#each benchmarkGroups as group}
            <a href="/{group.Name.toLowerCase().replaceAll(' ', '-')}"
               class="card card-compact bg-base-200 bg-opacity-30 border-1 border-black card-bordered shadow-md">
                <div class="card-body">
                    <h2 class="card-title"><a href="/{group.Name.toLowerCase().replaceAll(' ', '-')}">{group.Name}</a>
                    </h2>
                    <div>
                        <p>{group.Headline}</p>
                        {#if group.Benchmarks.length > 1}
                            <h3 class="header">Implementations:</h3>
                            <ul>
                                {#each group.Benchmarks as benchmark}
                                    <li>{benchmark.Name}</li>
                                {/each}
                            </ul>
                        {/if}
                        {#if Array.from(new Set(group.Benchmarks[0].Variations.map(v => v.Name))).length > 1}
                            <h3 class="header">Functions:</h3>
                            <ul>
                                {#each Array.from(new Set(group.Benchmarks[0].Variations.map(v => v.Name))) as v}
                                    <li>{v.toLowerCase()}</li>
                                {/each}
                            </ul>
                        {/if}
                    </div>
                </div>
            </a>
        {/each}
    </div>
</section>

<style lang="scss">
  hgroup {
    @apply mb-8;
  }

  .header {
    @apply text-sm mb-0 font-medium;
  }

  .card-title {
    @apply mt-0 font-bold;
  }
</style>
