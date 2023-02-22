<script>
    import { createEventDispatcher } from 'svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Delete user',
            text: 'Do your realy want to delete this user?',
            delete: 'Delete user',
            cancel: 'Cancel',
        },
        nl: {
            title: 'Verwijder gebruiker',
            text: 'Weet je zeker dat je deze gebruiker wilt verwijderen?',
            delete: 'Verwijder gebruiker',
            cancel: 'Annuleren',
        }
    };
    const t = (key) => lang[$language][key];

    // Props
    export let token;
    export let user;

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
    async function deleteUser() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${user.id}`, {
            method: 'DELETE',
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        if (response.status == 200) {
            close();
            dispatch('deleteUser', { user });
        }
    }
</script>

<form class="modal" class:is-active={isOpen} on:submit|preventDefault={deleteUser} style="z-index: 99999;">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">{t('title')}</p>
            <button type="button" class="delete" aria-label="close" on:click={close} />
        </header>
        <section class="modal-card-body">
            <p>{t('text')}</p>
        </section>
        <footer class="modal-card-foot">
            <button type="submit" class="button is-danger">{t('delete')}</button>
            <button class="button" on:click|preventDefault={close}>{t('cancel')}</button>
        </footer>
    </div>
</form>
