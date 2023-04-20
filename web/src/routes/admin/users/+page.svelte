<script>
    import { tick, onMount } from 'svelte';
    import CreateModal from '../../../components/modals/users/create-modal.svelte';
    import EditModal from '../../../components/modals/users/edit-modal.svelte';
    import DeleteModal from '../../../components/modals/delete-modal.svelte';
    import { language } from '../../../stores.js';

    // Language strings
    const lang = {
        en: {
            title: 'Users - Admin - BassieMusic',
            create: 'Create new user',
            index: '#',
            username: 'Username',
            email: 'Email address',
            role: 'Role',
            role_normal: 'Normal',
            role_admin: 'Admin',
            actions: 'Actions',
            edit: 'Edit',
            delete: 'Delete',
            user: 'user',
        },
        nl: {
            title: 'Gebruikers - Admin - BassieMusic',
            create: 'Maake nieuwe gebruiker aan',
            index: '#',
            username: 'Gebruikersnaam',
            email: 'Email adres',
            role: 'Rol',
            role_normal: 'Normaal',
            role_admin: 'Admin',
            actions: 'Acties',
            edit: 'Verander',
            delete: 'Verwijder',
            user: 'gebruiker',
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;
    let { token, users } = data;
    let selectedUser = null;
    let createModal;
    let editModal;
    let deleteModal;

    // Methods
    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            }
        );
        const { data: newUsers, pagination } = await response.json();
        users = [...users, ...newUsers];
        if (users.length < pagination.total) {
            fetchPage(page + 1);
        }
    }
    if (users.length < data.total) {
        onMount(() => {
            fetchPage(2);
        });
    }
</script>

<svelte:head>
    <title>{t('title')}</title>
</svelte:head>

<div class="columns">
    <div class="column">
        <h2 class="title">Admin Users</h2>
    </div>
    <div class="column">
        <div class="buttons is-pulled-right">
            <button class="button is-link" on:click={createModal.open()}>{t('create')}</button>
        </div>
    </div>
</div>

<table class="table" style="width: 100%;">
    <thead>
        <th style="width: 10%;">{t('index')}</th>
        <th style="width: 25%;">{t('username')}</th>
        <th style="width: 30%;">{t('email')}</th>
        <th style="width: 15%;">{t('role')}</th>
        <th style="width: 20%;">{t('actions')}</th>
    </thead>
    <tbody>
        {#each users as user, index}
            <tr>
                <td>{index + 1}</td>
                <td>{user.username}</td>
                <td>{user.email}</td>
                <td>{user.role === 'normal' ? t('role_normal') : t('role_admin')}</td>
                <td>
                    <div class="buttons">
                        <button
                            class="button is-link is-small"
                            on:click={async () => {
                                selectedUser = user;
                                await tick();
                                editModal.open();
                            }}
                        >
                            {t('edit')}
                        </button>

                        <button
                            class="button is-danger is-small"
                            on:click={async () => {
                                selectedUser = user;
                                await tick();
                                deleteModal.open();
                            }}
                        >
                            {t('delete')}
                        </button>
                    </div>
                </td>
            </tr>
        {/each}
    </tbody>
</table>

<CreateModal
    bind:this={createModal}
    {token}
    on:create={(event) => {
        users = [...users, event.detail.user].sort((a, b) =>
            a.username.toLowerCase().localeCompare(b.username.toLowerCase())
        );
    }}
/>

{#if selectedUser !== null}
    <EditModal
        bind:this={editModal}
        {token}
        user={selectedUser}
        on:update={(event) => {
            users = users.map((user) => {
                if (user.id === event.detail.user.id) return event.detail.user;
                return user;
            });
        }}
    />

    <DeleteModal
        bind:this={deleteModal}
        {token}
        item={selectedUser}
        itemRoute="users"
        itemLabel={t('user')}
        on:delete={() => {
            users = users.filter((user) => user.id !== selectedUser.id);
        }}
    />
{/if}
