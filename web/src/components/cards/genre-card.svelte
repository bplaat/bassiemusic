<script>
    import { language } from '../../stores.js';

    // Language strings
    const lang = {
        en: {
            image_alt: 'Image of genre $1',
        },
        nl: {
            image_alt: 'Afbeelding van genre $1',
        },
    };
    const t = (key, p1) => lang[$language][key].replace('$1', p1);

    // Props
    export let genre;
</script>

<a class="card" href="/genres/{genre.id}">
    <div class="card-image" style="aspect-ratio: 1;">
        <img src={genre.medium_image} alt={t('image_alt', genre.name)} loading="lazy" />
    </div>
    <div class="card-content">
        <h3 class="title is-6 ellipsis">{genre.name}</h3>
    </div>
</a>
