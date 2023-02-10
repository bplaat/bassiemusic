import { redirect } from '@sveltejs/kit';

export async function isGuestMiddleware({ fetch, cookies }) {
    if (cookies.get('token') != null) {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`,
            },
        });
        if (response.status == 200) {
            throw redirect(307, '/');
        }
    }
}

export async function isAuthedMiddleware({ fetch, cookies }) {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    if (response.status !== 200) {
        cookies.delete('token', {
            path: '/',
        });
        throw redirect(307, '/auth/login');
    }
    const { user } = await response.json();
    return user;
}
