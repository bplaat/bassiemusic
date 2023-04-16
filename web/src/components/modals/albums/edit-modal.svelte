<script>
    import { createEventDispatcher } from 'svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            header: 'Edit album',
            title: 'Title',
            type: 'Type',
            type_album: 'Album',
            type_ep: 'EP',
            type_single: 'Single',
            released_at: 'Released at',
            explicit: 'Explicit lyrics',
            explicit_description: 'This album contains explicit lyrics',
            deezer_id: 'Deezer ID',
            edit: 'Edit album',
            cancel: 'Cancel',
        },
        nl: {
            header: 'Verander album',
            title: 'Titel',
            type: 'Type',
            type_album: 'Album',
            type_ep: 'EP',
            type_single: 'Single',
            released_at: 'Released at',
            explicit: 'Expliciete teksten',
            explicit_description: 'Dit album bevat expliciete teksten',
            edit: 'Verander album',
            cancel: 'Annuleren',
        },
    };
    const t = (key) => lang[$language][key];

    // Props
    export let token;
    export let album;
    let errors = {};

    // State
    let isOpen = false;

    // Methods
    export function open() {
        isOpen = true;
    }

    export function close() {
        isOpen = false;
    }

    const dispatch = createEventDispatcher();
    async function editAlbum() {
        // Edit album
        const response = await fetch(`${import.meta.env.VITE_API_URL}/albums/${album.id}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                title: album.title,
                type: album.type,
                released_at: album.released_at,
                explicit: album.explicit,
                deezer_id: album.deezer_id,
            }),
        });

        if (response.status == 200) {
            const updatedAlbum = await response.json();
            album.title = updatedAlbum.title;
            album.type = updatedAlbum.type;
            album.released_at = updatedAlbum.released_at;
            album.explicit = updatedAlbum.explicit;
            album.deezer_id = updatedAlbum.deezer_id;
            close();
            dispatch('update', { album });
        } else {
            const data = await response.json();
            errors = data.errors;
        }
    }
</script>

<form class="modal" class:is-active={isOpen} on:submit|preventDefault={editAlbum}>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">{t('header')}</p>
            <button type="button" class="delete" aria-label="close" on:click={close} />
        </header>
        <section class="modal-card-body">
            <div class="field">
                <label class="label" for="albums_edit_title">{t('title')}</label>
                <div class="control">
                    <input
                        class="input"
                        class:is-danger={'title' in errors}
                        type="text"
                        id="albums_edit_title"
                        bind:value={album.title}
                        required
                    />
                </div>
            </div>

            <div class="field">
                <label class="label" for="albums_edit_type">{t('type')}</label>
                <div class="control">
                    <div class="select is-fullwidth">
                        <select id="albums_edit_type" bind:value={album.type} required>
                            <option value="album">{t('type_album')}</option>
                            <option value="ep">{t('type_ep')}</option>
                            <option value="single">{t('type_single')}</option>
                        </select>
                    </div>
                </div>
            </div>

            <div class="field">
                <label class="label" for="albums_edit_released_at">{t('released_at')}</label>
                <div class="control">
                    <input
                        class="input"
                        class:is-danger={'released_at' in errors}
                        type="text"
                        id="albums_edit_released_at"
                        bind:value={album.released_at}
                        required
                    />
                </div>
            </div>

            <div class="field">
                <label class="label" for="artists_edit_explicit">{t('explicit')}</label>
                <label class="checkbox" for="artists_edit_explicit">
                    <input type="checkbox" id="artists_edit_explicit" bind:checked={album.explicit} />
                    {t('explicit_description')}
                </label>
            </div>

            <div class="field">
                <label class="label" for="albums_edit_deezer_id">{t('deezer_id')}</label>
                <div class="control">
                    <input
                        class="input"
                        class:is-danger={'deezer_id' in errors}
                        type="number"
                        id="albums_edit_deezer_id"
                        bind:value={album.deezer_id}
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
