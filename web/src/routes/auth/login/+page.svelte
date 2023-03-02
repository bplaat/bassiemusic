<script>
    import { AUTH_TOKEN_EXPIRES_TIMEOUT } from '../../../consts.js';

    export let data;

    let logon = '';
    let password = '';
    let errors = {};

    async function login() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/login`, {
            method: 'POST',
            body: new URLSearchParams({
                logon,
                password,
            }),
        });
        const { success, token } = await response.json();
        if (success) {
            document.cookie = `token=${token}; path=/; samesite=strict; expires=${new Date(
                Date.now() + AUTH_TOKEN_EXPIRES_TIMEOUT
            ).toUTCString()}`;
            if (data.continueUrl != undefined) {
                window.location = data.continueUrl;
            } else {
                window.location = '/';
            }
        } else {
            errors = {
                logon: 'Wrong username, email address or password',
            };
        }
    }
</script>

<svelte:head>
    <title>Login - BassieMusic</title>
</svelte:head>

<div class="flex" />

<form class="box content" style="width: 40rem; max-width: 100%;" on:submit|preventDefault={login}>
    <h2 class="title is-4">Login with you BassieMusic account</h2>
    <p>
        <i>
            You need a BassieMusic account to listen to music<br />
            If you don't have an account ask the admins to create one for you
        </i>
    </p>

    <div class="field">
        <label class="label" for="logon">Username or email address</label>
        <div class="control">
            <!-- svelte-ignore a11y-autofocus -->
            <input
                class="input"
                class:is-danger={errors.logon}
                type="text"
                id="logon"
                bind:value={logon}
                autofocus
                required
            />
        </div>
        {#if errors.logon}
            <p class="help is-danger">{errors.logon}</p>
        {/if}
    </div>

    <div class="field">
        <label class="label" for="password">Password</label>
        <div class="control">
            <input
                class="input"
                class:is-danger={errors.logon}
                type="password"
                id="password"
                bind:value={password}
                required
            />
        </div>
    </div>

    <div class="field">
        <div class="control">
            <button class="button is-link" type="submit">Login</button>
        </div>
    </div>
</form>

<div class="flex" />

<p>
    Made with
    <svg class="icon is-inline is-colored" viewBox="0 0 24 24" style="width: 16px; height: 16px;" title="love">
        <path
            fill="#f14668"
            d="M12,21.35L10.55,20.03C5.4,15.36 2,12.27 2,8.5C2,5.41 4.42,3 7.5,3C9.24,3 10.91,3.81 12,5.08C13.09,3.81 14.76,3 16.5,3C19.58,3 22,5.41 22,8.5C22,12.27 18.6,15.36 13.45,20.03L12,21.35Z"
        />
    </svg>
    by <a href="https://www.plaatsoft.nl/" target="_blank" rel="noreferrer">PlaatSoft</a>
</p>

<style>
    :global(.section) {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        min-height: 100%;
    }
</style>
