import { error } from '@sveltejs/kit';

export async function load({ locals, url, fetch, cookies, params }) {
    // Fetch artist
    const response = await fetch(`${import.meta.env.VITE_API_URL}/artists/${params.id}`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    if (response.status === 404) {
        throw error(404, 'Not Found');
    }
    const artist = await response.json();

    // Filter artist albums
    const filterAlbumsBy = url.searchParams.get('albums_filter') || 'all';
    const filteredAlbums = artist.albums.filter((album) => {
        if (filterAlbumsBy === 'all') return true;
        if (filterAlbumsBy === 'album') return album.type === 'album';
        if (filterAlbumsBy === 'ep') return album.type === 'ep';
        if (filterAlbumsBy === 'single') return album.type === 'single';
    });

    // Return values
    return {
        token: cookies.get('token'),
        authUser: locals.authUser,
        artist,
        filterAlbumsBy,
        filteredAlbums,
    };
}
