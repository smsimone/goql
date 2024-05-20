<script lang="ts">
  import { onMount } from "svelte";
  import { useNavigate } from "svelte-navigator";
  import { GetAvailableTables } from "../../wailsjs/go/main/App";

  const navigate = useNavigate();

  interface Table {
    schemaname: string;
    tablename: string;
  }

  let tableNames: Table[] = [];

  onMount(async () => {
    loadTables();
  });

  async function loadTables() {
    try {
      const data = await GetAvailableTables().catch((err) => {
        console.error(err);
      });

      tableNames = data;
    } catch (err) {
      console.error("failed to get tables", err);
    }
  }

  async function showData(table: Table) {
    const route = `/tables/${table.schemaname}/${table.tablename}`;
    navigate(route);
    console.log("navigated to", route);
  }
</script>

<main>
  {#if tableNames.length === 0}
    <div>Tables view</div>
    <button on:click={loadTables}>Load all tables</button>
  {:else}
    <div>Tables</div>
    <table class="data-table">
      <tr>
        <th>Schema name</th>
        <th>Table name</th>
        <th></th>
      </tr>
      {#each tableNames as table}
        <tr>
          <td class="cell">{table.schemaname}</td>
          <td class="cell">{table.tablename}</td>
          <td><button on:click={() => showData(table)}>Show data</button></td>
        </tr>
      {/each}
    </table>
  {/if}
</main>

<style>
  .cell {
    border: 2px solid black;
  }
</style>
