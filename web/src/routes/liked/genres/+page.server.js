export async function load({ locals, url, fetch, cookies }) {
    // Fetch liked genres first page
    const sortBy = url.searchParams.get('sort_by') || 'liked_at_desc';
    const response = await fetch(
        `${import.meta.env.VITE_API_URL}/users/${locals.authUser.id}/liked_genres?${new URLSearchParams({
            sort_by: sortBy,
            page: 1,
        })}`,
        {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`,
            },
        }
    );
    const { data: genres, pagination } = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        genres,
        sortBy,
        total: pagination.total,
    };
}
