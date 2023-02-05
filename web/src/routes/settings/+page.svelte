<script>
    import Cookies from "js-cookie";

    export let data;
    const { authUser } = data;
    let newPassword = '';

    async function editDetails() {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${authUser.id}`,
            {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${Cookies.get("token")}`,
                },
                body: new URLSearchParams({
                    username: authUser.username,
                    email: authUser.email,
                    password: newPassword,
                    theme: authUser.theme,
                }),
            }
        );
        if (response.status == 200) {
            window.location = "/settings";
        } else {
            alert("Error!");
        }
    }
</script>

<svelte:head>
    <title>Settings - BassieMusic</title>
</svelte:head>

<h2 class="title">Settings</h2>

<form on:submit|preventDefault={editDetails} class="box">
    <h3 class="title is-4">Change details</h3>

    <div class="columns">
        <div class="column">
            <div class="field">
                <label class="label" for="username">Username</label>
                <div class="control">
                    <input
                        class="input"
                        type="text"
                        id="username"
                        bind:value={authUser.username}
                        required
                    />
                </div>
            </div>
        </div>

        <div class="column">
            <div class="field">
                <label class="label" for="email">Email</label>
                <div class="control">
                    <input
                        class="input"
                        type="email"
                        id="email"
                        bind:value={authUser.email}
                        required
                    />
                </div>
            </div>
        </div>
    </div>

    <div class="field">
        <label class="label" for="edit-password"
            >Password (leave empty to not change)</label
        >
        <div class="control">
            <input
                class="input"
                type="password"
                id="edit-password"
                bind:value={newPassword}
            />
        </div>
    </div>

    <div class="field">
        <label class="label" for="theme">Theme</label>
        <div class="control">
            <div class="select is-fullwidth">
                <select id="theme" bind:value={authUser.theme} required>
                    <option value="system">System</option>
                    <option value="light">Light</option>
                    <option value="dark">Dark</option>
                </select>
            </div>
        </div>
    </div>

    <div class="field">
        <div class="buttons">
            <button type="submit" class="button is-link">Change details</button>
        </div>
    </div>
</form>
