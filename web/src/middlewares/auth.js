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

export async function isAuthedMiddleware({ url, fetch, cookies }) {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });

    if (response.status !== 200) {
        cookies.delete('token', {
            path: '/',
        });
        throw redirect(
            307,
            `/auth/login?${new URLSearchParams({
                continue: url.href,
            })}`
        );
    }
    const { user } = await response.json();
    return user;
}

export async function isAdminMiddleware({ url, fetch, cookies }) {
    const authUser = await isAuthedMiddleware({ url, fetch, cookies });
    if (authUser.role !== 'admin') {
        throw redirect(307, '/');
    }
    return authUser;
}
