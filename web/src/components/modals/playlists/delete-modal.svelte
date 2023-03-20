<script>
    import { createEventDispatcher } from 'svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            header: 'Delete playlist',
            text: 'Do your really want to delete this playlist?',
            delete: 'Delete playlist',
            cancel: 'Cancel',
        },
        nl: {
            header: 'Verwijder afspeellijst',
            text: 'Weet je zeker dat je deze afspeellijst wilt verwijderen?',
            delete: 'Verwijder afspeellijst',
            cancel: 'Annuleren',
        },
    };
    const t = (key) => lang[$language][key];

    // Props
    export let token;
    export let playlist;

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
    async function deletePlaylist() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/playlists/${playlist.id}`, {
            method: 'DELETE',
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        if (response.status == 200) {
            close();
            dispatch('deletePlaylist', { playlist });
        }
    }
</script>

<form class="modal" class:is-active={isOpen} on:submit|preventDefault={deletePlaylist} style="z-index: 99999;">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">{t('header')}</p>
            <button type="button" class="delete" aria-label="close" on:click={close} />
        </header>
        <section class="modal-card-body">
            <p>{t('text')}</p>
        </section>
        <footer class="modal-card-foot">
            <button type="submit" class="button is-danger">{t('delete')}</button>
            <button class="button" on:click|preventDefault={close}>{t('cancel')}</button>
        </footer>
    </div>
</form>
