<script>
    import { createEventDispatcher } from 'svelte';
    import { language } from '../../../../stores.js';

    // Language strings
    const lang = {
        en: {
            header: 'Edit user',
            username: 'Username',
            email: 'Email address',
            password: 'Password (leave empty to not change)',
            role: 'Role',
            role_normal: 'Normal',
            role_admin: 'Admin',
            language: 'Language',
            theme: 'Theme',
            theme_system: 'System',
            theme_light: 'Light',
            theme_dark: 'Dark',
            allow_explicit: 'Allow explicit content',
            allow_explicit_description: 'Allow playback of explicit content',
            edit: 'Edit user',
            cancel: 'Cancel',
        },
        nl: {
            header: 'Verander gebruiker',
            username: 'Gebruikersnaam',
            email: 'Email adres',
            password: 'Wachtwoord (laat leeg om niet te veranderen)',
            role: 'Rol',
            role_normal: 'Normaal',
            role_admin: 'Admin',
            language: 'Taal',
            theme: 'Thema',
            theme_system: 'Systeem',
            theme_light: 'Licht',
            theme_dark: 'Donker',
            allow_explicit: 'Sta expliciete inhoud toe',
            allow_explicit_description: 'Sta het afspelen van expliciete inhoud toe',
            edit: 'Verander gebruiker',
            cancel: 'Annuleren',
        },
    };
    const t = (key) => lang[$language][key];

    // Props
    export let token;
    export let user;

    // State
    let newPassword = '';
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
        const body = new URLSearchParams({
            username: user.username,
            email: user.email,
            role: user.role,
            language: user.language,
            theme: user.theme,
            allow_explicit: user.allow_explicit,
        });
        if (newPassword != '') body.set('password', newPassword);
        const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${user.id}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body,
        });
        if (response.status == 200) {
            const updatedUser = await response.json();
            close();
            dispatch('updateUser', { user: updatedUser });
        } else {
            alert(`Error: ${JSON.stringify(await response.json())}`);
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
                        <label class="label" for="edit-username">{t('username')}</label>
                        <div class="control">
                            <input class="input" type="text" id="edit-username" bind:value={user.username} required />
                        </div>
                    </div>
                </div>

                <div class="column">
                    <div class="field">
                        <label class="label" for="edit-email">{t('email')}</label>
                        <div class="control">
                            <input class="input" type="email" id="edit-email" bind:value={user.email} required />
                        </div>
                    </div>
                </div>
            </div>

            <div class="field">
                <label class="label" for="edit-password">{t('password')}</label>
                <div class="control">
                    <input class="input" type="password" id="edit-password" bind:value={newPassword} />
                </div>
            </div>

            <div class="columns">
                <div class="column">
                    <div class="field">
                        <label class="label" for="edit-role">{t('role')}</label>
                        <div class="control">
                            <div class="select is-fullwidth">
                                <select id="edit-role" bind:value={user.role} required>
                                    <option value="normal">{t('role_normal')}</option>
                                    <option value="admin">{t('role_admin')}</option>
                                </select>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="column">
                    <div class="field">
                        <label class="label" for="edit-language">{t('language')}</label>
                        <div class="control">
                            <div class="select is-fullwidth">
                                <select id="edit-language" bind:value={user.language} required>
                                    <option value="en">English</option>
                                    <option value="nl">Nederlands</option>
                                </select>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="column">
                    <div class="field">
                        <label class="label" for="edit-theme">{t('theme')}</label>
                        <div class="control">
                            <div class="select is-fullwidth">
                                <select id="edit-theme" bind:value={user.theme} required>
                                    <option value="system">{t('theme_system')}</option>
                                    <option value="light">{t('theme_light')}</option>
                                    <option value="dark">{t('theme_dark')}</option>
                                </select>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="field">
                <label class="label" for="edit-allow_explicit">{t('allow_explicit')}</label>
                <label class="checkbox" for="edit-allow_explicit">
                    <input type="checkbox" id="edit-allow_explicit" bind:checked={user.allow_explicit} />
                    {t('allow_explicit_description')}
                </label>
            </div>
        </section>
        <footer class="modal-card-foot">
            <button type="submit" class="button is-link">{t('edit')}</button>
            <button class="button" on:click|preventDefault={close}>{t('cancel')}</button>
        </footer>
    </div>
</form>
