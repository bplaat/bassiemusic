<script>
    import { createEventDispatcher } from 'svelte';
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            header: 'Delete $1',
            text: 'Do your really want to delete this $1? This cannot be undone!',
            delete: 'Delete $1',
            cancel: 'Cancel',
        },
        nl: {
            header: 'Verwijder $1',
            text: 'Weet je zeker dat je deze $1 wilt verwijderen? Dit kan niet ongedaan gemaakt worden!',
            delete: 'Verwijder $1',
            cancel: 'Annuleren',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // Props
    export let token;
    export let itemRoute;
    export let itemLabel;
    export let item;

    // State
    let isOpen = false;

    // Methods
    export function open() {
        isOpen = true;
    }

    export function close() {
        isOpen = false;
    }

    const dispatch = createEventDispatcher();
    async function deleteItem() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/${itemRoute}/${item.id}`, {
            method: 'DELETE',
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        if (response.status == 200) {
            close();
            dispatch('delete');
        }
    }
</script>

<form class="modal" class:is-active={isOpen} on:submit|preventDefault={deleteItem}>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">{t('header', itemLabel)}</p>
            <button type="button" class="delete" aria-label="close" on:click={close} />
        </header>
        <section class="modal-card-body">
            <p>{t('text', itemLabel)}</p>
        </section>
        <footer class="modal-card-foot">
            <button type="submit" class="button is-danger">{t('delete', itemLabel)}</button>
            <button class="button" on:click|preventDefault={close}>{t('cancel')}</button>
        </footer>
    </div>
</form>
