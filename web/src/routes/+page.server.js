export async function load({ locals, fetch, cookies }) {
    // Fetch user last played tracks
    const response = await fetch(`${import.meta.env.VITE_API_URL}/users/${locals.authUser.id}/played_tracks`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data: lastPlayedTracks } = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        lastPlayedTracks,
    };
}
