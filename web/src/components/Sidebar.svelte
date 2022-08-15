<script>
import { page } from '$app/stores';
</script>

<div class="sidebar box has-background-white-bis m-0">
    <h1 class="title is-4 mb-5"><a href="/">BassieMusic</a></h1>
    <div class="menu">
        <p class="menu-label">Library</p>
        <ul class="menu-list mb-5">
            <li>
                <a href="/" class:is-active="{$page.url.pathname == '/'}">
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24"><path d="M10,20V14H14V20H19V12H22L12,3L2,12H5V20H10Z" /></svg>
                    Home
                </a>
            </li>
            <li>
                <a href="/search" class:is-active="{$page.url.pathname.startsWith('/search')}">
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24"><path d="M9.5,3A6.5,6.5 0 0,1 16,9.5C16,11.11 15.41,12.59 14.44,13.73L14.71,14H15.5L20.5,19L19,20.5L14,15.5V14.71L13.73,14.44C12.59,15.41 11.11,16 9.5,16A6.5,6.5 0 0,1 3,9.5A6.5,6.5 0 0,1 9.5,3M9.5,5C7,5 5,7 5,9.5C5,12 7,14 9.5,14C12,14 14,12 14,9.5C14,7 12,5 9.5,5Z" /></svg>
                    Search
                </a>
            </li>
            <li>
                <a href="/artists" class:is-active="{$page.url.pathname.startsWith('/artists')}">
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24"><path d="M9,3A4,4 0 0,1 13,7H5A4,4 0 0,1 9,3M11.84,9.82L11,18H10V19A2,2 0 0,0 12,21A2,2 0 0,0 14,19V14A4,4 0 0,1 18,10H20L19,11L20,12H18A2,2 0 0,0 16,14V19A4,4 0 0,1 12,23A4,4 0 0,1 8,19V18H7L6.16,9.82C5.67,9.32 5.31,8.7 5.13,8H12.87C12.69,8.7 12.33,9.32 11.84,9.82M9,11A1,1 0 0,0 8,12A1,1 0 0,0 9,13A1,1 0 0,0 10,12A1,1 0 0,0 9,11Z" /></svg>
                    Artists
                </a>
            </li>
            <li>
                <a href="/albums" class:is-active="{$page.url.pathname.startsWith('/albums')}">
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24"><path d="M12,14C10.89,14 10,13.1 10,12C10,10.89 10.89,10 12,10C13.11,10 14,10.89 14,12A2,2 0 0,1 12,14M12,4A8,8 0 0,0 4,12A8,8 0 0,0 12,20A8,8 0 0,0 20,12A8,8 0 0,0 12,4Z" /></svg>
                    Albums
                </a>
            </li>
            <li>
                <a href="/tracks" class:is-active="{$page.url.pathname == '/tracks'}">
                    <svg class="icon is-inline mr-2" viewBox="0 0 24 24"><path d="M21,3V15.5A3.5,3.5 0 0,1 17.5,19A3.5,3.5 0 0,1 14,15.5A3.5,3.5 0 0,1 17.5,12C18.04,12 18.55,12.12 19,12.34V6.47L9,8.6V17.5A3.5,3.5 0 0,1 5.5,21A3.5,3.5 0 0,1 2,17.5A3.5,3.5 0 0,1 5.5,14C6.04,14 6.55,14.12 7,14.34V6L21,3Z" /></svg>
                    Tracks
                </a>
            </li>
        </ul>
        <p class="menu-label">Playlists</p>
        <ul class="menu-list">
            <li>Comming soon...</li>
        </ul>
    </div>
    <div class="flex"></div>
    <p>
        Made with
        <svg class="icon is-inline" viewBox="0 0 24 24" style="width: 16px; height: 16px;"><path d="M12,21.35L10.55,20.03C5.4,15.36 2,12.27 2,8.5C2,5.41 4.42,3 7.5,3C9.24,3 10.91,3.81 12,5.08C13.09,3.81 14.76,3 16.5,3C19.58,3 22,5.41 22,8.5C22,12.27 18.6,15.36 13.45,20.03L12,21.35Z" /></svg>
        by<br/>
        <a href="https://bplaat.nl" target="_blank">Bastiaan van der Plaat</a>
    </p>
</div>
