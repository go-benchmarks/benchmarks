
import j from "../../../benchmarks.json"

export type BenchmarkGroups = BenchmarkGroup[]


export interface BenchmarkGroup {
    Name: string
    Benchmarks: Benchmark[]
}

export interface Benchmark {
    Name: string
    N: number
    NsPerOp: number
    AllocedBytesPerOp: number
    AllocsPerOp: number
    MBPerS: number
    Measured: number
    Ord: number
    Variation: string
    CPUCount: number
    Performance: number
}

// getBenchmarkGroups returns the benchmark groups from the JSON file.
export function getBenchmarkGroups(): BenchmarkGroups {
    return j
}
