import { redirect } from '@sveltejs/kit';
import { isAuthedMiddleware } from '../../../middlewares/auth.js';

export async function load({ fetch, cookies }) {
    const authUser = await isAuthedMiddleware({ fetch, cookies });
    if (authUser.role !== 'admin') {
        throw redirect(307, '/');
    }

    return { authUser };
}
