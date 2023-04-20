<script>
    import { createEventDispatcher } from 'svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            header: 'Edit playlist',
            name: 'Name',
            public: 'Public playlist',
            public_description: 'Make this playlist public',
            edit: 'Edit playlist',
            cancel: 'Cancel',
        },
        nl: {
            header: 'Verander afspeellijst',
            name: 'Naam',
            public: 'Publieke afspeellijst',
            public_description: 'Maak deze afspeellijst publiekelijk',
            edit: 'Verander afspeellijst',
            cancel: 'Annuleren',
        },
    };
    const t = (key) => lang[$language][key];

    // Props
    export let token;
    export let playlist;

    // State
    let errors = {};
    let isOpen = false;

    // Methods
    export function open() {
        isOpen = true;
    }

    export function close() {
        isOpen = false;
    }

    const dispatch = createEventDispatcher();
    async function editPlaylist() {
        // Edit playlist
        const response = await fetch(`${import.meta.env.VITE_API_URL}/playlists/${playlist.id}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                name: playlist.name,
                public: playlist.public,
            }),
        });

        if (response.status == 200) {
            const updatedPlaylist = await response.json();
            playlist.name = updatedPlaylist.name;
            playlist.public = updatedPlaylist.public;
            close();
            dispatch('update', { playlist });
        } else {
            const data = await response.json();
            errors = data.errors;
        }
    }
</script>

<form class="modal" class:is-active={isOpen} on:submit|preventDefault={editPlaylist}>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">{t('header')}</p>
            <button type="button" class="delete" aria-label="close" on:click={close} />
        </header>
        <section class="modal-card-body">
            <div class="field">
                <label class="label" for="playlists_edit_name">{t('name')}</label>
                <div class="control">
                    <input
                        class="input"
                        class:is-danger={'name' in errors}
                        type="text"
                        id="playlists_edit_name"
                        bind:value={playlist.name}
                        required
                    />
                </div>
            </div>

            <div class="field">
                <label class="label" for="playlists_edit_public">{t('public')}</label>
                <label class="checkbox" for="playlists_edit_public">
                    <input type="checkbox" id="playlists_edit_public" bind:checked={playlist.public} />
                    {t('public_description')}
                </label>
            </div>
        </section>
        <footer class="modal-card-foot">
            <button type="submit" class="button is-link">{t('edit')}</button>
            <button class="button" on:click|preventDefault={close}>{t('cancel')}</button>
        </footer>
    </div>
</form>
