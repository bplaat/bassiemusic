<script>
    import { page } from '$app/stores';
    import { tick, onMount } from 'svelte';
    import LikeButton from './like-button.svelte';
    import DeleteModal from './modals/delete-modal.svelte';
    import { sidebar, musicPlayer, musicState, language } from '../stores.js';
    import { formatDuration } from '../filters.js';

    // Language strings
    const lang = {
        en: {
            index: '#',
            title: 'Title',
            duration: 'Duration',
            plays: 'Plays',
            album: 'Album',
            disk: 'Disk $1',
            play_track: 'Play track',
            cover_alt: 'Cover of album $1',
            explicit: 'Explicit lyrics',
            like: 'Like track',
            remove_like: 'Remove track like',
            track: 'track',
            options: 'Track options',

            add_queue: 'Add track to play queue',
            remove_queue: 'Remove track from play queue',
            go_to_album: 'Go to album',
            go_to_artist: 'Go to artist',
            add_to_playlist: 'Add to playlist',
            playlists_empty: 'You have no playlists',
            remove_from_playlist: 'Remove from playlist',
            delete: 'Delete track',
        },
        nl: {
            index: '#',
            title: 'Titel',
            duration: 'Duratie',
            plays: 'Plays',
            album: 'Album',
            disk: 'Disk $1',
            play_track: 'Speel track',
            cover_alt: 'Hoes van album $1',
            explicit: 'Expliciete songtekst',
            like: 'Like track',
            remove_like: 'Verwijder track like',
            track: 'track',
            options: 'Track opties',

            add_queue: 'Voeg track toe aan wachtrij',
            remove_queue: 'Verwijder track van wachtrij',
            go_to_album: 'Ga naar album',
            go_to_artist: 'Ga naar artiest',
            context_menu: 'Open context menu',
            add_to_playlist: 'Voeg toe aan afspeellijst',
            playlists_empty: 'Je hebt geen afspeellijsten',
            remove_from_playlist: 'Verwijder van afspeellijst',
            delete: 'Verwijder track',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // Props
    export let token;
    export let authUser;
    export let tracks;
    export let isAlbum = false;
    export let inPlaylist = null;
    export let isMusicQueue = false;

    // State
    $: isMultiDisk = tracks.find((track) => track.disk != 1) != null;
    let isContextmenuOpen = false;
    let contextmenu;
    let contextmenuTrack;
    let contextmenuPosition;
    let lastPlaylists = [];
    let deleteModal;

    // On mount
    if (isAlbum) {
        onMount(() => {
            const trackPosition = $page.url.hash.substring(1);
            if (trackPosition) {
                const trackRow = document.getElementById(trackPosition);
                if (trackRow != null) {
                    document.querySelector('.app').scrollTop = trackRow.offsetTop;
                }
            }
        });
    }

    // Methods
    function playTrack(track) {
        if (authUser.allow_explicit) {
            $musicPlayer.playTracks(
                tracks.filter((otherTrack) => otherTrack.music != null),
                track
            );
        } else {
            $musicPlayer.playTracks(
                tracks.filter((otherTrack) => otherTrack.music != null && !otherTrack.explicit),
                track
            );
        }
    }

    export function playFirstTrack() {
        let firstTrack = authUser.allow_explicit
            ? tracks.find((otherTrack) => otherTrack.music != null)
            : tracks.find((otherTrack) => otherTrack.music != null && !otherTrack.explicit);
        if (firstTrack != null) {
            playTrack(firstTrack);
        }
    }

    async function fetchPlaylists() {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users/${authUser.id}/playlists?sort_by=updated_at_desc&limit=10`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data } = await response.json();
        lastPlaylists = data;
    }

    async function openContextmenu(track, position, x, y) {
        isContextmenuOpen = true;
        contextmenuTrack = track;
        contextmenuPosition = position;
        fetchPlaylists();
        await tick();
        const app = document.querySelector('.app');
        contextmenu.style.left = `${
            x + contextmenu.offsetWidth + 32 >= app.offsetWidth ? x - contextmenu.offsetWidth : x
        }px`;
        contextmenu.style.top = `${
            y + contextmenu.offsetHeight >= app.scrollTop + app.offsetHeight ? y - contextmenu.offsetHeight : y
        }px`;
    }

    function likeTrack(track) {
        fetch(`${import.meta.env.VITE_API_URL}/tracks/${track.id}/like`, {
            method: track.liked ? 'DELETE' : 'PUT',
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        track.liked = !track.liked;
        tracks = tracks;
    }

    async function appendTrackToPlaylist(playlist) {
        await fetch(`${import.meta.env.VITE_API_URL}/playlists/${playlist.id}/tracks`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                track_id: contextmenuTrack.id,
            }),
        });
        $sidebar.updateLastPlaylists();
    }

    async function removeTrackFromPlaylist(position) {
        await fetch(`${import.meta.env.VITE_API_URL}/playlists/${inPlaylist.id}/tracks/${position}`, {
            method: 'DELETE',
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        $sidebar.updateLastPlaylists();
        tracks.splice(position - 1, 1);
        tracks = tracks;
    }
</script>

<svelte:window
    on:click={() => (isContextmenuOpen = false)}
    on:resize={() => (isContextmenuOpen = false)}
    on:wheel={() => (isContextmenuOpen = false)}
/>

<!-- Tracks table -->
<table class="table" style="width: 100%; table-layout: fixed;">
    <thead>
        {#if isAlbum}
            <th style="width: 10%;"><div class="track-index">{t('index')}</div></th>
            <th style="width: 50%;">{t('title')}</th>
            <th style="width: 20%;">{t('duration')}</th>
            <th class="is-hidden-mobile">{t('plays')}</th>
            <th style="width: calc(40px + .75em);" />
            <th style="width: calc(40px + .75em);" class:is-hidden-mobile={!isMusicQueue} />
        {:else}
            <th style="width: 10%;"><div class="track-index">{t('index')}</div></th>
            <th style="width: calc(64px + 1.5em);">{t('title')}</th>
            <th class="track-title" />
            <th style="width: 30%;" class="is-hidden-mobile">{t('album')}</th>
            <th style="width: 15%;">{t('duration')}</th>
            <th style="width: 15%;" class="is-hidden-mobile">{t('plays')}</th>
            <th style="width: calc(40px + .75em);" />
            <th style="width: calc(40px + .75em);" class:is-hidden-mobile={!isMusicQueue} />
        {/if}
    </thead>
    <tbody>
        {#each tracks as track, index}
            {#if isAlbum && isMultiDisk && (index == 0 || track.disk != tracks[index - 1].disk)}
                <tr>
                    <td>
                        <svg class="icon is-colored" viewBox="0 0 24 24">
                            <path
                                fill="#777"
                                d="M12,14C10.89,14 10,13.1 10,12C10,10.89 10.89,10 12,10C13.11,10 14,10.89 14,12A2,2 0 0,1 12,14M12,4A8,8 0 0,0 4,12A8,8 0 0,0 12,20A8,8 0 0,0 20,12A8,8 0 0,0 12,4Z"
                            />
                        </svg>
                    </td>
                    <td style="height: 64px; font-weight: 500; color: #777;">{t('disk', track.disk)}</td>
                    <td />
                    <td class="is-hidden-mobile" />
                    <td class="is-hidden-mobile" />
                    <td />
                </tr>
            {/if}

            <tr
                id={isAlbum ? `${track.disk}-${track.position}` : undefined}
                class="track-container"
                class:disabled={track.music == null || (!authUser.allow_explicit && track.explicit)}
                on:contextmenu={(event) =>
                    openContextmenu(
                        track,
                        index + 1,
                        event.clientX - document.querySelector('.app').offsetLeft,
                        event.clientY +
                            document.querySelector('.app').scrollTop -
                            document.querySelector('.app').offsetTop
                    )}
                on:dblclick|preventDefault={() => playTrack(track)}
                class:has-background-light={$musicState.track != undefined && $musicState.track.id == track.id}
            >
                <td>
                    <div class="track-index">{isAlbum ? track.position : index + 1}</div>
                    <button
                        class="button is-small track-play"
                        on:click={() => playTrack(track)}
                        title={t('play_track')}
                    >
                        <svg class="icon" viewBox="0 0 24 24">
                            <path d="M8,5.14V19.14L19,12.14L8,5.14Z" />
                        </svg>
                    </button>
                </td>
                {#if !isAlbum}
                    <td>
                        <a
                            href="/albums/{track.album.id}"
                            class="box has-image m-0 p-0"
                            style="width: 64px; height: 64px;"
                        >
                            <img
                                src={track.album.small_cover || '/images/album-default.svg'}
                                alt={t('cover_alt', track.album)}
                            />
                        </a>
                    </td>
                {/if}
                <td>
                    <p class="ellipsis mb-1" style="font-weight: 500;">
                        {#if isAlbum}
                            <!-- svelte-ignore a11y-invalid-attribute -->
                            <a href="#" on:click|preventDefault={playTrack(track)}>{track.title}</a>
                        {:else}
                            <a href="/albums/{track.album.id}#{track.disk}-{track.position}">{track.title}</a>
                        {/if}
                    </p>
                    <p class="ellipsis">
                        {#if track.explicit}
                            <span class="tag is-danger mr-1" title={t('explicit')}>E</span>
                        {/if}
                        {#each track.artists as artist}
                            <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
                        {/each}
                    </p>
                </td>
                {#if !isAlbum}
                    <td class="ellipsis is-hidden-mobile">
                        <a href="/albums/{track.album.id}">{track.album.title}</a>
                    </td>
                {/if}
                <td>{formatDuration(track.duration)}</td>
                <td class="is-hidden-mobile">{track.plays}</td>
                <td class="px-0 is-hidden-mobile">
                    <LikeButton token={token} item={track} itemRoute="tracks" itemLabel={t('track')} />
                </td>
                <td class="pl-0">
                    <button
                        class="button"
                        on:click|stopPropagation={(event) =>
                            openContextmenu(
                                track,
                                index + 1,
                                event.target.offsetLeft + event.target.offsetWidth - contextmenu.offsetWidth,
                                event.target.offsetTop + event.target.offsetHeight
                            )}
                        title={t('options')}
                    >
                        <svg class="icon" viewBox="0 0 24 24" style="pointer-events: none;">
                            <path
                                d="M12,16A2,2 0 0,1 14,18A2,2 0 0,1 12,20A2,2 0 0,1 10,18A2,2 0 0,1 12,16M12,10A2,2 0 0,1 14,12A2,2 0 0,1 12,14A2,2 0 0,1 10,12A2,2 0 0,1 12,10M12,4A2,2 0 0,1 14,6A2,2 0 0,1 12,8A2,2 0 0,1 10,6A2,2 0 0,1 12,4Z"
                            />
                        </svg>
                    </button>
                </td>
            </tr>
        {/each}
    </tbody>
</table>

<!-- Tracks context menu -->
<div
    bind:this={contextmenu}
    class="contextmenu dropdown-content"
    class:hidden={!isContextmenuOpen}
    style="position: absolute; z-index: 99999;"
>
    {#if contextmenuTrack != null}
        {#if !isMusicQueue}
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
                class="dropdown-item"
                href="#"
                on:click|preventDefault={() => $musicPlayer.addTrack(contextmenuTrack)}
                title={t('add_queue')}
            >
                {t('add_queue')}
            </a>
        {:else}
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
                class="dropdown-item"
                class:disabled={$musicState.track != undefined && $musicState.track.id == contextmenuTrack.id}
                href="#"
                on:click|preventDefault={() => $musicPlayer.removeTrack(contextmenuTrack)}
            >
                {t('remove_queue')}
            </a>
        {/if}

        <a class="dropdown-item" href="/albums/{contextmenuTrack.album.id}">
            {t('go_to_album')}
        </a>

        {#if contextmenuTrack.artists.length > 1}
            <div class="dropdown is-hoverable" style="width: 100%;">
                <div class="dropdown-trigger dropdown-item" style="width: 100%;">
                    {t('go_to_artist')}
                    <svg class="icon is-inline is-pulled-right" viewBox="0 0 24 24">
                        <path d="M8.59,16.58L13.17,12L8.59,7.41L10,6L16,12L10,18L8.59,16.58Z" />
                    </svg>
                </div>
                <div class="dropdown-menu" style="top: 0; left: 32px; width: 100%;">
                    <div class="dropdown-content">
                        {#each contextmenuTrack.artists as artist}
                            <a class="dropdown-item ellipsis" href="/artists/{artist.id}">
                                {artist.name}
                            </a>
                        {/each}
                    </div>
                </div>
            </div>
        {:else}
            <a class="dropdown-item" href="/artists/{contextmenuTrack.artists[0].id}">
                {t('go_to_artist')}
            </a>
        {/if}

        <hr class="dropdown-divider" />

        {#if !contextmenuTrack.liked}
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a class="dropdown-item" href="#" on:click|preventDefault={() => likeTrack(contextmenuTrack)}>
                {t('like')}
            </a>
        {:else}
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a class="dropdown-item" href="#" on:click|preventDefault={() => likeTrack(contextmenuTrack)}>
                {t('remove_like')}
            </a>
        {/if}

        <div class="dropdown is-hoverable" style="width: 100%;">
            <div class="dropdown-trigger dropdown-item" style="width: 100%;">
                {t('add_to_playlist')}
                <svg class="icon is-inline is-pulled-right" viewBox="0 0 24 24">
                    <path d="M8.59,16.58L13.17,12L8.59,7.41L10,6L16,12L10,18L8.59,16.58Z" />
                </svg>
            </div>
            <div class="dropdown-menu" style="top: auto; bottom: 0; left: 32px; width: 100%;">
                <div class="dropdown-content">
                    {#each lastPlaylists as playlist}
                        <!-- svelte-ignore a11y-invalid-attribute -->
                        <a
                            class="dropdown-item ellipsis"
                            href="#"
                            on:click|preventDefault={() => appendTrackToPlaylist(playlist, contextmenuTrack)}
                        >
                            {playlist.name}
                        </a>
                    {/each}
                    {#if lastPlaylists.length == 0}
                        <p class="dropdown-item"><i>{t('playlists_empty')}</i></p>
                    {/if}
                </div>
            </div>
        </div>

        {#if inPlaylist != null && (authUser.role == 'admin' || inPlaylist.user.id == authUser.id)}
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
                class="dropdown-item"
                href="#"
                on:click|preventDefault={() => removeTrackFromPlaylist(contextmenuPosition)}
            >
                {t('remove_from_playlist')}
            </a>
        {/if}

        {#if authUser.role == 'admin'}
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a class="dropdown-item" href="#" on:click|preventDefault={() => deleteModal.open()}>
                {t('delete')}
            </a>
        {/if}
    {/if}
</div>

{#if authUser.role == 'admin'}
    <DeleteModal
        bind:this={deleteModal}
        {token}
        item={contextmenuTrack}
        itemRoute="tracks"
        itemLabel={t('track')}
        on:delete={() => {
            tracks = tracks.filter((track) => track.id != contextmenuTrack.id);
        }}
    />
{/if}

<style>
    .track-title {
        width: 50%;
    }
    @media (min-width: 768px) {
        .track-title {
            width: 30%;
        }
    }

    .track-index {
        width: 30px;
        text-align: center;
    }
    .track-play,
    .track-container:hover .track-index {
        display: none;
    }
    .track-index,
    .track-container:hover .track-play {
        display: block;
    }
</style>
