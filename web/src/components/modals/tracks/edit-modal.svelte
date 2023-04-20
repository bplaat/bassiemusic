<script>
    import { createEventDispatcher } from 'svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            header: 'Edit track',
            title: 'Title',
            album_id: 'Album',
            disk: 'Disk',
            position: 'Position',
            explicit: 'Explicit lyrics',
            explicit_description: 'This track contains explicit lyrics',
            deezer_id: 'Deezer ID',
            youtube_id: 'YouTube ID',
            edit: 'Edit track',
            cancel: 'Cancel',
        },
        nl: {
            header: 'Verander track',
            title: 'Titel',
            album_id: 'Album',
            disk: 'Disk',
            position: 'Positie',
            explicit: 'Expliciete teksten',
            explicit_description: 'Dit track bevat expliciete teksten',
            deezer_id: 'Deezer ID',
            youtube_id: 'YouTube ID',
            edit: 'Verander track',
            cancel: 'Annuleren',
        },
    };
    const t = (key) => lang[$language][key];

    // Props
    export let token;
    export let track;

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
    async function editTrack() {
        // Edit track
        const response = await fetch(`${import.meta.env.VITE_API_URL}/tracks/${track.id}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                title: track.title,
                disk: track.disk,
                position: track.position,
                explicit: track.explicit,
                deezer_id: track.deezer_id,
                youtube_id: track.youtube_id,
            }),
        });

        if (response.status === 200) {
            const updatedTrack = await response.json();
            track.title = updatedTrack.title;
            track.disk = updatedTrack.disk;
            track.position = updatedTrack.position;
            track.explicit = updatedTrack.explicit;
            track.deezer_id = updatedTrack.deezer_id;
            track.youtube_id = updatedTrack.youtube_id;
            close();
            dispatch('update', { track });
        } else {
            const data = await response.json();
            errors = data.errors;
        }
    }
</script>

<form class="modal" class:is-active={isOpen} on:submit|preventDefault={editTrack}>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">{t('header')}</p>
            <button type="button" class="delete" aria-label="close" on:click={close} />
        </header>
        <section class="modal-card-body">
            <div class="field">
                <label class="label" for="tracks_edit_title">{t('title')}</label>
                <div class="control">
                    <input
                        class="input"
                        class:is-danger={'title' in errors}
                        type="text"
                        id="tracks_edit_title"
                        bind:value={track.title}
                        required
                    />
                </div>
            </div>

            <div class="field">
                <!-- svelte-ignore a11y-label-has-associated-control -->
                <label class="label">{t('album_id')}</label>
                <div class="control">TODO</div>
            </div>

            <div class="columns">
                <div class="column">
                    <div class="field">
                        <label class="label" for="tracks_edit_disk">{t('disk')}</label>
                        <div class="control">
                            <input
                                class="input"
                                class:is-danger={'disk' in errors}
                                type="number"
                                id="tracks_edit_disk"
                                bind:value={track.disk}
                                required
                            />
                        </div>
                    </div>
                </div>

                <div class="column">
                    <div class="field">
                        <label class="label" for="tracks_edit_position">{t('position')}</label>
                        <div class="control">
                            <input
                                class="input"
                                class:is-danger={'position' in errors}
                                type="number"
                                id="tracks_edit_position"
                                bind:value={track.position}
                                required
                            />
                        </div>
                    </div>
                </div>
            </div>

            <div class="field">
                <label class="label" for="artists_edit_explicit">{t('explicit')}</label>
                <label class="checkbox" for="artists_edit_explicit">
                    <input type="checkbox" id="artists_edit_explicit" bind:checked={track.explicit} />
                    {t('explicit_description')}
                </label>
            </div>

            <div class="field">
                <label class="label" for="tracks_edit_deezer_id">{t('deezer_id')}</label>
                <div class="control">
                    <input
                        class="input"
                        class:is-danger={'deezer_id' in errors}
                        type="number"
                        id="tracks_edit_deezer_id"
                        bind:value={track.deezer_id}
                        required
                    />
                </div>
            </div>

            <div class="field">
                <label class="label" for="tracks_edit_youtube_id">{t('youtube_id')}</label>
                <div class="control">
                    <input
                        class="input"
                        class:is-danger={'youtube_id' in errors}
                        type="text"
                        id="tracks_edit_youtube_id"
                        bind:value={track.youtube_id}
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
