import { isAuthedMiddleware } from '../../middlewares/auth.js';

export async function load({ cookies, fetch }) {
    const authUser = await isAuthedMiddleware({ fetch, cookies });

    const response = await fetch(`${import.meta.env.VITE_API_URL}/artists`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`
        }
    });
    const { data: artists, pagination } = await response.json();

    return { token: cookies.get('token'), authUser, artists, total: pagination.total };
}
