<script>
    import AlbumCard from '../components/cards/album-card.svelte';
    import ArtistCard from '../components/cards/artist-card.svelte';
    import TracksTable from '../components/tracks-table.svelte';
    import { language } from '../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Home - BassieMusic',
            hey: 'Hey, $1!',
            last_albums: 'Last played albums',
            last_artists: 'Last played artists',
            last_tracks: 'Last played tracks',
            empty: "You haven't listened to any tracks yet, use the sidebar to find something you like",
        },
        nl: {
            title: 'Home - BassieMusic',
            hey: 'Hoi, $1!',
            last_albums: 'Laast gespeelde albums',
            last_artists: 'Laast gespeelde artisten',
            last_tracks: 'Laast gespeelde tracks',
            empty: 'Je hebt nog geen nummers beluisterd, gebruik de zijbalk om iets te vinden dat je leuk vindt',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // Utils
    function uniques(items) {
        const uniques = {};
        for (const item of items) {
            if (!uniques[item.id]) {
                uniques[item.id] = item;
            }
        }
        return Object.values(uniques);
    }

    // State
    export let data;
    const { token, authUser, lastPlayedTracks } = data;

    $: lastPlayedAlbums = uniques(lastPlayedTracks.map((track) => track.album)).slice(0, 5);
    $: lastPlayedArtists = uniques(lastPlayedTracks.map((track) => track.artists).flat()).slice(0, 5);
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<h1 class="title">{t('hey', authUser.username)}</h1>

{#if lastPlayedTracks.length > 0}
    <h2 class="title is-4">{t('last_albums')}</h2>
    <div class="columns is-multiline is-mobile">
        {#each lastPlayedAlbums as album}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <AlbumCard {album} {token} {authUser} />
            </div>
        {/each}
    </div>

    <h2 class="title is-4 mt-5">{t('last_artists')}</h2>
    <div class="columns is-multiline is-mobile">
        {#each lastPlayedArtists as artist}
            <div class="column is-half-mobile is-one-third-tablet is-one-quarter-desktop is-one-fifth-widescreen">
                <ArtistCard {artist} {token} {authUser} />
            </div>
        {/each}
    </div>

    <h2 class="title is-4 mt-5">{t('last_tracks')}</h2>
    <TracksTable {token} {authUser} tracks={lastPlayedTracks} displayMax="5" />
{:else}
    <p>{t('empty')}</p>
{/if}
