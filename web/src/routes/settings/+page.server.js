import { isAuthedMiddleware } from '../../middlewares/auth.js';

export async function load({ cookies, fetch }) {
    const authUser = await isAuthedMiddleware({ fetch, cookies });
    return { token: cookies.get('token'), authUser };
}
