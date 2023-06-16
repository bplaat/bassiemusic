import { error } from '@sveltejs/kit';

export async function load({ locals, url, fetch, cookies, params }) {
    // Fetch artist
    let response = await fetch(`${import.meta.env.VITE_API_URL}/artists/${params.id}`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    if (response.status === 404) {
        throw error(404, 'Not Found');
    }
    const artist = await response.json();

    // Fetch tracks first page
    const sortBy = url.searchParams.get('sort_by') || 'plays_desc';
    response = await fetch(
        `${import.meta.env.VITE_API_URL}/artists/${artist.id}/tracks?${new URLSearchParams({
            sort_by: sortBy,
            page: 1,
        })}`,
        {
            headers: {
                Authorization: `Bearer ${cookies.get('token')}`,
            },
        }
    );

    const { data: tracks, pagination } = await response.json();

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        artist,
        tracks,
        sortBy,
        total: pagination.total,
    };
}
