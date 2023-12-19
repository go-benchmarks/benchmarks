<script lang="ts">
    import type {BenchmarkGroup} from "$lib/model";

    export let benchmarkGroup: BenchmarkGroup;

    // Extract unique variation names
    const uniqueVariationNames = benchmarkGroup.Benchmarks.flatMap(b => b.Variations)
        .map(v => v.Name)
        .filter((v, i, a) => a.indexOf(v) === i);

    // Function to calculate the average NsPerOp for each variation
    function calculateAverageNsPerOp(variationName: string) {
        let totalNsPerOp = 0;
        let count = 0;

        benchmarkGroup.Benchmarks.forEach(benchmark => {
            benchmark.Variations.forEach(variation => {
                if (variation.Name === variationName) {
                    totalNsPerOp += variation.NsPerOp;
                    count++;
                }
            });
        });

        return count > 0 ? totalNsPerOp / count : null;
    }

    // Determine the fastest benchmark for each unique variation
    let fastestBenchmarkForVariation = {};
    uniqueVariationNames.forEach(name => {
        let lowestAverage = Infinity;
        let fastestBenchmark = null;

        benchmarkGroup.Benchmarks.forEach(benchmark => {
            let totalNsPerOp = 0;
            let count = 0;

            benchmark.Variations.forEach(variation => {
                if (variation.Name === name) {
                    totalNsPerOp += variation.NsPerOp;
                    count++;
                }
            });

            let average = count > 0 ? totalNsPerOp / count : Infinity;
            if (average < lowestAverage) {
                lowestAverage = average;
                fastestBenchmark = benchmark.Name;
            }
        });

        fastestBenchmarkForVariation[name] = fastestBenchmark;
    });

    let benchmarksWithFastestVariations = benchmarkGroup.Benchmarks.map(benchmark => {
        let fastestVariations = uniqueVariationNames.filter(name =>
            fastestBenchmarkForVariation[name] === benchmark.Name
        );
        return {
            name: benchmark.Name,
            fastestVariations: fastestVariations
        };
    });

    let slowestBenchmarkForVariation = {};
    uniqueVariationNames.forEach(name => {
        let highestAverage = -Infinity;
        let slowestBenchmark = null;

        benchmarkGroup.Benchmarks.forEach(benchmark => {
            let totalNsPerOp = 0;
            let count = 0;

            benchmark.Variations.forEach(variation => {
                if (variation.Name === name) {
                    totalNsPerOp += variation.NsPerOp;
                    count++;
                }
            });

            let average = count > 0 ? totalNsPerOp / count : -Infinity;
            if (average > highestAverage) {
                highestAverage = average;
                slowestBenchmark = benchmark.Name;
            }
        });

        slowestBenchmarkForVariation[name] = slowestBenchmark;
    });

    let benchmarksWithVariationSpeeds = benchmarkGroup.Benchmarks.map(benchmark => {
        let fastestVariations = uniqueVariationNames.filter(name =>
            fastestBenchmarkForVariation[name] === benchmark.Name
        );
        let slowestVariations = uniqueVariationNames.filter(name =>
            slowestBenchmarkForVariation[name] === benchmark.Name
        );
        return {
            name: benchmark.Name,
            fastestVariations: fastestVariations,
            slowestVariations: slowestVariations
        };
    });
</script>

{#if benchmarksWithVariationSpeeds.length > 1}
    <h2 class="mt-16">Implementations</h2>
    <blockquote class="prose">Green marks the fastest implementation, red marks the slowest</blockquote>
    <nav class="mb-32">
        <ul>
            {#each benchmarksWithVariationSpeeds as benchmark}
                <li>
                    <a class="link-hover text-blue-400 text-xl"
                       href="#{benchmark.name.replaceAll(' ', '-').toLowerCase()}">
                        {benchmark.name}
                    </a>
                    <span class="text-green-600 text-xs">
                        {#each benchmark.fastestVariations as variation}
                            {#if variation !== "run"}
                                <span>{variation} </span>
                            {:else}
                                fastest
                            {/if}
                        {/each}
                    </span>
                    <span class="text-red-600 text-xs">
                        {#each benchmark.slowestVariations as variation}
                            {#if variation !== "run"}
                                <span>{variation} </span>
                            {:else}
                                slowest
                            {/if}
                        {/each}
                    </span>
                </li>
            {/each}
        </ul>
    </nav>
{/if}





