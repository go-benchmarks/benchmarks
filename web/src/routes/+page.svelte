<script lang="ts">

    import {getBenchmarkGroups} from "$lib/model";

    const benchmarkGroups = getBenchmarkGroups()
</script>

<ul>
{#each benchmarkGroups as group}
    <li>
        <a href="./benchmark/{group.Name}">{group.Name}</a>
    </li>
{/each}
</ul>
