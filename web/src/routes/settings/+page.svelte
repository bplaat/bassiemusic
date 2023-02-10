<script>
    export let data;
    let { token, authUser, currentSession, sessions } = data;

    // Change details
    let newPassword = '';

    async function changeDetails() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                username: authUser.username,
                email: authUser.email,
                password: newPassword,
                theme: authUser.theme,
            }),
        });
        if (response.status == 200) {
            window.location = '/settings';
        } else {
            alert('Error!');
        }
    }

    // Change avatar
    let avatarInput;

    async function changeAvatar() {
        if (!avatarInput.files[0]) {
            return;
        }

        const formData = new FormData();
        formData.append('avatar', avatarInput.files[0], avatarInput.files[0].name);

        const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}/avatar`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: formData,
        });
        if (response.status == 200) {
            window.location = '/settings';
        } else {
            alert('Error!');
        }
    }

    async function deleteAvatar() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}/avatar/delete`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        if (response.status == 200) {
            window.location = '/settings';
        } else {
            alert('Error!');
        }
    }

    // Sessions management
    async function fetchSessionsPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/sessions?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data: newSessions, pagination } = await response.json();
        sessions.push(...newSessions);
        sessions = sessions;
        if (sessions.length != pagination.total) {
            fetchPage(page + 1);
        }
    }
    if (sessions.length != data.sessionsTotal) {
        fetchSessionsPage(2);
    }

    async function revokeSession(session) {
        fetch(`${import.meta.env.VITE_API_URL}/sessions/${session.id}/revoke`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        if (currentSession.id == session.id) {
            document.cookie = `token=; expires=${new Date(0).toUTCString()}`;
            window.location = '/auth/login';
        } else {
            sessions = sessions.filter((otherSession) => otherSession.id != session.id);
        }
    }
</script>

<svelte:head>
    <title>Settings - BassieMusic</title>
</svelte:head>

<h2 class="title">Settings</h2>

<div class="columns">
    <!-- Change details -->
    <div class="column">
        <form on:submit|preventDefault={changeDetails} class="box">
            <h3 class="title is-4">Change details</h3>

            <div class="columns">
                <div class="column">
                    <div class="field">
                        <label class="label" for="username">Username</label>
                        <div class="control">
                            <input class="input" type="text" id="username" bind:value={authUser.username} required />
                        </div>
                    </div>
                </div>

                <div class="column">
                    <div class="field">
                        <label class="label" for="email">Email</label>
                        <div class="control">
                            <input class="input" type="email" id="email" bind:value={authUser.email} required />
                        </div>
                    </div>
                </div>
            </div>

            <div class="field">
                <label class="label" for="password">Password (leave empty to not change)</label>
                <div class="control">
                    <input class="input" type="password" id="password" bind:value={newPassword} />
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
    </div>

    <!-- Change avatar -->
    <div class="column">
        <form on:submit|preventDefault={changeAvatar} class="box">
            <h3 class="title is-4">Change avatar</h3>

            <div class="field">
                <label class="label" for="avatar">Avatar file</label>
                <div class="control">
                    <input class="input" type="file" id="avatar" bind:this={avatarInput} />
                    <p class="help">You can upload an squared .jpg or .png image</p>
                </div>
            </div>

            <div class="field">
                <div class="buttons">
                    <button type="submit" class="button is-link">Change avatar</button>

                    {#if authUser.avatar}
                        <button type="button" class="button is-danger" on:click|preventDefault={deleteAvatar}
                            >Delete avatar</button
                        >
                    {/if}
                </div>
            </div>
        </form>
    </div>
</div>

<!-- Sessions management -->
<div class="box">
    <h3 class="title is-4">Sessions management</h3>

    <div class="columns is-multiline is-mobile">
        {#each sessions as session}
            {#if new Date(session.expires_at).getTime() >= Date.now()}
                <div class="column is-half">
                    <div class="box content">
                        <h3 class="title is-4">
                            {session.client_name} on {session.client_os}
                            {#if currentSession.id == session.id}
                                <span class="tag is-link is-pulled-right">CURRENT</span>
                            {/if}
                        </h3>
                        <p>
                            With {session.ip} at
                            {#if session.ip_city != null && session.ip_country != null}
                                {session.ip_city}, {session.ip_country}
                            {:else}
                                Unknown location
                            {/if}
                        </p>
                        <p>
                            Logged in at: {new Date(session.created_at).toLocaleString()}
                        </p>
                        <p>
                            Expires at: {new Date(session.expires_at).toLocaleString()}
                        </p>
                        <div class="buttons">
                            <button class="button is-danger" on:click={() => revokeSession(session)}
                                >Revoke session</button
                            >
                        </div>
                    </div>
                </div>
            {/if}
        {/each}
    </div>
</div>
