<script>
    import { musicPlayer } from '../../stores.js';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            image_alt: 'Image of artist $1',
            sync: 'This aritist is synced, we will download automatic new albums'
        },
        nl: {
            image_alt: 'Afbeelding van artiest $1',
            sync: 'Deze artiest is gesynchroniseerd, we zullen automatisch nieuwe albums downloaden'
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // Props
    export let artist;
    export let token;
    export let authUser;

    // Methods
    async function fetchAndPlayTracks() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/artists/${artist.id}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        const completeArtist = await response.json();
        let playableTracks = authUser.allow_explicit
            ? completeArtist.top_tracks.filter((track) => track.music != null)
            : completeArtist.top_tracks.filter((track) => track.music != null && !track.explicit);
        if (playableTracks.length > 0) {
            $musicPlayer.playTracks(playableTracks, playableTracks[0]);
        }
    }
</script>

<a class="card has-image-play-button" href="/artists/{artist.id}">
    <div class="card-image has-image-tags">
        <figure class="image is-1by1">
            <img src={artist.medium_image || '/images/avatar-default.svg'} alt={t('image_alt', artist.name)} />
        </figure>
        <div class="image-tags">
            {#if artist.sync}
                <span class="tag px-2 py-1" style="height: auto;" title={t('sync')}>
                    <svg class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                        <path
                            d="M12,18A6,6 0 0,1 6,12C6,11 6.25,10.03 6.7,9.2L5.24,7.74C4.46,8.97 4,10.43 4,12A8,8 0 0,0 12,20V23L16,19L12,15M12,4V1L8,5L12,9V6A6,6 0 0,1 18,12C18,13 17.75,13.97 17.3,14.8L18.76,16.26C19.54,15.03 20,13.57 20,12A8,8 0 0,0 12,4Z"
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
        <h3 class="title is-6 ellipsis">{artist.name}</h3>
    </div>
</a>
