export async function load({ locals, fetch, cookies, request }) {
    // When a user is authed return values
    if (locals.authUser !== null) {
        if(locals.lastTrack == undefined) locals.lastTrack = null;
        return {
            token: cookies.get('token'),
            authUser: locals.authUser,
            agent: locals.agent,
            lastTrack: locals.lastTrack,
            lastTrackPosition: locals.lastTrackPosition,
            lastPlaylists: locals.lastPlaylists,
        };
    }

    // When guest get agent information
    const response = await fetch(`${import.meta.env.VITE_API_URL}/agent`, {
        headers: {
            'User-Agent': request.headers.get('user-agent'),
        },
    });
    const agent = await response.json();
    return { authUser: null, agent, lastTrack: null };
}
