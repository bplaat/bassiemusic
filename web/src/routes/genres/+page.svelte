<script>
    import Cookies from "js-cookie";
    import GenreCard from "../../components/genre-card.svelte";

    export let data;
    const { genres } = data;

    async function fetchPage(page) {
        const response = await fetch(
            `${import.meta.env.VITE_API_URL}/genres?${new URLSearchParams({
                page,
            })}`,
            {
                headers: {
                    Authorization: `Bearer ${Cookies.get("token")}`,
                },
            }
        );
        const { data: newGenres, pagination } = await response.json();
        genres.push(...newGenres);
        if (genres.length != pagination.total) {
            fetchPage(page + 1);
        }
    }
    if (genres.length != data.total) {
        fetchPage(2);
    }
</script>

<svelte:head>
    <title>Genres - BassieMusic</title>
</svelte:head>

<h2 class="title">Genres</h2>

<div class="columns is-multiline">
    {#each genres as genre}
        <div class="column is-one-fifth">
            <GenreCard {genre}/>
        </div>
    {/each}
</div>
