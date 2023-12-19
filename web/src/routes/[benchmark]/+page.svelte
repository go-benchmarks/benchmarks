<script lang="ts">
    import {page} from '$app/stores';
    import {
        type Benchmark,
        getBarChartDataByVariation,
        getBarChartDataByVariationAndRunCount,
        getBenchmarkGroups,
        getLineChartOptions,
    } from "$lib/model";
    import {onMount} from "svelte";
    import hljs from "highlight.js/lib/core";
    import * as goLang from "highlight.js/lib/languages/go";
    import {Bar} from "svelte-chartjs";
    import {Chart as ChartJS} from 'chart.js';
    import "chart.js/auto"
    import annotationPlugin from 'chartjs-plugin-annotation';
    import Comparison from "$lib/components/Comparison.svelte";
    import ImplementationsToC from "$lib/components/ImplementationsToC.svelte";
    import ImplementationAnalysis from "$lib/components/ImplementationAnalysis.svelte";

    ChartJS.register(annotationPlugin)

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
</script>

<svelte:head>
    <title>Go {benchmarkGroup.Name} Benchmark</title>
    <link rel="stylesheet"
          href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/atom-one-dark.min.css"/>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/go.min.js"></script>
</svelte:head>

<div class="px-8 sm:px-0">
    <section>
        <hgroup>
            <h1>Go {@html benchmarkGroup.Name} Benchmark</h1>
            <p>
                {@html benchmarkGroup.Headline}
            </p>
        </hgroup>

        <b>Description:</b>
        <p class="small-width">{@html benchmarkGroup.Description}</p>

        <ImplementationsToC {benchmarkGroup}/>
    </section>

    <section>
        <Comparison {benchmarkGroup}/>
    </section>

    {#if uniqueVariationNames.length > 1}
        <section>
            <hgroup>
                <h3>Different Run Count Per Function</h3>
                <p>Some algorithms get slower or faster over time. Especially when they populate the same data
                    struct.</p>
            </hgroup>
            {#each uniqueVariationNames as vari}
                <h4><code>{vari}</code></h4>
                <Bar height="{50}" options="{getLineChartOptions(true, false)}"
                     data="{getBarChartDataByVariationAndRunCount(benchmarks, vari)}"/>
            {/each}
        </section>

        <section>
            <hgroup>
                <h3>Different CPU Core Count Per Function</h3>
                <p>Some algorithms are influenced by the amount of cores a CPU has.</p>
            </hgroup>
            {#each uniqueVariationNames as vari}
                <h4><code>{vari}</code></h4>
                <Bar height="{50}" options="{getLineChartOptions(true, false)}"
                     data="{getBarChartDataByVariation(benchmarks, vari)}"/>
            {/each}
        </section>
    {/if}


    {#each benchmarks as benchmark}
        <section id="{benchmark.Name.replaceAll(' ', '-').toLowerCase()}">
            <ImplementationAnalysis {benchmarkGroup} {benchmark}/>
        </section>
    {/each}

    <section>
        <h2>Full Benchmark Code</h2>
        <pre><code class="language-go">{benchmarkGroup.Code}</code></pre>
    </section>

    <section>
        <h2>Full Benchmark Output</h2>
        <div class="overflow-x-auto max-h-96">
            <table class="table table-xs">
                <thead>
                <tr>
                    <th>Implementation</th>
                    <th>Function</th>
                    <th>Runs</th>
                    <th>CPU Core Count</th>
                    <th>ns/op</th>
                    <th>ops/sec</th>
                    <th>B/op</th>
                    <th>allocs/op</th>
                    <th>MB/s</th>
                </tr>
                </thead>
                <tbody>
                {#each benchmarks as benchmark}
                    {#each benchmark.Variations as variation}
                        <tr>
                            <td>{benchmark.Name}</td>
                            <td>{variation.Name}</td>
                            <td>{variation.N}</td>
                            <td>{variation.CPUCount}</td>
                            <td>{variation.NsPerOp}</td>
                            <td>{variation.OpsPerSec}</td>
                            <td>{variation.AllocedBytesPerOp}</td>
                            <td>{variation.AllocsPerOp}</td>
                            <td>{variation.MBPerS}</td>
                        </tr>
                    {/each}
                {/each}
                </tbody>
            </table>
        </div>
    </section>
</div>

<style lang="scss">
  .code-area {
    @apply my-6;
  }

  section {
    @apply mb-40;
  }

  .small-width {
    @apply max-w-3xl;
  }
</style>
