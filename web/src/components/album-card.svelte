<script>
export let album;
</script>

<a class="card" href="/albums/{album.id}">
    <div class="card-image" style="background-image: url({album.cover});">
        <div class="card-image-tags">
            {#if album.type == 'album'}
                <span class="tag">ALBUM</span>
            {/if}
            {#if album.type == 'ep'}
                <span class="tag">EP</span>
            {/if}
            {#if album.type == 'single'}
                <span class="tag">SINGLE</span>
            {/if}
            {#if album.explicit}
                <span class="tag is-danger">E</span>
            {/if}
        </div>
    </div>
    <div class="card-content">
        <h3 class="title is-6 mb-2">{album.title}</h3>
        <p>
            {#each album.artists as artist}
                <a href="/artists/{artist.id}" class="mr-2">{artist.name}</a>
            {/each}
        </p>
    </div>
</a>
