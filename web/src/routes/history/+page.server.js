export async function load({ locals, fetch, cookies }) {
    // Fetch played tracks first page
    const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${locals.authUser.id}/played_tracks?page=1`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data: tracks, pagination } = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        tracks,
        total: pagination.total,
    };
}
