import { redirect } from '@sveltejs/kit';

export async function load({ cookies, fetch }) {
    await isNotAuthedMiddleware({ cookies, fetch });
}
