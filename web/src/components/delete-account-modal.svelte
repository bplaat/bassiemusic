<script>
    import { language } from '../stores.js';

    // Language strings
    const lang = {
        en: {
            header: 'Delete your account',
            text: "Do your really want to delete your account this can't be undone?",
            delete: 'Delete account',
            cancel: 'Cancel',
        },
        nl: {
            header: 'Verwijder jouw account',
            text: 'Wil je echt je account verwijderen maar kan dit niet ongedaan worden gemaakt?',
            delete: 'Verwijder account',
            cancel: 'Annuleren',
        },
    };
    const t = (key) => lang[$language][key];

    // Props
    export let token;
    export let authUser;

    // State
    let isOpen = false;

    // Methods
    export function open() {
        isOpen = true;
    }

    export function close() {
        isOpen = false;
    }

    async function deleteAccount() {
        await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}`, {
            method: 'DELETE',
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        document.cookie = `token=; expires=${new Date(0).toUTCString()}`;
        window.location = '/auth/login';
    }
</script>

<form class="modal" class:is-active={isOpen} on:submit|preventDefault={deleteAccount} style="z-index: 99999;">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="modal-background" on:click={close} />
    <div class="modal-card">
        <header class="modal-card-head">
            <p class="modal-card-title">{t('header')}</p>
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
