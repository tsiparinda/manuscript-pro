<script>
  import Modal from '../components/Modal.svelte'

  let inputRef, formRef
  let cexFile = ''
  let complete = false
  let loading = false
  let title = ''

  let showModal = false,
    modalHeader,
    modalText

  const errorMessages = {
    bad_file_ext:
      'The submitted file did not have the corresponding .cex file extension.',
    bad_file_body: 'The submitted file could not be read.',
    file_not_found: 'The submitted file could not be read.',
    bad_cex_data:
      'The CEX data contained erroneous data and could not be processed.',
    null_title: 'Type the title!',
    unknown:
      'An unknown error occurred. This is not necessarily due to the uploaded CEX data. Please try again and log in again if necessary.',
  }

  async function handleSubmit(event) {
    event.preventDefault()
    if (!complete) {
      return
    }

    if (inputRef.files.length < 1) {
      modalHeader = 'Failure'
      modalText = errorMessages.file_not_found
      showModal = true
      return
    }
    //console.log('title', title.length)
    if (title.length == 0) {
      loading = false
      modalHeader = 'Failure'
      modalText = errorMessages.null_title
      showModal = true
      return
    }
    const file = inputRef.files[0]
    const formData = new FormData()
    formData.append('title', title)
    formData.append('files', file)

    let res
    loading = true
    try {
      res = await fetch('/api/v2/cexupload', {
        method: 'POST',
        body: formData,
        redirect: 'error',
        credentials: 'same-origin',
      })
    } catch (err) {
      loading = false
      modalHeader = 'Failure'
      modalText = '1.' + errorMessages.unknown
      showModal = true
      return
    }
    loading = false

    const data = await res.json()
    if (res.status !== 200) {
      try {
        if (typeof errorMessages[data.message] != 'undefined') {
          modalText = errorMessages[data.message]
        } else {
          modalHeader = 'Warning'
          modalText = data.message
        }
        showModal = true
        return
      } catch (err) {
        modalHeader = 'Failure'
        modalText = '3.' + errorMessages.unknown
        showModal = true
        return
      }
    } else {
      //  console.log('!!!!!!!!!!!!!!!!!!!!data', data)
      formRef.reset()
    }
    // const data = await res.json()
    // console.log('data', data.message)
    modalText = data.message
    modalHeader = 'Success'
    showModal = true
    window.location.href = '/view/' + data.data + '/default'
  }

  $: complete = !!cexFile

  function onButtonClick() {
    inputRef.addEventListener('change', (event) => {
      const file = event.target.files[0]
      if (file) {
        const fileName = getFileName(file.name)
        cexFile = fileName
      }
    })

    inputRef.click()
  }

  function getFileName(fullPath) {
    var startIndex =
      fullPath.indexOf('\\') >= 0
        ? fullPath.lastIndexOf('\\')
        : fullPath.lastIndexOf('/')
    var filename = fullPath.substring(startIndex)
    if (filename.indexOf('\\') === 0 || filename.indexOf('/') === 0) {
      filename = filename.substring(1)
    }
    return filename
  }

</script>

<style>
  .form-group {
    display: flex;
    flex-direction: column;
  }

  .input-file {
    display: none;
  }

  .custom-button {
    background-color: #4caf50;
    border: none;
    color: white;
    padding: 8px 16px;
    text-align: center;
    text-decoration: none;
    display: inline-block;
    font-size: 16px;
    margin: 4px 2px;
    cursor: pointer;
  }

  .modal {
    display: block;
  }

  .spinner-container {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .spinner {
    border: 4px solid rgba(0, 0, 0, 0.1);
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border-left-color: #09f;
    animation: spin 1s linear infinite;
    margin-top: 10px; /* Add margin-left property */
  }

  @keyframes spin {
    100% {
      transform: rotate(360deg);
    }
  }

</style>

<div class="container-fluid">
  <div class="row">
    <div class="col-4 mx-auto ">
      <form bind:this={formRef} on:submit={handleSubmit}>
        <div class="form-group text-center">
          <h5>Add a new collection (CEX file) to Brucheion.</h5>
        </div>
        <div class="form-group">
          <label for="title">Title:</label>
          <input
            type="text"
            class="form-control"
            id="title"
            bind:value={title} />
        </div>
        <div class="form-group">
          <label for="cex-file">CEX-File:</label>
          <input
            id="cex-file"
            type="file"
            accept=".cex"
            class="input-file"
            bind:value={cexFile}
            bind:this={inputRef} />
          <h6>{cexFile}</h6>
          <button class="custom-button" on:click={onButtonClick}>Browse</button>
        </div>
        <div class="form-group text-left">
          <button
            class="btn btn-primary is-success"
            type="submit"
            class:is-loading={loading}
            disabled={!complete || loading}
            on:click={handleSubmit}>
            Submit
          </button>
          {#if loading}
            <div class="spinner-container">
              <div class="spinner" />
            </div>
          {/if}
        </div>
      </form>
    </div>
  </div>
</div>

{#if showModal}
  <Modal bind:showModal>
    <h3 slot="header">{modalHeader}</h3>
    <body>
      <h4>{modalText}</h4>
    </body>
  </Modal>
{/if}
