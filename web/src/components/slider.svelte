<script>
    import { onMount, createEventDispatcher } from "svelte";

    export let maxValue;

    let thumb, slider, container;

    function seekToValue(value, tempMaxValue) {
        if (tempMaxValue != undefined) {
            maxValue = tempMaxValue;
        }
        let newThumbPosition = container.offsetWidth * (value / maxValue);
        thumb.style.left = newThumbPosition - thumb.offsetWidth / 2 + "px";
        slider.style.width = newThumbPosition + "px";
    }

    const dispatch = createEventDispatcher();
    function newValue(value) {
        dispatch("newValue", { value: value });
    }

    onMount(() => {
        let drag = false;
        function updateValue(event) {
            let newThumbPosition =
                event.pageX - container.getBoundingClientRect().left;
            if (
                newThumbPosition >= 0 &&
                newThumbPosition <= container.offsetWidth
            ) {
                thumb.style.left = `${
                    newThumbPosition - thumb.offsetWidth / 2
                }px`;
                slider.style.width = `${newThumbPosition}px`;
                newValue(
                    (parseInt(slider.style.width) / container.offsetWidth) *
                        maxValue
                );
            }
        }
        container.addEventListener("mousedown", (event) => {
            drag = true;
            updateValue(event);
        });
        thumb.addEventListener("mousedown", (event) => {
            drag = true;
            updateValue(event);
        });
        window.addEventListener("mousemove", (event) => {
            if (drag) updateValue(event);
        });
        window.addEventListener("mouseup", () => {
            drag = false;
        });
    });

    export { seekToValue };
</script>

<div class="slider-container" bind:this={container} style="flex: 1;">
    <div class="slider" bind:this={slider} />
    <div class="slider-thumb" bind:this={thumb} />
</div>
