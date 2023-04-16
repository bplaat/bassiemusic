<script>
    import { createEventDispatcher } from 'svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            header: 'Edit genre',
            name: 'Name',
            deezer_id: 'Deezer ID',
            edit: 'Edit genre',
            cancel: 'Cancel',
        },
        nl: {
            header: 'Verander genre',
            name: 'Naam',
            deezer_id: 'Deezer ID',
            edit: 'Verander genre',
            cancel: 'Annuleren',
        },
    };
    const t = (key) => lang[$language][key];

    // Props
    export let token;
    export let genre;
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
    async function editGenre() {
        // Edit genre
        const response = await fetch(`${import.meta.env.VITE_API_URL}/genres/${genre.id}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                name: genre.name,
                deezer_id: genre.deezer_id,
            }),
        });

        if (response.status == 200) {
            const updatedGenre = await response.json();
            genre.name = updatedGenre.name;
            genre.deezer_id = updatedGenre.deezer_id;
            close();
            dispatch('update', { genre });
        } else {
            const data = await response.json();
            errors = data.errors;
        }
    }
</script>

<form class="modal" class:is-active={isOpen} on:submit|preventDefault={editGenre}>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">{t('header')}</p>
            <button type="button" class="delete" aria-label="close" on:click={close} />
        </header>
        <section class="modal-card-body">
            <div class="field">
                <label class="label" for="genres_edit_name">{t('name')}</label>
                <div class="control">
                    <input
                        class="input"
                        class:is-danger={'name' in errors}
                        type="text"
                        id="genres_edit_name"
                        bind:value={genre.name}
                        required
                    />
                </div>
            </div>

            <div class="field">
                <label class="label" for="genres_edit_deezer_id">{t('deezer_id')}</label>
                <div class="control">
                    <input
                        class="input"
                        class:is-danger={'deezer_id' in errors}
                        type="number"
                        id="genres_edit_deezer_id"
                        bind:value={genre.deezer_id}
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
