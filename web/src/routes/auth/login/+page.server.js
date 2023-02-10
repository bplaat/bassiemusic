import { isGuestMiddleware } from '../../../middlewares/auth.js';

export async function load({ cookies, fetch }) {
    await isGuestMiddleware({ cookies, fetch });
}
