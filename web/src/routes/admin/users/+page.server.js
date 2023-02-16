import { isAdminMiddleware } from '../../../middlewares/auth.js';

export async function load({ fetch, cookies }) {
    const authUser = await isAdminMiddleware({ fetch, cookies });

    const response = await fetch(`${import.meta.env.VITE_API_URL}/users`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data: users, pagination } = await response.json();

    return {
        token: cookies.get('token'),
        authUser,
        users,
        total: pagination.total,
    };
}
