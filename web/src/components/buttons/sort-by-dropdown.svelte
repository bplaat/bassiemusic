<script>
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            sort_by: 'Sort by',
        },
        nl: {
            sort_by: 'Sorteer op',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let sortBy;
    export let options;
    let open = false;
</script>

<svelte:window on:click={() => (open = false)} on:resize={() => (open = false)} on:wheel={() => (open = false)} />

<div class="dropdown is-pulled-right-tablet is-right-tablet" class:is-active={open}>
    <div class="dropdown-trigger">
        <button class="button" on:click|stopPropagation={() => (open = !open)}>
            <svg class="icon mr-2" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path d="M18 21L14 17H17V7H14L18 3L22 7H19V17H22M2 19V17H12V19M2 13V11H9V13M2 7V5H6V7H2Z" />
            </svg>
            {t('sort_by')}
        </button>
    </div>
    <div class="dropdown-menu">
        <div class="dropdown-content">
            {#each Object.keys(options) as key}
                <a href="?sort_by={key}" class="dropdown-item" class:is-active={sortBy === key}>
                    {options[key]}
                </a>
            {/each}
        </div>
    </div>
</div>
