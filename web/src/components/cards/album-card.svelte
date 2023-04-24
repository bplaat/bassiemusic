<script>
    import { goto } from '$app/navigation';
    import { musicPlayer, language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            cover_alt: 'Cover of album $1',
            explicit: 'Explicit lyrics',
        },
        nl: {
            cover_alt: 'Hoes van album $1',
            explicit: 'Expliciete songtekst',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // Props
    export let authUser;
    export let album;
    export let token;

    // Methods
    async function fetchAndPlayTracks() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/albums/${album.id}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        const completeAlbum = await response.json();
        const tracks = (
            authUser.allow_explicit
                ? completeAlbum.tracks.filter((track) => track.music !== null)
                : completeAlbum.tracks.filter((track) => track.music !== null && !track.explicit)
        ).map((track) => {
            track.album = album;
            return track;
        });
        if (tracks.length > 0) {
            $musicPlayer.playTracks(tracks, tracks[0]);
        }
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
    class="card has-image-play-button"
    class:disabled={!authUser.allow_explicit && album.explicit}
    on:click={() => goto(`/albums/${album.id}`)}
>
    <div class="card-image has-image-tags">
        <figure class="image is-1by1">
            <img src={album.medium_cover || '/images/album-default.svg'} alt={t('cover_alt', album.title)} />
        </figure>
        <div class="image-tags">
            {#if album.type === 'album'}
                <span class="tag">ALBUM</span>
            {/if}
            {#if album.type === 'ep'}
                <span class="tag">EP</span>
            {/if}
            {#if album.type === 'single'}
                <span class="tag">SINGLE</span>
            {/if}
            {#if album.explicit}
                <span class="tag is-danger" title={t('explicit')}>E</span>
            {/if}
        </div>
        <button class="button image-play-button" on:click|stopPropagation={fetchAndPlayTracks}>
            <svg class="icon" viewBox="0 0 24 24">
                <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
            </svg>
        </button>
    </div>
    <div class="card-content">
        <h3 class="title is-6 mb-2 ellipsis"><a href="/albums/{album.id}">{album.title}</a></h3>
        <p class="ellipsis">
            {#each album.artists as artist}
                <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
            {/each}
        </p>
    </div>
</div>
