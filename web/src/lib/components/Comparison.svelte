<script lang="ts">
    import {
        type Benchmark,
        type BenchmarkGroup,
        getBarChartDataByCPUCountMulti,
        getBarChartDataByRunsMulti,
        getLineChartOptions, type Variation
    } from "$lib/model.js";
  import {Bar} from "svelte-chartjs";

  export let benchmarkGroup: BenchmarkGroup;

  let benchmarks: Benchmark[] = benchmarkGroup.Benchmarks;

  const uniqueVariationNames = benchmarks.flatMap(b => b.Variations).map(v => v.Name).filter((v, i, a) => a.indexOf(v) === i);

  let uniqueVariationNamesString: string;
  if (uniqueVariationNames.length === 1) {
      uniqueVariationNamesString = `<code>${uniqueVariationNames[0]}</code>`;
  } else {
      uniqueVariationNamesString = uniqueVariationNames.map(vn => `<code>${vn}</code>`).slice(0, -1).join(", ") + " and " + uniqueVariationNames.map(vn => `<code>${vn}</code>`).slice(-1);
  }
</script>

<section>
    <hgroup>
        <h2>Comparison</h2>
        <p>This section compares the performance of <code>{benchmarkGroup.Name}</code> implementations with the
            functions: {@html uniqueVariationNamesString}</p>
    </hgroup>
    <h3>Different Run Count</h3>
    <Bar height="{50}" options="{getLineChartOptions(true, false)}"
         data="{getBarChartDataByRunsMulti(benchmarks)}"/>
    <h3>Different CPU Core Count</h3>
    <Bar height="{50}" options="{getLineChartOptions(true, false)}"
         data="{getBarChartDataByCPUCountMulti(benchmarks)}"/>
</section>
