export async function load({ locals, url, fetch, cookies }) {
    // Fetch artists first page
    const sortBy = url.searchParams.get('sort_by') || 'name';
    const response = await fetch(
        `${import.meta.env.VITE_API_URL}/artists?${new URLSearchParams({
            sort_by: sortBy,
            page: 1,
        })}`,
        {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`,
            },
        }
    );
    const { data: artists, pagination } = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        artists,
        sortBy,
        total: pagination.total,
    };
}
