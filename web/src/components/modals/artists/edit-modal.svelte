<script>
    import { createEventDispatcher } from 'svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            header: 'Edit artist',
            name: 'Name',
            sync: 'Sync artist',
            sync_description: 'Automatic download new albums from this artist',
            deezer_id: 'Deezer ID',
            edit: 'Edit artist',
            cancel: 'Cancel',
        },
        nl: {
            header: 'Verander artiest',
            name: 'Naam',
            sync: 'Sync artiest',
            sync_description: 'Download automatisch nieuwe albums van deze artiest',
            deezer_id: 'Deezer ID',
            edit: 'Verander artiest',
            cancel: 'Annuleren',
        },
    };
    const t = (key) => lang[$language][key];

    // Props
    export let token;
    export let artist;

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
    async function editArtist() {
        // Edit artist
        const response = await fetch(`${import.meta.env.VITE_API_URL}/artists/${artist.id}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                name: artist.name,
                sync: artist.sync,
                deezer_id: artist.deezer_id,
            }),
        });

        if (response.status === 200) {
            const updatedArtist = await response.json();
            artist.name = updatedArtist.name;
            artist.sync = updatedArtist.sync;
            artist.deezer_id = updatedArtist.deezer_id;
            close();
            dispatch('update', { artist });
        } else {
            const data = await response.json();
            errors = data.errors;
        }
    }
</script>

<form class="modal" class:is-active={isOpen} on:submit|preventDefault={editArtist}>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">{t('header')}</p>
            <button type="button" class="delete" aria-label="close" on:click={close} />
        </header>
        <section class="modal-card-body">
            <div class="field">
                <label class="label" for="artists_edit_name">{t('name')}</label>
                <div class="control">
                    <input
                        class="input"
                        class:is-danger={'name' in errors}
                        type="text"
                        id="artists_edit_name"
                        bind:value={artist.name}
                        required
                    />
                </div>
            </div>

            <div class="field">
                <label class="label" for="artists_edit_sync">{t('sync')}</label>
                <label class="checkbox" for="artists_edit_sync">
                    <input type="checkbox" id="artists_edit_sync" bind:checked={artist.sync} />
                    {t('sync_description')}
                </label>
            </div>

            <div class="field">
                <label class="label" for="artists_edit_deezer_id">{t('deezer_id')}</label>
                <div class="control">
                    <input
                        class="input"
                        class:is-danger={'deezer_id' in errors}
                        type="number"
                        id="artists_edit_deezer_id"
                        bind:value={artist.deezer_id}
                        required
                    />
                </div>
            </div>
        </section>
        <footer class="modal-card-foot">
            <button type="submit" class="button is-link">{t('edit')}</button>
            <button class="button" on:click|preventDefault={close}>{t('cancel')}</button>
        </footer>
    </div>
</form>
