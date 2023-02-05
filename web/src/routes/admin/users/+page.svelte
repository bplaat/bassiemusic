<script>
    import Cookies from "js-cookie";
    import CreateModal from "../../../components/admin/users/create-modal.svelte";
    import EditModal from "../../../components/admin/users/edit-modal.svelte";
    import DeleteModal from "../../../components/admin/users/delete-modal.svelte";

    export let data;
    let { users } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/users?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${Cookies.get("token")}`,
                },
            }
        );
        const { data: newUsers, pagination } = await response.json();
        users.push(...newUsers);
        users = users;
        if (users.length != pagination.total) {
            fetchPage(page + 1);
        }
    }
    if (users.length != data.total) {
        fetchPage(2);
    }

    function updateUsers() {
        users = [];
        fetchPage(1);
    }

    let selectedUser = {};
    let createModal;
    let editModal;
    let deleteModal;

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
    <title>Users - Admin - BassieMusic</title>
</svelte:head>

<div class="columns">
    <div class="column">
        <h2 class="title">Admin User</h2>
    </div>
    <div class="column">
        <div class="buttons is-pulled-right">
            <button class="button is-link" on:click={createModal.open()}>Create new user</button>
        </div>
    </div>
</div>

<table class="table" style="width: 100%;">
    <thead>
        <th style="width: 10%;">#</th>
        <th style="width: 30%;">Username</th>
        <th style="width: 30%;">Email</th>
        <th style="width: 15%;">Role</th>
        <th style="width: 15%;">Actions</th>
    </thead>
    <tbody>
        {#each users as user, index}
            <tr>
                <td>{index + 1}</td>
                <td>{user.username}</td>
                <td>{user.email}</td>
                <td>{user.role == "normal" ? "Normal" : "Admin"}</td>
                <td>
                    <div class="buttons">
                        <button
                            class="button is-link is-small"
                            on:click={() => editUser(user)}>Edit</button
                        >
                        <button
                            class="button is-danger is-small"
                            on:click={() => deleteUser(user)}>Delete</button
                        >
                    </div>
                </td>
            </tr>
        {/each}
    </tbody>
</table>

<CreateModal
    bind:this={createModal}
    on:createUser={updateUsers}
/>

<EditModal
    bind:this={editModal}
    user={selectedUser}
    on:updateUser={updateUsers}
/>

<DeleteModal
    bind:this={deleteModal}
    user={selectedUser}
    on:deleteUser={updateUsers}
/>
