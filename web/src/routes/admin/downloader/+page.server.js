export async function load({ locals, fetch, cookies }) {
    // Get storage size
    const response = await fetch(`${import.meta.env.VITE_API_URL}/storage_size`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const storage = await response.json();

    // Return value
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        storage,
    };
}
