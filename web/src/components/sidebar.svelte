<script>
    import { page } from '$app/stores';
    import { goto } from '$app/navigation';
    import { language } from '../stores.js';

    // Language strings
    const lang = {
        en: {
            home: 'Home',

            library: 'Library',
            search: 'Search',
            artists: 'Artists',
            genres: 'Genres',
            albums: 'Albums',
            tracks: 'Tracks',
            liked: 'Liked',
            history: 'History',

            playlists: 'Playlists',
            playlists_your: 'Your playlists',
            playlist_create: 'Create playlist',
            playlist_image_alt: 'Image of playlist $1',
            playlist_untitled: 'Untitled playlist',

            admin: 'Admin',
            users: 'Users',
            downloader: 'Downloader',

            avatar_alt: 'Avatar of user $1',
            settings: 'Settings',
            logout: 'Logout',
            made_by: 'Made with $1 by $2',
            love: 'love',
        },
        nl: {
            home: 'Home',

            library: 'Bibliotheek',
            search: 'Zoeken',
            artists: 'Artiesten',
            genres: 'Genres',
            albums: 'Albums',
            tracks: 'Tracks',
            liked: 'Geliked',
            history: 'Geschiedenis',

            playlists: 'Afspeellijsten',
            playlists_your: 'Jouw afspeellijsten',
            playlist_create: 'Maak afspeellijst',
            playlist_image_alt: 'Afbeelding van afspeellijst $1',
            playlist_untitled: 'Naamloze afspeellijst',

            admin: 'Admin',
            users: 'Gebruikers',
            downloader: 'Downloader',

            avatar_alt: 'Avatar van gebruiker $1',
            settings: 'Instellingen',
            logout: 'Log uit',
            made_by: 'Gemaakt met $1<br/> door $2',
            love: 'liefde',
        },
    };
    const t = (key, p1, p2) => lang[$language][key].replace('$1', p1).replace('$2', p2);

    // Props
    export let token;
    export let authUser;
    export let lastPlaylists;

    // State
    let isOpen = false;

    // Methods
    export function open() {
        isOpen = true;
    }
    export function close() {
        isOpen = false;
    }

    export async function updateLastPlaylists() {
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

    function gotoLikedPage() {
        goto(`/liked/${localStorage.getItem('liked-tab') || 'tracks'}`);
    }

    async function createPlaylist() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/playlists`, {
            method: 'POST',
            headers: {
                Authorization: `Bearer ${token}`,
            },
            body: new URLSearchParams({
                name: t('playlist_untitled'),
                public: false,
            }),
        });
        const playlist = await response.json();
        updateLastPlaylists();
        goto(`/playlists/${playlist.id}`);
    }

    async function logout() {
        await fetch(`${import.meta.env.VITE_API_URL}/auth/logout`, {
            method: 'PUT',
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        document.cookie = `token=; expires=${new Date(0).toUTCString()}`;
        window.location = '/auth/login';
    }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div class="modal-background is-hidden-desktop" class:is-hidden={!isOpen} on:click={close} style="z-index: 200;" />

<div class="sidebar box has-background-white-bis m-0" class:is-open={isOpen}>
    <h1 class="title is-4 mb-5">
        <a href="/">BassieMusic</a>
        <button class="delete is-pulled-right is-hidden-desktop" on:click={close} />
    </h1>
    <div class="menu">
        <p class="menu-label">{t('library')}</p>
        <ul class="menu-list mb-5">
            <li>
                <a href="/" class:is-active={$page.url.pathname === '/'}>
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                        <path d="M10,20V14H14V20H19V12H22L12,3L2,12H5V20H10Z" />
                    </svg>
                    {t('home')}
                </a>
            </li>
            <li>
                <a href="/search" class:is-active={$page.url.pathname.startsWith('/search')}>
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                        <path
                            d="M9.5,3A6.5,6.5 0 0,1 16,9.5C16,11.11 15.41,12.59 14.44,13.73L14.71,14H15.5L20.5,19L19,20.5L14,15.5V14.71L13.73,14.44C12.59,15.41 11.11,16 9.5,16A6.5,6.5 0 0,1 3,9.5A6.5,6.5 0 0,1 9.5,3M9.5,5C7,5 5,7 5,9.5C5,12 7,14 9.5,14C12,14 14,12 14,9.5C14,7 12,5 9.5,5Z"
                        />
                    </svg>
                    {t('search')}
                </a>
            </li>
            <li>
                <a href="/artists" class:is-active={$page.url.pathname.startsWith('/artists')}>
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                        <path
                            d="M9,3A4,4 0 0,1 13,7H5A4,4 0 0,1 9,3M11.84,9.82L11,18H10V19A2,2 0 0,0 12,21A2,2 0 0,0 14,19V14A4,4 0 0,1 18,10H20L19,11L20,12H18A2,2 0 0,0 16,14V19A4,4 0 0,1 12,23A4,4 0 0,1 8,19V18H7L6.16,9.82C5.67,9.32 5.31,8.7 5.13,8H12.87C12.69,8.7 12.33,9.32 11.84,9.82M9,11A1,1 0 0,0 8,12A1,1 0 0,0 9,13A1,1 0 0,0 10,12A1,1 0 0,0 9,11Z"
                        />
                    </svg>
                    {t('artists')}
                </a>
            </li>
            <li>
                <a href="/genres" class:is-active={$page.url.pathname.startsWith('/genres')}>
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                        <path
                            d="M11,13.5V21.5H3V13.5H11M12,2L17.5,11H6.5L12,2M17.5,13C20,13 22,15 22,17.5C22,20 20,22 17.5,22C15,22 13,20 13,17.5C13,15 15,13 17.5,13Z"
                        />
                    </svg>
                    {t('genres')}
                </a>
            </li>
            <li>
                <a href="/albums" class:is-active={$page.url.pathname.startsWith('/albums')}>
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                        <path
                            d="M12,14C10.89,14 10,13.1 10,12C10,10.89 10.89,10 12,10C13.11,10 14,10.89 14,12A2,2 0 0,1 12,14M12,4A8,8 0 0,0 4,12A8,8 0 0,0 12,20A8,8 0 0,0 20,12A8,8 0 0,0 12,4Z"
                        />
                    </svg>
                    {t('albums')}
                </a>
            </li>
            <li>
                <a href="/tracks" class:is-active={$page.url.pathname === '/tracks'}>
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                        <path
                            d="M21,3V15.5A3.5,3.5 0 0,1 17.5,19A3.5,3.5 0 0,1 14,15.5A3.5,3.5 0 0,1 17.5,12C18.04,12 18.55,12.12 19,12.34V6.47L9,8.6V17.5A3.5,3.5 0 0,1 5.5,21A3.5,3.5 0 0,1 2,17.5A3.5,3.5 0 0,1 5.5,14C6.04,14 6.55,14.12 7,14.34V6L21,3Z"
                        />
                    </svg>
                    {t('tracks')}
                </a>
            </li>
            <li>
                <a href="/playlists" class:is-active={$page.url.pathname === '/playlists'}>
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                        <path
                            d="M15,6H3V8H15V6M15,10H3V12H15V10M3,16H11V14H3V16M17,6V14.18C16.69,14.07 16.35,14 16,14A3,3 0 0,0 13,17A3,3 0 0,0 16,20A3,3 0 0,0 19,17V8H22V6H17Z"
                        />
                    </svg>
                    {t('playlists')}
                </a>
            </li>
            <li>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a
                    href="#"
                    on:click|preventDefault={gotoLikedPage}
                    class:is-active={$page.url.pathname.startsWith('/liked')}
                >
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                        <path
                            d="M13.5,20C6.9,13.9 3.5,10.8 3.5,7.1C3.5,4 5.9,1.6 9,1.6C10.7,1.6 12.4,2.4 13.5,3.7C14.6,2.4 16.3,1.6 18,1.6C21.1,1.6 23.5,4 23.5,7.1C23.5,10.9 20.1,14 13.5,20M12,21.1C5.4,15.2 1.5,11.7 1.5,7C1.5,6.8 1.5,6.6 1.5,6.4C0.9,7.3 0.5,8.4 0.5,9.6C0.5,13.4 3.9,16.5 10.5,22.4L12,21.1Z"
                        />
                    </svg>
                    {t('liked')}
                </a>
            </li>
            <li>
                <a href="/history" class:is-active={$page.url.pathname === '/history'}>
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                        <path
                            d="M13.5,8H12V13L16.28,15.54L17,14.33L13.5,12.25V8M13,3A9,9 0 0,0 4,12H1L4.96,16.03L9,12H6A7,7 0 0,1 13,5A7,7 0 0,1 20,12A7,7 0 0,1 13,19C11.07,19 9.32,18.21 8.06,16.94L6.64,18.36C8.27,20 10.5,21 13,21A9,9 0 0,0 22,12A9,9 0 0,0 13,3"
                        />
                    </svg>
                    {t('history')}
                </a>
            </li>
        </ul>

        <p class="menu-label">{t('playlists')}</p>
        <ul class="menu-list mb-5">
            <li>
                <a href="/your_playlists" class:is-active={$page.url.pathname === '/your_playlists'}>
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                        <path
                            d="M15,6H3V8H15V6M15,10H3V12H15V10M3,16H11V14H3V16M17,6V14.18C16.69,14.07 16.35,14 16,14A3,3 0 0,0 13,17A3,3 0 0,0 16,20A3,3 0 0,0 19,17V8H22V6H17Z"
                        />
                    </svg>
                    {t('playlists_your')}
                </a>
            </li>
            <li>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a href="#" on:click|preventDefault={createPlaylist}>
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                        <path d="M3 16H10V14H3M18 14V10H16V14H12V16H16V20H18V16H22V14M14 6H3V8H14M14 10H3V12H14V10Z" />
                    </svg>
                    {t('playlist_create')}
                </a>
            </li>
            {#each lastPlaylists as playlist}
                <li>
                    <a
                        href="/playlists/{playlist.id}"
                        class="ellipsis"
                        class:is-active={$page.url.pathname === `/playlists/${playlist.id}`}
                    >
                        <img
                            class="icon is-inline mr-2"
                            style="background-color: #ccc; border-radius: 3px; box-shadow: 0 .5em 1em -0.125em rgba(10,10,10,.1), 0 0px 0 1px rgba(10,10,10,.02);"
                            src={playlist.small_image || '/images/album-default.svg'}
                            alt={t('playlist_image_alt', playlist.name)}
                        />
                        {playlist.name}
                    </a>
                </li>
            {/each}
        </ul>

        {#if authUser.role === 'admin'}
            <p class="menu-label">{t('admin')}</p>
            <ul class="menu-list mb-5">
                <li>
                    <a href="/admin/users" class:is-active={$page.url.pathname === '/admin/users'}>
                        <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                            <path
                                d="M16.5,12A2.5,2.5 0 0,0 19,9.5A2.5,2.5 0 0,0 16.5,7A2.5,2.5 0 0,0 14,9.5A2.5,2.5 0 0,0 16.5,12M9,11A3,3 0 0,0 12,8A3,3 0 0,0 9,5A3,3 0 0,0 6,8A3,3 0 0,0 9,11M16.5,14C14.67,14 11,14.92 11,16.75V19H22V16.75C22,14.92 18.33,14 16.5,14M9,13C6.67,13 2,14.17 2,16.5V19H9V16.75C9,15.9 9.33,14.41 11.37,13.28C10.5,13.1 9.66,13 9,13Z"
                            />
                        </svg>
                        {t('users')}
                    </a>
                </li>
                <li>
                    <a href="/admin/downloader" class:is-active={$page.url.pathname === '/admin/downloader'}>
                        <svg class="icon is-inline mr-2" viewBox="0 0 24 24">
                            <path
                                d="M2 12H4V17H20V12H22V17C22 18.11 21.11 19 20 19H4C2.9 19 2 18.11 2 17V12M12 15L17.55 9.54L16.13 8.13L13 11.25V2H11V11.25L7.88 8.13L6.46 9.55L12 15Z"
                            />
                        </svg>
                        {t('downloader')}
                    </a>
                </li>
            </ul>
        {/if}
    </div>

    <div class="flex" />

    <div class="media mb-5">
        <div class="media-left">
            <div class="box has-image m-0 p-0" style="width: 48px; height: 48px;">
                <img
                    src={authUser.small_avatar || '/images/avatar-default.svg'}
                    alt="Avatar of user {authUser.username}"
                />
            </div>
        </div>
        <div class="media-content">
            <p><b>{authUser.username}</b></p>
            <p>
                <a href="/settings" class="mr-2">{t('settings')}</a>
                <!-- svelte-ignore a11y-invalid-attribute -->
                <a href="#" on:click|preventDefault={logout}>{t('logout')}</a>
            </p>
        </div>
    </div>

    <p>
        {@html t(
            'made_by',
            `<svg class="icon is-inline is-colored" viewBox="0 0 24 24" style="width: 16px; height: 16px;" title=${t(
                'love'
            )}>
            <path fill="#f14668" d="M12,21.35L10.55,20.03C5.4,15.36 2,12.27 2,8.5C2,5.41 4.42,3 7.5,3C9.24,3 10.91,3.81 12,5.08C13.09,3.81 14.76,3 16.5,3C19.58,3 22,5.41 22,8.5C22,12.27 18.6,15.36 13.45,20.03L12,21.35Z" />
            </svg>`,
            `<a href="https://www.plaatsoft.nl/" target="_blank" rel="noreferrer">PlaatSoft</a>`
        )}
    </p>
</div>

<style>
    .sidebar {
        position: fixed;
        top: 0;
        left: -16.5rem;
        width: 16.5rem;
        height: 100%;
        z-index: 300;
        display: flex;
        flex-direction: column;
        border-radius: 0;
        overflow-y: auto;
        transition: left 0.1s ease-in-out;
    }
    .sidebar.is-open {
        left: 0;
    }
    :global(.app.is-resizing .sidebar) {
        transition: none !important;
    }
    :global(.app.is-macos-app .sidebar) {
        padding-top: calc(28px + 1.25rem) !important;
    }
    :global(.app.macos-is-fullscreen .sidebar) {
        padding-top: calc(1.25rem) !important;
    }

    .title > a {
        text-decoration: none !important;
    }

    @media (min-width: 1024px) {
        .sidebar {
            left: 0 !important;
            z-index: 100;
        }
        :global(.app.is-playing .sidebar) {
            height: calc(100% - 6rem) !important;
        }
    }
</style>
