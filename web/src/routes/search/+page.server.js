import { isAuthedMiddleware } from '../../middlewares/auth.js';

export async function load({ url, fetch, cookies }) {
    const authUser = await isAuthedMiddleware({ url, fetch, cookies });

    // Do search when query is given
    const query = url.searchParams.get('q') || '';
    let searchResult = null;
    if (query != '') {
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

    // Always load first genres page
    const response = await fetch(`${import.meta.env.VITE_API_URL}/genres`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data: genres, pagination } = await response.json();

    return {
        token: cookies.get('token'),
        authUser,
        query,
        searchResult,
        genres,
        genresTotal: pagination.total,
        genresPage: 2,
    };
}
