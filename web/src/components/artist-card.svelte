<script>
    export let artist;
</script>

<a class="card" href="/artists/{artist.id}">
    <div class="card-image">
        <img
            style="aspect-ratio: 1;"
            src={artist.medium_image}
            alt="Image of artist {artist.name}"
            loading="lazy"
        />
    </div>
    <div class="card-content">
        <h3 class="title is-6 ellipsis">{artist.name}</h3>
    </div>
</a>
