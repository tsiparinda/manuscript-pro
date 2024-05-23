<script>
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()
  let barElement = undefined,
    innerOffsetX = null
  $: resizing = innerOffsetX !== null

  function handleMouseDown(e) {
    innerOffsetX = e.x
    document.addEventListener('mousemove', handleMouseMove, false)
    dispatch('down', {})
  }

  function handleMouseMove(e) {
    // console.log('!!', e.x, innerOffsetX)
    dispatch('move', {
      x: e.x - innerOffsetX,
    })
  }

  function handleMouseUp() {
    innerOffsetX = null
    document.removeEventListener('mousemove', handleMouseMove, false)
    dispatch('up')
  }

</script>

<style>
  .bar {
    display: flex;
    align-items: center;
    justify-content: center;

    box-sizing: border-box;
    width: 12px;
    height: auto;
    padding: 4px;

    background-color: var(--toolbar-bg-color);
    border-left: 1px solid var(--toolbar-border-color);
    cursor: ew-resize;
    transition: background-color 120ms ease-in-out;
  }

  .bar.active {
    background-color: rgb(255, 255, 255);
  }

  .handle {
    flex-shrink: 0;
    flex-grow: 0;
    width: 4px;
    height: 100%;
    border-radius: 2px;
    background-color: rgb(100, 100, 100);
    transition: background-color 120ms ease-in-out;
  }

  .bar:hover .handle {
    background-color: rgb(0, 0, 0);
  }

  .resizing {
    cursor: ew-resize;
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
