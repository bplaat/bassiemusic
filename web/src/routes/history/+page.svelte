<script>
    import TracksTable from '../../components/tracks-table.svelte';
    import { lazyLoader } from '../../utils.js';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Play History - BassieMusic',
            header: 'Play History',
            empty: 'You have not listened to any tracks',
        },
        nl: {
            title: 'Speel Geschiedenis - BassieMusic',
            header: 'Speel Geschiedenis',
            empty: 'Je hebt naar geen enkel nummer geluisterd',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;

    // Lazy loader
    lazyLoader(
        data.total,
        () => data.tracks.length,
        async (page) => {
            const response = await fetch(
                `${import.meta.env.VITE_API_URL}/users/${data.authUser.id}/played_tracks?${new URLSearchParams({
                    page,
                })}`,
                {
                    headers: {
                        Authorization: `Bearer ${data.token}`,
                    },
                }
            );
            const { data: newTracks } = await response.json();
            data.tracks = [...data.tracks, ...newTracks];
        }
    );
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<h1 class="title">{t('header')}</h1>

{#if data.tracks.length > 0}
    <TracksTable token={data.token} authUser={data.authUser} tracks={data.tracks} />
{:else}
    <p><i>{t('empty')}</i></p>
{/if}
