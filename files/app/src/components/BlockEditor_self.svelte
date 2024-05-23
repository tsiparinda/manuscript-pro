<script>
  import { onMount } from 'svelte'
  import { isElementInViewport } from '../lib/utils.js' // Import the isElementInViewport function from a utils.js file
  import Modal from './Modal.svelte'

  export let id_col, bucket, citation
  let err,
    passagetext = 'Loading transcription...',
    passage,
    textarea,
    dataloaded = false,
    isVisible = false

  let isLoading = true

  let showModal = false,
    modalHeader,
    modalText

  $: textareaRows = countRows(passagetext)
  $: urn = bucket + citation
  $: if (!!textarea) {
    textarea.disabled = !isVisible || !dataloaded
  }
  // $: console.log('!!!!!!!!!!!!!!', id_col, bucket, citation)

  async function getPassage(id_col, urn) {
    const res = await fetch(`/api/v1/passage/${id_col}/${urn}`)
    if (res.ok) {
      const d = await res.json()
      passagetext = d.data.transcriptionLines.join('\r\n')
      dataloaded = true
      return d.data
    } else {
      throw new Error(res.body)
    }
  }

  async function SaveBlockEditorTrancription() {
    //console.log('passagetext: ', passagetext)
    isLoading = true
    let data = {
      colid: parseInt(id_col),
      passageid: urn,
      text: passagetext,
    }
    // console.log('passagetext: ', data)

    const res = await fetch('/api/v2/savepassagetext', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })
    // .then((res) => {
    //   console.log('Request complete! response:', res)
    //   if (!textarea.disabled) {
    //     getPassage(id_col, urn)
    //   }
    //   isLoading = false
    //   // location.reload(true)
    // })

    if (res.ok) {
      modalHeader = 'Success'
      modalText = 'Witness saved'
      showModal = true
      // showMetadata = false
    } else {
      // Handle error
      console.error(`Ingestion failed: HTTP ${res.status} ${await res.text()}`)
      modalHeader = 'Error'
      modalText = 'An error occurred in the save transcription process.'
      showModal = true
      return
    }
    if (!textarea.disabled) {
      getPassage(id_col, urn)
    }
    isLoading = false
  }

  function countRows(text) {
    if (!text) {
      return 1
    }
    return text.split(/\n/).length
  }

  onMount(() => {
    updateTextareaState()
    window.addEventListener('scroll', updateTextareaState)
    window.addEventListener('resize', updateTextareaState)
    //getPassage(id_col, urn)
  })

  function updateTextareaState() {
    if (textarea) {
      isVisible = isElementInViewport(textarea)
      if (isVisible && !dataloaded) {
        getPassage(id_col, urn)
      }
    }
  }

</script>

<style>
  textarea.loading {
    cursor: progress;
  }

  textarea {
    overflow: hidden;
    white-space: pre;
    resize: none;
    width: 98%;
  }

</style>

<textarea
  bind:this={textarea}
  class="dynamicTextarea specialKey {isLoading ? 'loading' : ''}"
  rows={textareaRows}
  name="transcription"
  bind:value={passagetext}
  on:change={SaveBlockEditorTrancription} />

{#if showModal}
  <Modal bind:showModal>
    <h3 slot="header">{modalHeader}</h3>

    <body>
      <h4>{modalText}</h4>
    </body>
  </Modal>
{/if}
