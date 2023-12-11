<script lang="ts">
    import type {PageData} from './$types';
    import {page} from '$app/stores';
    import {
        type Benchmark,
        convertBenchmarksToLineChartData, getBarChartData, getBarChartDataByVariation,
        getBenchmarkGroups,
        getLineChartOptions,
    } from "$lib/model";
    import {Bar, Line} from "svelte-chartjs";
    import Chart from 'chart.js/auto';

    export let data: PageData;

    const benchmarkSlug = $page.params.benchmark;

    const benchmarkGroup = getBenchmarkGroups().filter(group => group.Name === benchmarkSlug)[0]
    let benchmarks: Benchmark[] = benchmarkGroup.Benchmarks;

    const uniqueVariationNames = benchmarks.flatMap(b => b.Variations).map(v => v.Name).filter((v, i, a) => a.indexOf(v) === i);
    const uniqueVariationNamesString = uniqueVariationNames.map(vn => `<code>${vn}</code>`).slice(0, -1).join(", ") + " and " + uniqueVariationNames.map(vn => `<code>${vn}</code>`).slice(-1);
</script>

<div class="text-red-600 text-8xl">
    <h2>Comparison</h2>
    <p>This site compares the performance of <code>{benchmarkGroup.Name}</code> implementations with the
        functions: {@html uniqueVariationNamesString}</p>
    <Bar options="{getLineChartOptions()}" data="{getBarChartData(benchmarks)}"/>
    <h3>Comparisons per function</h3>
    {#each uniqueVariationNames as vari}
        <details>
            <summary>{vari}</summary>
            <Bar options="{getLineChartOptions()}" data="{getBarChartDataByVariation(benchmarks, vari)}"/>
        </details>
    {/each}

    <h2>Implementations</h2>
    <aside>
        <nav>
            <ul>
                {#each benchmarks as benchmark}
                    <li><a href="#{benchmark.Name.replaceAll(' ', '-').toLowerCase()}">{benchmark.Name}</a></li>
                {/each}
            </ul>
        </nav>
    </aside>

    {#each benchmarks as benchmark}
        <section id="{benchmark.Name.replaceAll(' ', '-').toLowerCase()}">
            <h3>{benchmark.Name}</h3>
            <Line height={50} options="{getLineChartOptions()}" data="{convertBenchmarksToLineChartData([benchmark])}"/>
        </section>
    {/each}
</div>

<style>
    aside {
        @apply text-center;
    }
</style>
