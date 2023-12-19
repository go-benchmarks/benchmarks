<script lang="ts">
    import VariationAnalysis from "$lib/components/VariationAnalysis.svelte";
    import type {Benchmark, BenchmarkGroup, Variation} from "$lib/model";
    import {onMount} from "svelte";
    import hljs from "highlight.js/lib/core";
    import * as goLang from "highlight.js/lib/languages/go";

    export let benchmarkGroup: BenchmarkGroup;
    export let benchmark: Benchmark;
    export let benchmarks = benchmarkGroup.Benchmarks


    const uniqueVariationNames = benchmarks.flatMap(b => b.Variations).map(v => v.Name).filter((v, i, a) => a.indexOf(v) === i);

    onMount(() => {
        hljs.registerLanguage('go', goLang.default);
        hljs.highlightAll();
    });

    const averagedBenchmarks = benchmarkGroup.Benchmarks.map(b => {
        let averagedVariations: any = {};

        uniqueVariationNames.forEach(variationName => {
            const variations = b.Variations.filter(v => v.Name === variationName);
            const averageVariation = variations.reduce((acc, v) => {
                acc.NsPerOp += v.NsPerOp;
                acc.AllocedBytesPerOp += v.AllocedBytesPerOp;
                acc.AllocsPerOp += v.AllocsPerOp;
                acc.MBPerS += v.MBPerS;
                acc.OpsPerSec += v.OpsPerSec;
                return acc;
            }, {
                Name: variationName,
                NsPerOp: 0,
                AllocedBytesPerOp: 0,
                AllocsPerOp: 0,
                MBPerS: 0,
                OpsPerSec: 0,
            } as Variation);

            averageVariation.NsPerOp /= variations.length;
            averageVariation.AllocedBytesPerOp /= variations.length;
            averageVariation.AllocsPerOp /= variations.length;
            averageVariation.MBPerS /= variations.length;
            averageVariation.OpsPerSec /= variations.length;

            averagedVariations[variationName] = averageVariation;
        });

        return {
            Benchmark: b,
            Variations: averagedVariations
        };
    });

    const averageBenchmark = averagedBenchmarks.filter(b => b.Benchmark.Name === benchmark.Name)[0];
    function calculatePercentageDifference(current, other) {
        if (current <= 0 || other <= 0) return 0; // Avoid division by zero
        return ((current - other) / other) * 100;
    }

    const comparisonResults = [];

    uniqueVariationNames.forEach(variationName => {
        averagedBenchmarks.forEach(b => {
            if (b.Benchmark.Name !== benchmark.Name) {
                const currentVariation = averageBenchmark.Variations[variationName];
                const comparisonVariation = b.Variations[variationName];

                if (currentVariation && comparisonVariation) {
                    const percentageDifference = calculatePercentageDifference(currentVariation.NsPerOp, comparisonVariation.NsPerOp);
                    let comparisonStatement;

                    if (percentageDifference > 0) {
                        // Current is slower
                        comparisonStatement = `<span class="font-bold">${benchmark.Name}</span> <code>${variationName}</code> is <span class="text-red-500">${percentageDifference.toFixed(2)}% slower</span> than <span class="font-bold">${b.Benchmark.Name}</span> <code>${variationName}</code>.`;
                    } else {
                        // Current is faster
                        comparisonStatement = `<span class="font-bold">${benchmark.Name}</span> <code>${variationName}</code> is <span class="text-green-500">${Math.abs(percentageDifference).toFixed(2)}% faster</span> than <span class="font-bold">${b.Benchmark.Name}</span> <code>${variationName}</code>.`;
                    }

                    comparisonResults.push(comparisonStatement);
                }
            }
        });
    });



</script>

<section>
    <h3 class="underline"><a
            href="#{benchmark.Name.replaceAll(' ', '-').toLowerCase()}">{benchmark.Name}</a>
    </h3>
    <p class="small-width">{@html benchmark.Description}</p>
    <table class="table">
        <thead>
        <tr>
            <th>Function</th>
            <th>ns/op</th>
            <th>ops/sec</th>
            <th>B/op</th>
            <th>allocs/op</th>
            <th>MB/s</th>
        </tr>
        </thead>
        <tbody>
        {#each Object.values(averageBenchmark.Variations) as variation}
            <tr>
                <td>{variation.Name}</td>
                <td>{variation.NsPerOp.toFixed(2)}</td>
                <td>{variation.OpsPerSec.toFixed(0)}</td>
                <td>{variation.AllocedBytesPerOp.toFixed(1)}</td>
                <td>{variation.AllocsPerOp.toFixed(0)}</td>
                <td>{variation.MBPerS.toFixed(1)}</td>
            </tr>
        {/each}
        </tbody>
    </table>

    <h4>Comparison</h4>
    <ul>
        {#each comparisonResults as comparisonStatement}
            <li>{@html comparisonStatement}</li>
        {/each}
    </ul>

    {#if benchmark.Code}
        <pre class="code-area"><code class="language-go">{benchmarkGroup.Constants}{benchmark.Code}</code></pre>
        <div tabindex="-1" class="collapse collapse-plus bg-[#282c34]">
            <input type="checkbox"/>
            <div class="collapse-title text-xl font-medium">
                Benchmark Code
            </div>
            <div class="collapse-content">
                <pre class="code-area"><code
                        class="language-go">{benchmarkGroup.Constants}{benchmark.BenchmarkCode}</code></pre>
            </div>
        </div>
    {:else}
        <pre class="code-area"><code
                class="language-go">{benchmarkGroup.Constants}{benchmark.BenchmarkCode}</code></pre>
    {/if}
</section>

{#each uniqueVariationNames as variationName}
    <section>
        <h4>{benchmark.Name} <code>{variationName}</code></h4>
        <table class="table">
            <thead>
            <tr>
                <th>ns/op</th>
                <th>ops/sec</th>
                <th>B/op</th>
                <th>allocs/op</th>
                <th>MB/s</th>
            </tr>
            </thead>
            <tbody>
            {#each Object.values(averageBenchmark.Variations) as variation}
                {#if variation.Name === variationName}
                    <tr>
                        <td>{variation.NsPerOp.toFixed(2)}</td>
                        <td>{variation.OpsPerSec.toFixed(0)}</td>
                        <td>{variation.AllocedBytesPerOp.toFixed(1)}</td>
                        <td>{variation.AllocsPerOp.toFixed(0)}</td>
                        <td>{variation.MBPerS.toFixed(1)}</td>
                    </tr>
                {/if}
            {/each}
            </tbody>
        </table>
        <VariationAnalysis {benchmark} {variationName}/>
    </section>
{/each}

<style lang="scss">
  .code-area {
    @apply my-6;
  }

  section {
    @apply mb-40;
  }
</style>
