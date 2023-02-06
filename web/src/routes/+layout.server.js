import { redirect } from "@sveltejs/kit";

export async function load({ cookies, fetch }) {
    if (cookies.get('token') != null) {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`
            }
        });
        const { success, user: authUser, last_track: lastTrack } = await response.json();
        if (!success) {
            throw redirect(307, '/auth/login');
        }
        return { token: cookies.get('token'), authUser, lastTrack };
    }
    return { authUser: null };
}
