import { redirect } from '@sveltejs/kit';

// This server hook is really the auth middleware of the web player
export async function handle({ event, resolve }) {
    // Locals authUser is used to check if user is authed or not
    event.locals.authUser = null;

    // Guest pages
    if (event.url.pathname === '/auth/login') {
        if (event.cookies.get('token') !== null) {
            const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
                headers: {
                    Authorization: `Bearer ${event.cookies.get('token')}`,
                },
            });
            if (response.status === 200) {
                throw redirect(307, '/');
            }
        }
        return await resolve(event);
    }

    // Authed pages
    if (event.cookies.get('token') !== null) {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
            headers: {
                Authorization: `Bearer ${event.cookies.get('token')}`,
            },
        });
        if (response.status !== 200) {
            event.cookies.delete('token', { path: '/' });
            throw redirect(307, `/auth/login?${new URLSearchParams({ continue: event.url.href })}`);
        }
        const { user, session_id, agent, last_track, last_track_position, last_playlists } = await response.json();
        event.locals.authUser = user;
        event.locals.authSessionId = session_id;
        event.locals.agent = agent;
        event.locals.lastTrack = last_track;
        event.locals.lastTrackPosition = last_track_position;
        event.locals.lastPlaylists = last_playlists;
    } else {
        throw redirect(307, `/auth/login?${new URLSearchParams({ continue: event.url.href })}`);
    }

    // Admin pages
    if (event.url.pathname.startsWith('/admin') && event.locals.authUser.role !== 'admin') {
        throw redirect(307, '/');
    }
    return await resolve(event);
}
