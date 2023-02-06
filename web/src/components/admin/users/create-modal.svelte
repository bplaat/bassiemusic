<script>
    import { createEventDispatcher } from "svelte";

    export let token;

    let user = {
        username: "",
        email: "",
        password: "",
        role: "normal",
        theme: "system",
    };

    let isOpen = false;
    export function open() {
        isOpen = true;
    }
    export function close() {
        isOpen = false;
    }

    const dispatch = createEventDispatcher();
    async function editUser() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/users`, {
            method: "POST",
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                username: user.username,
                email: user.email,
                password: user.password,
                role: user.role,
                theme: user.theme,
            }),
        });
        const createdUser = await response.json();
        if (response.status == 200) {
            close();
            dispatch("createUser", { user: createdUser });
        }
    }
</script>

<form
    class="modal"
    class:is-active={isOpen}
    on:submit|preventDefault={editUser}
    style="z-index: 99999;"
>
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">Create new user</p>
            <button
                type="button"
                class="delete"
                aria-label="close"
                on:click={close}
            />
        </header>
        <section class="modal-card-body">
            <div class="columns">
                <div class="column">
                    <div class="field">
                        <label class="label" for="create-username"
                            >Username</label
                        >
                        <div class="control">
                            <input
                                class="input"
                                type="text"
                                id="create-username"
                                bind:value={user.username}
                                required
                            />
                        </div>
                    </div>
                </div>

                <div class="column">
                    <div class="field">
                        <label class="label" for="create-email">Email</label>
                        <div class="control">
                            <input
                                class="input"
                                type="email"
                                id="create-email"
                                bind:value={user.email}
                                required
                            />
                        </div>
                    </div>
                </div>
            </div>

            <div class="field">
                <label class="label" for="create-password">Password</label>
                <div class="control">
                    <input
                        class="input"
                        type="password"
                        id="create-password"
                        bind:value={user.password}
                    />
                </div>
            </div>

            <div class="columns">
                <div class="column">
                    <div class="field">
                        <label class="label" for="create-role">Role</label>
                        <div class="control">
                            <div class="select is-fullwidth">
                                <select
                                    id="create-role"
                                    bind:value={user.role}
                                    required
                                >
                                    <option value="normal">Normal</option>
                                    <option value="admin">Admin</option>
                                </select>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="column">
                    <div class="field">
                        <label class="label" for="create-theme">Theme</label>
                        <div class="control">
                            <div class="select is-fullwidth">
                                <select
                                    id="create-theme"
                                    bind:value={user.theme}
                                    required
                                >
                                    <option value="system">System</option>
                                    <option value="light">Light</option>
                                    <option value="dark">Dark</option>
                                </select>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>
        <footer class="modal-card-foot">
            <button type="submit" class="button is-link">Create new user</button
            >
            <button class="button" on:click|preventDefault={close}
                >Cancel</button
            >
        </footer>
    </div>
</form>
