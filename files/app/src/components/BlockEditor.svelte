<script>
  export let passagetext
  export let transcriptionFontSize

  import { createEventDispatcher } from 'svelte'

  // function autoResize(textarea) {
  //   textarea.style.height = 'auto'
  //   textarea.style.height = textarea.scrollHeight + 'px'
  //   console.log('!!!!', textarea.scrollHeight)
  // }
  // on:input={autoResize(this)}

  const dispatch = createEventDispatcher()

  function updateValue(event) {
    const newpassagetext = event.target.value
    dispatch('update', newpassagetext)
  }

</script>

<style>
  .container {
    height: 100%;
    width: 100%;
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  textarea {
    height: 100%;
    width: 100%;
    resize: none;
    box-sizing: border-box;
  }

</style>

<div class="container">
  <textarea
    class="specialKey"
    rows="15"
    name="transcription"
    style="font-size: {transcriptionFontSize};"
    bind:value={passagetext}
    on:change={updateValue} />
</div>
<!-- 
<input class="specialKey" name="transcription" style="margin-left: 25px;  webkit-line-break:
normal; width: 98%; "> -->
