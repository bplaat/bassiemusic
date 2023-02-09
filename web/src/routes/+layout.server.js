import { redirect } from "@sveltejs/kit";

export async function load({ request, cookies, fetch }) {
    // When a token exist validate token
    if (cookies.get('token') != null) {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`
            }
        });
        if (response.status != 200) {
            cookies.delete('token', {
                path: '/'
            });
            throw redirect(307, '/auth/login');
        }
        const { user: authUser, agent, last_track: lastTrack, last_track_position: lastTrackPosition } = await response.json();
        return { token: cookies.get('token'), authUser, agent, lastTrack, lastTrackPosition };
    }

    // Just get user agent information
    const response = await fetch(`${import.meta.env.VITE_API_URL}/agent`, {
        headers: {
            'User-Agent': request.headers.get('user-agent')
        }
    });
    const agent = await response.json();
    return { authUser: null, agent };
}
