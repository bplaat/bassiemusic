import { authMiddleware } from '../../middlewares/auth.js';

export async function load({ cookies, fetch }) {
    const authUser = authMiddleware({ fetch, cookies });

    const response = await fetch(`${import.meta.env.VITE_API_URL}/artists`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`
        }
    });
    const { data: artists, pagination } = await response.json();

    return { authUser, artists, total: pagination.total };
}
