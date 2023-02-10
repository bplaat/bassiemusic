<script>
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
                Date.now() + 356 * 24 * 60 * 60 + 1000
            ).toUTCString()}`;
            window.location = '/';
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

<form on:submit|preventDefault={login}>
    <h2 class="title">Login with you BassieMusic account</h2>

    <div class="box">
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
    </div>
</form>
