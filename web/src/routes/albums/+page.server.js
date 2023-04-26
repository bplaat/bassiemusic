export async function load({ locals, url, fetch, cookies }) {
    // Fetch albums first page
    const sortBy = url.searchParams.get('sort_by') || 'title';
    const response = await fetch(
        `${import.meta.env.VITE_API_URL}/albums?${new URLSearchParams({
            sort_by: sortBy,
            page: 1,
        })}`,
        {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`,
            },
        }
    );
    const { data: albums, pagination } = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        albums,
        sortBy,
        total: pagination.total,
    };
}
