<script>
  import Modal from '../components/Modal.svelte'

  let inputRef, formRef
  let errorMessage = null
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

    console.log('title', title.length)
    if (title.length == 0) {
      modalHeader = 'Failure'
      modalText = errorMessages.null_title
      showModal = true
      return
    }
    const formData = new FormData()
    formData.append('title', title)

    let res
    try {
      res = await fetch('/api/v2/addcollection', {
        method: 'POST',
        body: formData,
        redirect: 'error',
        credentials: 'same-origin',
      })
    } catch (err) {
      modalHeader = 'Failure'
      modalText = errorMessages.unknown + '1'
      showModal = true
      return
    }
    const data = await res.json()
    if (res.status !== 200) {
      try {
        if (typeof errorMessages[data.Message] != 'undefined') {
          errorMessage = errorMessages[data.Message]
        } else {
          errorMessage = data.Message
        }
      } catch (err) {
        modalHeader = 'Failure'
        modalText = errorMessages.unknown + '3'
        showModal = true
        return
      }
    } else {
      formRef.reset()
    }

    modalHeader = 'Success'
    modalText = 'Collection added'
    showModal = true
    window.location.href = '/view/' + data.data + '/default'
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
          <h5>Add a new collection to Brucheion.</h5>
        </div>
        <div class="form-group">
          <label for="title">Title:</label>
          <input
            type="text"
            class="form-control"
            id="title"
            bind:value={title} />
        </div>
        <div class="form-group text-left">
          <button
            class="btn btn-primary is-success"
            type="submit"
            on:click={handleSubmit}>
            Submit
          </button>
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
