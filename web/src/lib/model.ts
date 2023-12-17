import j from "../../../benchmarks.json"
import showdown from "showdown"

export type BenchmarkGroups = BenchmarkGroup[]

export interface BenchmarkGroup {
    Name: string
    Headline: string
    Description: string
    Benchmarks: Benchmark[]
    Code: string
    Constants: string
}

export interface Benchmark {
    Name: string
    Description: string
    BenchmarkCode: string
    Code: string
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
    let benchmarkGroups: BenchmarkGroups = j;

    // Convert markdown to HTML
    benchmarkGroups.forEach(benchmarkGroup => {
        benchmarkGroup.Headline = markdownToHtml(benchmarkGroup.Headline);
        benchmarkGroup.Description = markdownToHtml(benchmarkGroup.Description);
        benchmarkGroup.Benchmarks.forEach(benchmark => {
            benchmark.Description = markdownToHtml(benchmark.Description);
        });
    });

    return benchmarkGroups;
}

function markdownToHtml(md: string): string {
    let converter = new showdown.Converter();
    return converter.makeHtml(md);
}

export function convertBenchmarksToCPUCountPerformanceLineChart(benchmarks: Benchmark[]) {
    // Make copy of benchmarks
    benchmarks = JSON.parse(JSON.stringify(benchmarks));

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
            let dataset = datasets.find(d => d.label === `${variation.Name}`);
            if (!dataset) {
                dataset = {
                    label: `${variation.Name}`,
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

export function convertBenchmarksToRunCountPerformanceLineChart(benchmarks: Benchmark[]) {
    // Make copy of benchmarks
    benchmarks = JSON.parse(JSON.stringify(benchmarks));

    // Prepare an array to hold datasets for each Benchmark
    let datasets: {
        label: string;
        data: { x: number; y: number; }[];
        fill: boolean;
        borderColor: string;
        tension: number;
    }[] = [];
    let ns = new Set<number>();

    // Only keep variations that have the highest CPU Count in each benchmark
    // For example: If there are variations for a benchmark with CPUCount=1, CPUCount=2, CPUCount=4, CPUCount=8, CPUCount=16, CPUCount=32
    // We only want to keep all variations with CPUCount=32
    benchmarks.forEach(benchmark => {
        const maxCPUCount = Math.max(...benchmark.Variations.map(variation => variation.CPUCount));
        benchmark.Variations = benchmark.Variations.filter(variation => variation.CPUCount === maxCPUCount);
    });

    benchmarks.forEach(benchmark => {
        benchmark.Variations.forEach(variation => {
            // Collect all unique N values
            ns.add(variation.N);

            // Find or create the dataset for this Benchmark and Variation combination
            let dataset = datasets.find(d => d.label === `${variation.Name}`);
            if (!dataset) {
                dataset = {
                    label: `${variation.Name}`,
                    data: [],
                    fill: false,
                    borderColor: getRandomColor(), // Function to generate random colors
                    tension: 0.1
                };
                datasets.push(dataset);
            }

            // Add the data point for this Variation
            dataset.data.push({
                x: variation.N,
                y: variation.OpsPerSec
            });
        });
    });

    // Sort the data points within each dataset by N
    datasets.forEach(dataset => {
        dataset.data.sort((a, b) => a.x - b.x);
    });

    // Convert the N set to an array and sort it
    const labels = Array.from(ns).sort((a, b) => a - b);


    return {
        labels: labels.map(l => l + " Runs"),
        datasets: datasets
    };
}

export function getLineChartOptions(isLogarithmic: boolean) {
    return {
        responsive: true,
        plugins: {
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
                type: isLogarithmic ? 'logarithmic' : 'linear', // Set scale type based on isLogarithmic
                ticks: {
                    callback: function (value: any) {
                        return value + ' ops/sec';
                    }
                }
            }
        }
    }
}

export function getBarChartDataByCPUCountMulti(benchmarks: Benchmark[]) {
    // Make copy of benchmarks
    benchmarks = JSON.parse(JSON.stringify(benchmarks));

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
            label: `${cpuCount} CPU Cores`,
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
            const dataset = datasets.find(d => d.label === `${variation.CPUCount} CPU Cores`);
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

export function getChartDataByCPUCount(benchmark: Benchmark) {
    // Make copy of benchmarks
    benchmark = JSON.parse(JSON.stringify(benchmark));

    // Create a unique list of all CPU counts across all variations
    const cpuCounts = benchmark.Variations
        .map(variation => variation.CPUCount)
        .filter((value, index, self) => self.indexOf(value) === index)
        .sort((a, b) => a - b); // Sort the CPU counts for consistent ordering

    // Create a unique list of all Run Counts across all variations
    const runCounts = benchmark.Variations
        .map(variation => variation.N)
        .filter((value, index, self) => self.indexOf(value) === index)
        .sort((a, b) => a - b); // Sort the Run Counts for consistent ordering

    // Initialize datasets.
    const datasets = runCounts.map(runCount => {
        return {
            label: `${runCount} Runs`,
            data: cpuCounts.map(cpuCount => {
                // Find the variation with this Run Count and CPU Count
                const variation = benchmark.Variations.find(v => v.N === runCount && v.CPUCount === cpuCount);
                // If a variation is found, return its OpsPerSec, otherwise return null
                return variation ? variation.OpsPerSec : null;
            }),
            backgroundColor: getRandomColor(), // Function to generate a color for each Run Count line
            borderColor: getRandomColor(), // Function to generate a color for each line
        };
    });

    // Labels for the chart is the CPU Core count
    const labels = cpuCounts.map(cpuCount => cpuCount + " CPU Cores");

    return {
        labels: labels, // x-axis labels for the benchmarks
        datasets: datasets, // datasets for each Run Count
    };
}

export function filterBenchmarkByVariationName(benchmark: Benchmark, variationName: string) {
    // Make copy of benchmarks
    benchmark = JSON.parse(JSON.stringify(benchmark));

    // Filter out benchmarks that have the specified variation name
    benchmark.Variations = benchmark.Variations.filter(variation => variation.Name === variationName);

    return benchmark;
}

export function getBarChartDataByRunsMulti(benchmarks: Benchmark[]) {
    // Make copy of benchmarks
    benchmarks = JSON.parse(JSON.stringify(benchmarks));

    // Get a list of N values across all variations
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

export function getChartDataByRuns(benchmark: Benchmark) {
    // Make copy of benchmarks
    benchmark = JSON.parse(JSON.stringify(benchmark));

    // Create a unique list of all Run Counts across all variations
    const runCounts = benchmark.Variations
        .map(variation => variation.N)
        .filter((value, index, self) => self.indexOf(value) === index)
        .sort((a, b) => a - b); // Sort the Run Counts for consistent ordering

    // Create a unique list of all CPU Counts across all variations
    const cpuCounts = benchmark.Variations
        .map(variation => variation.CPUCount)
        .filter((value, index, self) => self.indexOf(value) === index)
        .sort((a, b) => a - b); // Sort the CPU Counts for consistent ordering

    // Initialize datasets.
    const datasets = cpuCounts.map(cpuCount => {
        return {
            label: `${cpuCount} CPU Cores`,
            data: runCounts.map(n => {
                // Find the variation with this Run Count and CPU Count
                const variation = benchmark.Variations.find(v => v.N === n && v.CPUCount === cpuCount);
                // If a variation is found, return its OpsPerSec, otherwise return null
                return variation ? variation.OpsPerSec : null;
            }),
            backgroundColor: getRandomColor(), // Function to generate a color for each CPU Count line
            borderColor: getRandomColor(), // Function to generate a color for each line
        };
    });

    // Labels for the chart is the Run count
    const labels = runCounts.map(n => n + " Runs");

    return {
        labels: labels, // x-axis labels for the benchmarks
        datasets: datasets, // datasets for each CPU Count
    };
}

export function getBarChartDataByVariation(benchmarks: Benchmark[], variationName: string) {
    // Make copy of benchmarks
    benchmarks = JSON.parse(JSON.stringify(benchmarks));

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
            label: `${cpuCount} CPU Cores`,
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
    // Make copy of benchmarks
    benchmarks = JSON.parse(JSON.stringify(benchmarks));

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
