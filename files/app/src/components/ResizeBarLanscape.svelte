<script>
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()
  let barElement = undefined,
    innerOffsetY = null
  $: resizing = innerOffsetY !== null

  function handleMouseDown(e) {
    innerOffsetY = e.y
    document.addEventListener('mousemove', handleMouseMove, false)
    dispatch('down', {})
  }

  function handleMouseMove(e) {
    dispatch('move', {
      y: e.y - innerOffsetY,
    })
  }

  function handleMouseUp() {
    innerOffsetY = null
    document.removeEventListener('mousemove', handleMouseMove, false)
    dispatch('up')
  }

</script>

<style>
  .bar {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

    box-sizing: border-box;
    width: 100%;
    height: 12px;
    padding: 4px;

    background-color: var(--toolbar-bg-color);
    border-top: 1px solid var(--toolbar-border-color);
    cursor: ns-resize;
    transition: background-color 120ms ease-in-out;
  }

  .bar.active {
    background-color: rgb(255, 255, 255);
  }

  .handle {
    flex-shrink: 0;
    flex-grow: 0;
    width: 32px;
    height: 4px;
    border-radius: 2px;
    background-color: rgb(100, 100, 100);
    transition: background-color 120ms ease-in-out;
  }

  .bar:hover .handle {
    background-color: rgb(0, 0, 0);
  }

  .resizing {
    cursor: ns-resize;
  }

</style>

<svelte:window on:mouseup={handleMouseUp} />
<svelte:body class:resizing />

<div
  class="bar resizing"
  bind:this={barElement}
  on:mousedown={handleMouseDown}
  class:active={resizing}>
  <div class="handle" />
</div>
