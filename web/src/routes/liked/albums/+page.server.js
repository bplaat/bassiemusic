import { isAuthedMiddleware } from '../../../middlewares/auth.js';

export async function load({ cookies, fetch }) {
    const authUser = await isAuthedMiddleware({ fetch, cookies });

    const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}/liked_albums`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`
        }
    });
    const { data: albums, pagination } = await response.json();

    return { token: cookies.get('token'), authUser, albums, total: pagination.total };
}
