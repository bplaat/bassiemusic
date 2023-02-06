<script>
    export let data;
    const { token, authUser } = data;
    let newPassword = "";

    async function changeDetails() {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${authUser.id}`,
            {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${token}`,
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

    let avatarInput;

    async function changeAvatar() {
        if (!avatarInput.files[0]) {
            return;
        }

        const formData = new FormData();
        formData.append(
            "avatar",
            avatarInput.files[0],
            avatarInput.files[0].name
        );

        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${authUser.id}/avatar`,
            {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${token}`,
                },
                body: formData,
            }
        );
        if (response.status == 200) {
            window.location = "/settings";
        } else {
            alert("Error!");
        }
    }

    async function deleteAvatar() {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${
                authUser.id
            }/avatar/delete`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
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
                <label class="label" for="password"
                    >Password (leave empty to not change)</label
                >
                <div class="control">
                    <input
                        class="input"
                        type="password"
                        id="password"
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
                    <button type="submit" class="button is-link"
                        >Change details</button
                    >
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
                    <input
                        class="input"
                        type="file"
                        id="avatar"
                        bind:this={avatarInput}
                    />
                    <p class="help">You can upload an squared .jpg or .png image</p>
                </div>
            </div>

            <div class="field">
                <div class="buttons">
                    <button type="submit" class="button is-link"
                        >Change avatar</button
                    >

                    {#if authUser.avatar}
                        <button
                            type="button"
                            class="button is-danger"
                            on:click|preventDefault={deleteAvatar}
                            >Delete avatar</button
                        >
                    {/if}
                </div>
            </div>
        </form>
    </div>
</div>
