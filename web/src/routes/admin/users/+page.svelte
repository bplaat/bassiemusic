<script>
    import Cookies from "js-cookie";

    export let data;
    const { users } = data;

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
        if (users.length != pagination.total) {
            fetchPage(page + 1);
        }
    }
    if (users.length != data.total) {
        fetchPage(2);
    }
</script>

<h1 class="title">Admin Users</h1>

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
                <td>{user.role == 'normal' ? 'Normal' : 'Admin'}</td>
                <td>
                    <div class="buttons">
                        <button class="button is-link is-small">Edit</button>
                        <button class="button is-danger is-small">Delete</button>
                    </div>
                </td>
            </tr>
        {/each}
    </tbody>
</table>
