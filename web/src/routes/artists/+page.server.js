import { isAuthedMiddleware } from '../../middlewares/auth.js';

export async function load({ url, fetch, cookies }) {
    const authUser = await isAuthedMiddleware({ url, fetch, cookies });

    const sortBy = url.searchParams.get('sort_by') || 'name';
    const response = await fetch(
        `${import.meta.env.VITE_API_URL}/artists?${new URLSearchParams({
            sort_by: sortBy,
        })}`,
        {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`,
            },
        }
    );
    const { data: artists, pagination } = await response.json();

    return {
        token: cookies.get('token'),
        authUser,
        artists,
        sortBy,
        total: pagination.total,
    };
}
