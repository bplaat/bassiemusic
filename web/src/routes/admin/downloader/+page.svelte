<script>
    import { onMount, onDestroy } from 'svelte';
    import { WEBSOCKET_RECONNECT_TIMEOUT } from '../../../consts.js';
    import { formatBytes } from '../../../filters.js';
    import { language } from '../../../stores.js';
    import DeleteModal from '../../../components/modals/delete-modal.svelte';

    // Language strings
    const lang = {
        en: {
            title: 'Downloader - Admin - BassieMusic',
            header: 'Admin Downloader',

            storage_size: 'Storage folder size',
            storage_used: 'Used: $1',
            storage_max: 'Max: $1',

            download_task_label: 'download task',
            download_tasks: 'Current download tasks',
            download_tasks_index: '#',
            download_tasks_type: 'Type',
            download_tasks_deezer_artist: 'Deezer Artist',
            download_tasks_deezer_album: 'Deezer Album',
            download_tasks_youtube_track: 'Youtube Track',
            download_tasks_update_deezer_artist: 'Update artist tracks',
            download_tasks_update_deezer_album: 'Update album tracks',
            download_tasks_display_name: 'Name',
            download_tasks_status: 'Status',
            download_tasks_status_pending: 'Pending',
            download_tasks_empty: 'There are no current download tasks',
            download_task_delete: 'Delete download task',

            search_header: 'Search and download albums and artists',
            query_placeholder: 'Find an album or artist...',
            search: 'Search',
            albums: 'Albums',
            cover_alt: 'Cover of album $1',
            add_album: 'Add album to BassieMusic',
            albums_empty: "Can't find any albums on Deezer",
            artists: 'Artists',
            image_alt: 'Image of artist $1',
            add_artist: 'Add artist to BassieMusic',
            artists_empty: "Can't find any artists on Deezer",
        },
        nl: {
            title: 'Downloader - Admin - BassieMusic',
            header: 'Admin Downloader',

            storage_size: 'Storage folder groote',
            storage_used: 'Gebruikt: $1',
            storage_max: 'Max: $1',

            download_task_label: 'download taak',
            download_tasks: 'Huidge download taken',
            download_tasks_index: '#',
            download_tasks_type: 'Type',
            download_tasks_deezer_artist: 'Deezer Artiest',
            download_tasks_deezer_album: 'Deezer Album',
            download_tasks_youtube_track: 'Youtube Track',
            download_tasks_update_deezer_artist: 'Update artiest nummers',
            download_tasks_update_deezer_album: 'Update album nummers',
            download_tasks_display_name: 'Naam',
            download_tasks_status: 'Status',
            download_tasks_status_pending: 'Wachtend',
            download_tasks_empty: 'Er zijn geen huidige download taken',
            download_task_delete: 'Verwijder download opdracht',

            search_header: 'Zoek en download albums en artisten',
            query_placeholder: 'Vind een album of artist...',
            search: 'Zoeken',
            albums: 'Albums',
            cover_alt: 'Hoes van album $1',
            add_album: 'Voeg album toe aan BassieMusic',
            albums_empty: 'Kan geen albums vinden op Deezer',
            artists: 'Artiesten',
            image_alt: 'Afbeelding van artist $1',
            add_artist: 'Voeg artist toe aan BassieMusic',
            artists_empty: 'Kan geen artisten vinden op Deezer',
            tasks_header: 'Huidige opdrachten',
        },
    };
    const t = (key, p1 = '') => lang[$language][key].replace('$1', p1);

    // State
    export let data;
    let query = '';
    let results = false;
    let albums = [];
    let artists = [];
    let deleteModal;
    let selectedTask;

    // Methods
    async function search() {
        if (query === '') {
            results = false;
            albums = [];
            artists = [];
            return;
        }

        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/deezer_search?${new URLSearchParams({
                q: query,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${data.token}`,
                },
            }
        );
        const result = await response.json();
        results = true;
        albums = result.albums;
        artists = result.artists;
    }

    async function downloadAlbum(album) {
        await fetch(`${import.meta.env.VITE_API_URL}/download/album`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${data.token}`,
            },
            body: new URLSearchParams({
                deezer_id: album.id,
                display_name: album.title,
            }),
        });
        albums = albums.filter((otherAlbum) => otherAlbum.id !== album.id);
    }

    async function downloadArtist(artist) {
        await fetch(`${import.meta.env.VITE_API_URL}/download/artist`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${data.token}`,
            },
            body: new URLSearchParams({
                deezer_id: artist.id,
                display_name: artist.name,
            }),
        });
        artists = artists.filter((otherArtist) => otherArtist.id !== artist.id);
    }

    // Download Tasks Logger
    let ws;
    let connected = false;
    let tasks = [];
    function websocketConnect() {
        ws = new WebSocket(import.meta.env.VITE_WEBSOCKET_URL);
        ws.onopen = () => {
            connected = true;
            ws.send(JSON.stringify({ type: 'auth.validate', data: { token: data.token } }));
            ws.send(JSON.stringify({ type: 'download_tasks.init' }));
        };
        ws.onmessage = (event) => {
            const { type, data } = JSON.parse(event.data);

            // Download Tasks messages
            if (type === 'download_tasks.init.response') {
                tasks = data;
            }
            if (type === 'download_tasks.create') {
                tasks = [...tasks, data];
            }
            if (type === 'download_tasks.update') {
                tasks = tasks.map((task) => {
                    if (task.id === data.id) return data;
                    return task;
                });
            }
            if (type === 'download_tasks.delete') {
                tasks = tasks.filter((task) => task.id !== data.id);
            }
        };
        ws.disconnect = () => {
            connected = false;
            setTimeout(websocketConnect, WEBSOCKET_RECONNECT_TIMEOUT);
        };
    }
    onMount(() => {
        websocketConnect();
    });
    onDestroy(() => {
        if (!connected) return;
        ws.close();
    });
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<h1 class="title">{t('header')}</h1>

<div class="box">
    <h2 class="title is-4">{t('storage_size')}</h2>
    <progress class="progress is-link" value={data.storage.used} max={data.storage.max}>
        {((data.storage.used / data.storage.max) * 100).toFixed(2)}%
    </progress>
    <p>
        <span class="mr-3">{t('storage_used', formatBytes(data.storage.used))}</span>
        <span>{t('storage_max', formatBytes(data.storage.max))}</span>
    </p>
</div>

<div class="box">
    <h2 class="title is-4">{t('download_tasks')}</h2>
    {#if tasks.length > 0}
        <table class="table">
            <thead>
                <th style="width: 10%;">{t('download_tasks_index')}</th>
                <th style="width: 35%;">{t('download_tasks_type')}</th>
                <th style="width: 35%;">{t('download_tasks_display_name')}</th>
                <th style="width: 30%;">{t('download_tasks_status')}</th>
                <th style="width: calc(40px + .75em);" />
            </thead>
            <tbody>
                {#each tasks as task, index}
                    <tr>
                        <td>{index + 1}</td>
                        <td class="ellipsis">
                            {#if task.type === 'deezer_artist'}
                                {t('download_tasks_deezer_artist')}
                            {/if}
                            {#if task.type === 'deezer_album'}
                                {t('download_tasks_deezer_album')}
                            {/if}
                            {#if task.type === 'youtube_track'}
                                {t('download_tasks_youtube_track')}
                            {/if}
                            {#if task.type === 'update_deezer_artist'}
                                {t('download_tasks_update_deezer_artist')}
                            {/if}
                            {#if task.type === 'update_deezer_album'}
                                {t('download_tasks_update_deezer_album')}
                            {/if}
                        </td>
                        <td class="ellipsis" style="font-weight: 500;">{task.display_name}</td>
                        <td>
                            {#if task.status === 'downloading'}
                                <progress class="progress is-link" style="width: 100%;" value={task.progress} max="100">
                                    {task.progress}%
                                </progress>
                            {:else}
                                <span class="ellipsis">{t('download_tasks_status_pending')}</span>
                            {/if}
                        </td>
                        <td>
                            <button
                                class="button"
                                on:click={() => {
                                    selectedTask = task;
                                    deleteModal.open();
                                }}
                                title={t('download_task_delete')}
                                disabled={task.status !== 'pending'}
                            >
                                <svg class="icon" viewBox="0 0 24 24">
                                    <path
                                        d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"
                                    />
                                </svg>
                            </button>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {:else}
        <p>{t('download_tasks_empty')}</p>
    {/if}
</div>

<div class="box">
    <h2 class="title is-4">{t('search_header')}</h2>

    <form on:submit|preventDefault={search} class="field has-addons">
        <div class="control" style="width: 100%;">
            <input class="input" type="text" bind:value={query} placeholder={t('query_placeholder')} />
        </div>
        <div class="control">
            <button type="submit" class="button is-link">{t('search')}</button>
        </div>
    </form>

    {#if results}
        <div class="columns mt-5">
            <div class="column is-half">
                <h2 class="title is-4">{t('albums')}</h2>
                {#each albums as album}
                    <div class="media">
                        <div class="media-left">
                            <div class="box m-0 p-0" style="width: 48px; height: 48px;">
                                <img src={album.cover_medium} alt={t('cover_alt', album.title)} />
                            </div>
                        </div>
                        <div class="media-content" style="min-width: 0;">
                            <p class="ellipsis" style="font-weight: 500;">{album.title}</p>
                            <p class="ellipsis">{album.artist.name}</p>
                        </div>
                        <div class="media-right">
                            <button class="button is-link" on:click={() => downloadAlbum(album)} title={t('add_album')}>
                                <svg class="icon" viewBox="0 0 24 24">
                                    <path d="M19,13H13V19H11V13H5V11H11V5H13V11H19V13Z" />
                                </svg>
                            </button>
                        </div>
                    </div>
                {/each}
                {#if albums.length === 0}
                    <p><i>{t('albums_empty')}</i></p>
                {/if}
            </div>

            <div class="column is-half">
                <h2 class="title is-4">{t('artists')}</h2>
                {#each artists as artist}
                    <div class="media">
                        <div class="media-left">
                            <div class="box m-0 p-0" style="width: 48px; height: 48px;">
                                <img src={artist.picture_medium} alt={t('image_alt', artist.name)} />
                            </div>
                        </div>
                        <div class="media-content">
                            <p class="ellipsis" style="font-weight: 500;">{artist.name}</p>
                        </div>
                        <button class="button is-link" on:click={() => downloadArtist(artist)} title={t('add_artist')}>
                            <svg class="icon" viewBox="0 0 24 24">
                                <path d="M19,13H13V19H11V13H5V11H11V5H13V11H19V13Z" />
                            </svg>
                        </button>
                    </div>
                {/each}
                {#if artists.length === 0}
                    <p><i>{t('artists_empty')}</i></p>
                {/if}
            </div>
        </div>
    {/if}
</div>

<DeleteModal
    bind:this={deleteModal}
    token={data.token}
    item={selectedTask}
    itemRoute="download"
    itemLabel={t('download_task_label')}
/>
