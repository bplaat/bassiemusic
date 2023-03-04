<script>
    import DeleteAccountModal from '../../components/delete-account-modal.svelte';
    import { lazyLoader } from '../../utils.js';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Settings - BassieMusic',
            header: 'Settings',

            change_details: 'Change details',
            username: 'Username',
            email: 'Email address',
            password: 'Password (leave empty to not change)',
            language: 'Language',
            theme: 'Theme',
            theme_system: 'System',
            theme_light: 'Light',
            theme_dark: 'Dark',
            allow_explicit: 'Allow explicit content',
            allow_explicit_description: 'Allow playback of explicit content',

            change_avatar: 'Change avatar',
            avatar_file: 'Avatar file',
            avatar_file_help: 'You can upload an squared .jpg or .png image',
            delete_avatar: 'Delete avatar',

            delete_account: 'Delete account',

            sessions: 'Sessions management',
            current: 'CURRENT',
            location: 'With $1 at $2',
            unknown_location: 'Unknown location',
            logged_in_at: 'Logged in at: $1',
            expires_at: 'Expires at: $1',
            revoke_session: 'Revoke session',
        },
        nl: {
            title: 'Instellingen - BassieMusic',
            header: 'Instellingen',

            change_details: 'Verander gegevens',
            username: 'Gebruikersnaam',
            email: 'Email adres',
            password: 'Wachtwoord (laat leeg om niet te veranderen)',
            language: 'Taal',
            theme: 'Thema',
            theme_system: 'Systeem',
            theme_light: 'Licht',
            theme_dark: 'Donker',
            allow_explicit: 'Sta expliciete inhoud toe',
            allow_explicit_description: 'Sta het afspelen van expliciete inhoud toe',

            change_avatar: 'Verander avatar',
            avatar_file: 'Avatar bestand',
            avatar_file_help: 'U kunt een vierkante .jpg of .png-afbeelding uploaden',
            delete_avatar: 'Verwijder avatar',

            delete_account: 'Verwijder account',

            sessions: 'Sessie beheer',
            current: 'HUIDIGE',
            location: 'Met $1 in $2',
            unknown_location: 'Onbekende locatie',
            logged_in_at: 'Ingelogt op: $1',
            expires_at: 'Verloopt op: $1',
            revoke_session: 'Sessie intrekken',
        },
    };
    const t = (key, p1, p2) => lang[$language][key].replace('$1', p1).replace('$2', p2);

    // State
    export let data;
    let { token, authUser, currentSessionId, sessions } = data;
    let deleteAccountModal;

    // Change details
    let newPassword = '';

    async function changeDetails() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                username: authUser.username,
                email: authUser.email,
                password: newPassword,
                language: authUser.language,
                theme: authUser.theme,
                allow_explicit: authUser.allow_explicit,
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
            method: 'PUT',
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
        const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}/avatar`, {
            method: 'DELETE',
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
    lazyLoader(
        data.sessionsTotal,
        () => data.sessions.length,
        async (page) => {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/users/${authUser.id}/active_sessions?${new URLSearchParams({
                    page,
                })}`,
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );
            const { data: newSessions } = await response.json();
            sessions = [...sessions, ...newSessions];
        }
    );

    async function revokeSession(session) {
        await fetch(`${import.meta.env.VITE_API_URL}/sessions/${session.id}/revoke`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        if (currentSessionId == session.id) {
            document.cookie = `token=; expires=${new Date(0).toUTCString()}`;
            window.location = '/auth/login';
        } else {
            sessions = sessions.filter((otherSession) => otherSession.id != session.id);
        }
    }
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<h2 class="title">{t('header')}</h2>

<div class="columns">
    <!-- Change details -->
    <div class="column">
        <form on:submit|preventDefault={changeDetails} class="box">
            <h3 class="title is-4">{t('change_details')}</h3>

            <div class="columns">
                <div class="column">
                    <div class="field">
                        <label class="label" for="username">{t('username')}</label>
                        <div class="control">
                            <input class="input" type="text" id="username" bind:value={authUser.username} required />
                        </div>
                    </div>
                </div>

                <div class="column">
                    <div class="field">
                        <label class="label" for="email">{t('email')}</label>
                        <div class="control">
                            <input class="input" type="email" id="email" bind:value={authUser.email} required />
                        </div>
                    </div>
                </div>
            </div>

            <div class="field">
                <label class="label" for="password">{t('password')}</label>
                <div class="control">
                    <input class="input" type="password" id="password" bind:value={newPassword} />
                </div>
            </div>

            <div class="columns">
                <div class="column">
                    <div class="field">
                        <label class="label" for="language">{t('language')}</label>
                        <div class="control">
                            <div class="select is-fullwidth">
                                <select id="language" bind:value={authUser.language} required>
                                    <option value="en">English</option>
                                    <option value="nl">Nederlands</option>
                                </select>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="column">
                    <div class="field">
                        <label class="label" for="theme">{t('theme')}</label>
                        <div class="control">
                            <div class="select is-fullwidth">
                                <select id="theme" bind:value={authUser.theme} required>
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
                <label class="label" for="allow_explicit">{t('allow_explicit')}</label>
                <label class="checkbox" for="allow_explicit">
                    <input type="checkbox" id="allow_explicit" bind:checked={authUser.allow_explicit} />
                    {t('allow_explicit_description')}
                </label>
            </div>

            <div class="field">
                <div class="buttons">
                    <button type="submit" class="button is-link">{t('change_details')}</button>
                </div>
            </div>
        </form>
    </div>

    <div class="column">
        <!-- Change avatar -->
        <form on:submit|preventDefault={changeAvatar} class="box">
            <h3 class="title is-4">{t('change_avatar')}</h3>

            <div class="field">
                <label class="label" for="avatar">{t('avatar_file')}</label>
                <div class="control">
                    <input class="input" type="file" id="avatar" accept="*.jpg,*.jpeg,*.png" bind:this={avatarInput} />
                    <p class="help">{t('avatar_file_help')}</p>
                </div>
            </div>

            <div class="field">
                <div class="buttons">
                    <button type="submit" class="button is-link">{t('change_avatar')}</button>

                    {#if authUser.small_avatar}
                        <button type="button" class="button is-danger" on:click|preventDefault={deleteAvatar}>
                            {t('delete_avatar')}
                        </button>
                    {/if}
                </div>
            </div>
        </form>

        <!-- Change avatar -->
        <div class="box">
            <h3 class="title is-4">{t('delete_account')}</h3>

            <div class="buttons">
                <button class="button is-danger" on:click={() => deleteAccountModal.open()}>
                    {t('delete_account')}
                </button>
            </div>
        </div>
    </div>
</div>

<!-- Sessions management -->
<div class="box">
    <h3 class="title is-4">{t('sessions')}</h3>

    <div class="columns is-multiline">
        {#each sessions as session}
            <div class="column is-half">
                <div class="box content">
                    <h3 class="title is-4">
                        {session.client_name} on {session.client_os}
                        {#if currentSessionId == session.id}
                            <span class="tag is-link is-pulled-right">{t('current')}</span>
                        {/if}
                    </h3>
                    <p>
                        {t(
                            'location',
                            session.ip,
                            session.ip_city != undefined && session.ip_country != undefined
                                ? `${session.ip_city}, ${session.ip_country.toUpperCase()}`
                                : t('unknown_location')
                        )}
                    </p>
                    <p>
                        {t('logged_in_at', new Date(session.created_at).toLocaleString(authUser.lang))}
                    </p>
                    <p>
                        {t('expires_at', new Date(session.expires_at).toLocaleString(authUser.lang))}
                    </p>
                    <div class="buttons">
                        <button
                            class="button is-danger"
                            on:click={() => revokeSession(session)}
                            title="This will revoke this active session, you will be logged out of this browser or client"
                        >
                            {t('revoke_session')}
                        </button>
                    </div>
                </div>
            </div>
        {/each}
    </div>
</div>

<DeleteAccountModal bind:this={deleteAccountModal} {token} {authUser} />
