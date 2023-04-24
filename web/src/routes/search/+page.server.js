export async function load({ locals, url, fetch, cookies }) {
    // Do search when query is given
    const query = url.searchParams.get('q') || '';
    let searchResult = null;
    if (query !== '') {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/search?${new URLSearchParams({
                q: query,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${cookies.get('token')}`,
                },
            }
        );
        searchResult = await response.json();
    }

    // Always fetch genres first page
    const response = await fetch(`${import.meta.env.VITE_API_URL}/genres?page=1`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data: genres, pagination } = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        query,
        searchResult,
        genres,
        genresTotal: pagination.total,
        genresPage: 2,
    };
}
