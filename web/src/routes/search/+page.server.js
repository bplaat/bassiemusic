import { isAuthedMiddleware } from '../../middlewares/auth.js';

export async function load({ url, cookies, fetch }) {
    const authUser = await isAuthedMiddleware({ fetch, cookies });

    const response = await fetch(`${import.meta.env.VITE_API_URL}/genres`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data: genres, pagination } = await response.json();

    return {
        token: cookies.get('token'),
        authUser,
        genres,
        query: url.searchParams.get('q') || undefined,
        total: pagination.total,
    };
}
