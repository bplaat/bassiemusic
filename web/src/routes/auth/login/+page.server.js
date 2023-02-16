import { isGuestMiddleware } from '../../../middlewares/auth.js';

export async function load({ url, fetch, cookies }) {
    await isGuestMiddleware({ fetch, cookies });

    return {
        continueUrl: url.searchParams.get('continue') || undefined,
    };
}
