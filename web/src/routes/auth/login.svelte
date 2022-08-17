<script>
import { API_URL } from '../../config.js';

export let email = '';
export let password = '';

async function login() {
    const response = await fetch(`${API_URL}/auth/login?email=${encodeURIComponent(email)}&password=${encodeURIComponent(password)}`, {
        headers: {
            Authorization: 'Bearer ' + localStorage.token
        }
    });
    const data = await response.json();
    if (data.success) {
        alert('Authed!');
        localStorage.token = data.token;
    } else {
        alert(data.message);
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
            <label class="label" for="email">Email</label>
            <div class="control">
                <input class="input" type="email" id="email" bind:value={email} autofocus required>
            </div>
        </div>

        <div class="field">
            <label class="label" for="password">Password</label>
            <div class="control">
                <input class="input" type="password" id="password" bind:value={password} required>
            </div>
        </div>

        <div class="field">
            <div class="control">
                <button class="button is-link" type="submit">Login</button>
            </div>
        </div>
    </div>
</form>
