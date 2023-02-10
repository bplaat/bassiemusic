import { isNotAuthedMiddleware } from '../../middlewares/auth.js';

export async function load({ cookies, fetch }) {
    await isNotAuthedMiddleware({ cookies, fetch });
}
