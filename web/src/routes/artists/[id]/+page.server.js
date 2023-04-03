import { error } from '@sveltejs/kit';
import { isAuthedMiddleware } from '../../../middlewares/auth.js';

export async function load({ url, fetch, cookies, params }) {
    const authUser = await isAuthedMiddleware({ url, fetch, cookies });

    const response = await fetch(`${import.meta.env.VITE_API_URL}/artists/${params.id}`, {
        headers: {
            Authorization: `Bearer ${cookies.get('token')}`,
        },
    });
    if (response.status == 404) {
        throw error(404, 'Not Found');
    }
    const artist = await response.json();

    const filterAlbumsBy = url.searchParams.get('albums_filter') || 'all';
    const filteredAlbums = artist.albums.filter((album) => {
        if (filterAlbumsBy == 'all') return true;
        if (filterAlbumsBy == 'album') return album.type == 'album';
        if (filterAlbumsBy == 'ep') return album.type == 'ep';
        if (filterAlbumsBy == 'single') return album.type == 'single';
    });

    return { token: cookies.get('token'), authUser, artist, filterAlbumsBy, filteredAlbums };
}
