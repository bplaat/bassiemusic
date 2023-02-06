<script>
    export let artist;
</script>

<a class="card" href="/artists/{artist.id}">
    <div class="card-image" style="background-image: url({artist.medium_image});" />
    <div class="card-content">
        <h3 class="title is-6">{artist.name}</h3>
    </div>
</a>
