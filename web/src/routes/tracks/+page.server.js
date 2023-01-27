import { authMiddleware } from '../../middlewares/auth.js';

export async function load({ cookies, fetch }) {
    const authUser = authMiddleware({ fetch, cookies });

    const response = await fetch(`${import.meta.env.VITE_API_URL}/tracks?limit=50`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`
        }
    });
    const { data: tracks, pagination } = await response.json();

    return { authUser, tracks, total: pagination.total };
}
