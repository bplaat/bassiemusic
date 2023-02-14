<script>
    export let genre;
</script>

<a class="card" href="/genres/{genre.id}">
    <div class="card-image" style="aspect-ratio: 1;">
        <img src={genre.medium_image} alt="Image of genre {genre.name}" loading="lazy" />
    </div>
    <div class="card-content">
        <h3 class="title is-6 ellipsis">{genre.name}</h3>
    </div>
</a>
