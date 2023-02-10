<script>
    import { createEventDispatcher } from "svelte";

    export let maxValue;
    export let style;

    let thumb, slider, container;

    export function seekToValue(value, tempMaxValue) {
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

    let drag = false;

    function updateValue(event) {
        let newThumbPosition =
            event.pageX - container.getBoundingClientRect().left;
        if (
            newThumbPosition >= 0 &&
            newThumbPosition <= container.offsetWidth
        ) {
            thumb.style.left = `${newThumbPosition - thumb.offsetWidth / 2}px`;
            slider.style.width = `${newThumbPosition}px`;
            newValue(
                (parseInt(slider.style.width) / container.offsetWidth) *
                    maxValue
            );
        }
    }

    function onMousedown(event) {
        drag = true;
        updateValue(event);
    }

    function onMousemove(event) {
        if (drag) updateValue(event);
    }

    function onMouseup() {
        drag = false;
    }
</script>

<svelte:window on:mousemove|preventDefault={onMousemove} on:mouseup|preventDefault={onMouseup} />

<div class="slider-container" style="{style};" bind:this={container} on:mousedown|preventDefault={onMousedown}>
    <div class="slider" bind:this={slider} />
    <div class="slider-thumb" bind:this={thumb} on:mousedown|preventDefault={onMousedown} />
</div>
