<script lang="ts">
    import {page} from '$app/stores';
    import {
        type Benchmark,
        type BenchmarkGroup,
        convertBenchmarksToCPUCountPerformanceLineChart,
        convertBenchmarksToRunCountPerformanceLineChart,
        filterBenchmarkByVariationName,
        getBarChartDataByCPUCountMulti,
        getBarChartDataByRunsMulti,
        getBarChartDataByVariation,
        getBarChartDataByVariationAndRunCount,
        getBenchmarkGroups, getChartDataByCPUCount, getChartDataByRuns,
        getLineChartOptions,
    } from "$lib/model";
    import {Bar, Line} from "svelte-chartjs";
    import {pageTitle} from "$lib/store";
    import Chart from "chart.js/auto";

    const benchmarkSlug = $page.params.benchmark;
    const benchmarkGroup = getBenchmarkGroups().filter(group => group.Name === benchmarkSlug)[0];

    if (!benchmarkGroup) {
        throw new Error(`Benchmark Group ${benchmarkSlug} not found`);
    }

    $pageTitle = benchmarkGroup.Name;

    let benchmarks: Benchmark[] = benchmarkGroup.Benchmarks;

    const uniqueVariationNames = benchmarks.flatMap(b => b.Variations).map(v => v.Name).filter((v, i, a) => a.indexOf(v) === i);
    const uniqueVariationNamesString = uniqueVariationNames.map(vn => `<code>${vn}</code>`).slice(0, -1).join(", ") + " and " + uniqueVariationNames.map(vn => `<code>${vn}</code>`).slice(-1);
</script>

<div>
    <hgroup>
        <h2>Comparison</h2>
        <p>This site compares the performance of <code>{benchmarkGroup.Name}</code> implementations with the
            functions: {@html uniqueVariationNamesString}</p>
    </hgroup>
    <h3>Different Run Count</h3>
    <Bar height="{50}" options="{getLineChartOptions(true)}" data="{getBarChartDataByRunsMulti(benchmarks)}"/>
    <h3>Different CPU Core Count</h3>
    <Bar height="{50}" options="{getLineChartOptions(true)}" data="{getBarChartDataByCPUCountMulti(benchmarks)}"/>

    <hgroup>
        <h3>Different Run Count</h3>
        <p>Some algorithms get slower or faster over time. Especially when they populate the same data struct.</p>
    </hgroup>
    {#each uniqueVariationNames as vari}
        <h4>{vari}</h4>
        <Bar height="{50}" options="{getLineChartOptions(true)}" data="{getBarChartDataByVariationAndRunCount(benchmarks, vari)}"/>
    {/each}


    <hgroup>
        <h3>Different CPU Core Count</h3>
        <p>Some algorithms are influenced by the amount of cores a CPU has.</p>
    </hgroup>
    {#each uniqueVariationNames as vari}
        <h4>{vari}</h4>
        <Bar height="{50}" options="{getLineChartOptions(true)}" data="{getBarChartDataByVariation(benchmarks, vari)}"/>
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

            {#each uniqueVariationNames as vari}
                <h4>{vari}</h4>
                <h5>By Run Count</h5>
                <Line height="{50}" options="{getLineChartOptions(false)}" data="{getChartDataByRuns(filterBenchmarkByVariationName(benchmark, vari))}"/>
                <h5>By CPU Core Count</h5>
                <Line height="{50}" options="{getLineChartOptions(false)}" data="{getChartDataByCPUCount(filterBenchmarkByVariationName(benchmark, vari))}"/>
            {/each}
        </section>
    {/each}
</div>

<style>
    aside {
        @apply text-center;
    }
</style>
