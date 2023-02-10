import { isAuthedMiddleware } from '../../middlewares/auth.js';

export async function load({ cookies, fetch }) {
    const authUser = await isAuthedMiddleware({ fetch, cookies });

    const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}/played_tracks`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data: tracks, pagination } = await response.json();

    return {
        token: cookies.get('token'),
        authUser,
        tracks,
        total: pagination.total,
    };
}
