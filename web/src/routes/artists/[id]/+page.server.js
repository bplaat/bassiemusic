import { error } from '@sveltejs/kit';
import { authMiddleware } from '../../../middlewares/auth.js';

export async function load({ fetch, cookies, params }) {
    const authUser = authMiddleware({ fetch, cookies });

    const response = await fetch(`${import.meta.env.VITE_API_URL}/artists/${params.id}`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`
        }
    });
    if (response.status === 404) {
        throw error(404, 'Not Found');
    }
    const artist = await response.json();

    return { authUser, artist };
}
