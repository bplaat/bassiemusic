import { isAuthedMiddleware } from '../middlewares/auth.js';

export async function load({ cookies, fetch }) {
    if (cookies.get('token') != null) {
        const authUser = await isAuthedMiddleware({ fetch, cookies });

        const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}/played_tracks`, {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`,
            },
        });
        const { data: lastPlayedTracks } = await response.json();

        return {
            token: cookies.get('token'),
            authUser,
            lastPlayedTracks,
        };
    }
}
