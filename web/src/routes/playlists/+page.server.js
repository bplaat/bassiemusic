export async function load({ locals, url, fetch, cookies }) {
    // First playlists first page
    const sortBy = url.searchParams.get('sort_by') || 'updated_at_desc';
    const response = await fetch(
        `${import.meta.env.VITE_API_URL}/playlists?${new URLSearchParams({
            sort_by: sortBy,
            page: 1,
        })}`,
        {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`,
            },
        }
    );
    const { data: playlists, pagination } = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        playlists,
        sortBy,
        total: pagination.total,
    };
}
