import { redirect } from "@sveltejs/kit";

export async function load({ cookies, fetch }) {
    if (cookies.get('token') != null) {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/validate`, {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`
            }
        });
        const { success } = await response.json();
        if (success) {
            throw redirect(307, '/');
        }
    }
}
