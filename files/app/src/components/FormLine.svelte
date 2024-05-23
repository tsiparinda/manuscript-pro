<script>
  export let id = ''
  export let label = ''
  export let offset = false
  export let info = ''

</script>

<style>
  .offset {
    margin-top: 25px;
  }

  .label-wrapper {
    display: flex;
    align-items: flex-start; /* Aligns label to the top */
  }

  .line-label {
    margin-right: 10px;
  }

  [data-title]:hover:after {
    opacity: 1;
    transition: all 0.1s ease 0.5s;
    visibility: visible;
  }
  [data-title]:after {
    content: attr(data-title);
    position: absolute;
    bottom: -1.6em;
    left: 100%;
    padding: 4px 4px 4px 8px;
    color: #222;
    white-space: nowrap;
    -moz-border-radius: 5px;
    -webkit-border-radius: 5px;
    border-radius: 5px;
    -moz-box-shadow: 0px 0px 4px #222;
    -webkit-box-shadow: 0px 0px 4px #222;
    box-shadow: 0px 0px 4px #222;
    background-image: -moz-linear-gradient(top, #f8f8f8, #cccccc);
    background-image: -webkit-gradient(
      linear,
      left top,
      left bottom,
      color-stop(0, #f8f8f8),
      color-stop(1, #cccccc)
    );
    background-image: -webkit-linear-gradient(top, #f8f8f8, #cccccc);
    background-image: -moz-linear-gradient(top, #f8f8f8, #cccccc);
    background-image: -ms-linear-gradient(top, #f8f8f8, #cccccc);
    background-image: -o-linear-gradient(top, #f8f8f8, #cccccc);
    opacity: 0;
    z-index: 99999;
    visibility: hidden;
  }
  [data-title] {
    position: relative;
  }

</style>

<div class="field is-horizontal" class:offset>
  <div class="field-label is-normal label-wrapper">
    {#if label}
      <label class="label offset line-label" for={id}>{label}
        {#if info}<i class="fa fa-question-circle-o" data-title={info} />{/if}
      </label>
    {:else}
      <div class="line-label" />
    {/if}
  </div>
  <div class="field-body">
    <div class="field">
      <div class="control">
        <slot />
      </div>
    </div>
  </div>
</div>
