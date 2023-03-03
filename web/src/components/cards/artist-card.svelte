<script>
    import { musicPlayer } from '../../stores.js';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            image_alt: 'Image of artist $1',
        },
        nl: {
            image_alt: 'Afbeelding van artiest $1',
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
    <div class="card-image" style="aspect-ratio: 1;">
        <img src={artist.medium_image} alt={t('image_alt', artist.name)} />
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
