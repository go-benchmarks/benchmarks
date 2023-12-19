<script lang="ts">
    import {
        type Benchmark,
        filterBenchmarkByVariationName,
        getChartDataByCPUCount,
        getChartDataByRuns,
        getLineChartOptions
    } from "$lib/model";
    import {Line} from "svelte-chartjs";

    export let benchmark: Benchmark
    export let variationName: string
</script>

<h5>By Run Count</h5>
<Line height="{50}" options="{getLineChartOptions(false, true)}"
      data="{getChartDataByRuns(filterBenchmarkByVariationName(benchmark, variationName))}"/>
<h5>By CPU Core Count</h5>
<Line height="{50}" options="{getLineChartOptions(false, true)}"
      data="{getChartDataByCPUCount(filterBenchmarkByVariationName(benchmark, variationName))}"/>

