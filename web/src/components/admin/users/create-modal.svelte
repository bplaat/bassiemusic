<script>
    import { createEventDispatcher } from 'svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            header: 'Create new user',
            username: 'Username',
            email: 'Email address',
            password: 'Password',
            role: 'Role',
            role_normal: 'Normal',
            role_admin: 'Admin',
            create: 'Create new user',
            cancel: 'Cancel',
        },
        nl: {
            header: 'Maak nieuw gebruiker aan',
            username: 'Gebruikersnaam',
            email: 'Email adres',
            password: 'Wachtwoord',
            role: 'Rol',
            role_normal: 'Normaal',
            role_admin: 'Admin',
            create: 'Maak nieuwe gebruiker aan',
            cancel: 'Annuleren',
        },
    };
    const t = (key) => lang[$language][key];

    // Props
    export let token;

    // State
    let user = {
        username: '',
        email: '',
        password: '',
        role: 'normal',
        theme: 'dark',
    };
    let isOpen = false;

    // Methods
    export function open() {
        isOpen = true;
    }

    export function close() {
        isOpen = false;
    }

    const dispatch = createEventDispatcher();
    async function editUser() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/users`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                username: user.username,
                email: user.email,
                password: user.password,
                role: user.role,
            }),
        });
        if (response.status == 200) {
            const createdUser = await response.json();
            close();
            dispatch('createUser', { user: createdUser });
        }
    }
</script>

<form class="modal" class:is-active={isOpen} on:submit|preventDefault={editUser} style="z-index: 99999;">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">{t('header')}</p>
            <button type="button" class="delete" aria-label="close" on:click={close} />
        </header>
        <section class="modal-card-body">
            <div class="columns">
                <div class="column">
                    <div class="field">
                        <label class="label" for="create-username">{t('username')}</label>
                        <div class="control">
                            <input class="input" type="text" id="create-username" bind:value={user.username} required />
                        </div>
                    </div>
                </div>

                <div class="column">
                    <div class="field">
                        <label class="label" for="create-email">{t('email')}</label>
                        <div class="control">
                            <input class="input" type="email" id="create-email" bind:value={user.email} required />
                        </div>
                    </div>
                </div>
            </div>

            <div class="field">
                <label class="label" for="create-password">{t('password')}</label>
                <div class="control">
                    <input class="input" type="password" id="create-password" bind:value={user.password} />
                </div>
            </div>

            <div class="field">
                <label class="label" for="create-role">{t('role')}</label>
                <div class="control">
                    <div class="select is-fullwidth">
                        <select id="create-role" bind:value={user.role} required>
                            <option value="normal">{t('role_normal')}</option>
                            <option value="admin">{t('role_admin')}</option>
                        </select>
                    </div>
                </div>
            </div>
        </section>
        <footer class="modal-card-foot">
            <button type="submit" class="button is-link">{t('create')}</button>
            <button class="button" on:click|preventDefault={close}>{t('cancel')}</button>
        </footer>
    </div>
</form>
