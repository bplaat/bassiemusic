<script>
    import { formatBytes } from '../../../filters.js';
    import { language } from '../../../stores.js';
    import { onMount } from 'svelte';

    // Language strings
    const lang = {
        en: {
            title: 'Downloader - Admin - BassieMusic',
            header: 'Admin Downloader',

            index: '#',
            type: "Type",
            display_name: "Name",
            status: "Status",
            progress: "progress",

            storage_size: 'Storage folder size',
            storage_used: 'Used: $1',
            storage_max: 'Max: $1',

            deezer_artist: "Artist",
            deezer_album: "Album",
            pending: "Pending",
            downloading: "Downloading",
            empty_tasks: "There are no download tasks",

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
            tasks_header: 'Current tasks',
        },
        nl: {
            title: 'Downloader - Admin - BassieMusic',
            header: 'Admin Downloader',

            index: '#',
            type: "Type",
            display_name: "Naam",
            status: "Status",
            progress: "Progressie",

            storage_size: 'Storage folder groote',
            storage_used: 'Gebruikt: $1',
            storage_max: 'Max: $1',

            deezer_artist: "Artiest",
            deezer_album: "Album",
            pending: "Wachten",
            downloading: "Aan het downloaden",
            empty_tasks: "Er zijn geen download opdrachten",

            search_header: 'Zoek en download albums en artisten',
            query_placeholder: 'Vind een album of artist...',
            search: 'Zoeken',
            albums: 'Albums',
            cover_alt: 'Hoes van album $1',
            add_album: 'Voeg album toe aan BassieMusic',
            albums_empty: 'Kan geen albums vinden op Deezer',
            artists: 'Artisten',
            image_alt: 'Afbeelding van artist $1',
            add_artist: 'Voeg artist toe aan BassieMusic',
            artists_empty: 'Kan geen artisten vinden op Deezer',
            tasks_header: 'Huidige opdrachten',
        },
    };
    const t = (key, p1 = '') => lang[$language][key].replace('$1', p1);

    // State
    export let data;
    let token = data.token;
    let query = '';
    let results = false;
    let albums = [];
    let artists = [];

    // Methods
    async function search() {
        if (query == '') {
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
                    Authorization: `Bearer ${token}`,
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
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                deezer_id: album.id,
                display_name: album.title,
            }),
        });
        albums = albums.filter((otherAlbum) => otherAlbum.id != album.id);
    }

    async function downloadArtist(artist) {
        await fetch(`${import.meta.env.VITE_API_URL}/download/artist`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                deezer_id: artist.id,
                display_name: artist.name,
            }),
        });
        artists = artists.filter((otherArtist) => otherArtist.id != artist.id);
    }

    // Logger
    function removeAt(arr, index) {
        var j = 0;
        var arr2 = [];
        for (var i = 0; i < arr.length; i++) {
            if (i != index) {
            arr2[j] = arr[i];
            j++;
            }
        }
        return arr2
    }

    let ws;
    let tasks = [];
    onMount(() => {
        ws = new WebSocket(import.meta.env.VITE_WEBSOCKET_URL);
        ws.onopen = () => {
            ws.onmessage = (event) => {
                let data = JSON.parse(event.data);

                if(data.type == 'allTasks'){
                    tasks = data['data']
                }else if(data.type == 'newTask'){
                    tasks.push(data['data'])
                }else if(data.type == 'taskUpdate'){
                    let newTasks = [];
                    tasks.forEach((task) => {
                        if(task.id == data['data'].id){
                            task = data['data']
                        }
                        newTasks.push(task)
                    })

                    tasks = newTasks;
                }else if(data.type == 'taskDelete'){
                    let i = 0;
                    tasks.forEach((task) => {
                        if(task.id == data['data'].id){
                            tasks = removeAt(tasks, i)
                        }
                        i += 1
                    })
                }

                tasks = tasks;
            }

            ws.send(JSON.stringify({type: 'auth', token}));
        };
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
    <h2 class="title is-4">{t('tasks_header')}</h2>
    {#key tasks}
        {#if tasks.length != 0}
            <table class="table" style="width: 100%; table-layout: fixed;">
                <thead>
                    <th style="width: 10%;">{t('index')}</th>
                    <th style="width: 20%;">{t('type')}</th>
                    <th style="width: 20%;">{t('display_name')}</th>
                    <th style="width: 20%; text-align: center;">{t('status')}</th>
                </thead>
                <tbody>
                    {#each tasks as task, index}
                        <tr>
                            <td>
                                <div>{index + 1}</div>
                            </td>
                            <td>
                                <p class="ellipsis">{t(task.type)}</p>
                            </td>
                            <td>
                                <p class="ellipsis mb-1" style="font-weight: 500;">{task.display_name}</p>
                            </td>
                            <td style="text-align: center;">
                                {#if task.status == "downloading"}
                                    <progress class="progress is-link" value={task.progress} max=100>
                                        {task.progress}%
                                    </progress>
                                {:else}
                                    <p class="ellipsis">{t(task.status)}</p>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        {:else}
            <p>{t('empty_tasks')}</p>
        {/if}
    {/key}
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
                {#if albums.length == 0}
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
                {#if artists.length == 0}
                    <p><i>{t('artists_empty')}</i></p>
                {/if}
            </div>
        </div>
    {/if}
</div>
