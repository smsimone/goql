<script lang="ts">
    import { onMount } from "svelte";
    import type { NavigateFn } from "svelte-navigator";
    import {
        AddConfiguration,
        GetConfiguration,
        GetCurrentConfigurations,
        SetActiveConnection,
        TestConnection,
    } from "../../wailsjs/go/main/App";
    import type { configuration } from "../../wailsjs/go/models";
    import logo from "../assets/images/logo-universal.png";
    import InputComponent from "../components/InputComponent.svelte";

    export let location: string;
    export let navigate: NavigateFn;

    let configurationNames: string[] = [];

    let data: configuration.DatabaseConnection = {
        name: "",
        url: "",
        port: 5432,
        username: "",
        password: "",
        database: "",
        id: -1,
    };

    $: data.port = Number(tempPort);

    let tempPort: string = "5432";

    onMount(() => {
        getConfigurations();
    });

    /**
     * Gets all the configurations saved
     */
    function getConfigurations() {
        GetCurrentConfigurations().then((names) => {
            configurationNames = names;
        });
    }

    /**
     * Loads a single configuration given its name
     */
    async function loadConfiguration(name: string) {
        const conf = await GetConfiguration(name);
        data = {
            name: conf.name,
            password: conf.password,
            port: conf.port,
            username: conf.username,
            url: conf.url,
            database: conf.database,
            id: conf.id,
        };
    }

    async function addConnection() {
        await AddConfiguration({
            name: data.name,
            url: data.url,
            port: data.port,
            username: data.username,
            password: data.password,
            database: data.database,
        });
        getConfigurations();
        data = {
            name: "",
            url: "",
            port: 5432,
            username: "",
            password: "",
            database: "",
            id: -1,
        };
    }

    async function connect() {
        if (await testConnection()) {
            const val = await SetActiveConnection(data.id)
                .catch((err) => {
                    alert(err);
                    return false;
                })
                .then((_) => true);
            if (val) {
                navigate("/tables");
            }
        }
    }

    async function testConnection(): Promise<boolean> {
        return await TestConnection(data)
            .catch((err) => {
                console.error("failed to connect to database", err);
                return false;
            })
            .then((_) => true);
    }
</script>

<main>
    <img alt="Wails logo" id="logo" src={logo} />
    {#if configurationNames.length !== 0}
        <div class="dropdown">
            <button class="dropbtn">Configurations</button>
            <div class="dropdown-content">
                {#each configurationNames as name}
                    <p
                        on:click={() => {
                            loadConfiguration(name);
                        }}
                        on:keydown={() => {
                            loadConfiguration(name);
                        }}
                    >
                        {name}
                    </p>
                {/each}
            </div>
        </div>
    {/if}
    <div class="input-box">
        <InputComponent
            placeholder="Connection name"
            label="name"
            bind:value={data.name}
        />
        <InputComponent
            placeholder="postgres://localhost"
            label="URL:"
            bind:value={data.url}
        />
        <InputComponent
            placeholder="5432"
            label="Port:"
            bind:value={tempPort}
        />
        <InputComponent
            placeholder="postgres"
            label="Database:"
            bind:value={data.database}
        />
        <InputComponent
            placeholder="postgres"
            label="Username:"
            bind:value={data.username}
        />
        <InputComponent
            placeholder="super_secret_password"
            label="Password:"
            type="password"
            bind:value={data.password}
        />
        <div class="vert-center">
            <div class="button-row">
                <button class="btn" on:click={addConnection}>Save</button>
                <button class="btn" on:click={testConnection}
                    >Test connection</button
                >
                <button class="btn" on:click={connect}>Connect</button>
            </div>
        </div>
    </div>
</main>

<style>
    .button-row {
        padding-top: 10px;
        display: flex;
        width: 30%;
        align-items: center;
        justify-content: space-evenly;
    }

    .vert-center {
        justify-content: center;
        width: 100%;
        display: flex;
        height: fit-content;
    }

    #logo {
        display: block;
        width: 30%;
        height: 30%;
        margin: auto;
        padding: 10% 0 0;
        background-position: center;
        background-repeat: no-repeat;
        background-size: 100% 100%;
        background-origin: content-box;
    }

    .btn {
        text-align: center;
        margin: auto;
    }

    .input-box {
        flex-direction: column;
        display: flex;
        align-content: center;
        justify-content: center;
    }

    .dropbtn {
        background-color: #04aa6d;
        color: white;
        padding: 16px;
        font-size: 16px;
        border: none;
    }

    /* The container <div> - needed to position the dropdown content */
    .dropdown {
        position: relative;
        display: inline-block;
    }

    /* Dropdown Content (Hidden by Default) */
    .dropdown-content {
        display: none;
        position: absolute;
        background-color: #f1f1f1;
        min-width: 160px;
        box-shadow: 0px 8px 16px 0px rgba(0, 0, 0, 0.2);
        z-index: 1;
    }

    .dropdown-content {
        display: none;
        position: absolute;
        background-color: #f1f1f1;
        min-width: 160px;
        box-shadow: 0px 8px 16px 0px rgba(0, 0, 0, 0.2);
        z-index: 1;
    }

    /* Links inside the dropdown */
    .dropdown-content p {
        color: black;
        padding: 12px 16px;
        text-decoration: none;
        display: block;
    }

    /* Change color of dropdown links on hover */
    .dropdown-content p:hover {
        background-color: #ddd;
    }

    /* Show the dropdown menu on hover */
    .dropdown:hover .dropdown-content {
        display: block;
    }

    /* Change the background color of the dropdown button when the dropdown content is shown */
    .dropdown:hover .dropbtn {
        background-color: #3e8e41;
    }
</style>
