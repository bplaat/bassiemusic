<script>
    import { createEventDispatcher } from 'svelte';

    export let token;
    export let user;

    let isOpen = false;
    export function open() {
        isOpen = true;
    }
    export function close() {
        isOpen = false;
    }

    const dispatch = createEventDispatcher();
    async function deleteUser() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${user.id}/delete`, {
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
            <p class="modal-card-title">Delete user</p>
            <button type="button" class="delete" aria-label="close" on:click={close} />
        </header>
        <section class="modal-card-body">
            <p>Do your realy want to delete this user?</p>
        </section>
        <footer class="modal-card-foot">
            <button type="submit" class="button is-danger">Delete user</button>
            <button class="button" on:click|preventDefault={close}>Cancel</button>
        </footer>
    </div>
</form>
