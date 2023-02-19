import { redirect } from '@sveltejs/kit';

export async function load({ url, fetch, cookies }) {
    // Validate token
    const validateResponse = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    if (validateResponse.status !== 200) {
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
    const { user: authUser, session_id: currentSessionId } = await validateResponse.json();

    // Fetch sessions
    const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}/sessions`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data: sessions, pagination } = await response.json();

    return {
        token: cookies.get('token'),
        authUser,
        currentSessionId,
        sessions,
        sessionsTotal: pagination.total,
    };
}
