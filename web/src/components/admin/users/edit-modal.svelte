<script>
    import { createEventDispatcher } from "svelte";

    export let token;
    export let user;
    let newPassword = '';

    let isOpen = false;
    export function open() {
        isOpen = true;
    }
    export function close() {
        isOpen = false;
    }

    const dispatch = createEventDispatcher();
    async function editUser() {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${user.id}`,
            {
                method: "POST",
                headers: {
                    Authorization: `Bearer ${token}`,
                },
                body: new URLSearchParams({
                    username: user.username,
                    email: user.email,
                    password: newPassword,
                    role: user.role,
                    theme: user.theme
                }),
            }
        );
        const updatedUser = await response.json();
        if (response.status == 200) {
            close();
            dispatch("updateUser", { user: updatedUser });
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
            <p class="modal-card-title">Edit user</p>
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
                        <label class="label" for="edit-username"
                            >Username</label
                        >
                        <div class="control">
                            <input
                                class="input"
                                type="text"
                                id="edit-username"
                                bind:value={user.username}
                                required
                            />
                        </div>
                    </div>
                </div>

                <div class="column">
                    <div class="field">
                        <label class="label" for="edit-email">Email</label>
                        <div class="control">
                            <input
                                class="input"
                                type="email"
                                id="edit-email"
                                bind:value={user.email}
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

            <div class="columns">
                <div class="column">
                    <div class="field">
                        <label class="label" for="edit-role">Role</label>
                        <div class="control">
                            <div class="select is-fullwidth">
                                <select
                                    id="edit-role"
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
                        <label class="label" for="edit-theme">Theme</label>
                        <div class="control">
                            <div class="select is-fullwidth">
                                <select
                                    id="edit-theme"
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
            <button type="submit" class="button is-link">Edit user</button>
            <button class="button" on:click|preventDefault={close}
                >Cancel</button
            >
        </footer>
    </div>
</form>
