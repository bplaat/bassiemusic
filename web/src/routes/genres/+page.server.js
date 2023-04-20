export async function load({ locals, url, fetch, cookies }) {
    // Fetch genres first page
    const sortBy = url.searchParams.get('sort_by') || 'name';
    const response = await fetch(
        `${import.meta.env.VITE_API_URL}/genres?${new URLSearchParams({
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
