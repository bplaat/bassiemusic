<script>
    export let genre;
</script>

<a class="card" href="/genres/{genre.id}">
    <div
        class="card-image"
        style="background-image: url({genre.image});"
    />
    <div class="card-content">
        <h3 class="title is-6">{genre.name}</h3>
    </div>
</a>