import { isGuestMiddleware } from '../../middlewares/auth.js';

export async function load({ fetch, cookies }) {
    await isGuestMiddleware({ fetch, cookies });
}
