import { redirect } from '@sveltejs/kit';

export async function authMiddleware({ fetch, cookies }) {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`
        }
    });
    if (response.status !== 200) {
        throw redirect(307, '/auth/login');
    }
    const { user } = await response.json();
    return user;
}
