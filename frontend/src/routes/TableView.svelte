<script lang="ts">
    import { onMount } from "svelte";
    import { GetTableData, UpdateValue } from "../../wailsjs/go/main/App";
    import type { database } from "wailsjs/go/models";

    export let tableName: string | undefined | null;
    export let schemaName: string | undefined | null;

    let columns: string[] | undefined;
    let rows: database.RowData[] | undefined;
    let primaryKeys: string[] | undefined;

    onMount(() => {
        console.log(schemaName, tableName);

        GetTableData(schemaName, tableName).then((val) => {
            columns = val.columns;
            rows = val.rows;
            primaryKeys = val.primay_key;
        });
    });

    async function updateValue(rowId: string, value: string, column: number) {
        // UpdateValue({});
        console.table({
            tableName,
            schemaName,
            rowId,
            value,
            column,
        });

        if (false)
            UpdateValue({
                value,
                table: tableName,
                field: columns[column],
                where: "",
                schema: schemaName,
            });
    }
</script>

<main>
    <div>
        <button on:click={() => history.back()}>Back</button>
        <h1>Table {schemaName}.{tableName}</h1>
    </div>
    {#if columns}
        <table class="table">
            <tr>
                {#each columns as col}
                    <th
                        class={`cell ${primaryKeys?.includes(col) ?? false ? "primary" : ""}`}
                        ><div>{col}</div></th
                    >
                {/each}
            </tr>
            {#if rows}
                {#each rows as row}
                    <tr>
                        {#each row.columns as col, idx}
                            <td
                                class="cell"
                                on:keydown={() => {}}
                                on:click={() => {
                                    updateValue(
                                        row.columns[0].value,
                                        "prova",
                                        idx,
                                    );
                                }}>{col.value}</td
                            >
                        {/each}
                    </tr>
                {/each}
            {/if}
        </table>
    {/if}
</main>

<style>
    .table {
        width: 100%;
    }

    .cell {
        border: solid black 2px;
        border-radius: 3px;
    }

    .primary {
        background: red;
    }
</style>
