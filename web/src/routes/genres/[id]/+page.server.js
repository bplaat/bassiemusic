import { error } from '@sveltejs/kit';

export async function load({ locals, fetch, cookies, params }) {
    // Fetch genre
    const response = await fetch(`${import.meta.env.VITE_API_URL}/genres/${params.id}`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    if (response.status === 404) {
        throw error(404, 'Not Found');
    }
    const genre = await response.json();

    // Fetch genre albums first page
    const albumsResponse = await fetch(`${import.meta.env.VITE_API_URL}/genres/${params.id}/albums`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    const { data, pagination } = await albumsResponse.json();
    genre.albums = data;

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        genre,
        albumsTotal: pagination.total,
        albumsPage: 2,
    };
}
