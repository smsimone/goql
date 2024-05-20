<script lang="ts">
    export let value: string;
    export let label: string;
    export let placeholder: string | undefined | null;
    export let type: "text" | "password" = "text";
</script>

<main>
    <div class="full-row">
        <div class="input-line">
            <p class="label">{label}</p>
            {#if type === "password"}
                <input
                    type="password"
                    bind:value
                    class="input"
                    placeholder={placeholder ?? ""}
                />
            {:else}
                <input
                    bind:value
                    class="input"
                    placeholder={placeholder ?? ""}
                />
            {/if}
        </div>
    </div>
</main>

<style>
    .input-line {
        display: flex;
        flex-direction: row;
        justify-content: center;
        width: 30%;
        align-items: center;
        justify-content: space-between;
        height: fit-content;
    }

    .full-row {
        width: 100%;
        justify-content: center;
        display: flex;
        height: fit-content;
    }

    .input {
        border-radius: 5px;
    }
</style>
