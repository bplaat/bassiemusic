import { isAuthedMiddleware } from '../../../middlewares/auth.js';

export async function load({ url, fetch, cookies }) {
    const authUser = await isAuthedMiddleware({ url, fetch, cookies });

    const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}/liked_artists`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data: artists, pagination } = await response.json();

    return {
        token: cookies.get('token'),
        authUser,
        artists,
        total: pagination.total,
    };
}
