import { isAuthedMiddleware } from '../../../middlewares/auth.js';

export async function load({ url, fetch, cookies }) {
    const authUser = await isAuthedMiddleware({ url, fetch, cookies });

    const sortBy = url.searchParams.get('sort_by') || 'liked_at_desc';
    const response = await fetch(
        `${import.meta.env.VITE_API_URL}/users/${authUser.id}/liked_genres?${new URLSearchParams({
            sort_by: sortBy,
        })}`,
        {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`,
            },
        }
    );
    const { data: genres, pagination } = await response.json();

    return {
        token: cookies.get('token'),
        authUser,
        genres,
        sortBy,
        total: pagination.total,
    };
}
