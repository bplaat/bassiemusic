export async function load({ locals, fetch, cookies }) {
    // Fetch user active sessions first page
    const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${locals.authUser.id}/active_sessions`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data: sessions, pagination } = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        authSessionId: locals.authSessionId,
        sessions,
        sessionsTotal: pagination.total,
    };
}
