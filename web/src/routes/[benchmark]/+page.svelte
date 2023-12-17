<script lang="ts">
    import {page} from '$app/stores';
    import {
        type Benchmark,
        filterBenchmarkByVariationName,
        getBarChartDataByCPUCountMulti,
        getBarChartDataByRunsMulti,
        getBarChartDataByVariation,
        getBarChartDataByVariationAndRunCount,
        getBenchmarkGroups,
        getChartDataByCPUCount,
        getChartDataByRuns,
        getLineChartOptions,
    } from "$lib/model";
    import {Bar, Line} from "svelte-chartjs";
    import "chart.js/auto";
    import {onMount} from "svelte";
    import hljs from "highlight.js/lib/core";
    import * as goLang from "highlight.js/lib/languages/go";

    onMount(() => {
        hljs.registerLanguage('go', goLang.default);
        hljs.highlightAll();
    });

    const benchmarkSlug = $page.params.benchmark;
    const benchmarkGroup = getBenchmarkGroups().filter(group => group.Name.toLowerCase().replaceAll(" ", "-") === benchmarkSlug)[0];

    if (!benchmarkGroup) {
        throw new Error(`Benchmark Group ${benchmarkSlug} not found`);
    }

    let benchmarks: Benchmark[] = benchmarkGroup.Benchmarks;

    const uniqueVariationNames = benchmarks.flatMap(b => b.Variations).map(v => v.Name).filter((v, i, a) => a.indexOf(v) === i);


    let uniqueVariationNamesString: string;
    if (uniqueVariationNames.length === 1) {
        uniqueVariationNamesString = `<code>${uniqueVariationNames[0]}</code>`;
    } else {
        uniqueVariationNamesString = uniqueVariationNames.map(vn => `<code>${vn}</code>`).slice(0, -1).join(", ") + " and " + uniqueVariationNames.map(vn => `<code>${vn}</code>`).slice(-1);
    }
</script>

<svelte:head>
    <title>Go {benchmarkGroup.Name} Benchmarks</title>
    <link rel="stylesheet"
          href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/atom-one-dark.min.css"/>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/go.min.js"></script>
</svelte:head>

<div>
    <hgroup>
        <h1>Go {benchmarkGroup.Name} Benchmark</h1>
        <p><a target="_blank" class="link text-blue-400" href="//github.com/go-benchmarks/benchmarks/tree/main/benchmarks/{benchmarkGroup.Name.toLowerCase().replaceAll(' ', '-')}">github.com/go-benchmarks/benchmarks/tree/main/benchmarks/{benchmarkGroup.Name.toLowerCase().replaceAll(' ', '-')}</a>
        </p>
    </hgroup>
    <br/>
    <b>Description:</b>
    <p class="description">{benchmarkGroup.Description}</p>

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
        <Bar height="{50}" options="{getLineChartOptions(true)}"
             data="{getBarChartDataByVariationAndRunCount(benchmarks, vari)}"/>
    {/each}


    <hgroup>
        <h3>Different CPU Core Count</h3>
        <p>Some algorithms are influenced by the amount of cores a CPU has.</p>
    </hgroup>
    {#each uniqueVariationNames as vari}
        <h4>{vari}</h4>
        <Bar height="{50}" options="{getLineChartOptions(true)}" data="{getBarChartDataByVariation(benchmarks, vari)}"/>
    {/each}


    <h2 class="mt-64">Implementations</h2>
    <nav class="mb-32">
        <ul>
            {#each benchmarks as benchmark}
                <li><a class="link-hover text-blue-400 text-xl"
                       href="#{benchmark.Name.replaceAll(' ', '-').toLowerCase()}">{benchmark.Name}</a></li>
            {/each}
        </ul>
    </nav>

    {#each benchmarks as benchmark}
        <section id="{benchmark.Name.replaceAll(' ', '-').toLowerCase()}">
            <h3 class="underline"><a href="#{benchmark.Name.replaceAll(' ', '-').toLowerCase()}">{benchmark.Name}</a>
            </h3>
            <p class="description">{benchmark.Description}</p>
            <b>Code:</b>
            <pre><code class="language-go">{benchmarkGroup.Constants}{benchmark.Code}</code></pre>

            {#each uniqueVariationNames as vari}
                <h4>{vari}</h4>
                <h5>By Run Count</h5>
                <Line height="{50}" options="{getLineChartOptions(false)}"
                      data="{getChartDataByRuns(filterBenchmarkByVariationName(benchmark, vari))}"/>
                <h5>By CPU Core Count</h5>
                <Line height="{50}" options="{getLineChartOptions(false)}"
                      data="{getChartDataByCPUCount(filterBenchmarkByVariationName(benchmark, vari))}"/>
            {/each}
        </section>
    {/each}

    <h2>Full Benchmark Code</h2>
    <pre><code class="language-go">{benchmarkGroup.Code}</code></pre>
</div>

<style lang="scss">
  .description {
    @apply max-w-3xl;
  }
</style>
