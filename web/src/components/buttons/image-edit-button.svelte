<script>
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            artist_image_alt: 'Image of artist $1',
            artist_sync: 'This aritist is synced, we will download automatic new albums',

            genre_image_alt: 'Image of genre $1',

            album_cover_alt: 'Cover of album $1',
            album_album: 'Album',
            album_ep: 'EP',
            album_single: 'Single',
            album_explicit: 'Explicit lyrics',

            playlist_image_alt: 'Image of playlist $1',
            playlist_public: 'Public playlist',

            edit_image: 'Edit image',
        },
        nl: {
            artist_image_alt: 'Afbeelding van artiest $1',
            artist_sync: 'Deze artiest is gesynchroniseerd, we zullen automatisch nieuwe albums downloaden',

            genre_image_alt: 'Afbeelding van genre $1',

            album_cover_alt: 'Cover of album $1',
            album_album: 'Album',
            album_ep: 'EP',
            album_single: 'Single',
            album_explicit: 'Explicit lyrics',

            playlist_image_alt: 'Afbeelding van afspeellijst $1',
            playlist_public: 'Publieke afspeellijst',

            edit_image: 'Verander afbeelding',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // State
    export let token;
    export let item;
    export let itemRoute;
    export let editable;

    // Methods
    function editImage() {
        const imageInput = document.createElement('input');
        imageInput.type = 'file';
        imageInput.accept = '*.jpg,*.jpeg,*.png';
        imageInput.addEventListener('change', async () => {
            const body = new FormData();
            if (itemRoute === 'artists' || itemRoute === 'genres' || itemRoute === 'playlists') {
                body.set('image', imageInput.files[0], imageInput.files[0].name);
            }
            if (itemRoute === 'albums') {
                body.set('cover', imageInput.files[0], imageInput.files[0].name);
            }
            const response = await fetch(`${import.meta.env.VITE_API_URL}/${itemRoute}/${item.id}`, {
                method: 'PUT',
                headers: {
                    Authorization: `Bearer ${token}`,
                },
                body,
            });
            if (response.status === 200) {
                const updatedItem = await response.json();
                if (itemRoute === 'artists' || itemRoute === 'genres' || itemRoute === 'playlists') {
                    item.small_image = updatedItem.small_image;
                    item.medium_image = updatedItem.medium_image;
                }
                if (itemRoute === 'artists' || itemRoute === 'genres') {
                    item.large_image = updatedItem.large_image;
                }
                if (itemRoute === 'albums') {
                    item.small_cover = updatedItem.small_cover;
                    item.medium_cover = updatedItem.medium_cover;
                    item.large_cover = updatedItem.large_cover;
                }
            }
        });
        imageInput.dispatchEvent(new MouseEvent('click'));
    }

    async function deleteImage() {
        const body = new URLSearchParams();
        if (itemRoute === 'artists' || itemRoute === 'genres' || itemRoute === 'playlists') {
            body.set('image', '');
        }
        if (itemRoute === 'albums') {
            body.set('cover', '');
        }
        await fetch(`${import.meta.env.VITE_API_URL}/${itemRoute}/${item.id}`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body,
        });
        if (itemRoute === 'artists' || itemRoute === 'genres' || itemRoute === 'playlists') {
            item.small_image = null;
            item.medium_image = null;
        }
        if (itemRoute === 'artists' || itemRoute === 'genres') {
            item.large_image = null;
        }
        if (itemRoute === 'albums') {
            item.small_cover = null;
            item.medium_cover = null;
            item.large_cover = null;
        }
    }
</script>

<div class="image-edit-button box has-image p-0" class:is-editable={editable}>
    <figure class="image is-1by1">
        {#if itemRoute === 'artists'}
            <img src={item.large_image || '/images/avatar-default.svg'} alt={t('artist_image_alt', item.name)} />
        {/if}
        {#if itemRoute === 'genres'}
            <img src={item.large_image || '/images/album-default.svg'} alt={t('genre_image_alt', item.name)} />
        {/if}
        {#if itemRoute === 'albums'}
            <img src={item.large_cover || '/images/album-default.svg'} alt={t('album_cover_alt', item.title)} />
        {/if}
        {#if itemRoute === 'playlists'}
            <img src={item.medium_image || '/images/album-default.svg'} alt={t('playlist_image_alt', item.name)} />
        {/if}
    </figure>

    <div class="image-tags">
        {#if itemRoute === 'artists' && item.sync}
            <span class="tag px-2 py-1" style="height: auto;" title={t('artist_sync')}>
                <svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                    <path
                        d="M12,18A6,6 0 0,1 6,12C6,11 6.25,10.03 6.7,9.2L5.24,7.74C4.46,8.97 4,10.43 4,12A8,8 0 0,0 12,20V23L16,19L12,15M12,4V1L8,5L12,9V6A6,6 0 0,1 18,12C18,13 17.75,13.97 17.3,14.8L18.76,16.26C19.54,15.03 20,13.57 20,12A8,8 0 0,0 12,4Z"
                    />
                </svg>
            </span>
        {/if}

        {#if itemRoute === 'albums'}
            {#if item.type === 'album'}
                <span class="tag" style="text-transform: uppercase;">{t('album_album')}</span>
            {/if}
            {#if item.type === 'ep'}
                <span class="tag" style="text-transform: uppercase;">{t('album_ep')}</span>
            {/if}
            {#if item.type === 'single'}
                <span class="tag" style="text-transform: uppercase;">{t('album_single')}</span>
            {/if}
            {#if item.explicit}
                <span class="tag is-danger" title={t('album_explicit')}>E</span>
            {/if}
        {/if}

        {#if itemRoute === 'playlists' && item.public}
            <span class="tag px-2 py-1" style="height: auto;" title={t('playlist_public')}>
                <svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                    <path
                        d="M16.36,14C16.44,13.34 16.5,12.68 16.5,12C16.5,11.32 16.44,10.66 16.36,10H19.74C19.9,10.64 20,11.31 20,12C20,12.69 19.9,13.36 19.74,14M14.59,19.56C15.19,18.45 15.65,17.25 15.97,16H18.92C17.96,17.65 16.43,18.93 14.59,19.56M14.34,14H9.66C9.56,13.34 9.5,12.68 9.5,12C9.5,11.32 9.56,10.65 9.66,10H14.34C14.43,10.65 14.5,11.32 14.5,12C14.5,12.68 14.43,13.34 14.34,14M12,19.96C11.17,18.76 10.5,17.43 10.09,16H13.91C13.5,17.43 12.83,18.76 12,19.96M8,8H5.08C6.03,6.34 7.57,5.06 9.4,4.44C8.8,5.55 8.35,6.75 8,8M5.08,16H8C8.35,17.25 8.8,18.45 9.4,19.56C7.57,18.93 6.03,17.65 5.08,16M4.26,14C4.1,13.36 4,12.69 4,12C4,11.31 4.1,10.64 4.26,10H7.64C7.56,10.66 7.5,11.32 7.5,12C7.5,12.68 7.56,13.34 7.64,14M12,4.03C12.83,5.23 13.5,6.57 13.91,8H10.09C10.5,6.57 11.17,5.23 12,4.03M18.92,8H15.97C15.65,6.75 15.19,5.55 14.59,4.44C16.43,5.07 17.96,6.34 18.92,8M12,2C6.47,2 2,6.5 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z"
                    />
                </svg>
            </span>
        {/if}
    </div>

    <!-- svelte-ignore a11y-click-events-have-key-events -->
    {#if editable}
        <div class="image-edit-button-overlay" on:click={editImage}>
            {#if item.small_image !== null || item.small_cover !== null}
                <div class="image-tags">
                    <!-- svelte-ignore a11y-click-events-have-key-events -->
                    <span class="delete" on:click|stopPropagation={deleteImage} />
                </div>
            {/if}

            <svg class="icon" viewBox="0 0 24 24" style="width: 48px; height: 48px;">
                <path
                    d="M20.71,7.04C21.1,6.65 21.1,6 20.71,5.63L18.37,3.29C18,2.9 17.35,2.9 16.96,3.29L15.12,5.12L18.87,8.87M3,17.25V21H6.75L17.81,9.93L14.06,6.18L3,17.25Z"
                />
            </svg>
            <div class="mt-1">{t('edit_image')}</div>
        </div>
    {/if}
</div>

<style>
    .image-edit-button {
        position: relative;
    }

    .image-edit-button-overlay {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        flex-direction: column;
        justify-content: center;
        background-color: rgba(10, 10, 10, 0.86);
        visibility: hidden;
        cursor: pointer;
    }
    .image-edit-button.is-editable:hover > .image-tags {
        visibility: hidden;
    }
    .image-edit-button.is-editable:hover > .image-edit-button-overlay {
        visibility: visible;
    }
</style>
