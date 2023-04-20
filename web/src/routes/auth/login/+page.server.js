export async function load({ url }) {
    // Return continue query variable
    return { continueUrl: url.searchParams.get('continue') };
}
