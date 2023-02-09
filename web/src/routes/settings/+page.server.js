import { isAuthedMiddleware } from '../../middlewares/auth.js';

export async function load({ cookies, fetch }) {
    // Validate token
    const validateResponse = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`
        }
    });
    if (validateResponse.status !== 200) {
        throw redirect(307, '/auth/login');
    }
    const { user: authUser, session: currentSession } = await validateResponse.json();

    // Fetch sessions
    const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${authUser.id}/sessions`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`
        }
    });
    const { data: sessions, pagination } = await response.json();

    return { token: cookies.get('token'), authUser, currentSession, sessions, sessionsToal: pagination.total };
}
