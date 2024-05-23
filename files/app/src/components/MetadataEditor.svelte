<script>
  import { onMount } from 'svelte'
  import Message from '../components/Message.svelte'
  import Modal from './Modal.svelte'

  export let passage,
    savedata,
    canEditMetadata = false

  let formData = {
    urn: '',
    citationScheme: '',
    groupName: '',
    workTitle: '',
    versionLabel: '',
    exemplarLabel: '',
    online: '',
    language: '',
  }

  let showModal = false,
    modalHeader,
    modalText

  onMount(() => {
    formData = { ...passage.catalog }
  })

  $: handleSubmit(savedata)

  async function handleSubmit() {
    if (!savedata) return
    // formData.colid= parseInt(id_col)
    // console.log('!!!!!', JSON.stringify(formData))
    const response = await fetch('/api/v2/savemetadata', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(formData),
    })

    if (response.ok) {
      // Handle successful POST action
      passage.catalog = { ...formData }

      modalHeader = 'Success'
      modalText = 'Metadata saved'
      showModal = true
      // showMetadata = false
    } else {
      // Handle error
      modalHeader = 'Error'
      modalText = 'Metadata not saved'
      showModal = true
    }
    savedata = false
  }

</script>

<style>
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

<form on:submit|preventDefault={handleSubmit}>
  <dl>
    <dt>
      Work URN
      <i
        class="fa fa-question-circle-o"
        data-title="This field is non-editable" />
    </dt>
    <dd>
      <input
        type="text"
        class="form-control"
        bind:value={formData.urn}
        readonly />
    </dd>
    <dt>
      Scheme
      <i
        class="fa fa-question-circle-o"
        data-title="Citation scheme of the work for example, 'a.b.c' could resemble 'adhyāya / āhnika / sūtra'" />
    </dt>
    <dd>
      <input
        type="text"
        class="form-control"
        bind:value={formData.citationScheme}
        readonly={!canEditMetadata} />
    </dd>
    <dt>
      Workgroup
      <i
        class="fa fa-question-circle-o"
        data-title="Workgroup in natural language" />
    </dt>
    <dd>
      <input
        type="text"
        class="form-control"
        bind:value={formData.groupName}
        readonly={!canEditMetadata} />
    </dd>
    <dt>
      Title
      <i
        class="fa fa-question-circle-o"
        data-title="Title in natural language" />
    </dt>
    <dd>
      <input
        type="text"
        class="form-control"
        bind:value={formData.workTitle}
        readonly={!canEditMetadata} />
    </dd>
    <dt>
      Version
      <i
        class="fa fa-question-circle-o"
        data-title="Version in natural language" />
    </dt>
    <dd>
      <input
        type="text"
        class="form-control"
        bind:value={formData.versionLabel}
        readonly={!canEditMetadata} />
    </dd>
    <dt>
      Exemplar
      <i
        class="fa fa-question-circle-o"
        data-title="Exemplar in natural language" />
    </dt>
    <dd>
      <input
        type="text"
        class="form-control"
        bind:value={formData.exemplarLabel}
        readonly={!canEditMetadata} />
    </dd>
    <dt>
      Online
      <i
        class="fa fa-question-circle-o"
        data-title="Boolean type; usually 'true' or 'false'" />
    </dt>
    <dd>
      <input
        type="text"
        class="form-control"
        bind:value={formData.online}
        readonly={!canEditMetadata} />
    </dd>
    <dt>
      Language
      <i class="fa fa-question-circle-o" data-title="Language ID" />
    </dt>
    <dd>
      <input
        type="text"
        class="form-control"
        bind:value={formData.language}
        readonly={!canEditMetadata} />
    </dd>
  </dl>
  <!-- <button type="submit">Save</button> -->
</form>

{#if showModal}
  <Modal bind:showModal>
    <h3 slot="header">{modalHeader}</h3>
    <body>
      <h4>{modalText}</h4>
    </body>
  </Modal>
{/if}
