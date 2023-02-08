<script>
    import { onMount, createEventDispatcher  } from 'svelte';

    export let maxValue;

    let thumb,
        slider,
        container;

    function seekToValue(value) {
        let newThumbPosition = container.offsetWidth * (value / maxValue);
        thumb.style.left = newThumbPosition - thumb.offsetWidth / 2 + "px";
        slider.style.width = newThumbPosition + "px";
    }

    const dispatch = createEventDispatcher();
    function newValue(value){
        dispatch("newValue", { value: value });
    }

    onMount(() => {
        thumb = document.querySelector(".slider-thumb");
        slider = document.querySelector(".slider");
        container = document.querySelector(".slider-container");
        container.addEventListener("click", function(event) {
            let newThumbPosition = event.pageX - container.getBoundingClientRect().left;
            if (newThumbPosition >= 0 && newThumbPosition <= container.offsetWidth) {
                thumb.style.left = newThumbPosition - thumb.offsetWidth / 2 + "px";
                slider.style.width = newThumbPosition + "px";

                newValue(((parseInt(slider.style.width) / container.offsetWidth)) * maxValue)
            }
        });

        thumb.addEventListener("mousedown", function(event) {
            let thumbPosition = event.pageX - container.getBoundingClientRect().left;
            thumb.classList.add("active");
            let mouseMoveHandler = function(event) {
                let newThumbPosition = event.pageX - container.getBoundingClientRect().left;

                if (newThumbPosition >= 0 && newThumbPosition <= container.offsetWidth) {
                    thumbPosition = newThumbPosition;
                    thumb.style.left = thumbPosition - thumb.offsetWidth / 2 + "px";
                    slider.style.width = thumbPosition + "px";
                    
                    newValue(((parseInt(slider.style.width) / container.offsetWidth)) * maxValue)
                }
            };

            let mouseUpHandler = function() {
                thumb.classList.remove("active");
                document.removeEventListener("mousemove", mouseMoveHandler);
                document.removeEventListener("mouseup", mouseUpHandler);
            };
            document.addEventListener("mousemove", mouseMoveHandler);
            document.addEventListener("mouseup", mouseUpHandler);
        });
    });

    export {seekToValue}
</script>

<div class="slider-container" style="flex:1;">
    <div class="slider"></div>
    <div class="slider-thumb"></div>
</div>