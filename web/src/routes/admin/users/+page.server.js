export async function load({ locals, fetch, cookies }) {
    // Fetch users first page
    const response = await fetch(`${import.meta.env.VITE_API_URL}/users?page=1`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data: users, pagination } = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        users,
        total: pagination.total,
    };
}
