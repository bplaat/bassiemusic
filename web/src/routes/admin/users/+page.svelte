<script>
    import { onMount } from 'svelte';
    import CreateModal from '../../../components/admin/users/create-modal.svelte';
    import EditModal from '../../../components/admin/users/edit-modal.svelte';
    import DeleteModal from '../../../components/admin/users/delete-modal.svelte';
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
        },
    };
    const t = (key) => lang[$language][key];

    // State
    export let data;
    let { token, users } = data;
    let selectedUser = {};
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

    function updateUsers() {
        users = [];
        fetchPage(1);
    }

    function editUser(user) {
        selectedUser = user;
        editModal.open();
    }

    function deleteUser(user) {
        selectedUser = user;
        deleteModal.open();
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
                <td>{user.role == 'normal' ? t('role_normal') : t('role_admin')}</td>
                <td>
                    <div class="buttons">
                        <button class="button is-link is-small" on:click={() => editUser(user)}>{t('edit')}</button>
                        <button class="button is-danger is-small" on:click={() => deleteUser(user)}>
                            {t('delete')}
                        </button>
                    </div>
                </td>
            </tr>
        {/each}
    </tbody>
</table>

<CreateModal bind:this={createModal} {token} on:createUser={updateUsers} />

<EditModal bind:this={editModal} {token} user={selectedUser} on:updateUser={updateUsers} />

<DeleteModal bind:this={deleteModal} {token} user={selectedUser} on:deleteUser={updateUsers} />
