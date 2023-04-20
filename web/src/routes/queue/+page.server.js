export async function load({ locals, cookies }) {
    // Return token and auth user
    return { token: cookies.get('token'), authUser: locals.authUser };
}
