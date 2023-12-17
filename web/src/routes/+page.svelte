<script lang="ts">
    import {getBenchmarkGroups} from "$lib/model";

    const benchmarkGroups = getBenchmarkGroups()
</script>

<svelte:head>
    <title>Go Benchmarks</title>
</svelte:head>

<section>
    <div class="flex flex-wrap gap-8">
        {#each benchmarkGroups as group}
            <a href="/{group.Name.toLowerCase().replaceAll(' ', '-')}"
               class="card w-1/4 flex-grow card-compact bg-base-200 border-1 border-black card-bordered shadow-md">
                <div class="card-body">
                    <h2 class="card-title"><a href="/{group.Name.toLowerCase().replaceAll(' ', '-')}">{group.Name}</a></h2>
                    <h3>Functions:</h3>
                    <ul>
                        {#each Array.from(new Set(group.Benchmarks[0].Variations.map(v => v.Name))) as v}
                            <li>{v}</li>
                        {/each}
                    </ul>
                    {#if group.Benchmarks.length > 1}
                        <h3>Implementations:</h3>
                        <ul>
                            {#each group.Benchmarks as benchmark}
                                <li>{benchmark.Name}</li>
                            {/each}
                        </ul>
                    {/if}
                </div>
            </a>
        {/each}
    </div>
</section>
