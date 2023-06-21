<script>
    import { musicPlayer, language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            image_alt: 'Image of playlist $1',
            public: 'Public playlist',
        },
        nl: {
            image_alt: 'Afbeelding van afspeellijst $1',
            public: 'Publieke afspeellijst',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // Props
    export let authUser;
    export let playlist;
    export let token;

    // Methods
    async function fetchAndPlayTracks() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/playlists/${playlist.id}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        const completePlaylist = await response.json();
        const tracks = authUser.allow_explicit
            ? completePlaylist.tracks.filter((track) => track.music !== null)
            : completePlaylist.tracks.filter((track) => track.music !== null && !track.explicit);
        if (tracks.length > 0) {
            $musicPlayer.playTracks(tracks, tracks[0]);
        }
    }
</script>

<a class="card has-image-play-button" href="/playlists/{playlist.id}">
    <div class="card-image has-image-tags">
        <figure class="image is-1by1">
            <img src={playlist.medium_image || '/images/album-default.svg'} alt={t('image_alt', playlist.name)} />
        </figure>
        <div class="image-tags">
            {#if playlist.public}
                <span class="tag px-2 py-1" style="height: auto;" title={t('public')}>
                    <svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                        <path
                            d="M16.36,14C16.44,13.34 16.5,12.68 16.5,12C16.5,11.32 16.44,10.66 16.36,10H19.74C19.9,10.64 20,11.31 20,12C20,12.69 19.9,13.36 19.74,14M14.59,19.56C15.19,18.45 15.65,17.25 15.97,16H18.92C17.96,17.65 16.43,18.93 14.59,19.56M14.34,14H9.66C9.56,13.34 9.5,12.68 9.5,12C9.5,11.32 9.56,10.65 9.66,10H14.34C14.43,10.65 14.5,11.32 14.5,12C14.5,12.68 14.43,13.34 14.34,14M12,19.96C11.17,18.76 10.5,17.43 10.09,16H13.91C13.5,17.43 12.83,18.76 12,19.96M8,8H5.08C6.03,6.34 7.57,5.06 9.4,4.44C8.8,5.55 8.35,6.75 8,8M5.08,16H8C8.35,17.25 8.8,18.45 9.4,19.56C7.57,18.93 6.03,17.65 5.08,16M4.26,14C4.1,13.36 4,12.69 4,12C4,11.31 4.1,10.64 4.26,10H7.64C7.56,10.66 7.5,11.32 7.5,12C7.5,12.68 7.56,13.34 7.64,14M12,4.03C12.83,5.23 13.5,6.57 13.91,8H10.09C10.5,6.57 11.17,5.23 12,4.03M18.92,8H15.97C15.65,6.75 15.19,5.55 14.59,4.44C16.43,5.07 17.96,6.34 18.92,8M12,2C6.47,2 2,6.5 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z"
                        />
                    </svg>
                </span>
            {/if}
        </div>

        <button class="button image-play-button" on:click|preventDefault={fetchAndPlayTracks}>
            <svg class="icon" viewBox="0 0 24 24">
                <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
            </svg>
        </button>
    </div>
    <div class="card-content">
        <h3 class="title is-6 ellipsis">{playlist.name}</h3>
        <p class="ellipsis">
            {#each playlist.owners as owner}
                <span class="mr-2">{owner.username}</span>
            {/each}
        </p>
    </div>
</a>
