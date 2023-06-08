export function load({ url }) {
    return { sortBy: url.searchParams.get('sort_by') || 'name' };
}
