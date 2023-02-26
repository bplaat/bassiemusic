<script>
    import TracksTable from '../../components/tracks-table.svelte';
    import { musicState } from '../../stores';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Play Queue - BassieMusic',
            back: 'Go back one page',
            header: 'Play Queue',
            empty: 'The music player queue is empty',
        },
        nl: {
            title: 'Wachtrij - BassieMusic',
            back: 'Ga een pagina terug',
            header: 'Wachtrij',
            empty: 'De muziek wachtrij is leeg',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;
    const { token, authUser } = data;
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<div class="buttons">
    <button class="button" on:click={() => history.back()} title={t('back')}>
        <svg class="icon" viewBox="0 0 24 24">
            <path d="M20,11V13H8L13.5,18.5L12.08,19.92L4.16,12L12.08,4.08L13.5,5.5L8,11H20Z" />
        </svg>
    </button>
</div>

<h2 class="title">{t('header')}</h2>
{#if $musicState.queue.length > 0}
    <TracksTable {token} {authUser} isMusicQueue={true} tracks={$musicState.queue} />
{:else}
    <p><i>{t('empty')}</i></p>
{/if}
