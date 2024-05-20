<script lang="ts">
    import { Route, Router } from "svelte-navigator";
    import DatabaseTables from "./routes/DatabaseTables.svelte";
    import Homepage from "./routes/Homepage.svelte";
    import NotFound from "./routes/NotFound.svelte";
    import TableView from "./routes/TableView.svelte";
</script>

<main>
    <Router>
        <div>
            <Route path="/" component={Homepage} />
            <Route path="/tables/*">
                <Route path="/" component={DatabaseTables} />
                <Route path=":schema/:table" let:params>
                    <TableView
                        tableName={params.table}
                        schemaName={params.schema}
                    />
                </Route>
            </Route>
            <Route component={NotFound} />
        </div>
    </Router>
</main>
