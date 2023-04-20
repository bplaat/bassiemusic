export async function load({ locals, url, fetch, cookies }) {
    // Fetch tracks first page
    const sortBy = url.searchParams.get('sort_by') || 'plays_desc';
    const response = await fetch(
        `${import.meta.env.VITE_API_URL}/tracks?${new URLSearchParams({
            sort_by: sortBy,
            page: 1,
        })}`,
        {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`,
            },
        }
    );
    const { data: tracks, pagination } = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        tracks,
        sortBy,
        total: pagination.total,
    };
}
