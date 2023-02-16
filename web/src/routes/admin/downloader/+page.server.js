import { isAdminMiddleware } from '../../../middlewares/auth.js';

export async function load({ fetch, cookies }) {
    const authUser = await isAdminMiddleware({ fetch, cookies });

    // Get storage size
    const response = await fetch(`${import.meta.env.VITE_API_URL}/storage_size`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const storage = await response.json();

    return { token: cookies.get('token'), authUser, storage };
}
