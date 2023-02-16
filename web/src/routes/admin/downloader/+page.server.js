import { redirect } from '@sveltejs/kit';
import { isAuthedMiddleware } from '../../../middlewares/auth.js';

export async function load({ fetch, cookies }) {
    const authUser = await isAuthedMiddleware({ fetch, cookies });
    if (authUser.role !== 'admin') {
        throw redirect(307, '/');
    }

    // Get storage size
    const response = await fetch(`${import.meta.env.VITE_API_URL}/storage_size`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const storage = await response.json();

    return { token: cookies.get('token'), authUser, storage };
}
