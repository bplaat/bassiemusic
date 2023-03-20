<script>
    import { createEventDispatcher } from 'svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            header: 'Edit playlist',
            name: 'Name',
            public: 'Public playlist',
            image: 'Image (leave empty to not change)',
            image_help: 'You can upload an squared .jpg or .png image',
            delete_image: 'Delete image',
            public_description: 'Make this playlist public',
            edit: 'Edit playlist',
            cancel: 'Cancel',
        },
        nl: {
            header: 'Verander afspeellijst',
            name: 'Name',
            image: 'Afbeelding (laat leeg om niet te veranderen)',
            image_help: 'U kunt een vierkante .jpg of .png-afbeelding uploaden',
            delete_image: 'Verwijder afbeelding',
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
    let imageInput;

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
    async function editPlaylist() {
        // Change playlist image
        if (imageInput.files[0]) {
            const formData = new FormData();
            formData.append('image', imageInput.files[0], imageInput.files[0].name);
            await fetch(`${import.meta.env.VITE_API_URL}/playlists/${playlist.id}/image`, {
                method: 'POST',
                headers: {
                    Authorization: `Bearer ${token}`,
                },
                body: formData,
            });
        }

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
            updatedPlaylist.user = playlist.user;
            updatedPlaylist.tracks = playlist.tracks;
            close();
            imageInput.value = '';
            dispatch('updatePlaylist', { playlist: updatedPlaylist });
        }
    }

    async function deleteImage() {
        await fetch(`${import.meta.env.VITE_API_URL}/playlists/${playlist.id}/image`, {
            method: 'DELETE',
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        playlist.small_image = null;
        playlist.medium_image = null;
        dispatch('updatePlaylist', { playlist: playlist });
    }
</script>

<form class="modal" class:is-active={isOpen} on:submit|preventDefault={editPlaylist} style="z-index: 99999;">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">{t('header')}</p>
            <button type="button" class="delete" aria-label="close" on:click={close} />
        </header>
        <section class="modal-card-body">
            <div class="field">
                <label class="label" for="edit-name">{t('name')}</label>
                <div class="control">
                    <input class="input" type="text" id="edit-name" bind:value={playlist.name} required />
                </div>
            </div>

            <div class="field">
                <label class="label" for="edit-image">{t('image')}</label>
                <div class="control">
                    <input
                        class="input"
                        type="file"
                        id="edit-image"
                        accept="*.jpg,*.jpeg,*.png"
                        bind:this={imageInput}
                    />
                    <p class="help">{t('image_help')}</p>
                </div>
            </div>

            {#if playlist.small_image}
                <div class="field">
                    <div class="buttons">
                        <button type="button" class="button is-danger" on:click|preventDefault={deleteImage}>
                            {t('delete_image')}
                        </button>
                    </div>
                </div>
            {/if}

            <div class="field">
                <label class="label" for="edit-public">{t('public')}</label>
                <label class="checkbox" for="edit-public">
                    <input type="checkbox" id="edit-public" bind:checked={playlist.public} />
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
