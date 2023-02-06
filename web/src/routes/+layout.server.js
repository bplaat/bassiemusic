import { redirect } from "@sveltejs/kit";

export async function load({ cookies, fetch }) {
    if (cookies.get('token') != null) {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`
            }
        });
        const { user: authUser, last_track: lastTrack, last_track_position: lastTrackPosition } = await response.json();
        if (response.status != 200) {
            throw redirect(307, '/auth/login');
        }
        return { token: cookies.get('token'), authUser, lastTrack, lastTrackPosition };
    }
    return { authUser: null };
}
