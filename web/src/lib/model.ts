import j from "../../../benchmarks.json"

export type BenchmarkGroups = BenchmarkGroup[]

export interface BenchmarkGroup {
    Name: string
    Description: string
    Benchmarks: Benchmark[]
}

export interface Benchmark {
    Name: string
    Description: string
    Variations: Variation[]
}

export interface Variation {
    N: number
    NsPerOp: number
    AllocedBytesPerOp: number
    AllocsPerOp: number
    MBPerS: number
    Measured: number
    Ord: number
    Name: string
    CPUCount: number
    OpsPerSec: number
}


// getBenchmarkGroups returns the benchmark groups from the JSON file.
export function getBenchmarkGroups(): BenchmarkGroups {
    return j
}

export function convertBenchmarksTo() {

}

export function convertBenchmarksToLineChartData(benchmarks: Benchmark[]) {
// Prepare an array to hold datasets for each Benchmark
    let datasets: {
        label: string;
        data: { x: number; y: number; }[];
        fill: boolean;
        borderColor: string;
        tension: number;
    }[] = [];
    let cpuCounts = new Set<number>();

    // Only keep variations that have the highest N in each benchmark
    // For example: If there are 8 variations for a benchmark with N=1, N=2, N=4, N=8, N=16, N=32, N=64, N=128
    // We only want to keep all variations with N=128
    benchmarks.forEach(benchmark => {
        const maxN = Math.max(...benchmark.Variations.map(variation => variation.N));
        benchmark.Variations = benchmark.Variations.filter(variation => variation.N === maxN);
    });

    benchmarks.forEach(benchmark => {
        benchmark.Variations.forEach(variation => {
            // Collect all unique CPU counts
            cpuCounts.add(variation.CPUCount);

            // Find or create the dataset for this Benchmark and Variation combination
            let dataset = datasets.find(d => d.label === `${benchmark.Name} | ${variation.Name}`);
            if (!dataset) {
                dataset = {
                    label: `${benchmark.Name} | ${variation.Name}`,
                    data: [],
                    fill: false,
                    borderColor: getRandomColor(), // Function to generate random colors
                    tension: 0.1
                };
                datasets.push(dataset);
            }

            // Add the data point for this Variation
            dataset.data.push({
                x: variation.CPUCount,
                y: variation.OpsPerSec
            });
        });
    });

    // Sort the data points within each dataset by CPUCount
    datasets.forEach(dataset => {
        dataset.data.sort((a, b) => a.x - b.x);
    });

    // Convert the CPU counts set to an array and sort it
    const labels = Array.from(cpuCounts).sort((a, b) => a - b);

    return {
        labels: labels.map(l => l + " CPU Cores"),
        datasets: datasets
    };
}

export function getLineChartOptions() {
    return {
        plugins: {
            zoom: {
                zoom: {
                    wheel: {
                        enabled: true,
                    },
                    pinch: {
                        enabled: true
                    },
                    mode: 'xy',
                }
            },
            tooltip: {
                callbacks: {
                    label: function (context: any) {
                        let label = context.dataset.label || '';

                        if (label) {
                            label += ': ';
                        }
                        if (context.parsed.y !== null) {
                            // Calculate ns/op from ops/sec
                            label += context.parsed.y.toFixed(0) + ' ops/sec (~' + (1000000000 / context.parsed.y).toFixed(2) + ' ns/op)';
                        }
                        return label;
                    }
                }
            },
        },
        scales: {
            y: {
                ticks: {
                    callback: function (value: any) {
                        return value + ' ops/sec';
                    }
                }
            }
        }
    }
}

export function getBarChartDataByCPUCount(benchmarks: Benchmark[]) {
    // First, we need to create a unique list of all CPU counts across all variations
    const cpuCounts = benchmarks
        .flatMap(benchmark => benchmark.Variations)
        .map(variation => variation.CPUCount)
        .filter((value, index, self) => self.indexOf(value) === index)
        .sort((a, b) => a - b); // Sort the CPU counts for consistent ordering

    // Create labels for the chart based on the names of the benchmarks and their variations
    const labels = benchmarks.flatMap(benchmark =>
        benchmark.Variations.map(variation => `${benchmark.Name} | ${variation.Name}`)
    ).filter((value, index, self) => self.indexOf(value) === index); // Ensure uniqueness

    // Initialize datasets, one for each CPU count
    const datasets = cpuCounts.map(cpuCount => {
        return {
            label: `CPU ${cpuCount}`,
            data: new Array(labels.length).fill(null), // Initialize with null which will be replaced with actual data
            backgroundColor: getRandomColor(), // Function to generate a color for each CPU count stack
        };
    });

    // Populate the datasets with NsPerOp values
    benchmarks.forEach(benchmark => {
        benchmark.Variations.forEach(variation => {
            // Find the index of the label corresponding to this variation
            const labelIndex = labels.indexOf(`${benchmark.Name} | ${variation.Name}`);
            // Find the dataset corresponding to this CPU count
            const dataset = datasets.find(d => d.label === `CPU ${variation.CPUCount}`);
            if (dataset) {
                // Replace the null with the actual NsPerOp value for the corresponding CPU count and variation
                dataset.data[labelIndex] = variation.OpsPerSec;
            }
        });
    });

    return {
        labels: labels, // x-axis labels for the variations
        datasets: datasets, // datasets for each CPU count
    };
}

export function getBarChartDataByRuns(benchmarks: Benchmark[]) {
    // First, we need to create a unique list of all CPU counts across all variations
    // const cpuCounts = benchmarks
    //     .flatMap(benchmark => benchmark.Variations)
    //     .map(variation => variation.CPUCount)
    //     .filter((value, index, self) => self.indexOf(value) === index)
    //     .sort((a, b) => a - b); // Sort the CPU counts for consistent ordering

    // Get a unique list of N values across all variations
    const ns = benchmarks
        .flatMap(benchmark => benchmark.Variations)
        .map(variation => variation.N)
        .filter((value, index, self) => self.indexOf(value) === index)
        .sort((a, b) => a - b); // Sort the N values for consistent ordering

    // Create labels for the chart based on the names of the benchmarks and their variations
    const labels = benchmarks.flatMap(benchmark =>
        benchmark.Variations.map(variation => `${benchmark.Name} | ${variation.Name}`)
    ).filter((value, index, self) => self.indexOf(value) === index); // Ensure uniqueness

    // Initialize datasets, one for each CPU count
    const datasets = ns.map(n => {
        return {
            label: `${n} Runs`,
            data: new Array(labels.length).fill(null), // Initialize with null which will be replaced with actual data
            backgroundColor: getRandomColor(), // Function to generate a color for each CPU count stack
        };
    });

    // Populate the datasets with NsPerOp values
    benchmarks.forEach(benchmark => {
        benchmark.Variations.forEach(variation => {
            // Find the index of the label corresponding to this variation
            const labelIndex = labels.indexOf(`${benchmark.Name} | ${variation.Name}`);
            // Find the dataset corresponding to this CPU count
            const dataset = datasets.find(d => d.label === `${variation.N} Runs`);
            if (dataset) {
                // Replace the null with the actual NsPerOp value for the corresponding CPU count and variation
                dataset.data[labelIndex] = variation.OpsPerSec;
            }
        });
    });

    return {
        labels: labels, // x-axis labels for the variations
        datasets: datasets, // datasets for each CPU count
    };
}

export function getBarChartDataByVariation(benchmarks: Benchmark[], variationName: string) {
    // Filter out benchmarks that have the specified variation name
    const filteredBenchmarks = benchmarks.map(benchmark => ({
        ...benchmark,
        Variations: benchmark.Variations.filter(variation => variation.Name === variationName)
    })).filter(benchmark => benchmark.Variations.length > 0);

    // Create a unique list of all CPU counts for the specified variation
    const cpuCounts = filteredBenchmarks
        .flatMap(benchmark => benchmark.Variations)
        .map(variation => variation.CPUCount)
        .filter((value, index, self) => self.indexOf(value) === index)
        .sort((a, b) => a - b); // Sort the CPU counts for consistent ordering

    // Initialize datasets, one for each CPU count
    const datasets = cpuCounts.map(cpuCount => {
        const data = filteredBenchmarks.map(benchmark => {
            const variation = benchmark.Variations.find(v => v.CPUCount === cpuCount);
            return variation ? variation.OpsPerSec : null;
        });
        return {
            label: `CPU ${cpuCount}`,
            data: data,
            backgroundColor: getRandomColor(), // Function to generate a color for each CPU count stack
        };
    });

    // Labels for the chart are the benchmark names
    const labels = filteredBenchmarks.map(benchmark => benchmark.Name);

    return {
        labels: labels, // x-axis labels for the benchmarks
        datasets: datasets, // datasets for each CPU count for the specified variation
    };
}

export function getBarChartDataByVariationAndRunCount(benchmarks: Benchmark[], variationName: string) {
    // Filter out benchmarks that have the specified variation name
    const filteredBenchmarks = benchmarks.map(benchmark => ({
        ...benchmark,
        Variations: benchmark.Variations.filter(variation => variation.Name === variationName)
    })).filter(benchmark => benchmark.Variations.length > 0);

    // Get a unique list of N values across all variations
    const ns = filteredBenchmarks
        .flatMap(benchmark => benchmark.Variations)
        .map(variation => variation.N)
        .filter((value, index, self) => self.indexOf(value) === index)
        .sort((a, b) => a - b); // Sort the N values for consistent ordering

    // Initialize datasets, one for each CPU count
    const datasets = ns.map(n => {
        const data = filteredBenchmarks.map(benchmark => {
            const variation = benchmark.Variations.find(v => v.N === n);
            return variation ? variation.OpsPerSec : null;
        });
        return {
            label: `${n} Runs`,
            data: data,
            backgroundColor: getRandomColor(), // Function to generate a color for each CPU count stack
        };
    });

    // Labels for the chart are the benchmark names
    const labels = filteredBenchmarks.map(benchmark => benchmark.Name);

    return {
        labels: labels, // x-axis labels for the benchmarks
        datasets: datasets, // datasets for each CPU count for the specified variation
    };
}

// Utility function to generate random light colors for the chart
function getRandomColor() {
    const r = Math.floor(Math.random() * 200 + 55);
    const g = Math.floor(Math.random() * 200 + 55);
    const b = Math.floor(Math.random() * 200 + 55);
    return `rgb(${r}, ${g}, ${b})`;
}
