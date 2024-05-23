<script>
  import Objects from '../components/Objects.svelte'
  import jQuery from 'jquery'
  export let id_col

  // console.info('from routes share collections colid: ', id_col)

  async function getShareCollection(id_col) {
    const res = await fetch(`/api/v2/sharecollection/${id_col}`)
    if (res.ok) {
      const d = await res.json()
      console.log(d)
      return d
    } else {
      throw new Error(res.body)
    }
  }

  async function getUser() {
    const res = await fetch(`/api/v2/users/0`)
    if (res.ok) {
      const d = await res.json()
      return d.data
    } else {
      throw new Error(res.body)
    }
  }

  let results = []
  let possibleResults = []
  let peopleWithAccess = []
  let access = false
  let readOnly = false
  let makePublic = false

  let searchInput = {
    user_name: '',
  }
  let isFocused = false
  const onFocus = () => (isFocused = true)
  const onBlur = () => (isFocused = false)
  let collection, colUsers, err
  ;(async () => {
    const result = Promise.all([getShareCollection(id_col), getUser()])
      .then(([c, u]) => {
        collection = c.collection
        colUsers = c.colusers
        possibleResults = u
        makePublic = collection.is_public
      })
      .catch((e) => (err = e))

    await result
    console.log('!!!', colUsers)
    var cU = (colUsers ?? []).map((x) => x.user_id)
    peopleWithAccess = possibleResults.filter((x) => cU.includes(x.user_id))
    //mix colusers with user
    peopleWithAccess = peopleWithAccess.map((person) => {
      const colUser = colUsers.find((x) => x.user_id === person.user_id)
      if (colUser) {
        return { ...person, ...colUser }
      }
      return person
    })
  })()

  const typeahead = () => {
    let resultsIncludes = possibleResults.filter((possibleResults) =>
      possibleResults.user_name
        .toLowerCase()
        .includes(searchInput.user_name.toLowerCase())
    )
    let resultsStartWith = possibleResults.filter((possibleResults) =>
      possibleResults.user_name
        .toLowerCase()
        .startsWith(searchInput.user_name.toLowerCase())
    )
    results = resultsStartWith.sort().concat(resultsIncludes.sort())
    results = [...new Set(results)]
  }

  const newSearchInput = (oneResult) => {
    searchInput = oneResult
  }

  function handleDeleteColUser(userid) {
    let data = {
      colid: collection.id,
      userid: userid,
    }
    fetch('/api/v2/collectionsuser', {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    }).then((res) => {
      // console.log('Request complete! response:', res)
      location.reload(true)
    })
  }

  function handleSubmit() {
    makePublic = !makePublic

    collection.is_public = makePublic

    let selUser = {
      col_id: collection.id,
      user_id: searchInput.user_id,
      is_write: access,
    }

    colUsers = (colUsers ?? []).concat(selUser)

    let data = {
      collection: collection,
      colusers: colUsers,
    }

    // console.log(data)

    fetch('/api/v2/sharecollection', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    }).then((res) => {
      // console.log('Request complete! response:', res)
      location.reload(true)
    })
  }

</script>

<style>
  .typeahead {
    position: relative;
  }
  .ip {
    margin-bottom: 0;
    box-shadow: rgba(0, 0, 0, 0.16) 0px 3px 6px, rgba(0, 0, 0, 0.23) 0px 3px 6px;
  }
  .ip:hover {
    box-shadow: rgba(0, 0, 0, 0.16) 0px 4px 8px, rgba(0, 0, 0, 0.23) 0px 4px 8px;
  }
  input[type='text'].ip {
    width: 60%;
    display: block;
    text-align: center;
    border: 3px solid black;
    border-radius: 0;
    margin-right: auto;
    transition: width 0.4s ease-in-out;
  }
  input[type='text'].ip:focus {
    background-color: whitesmoke;
    outline: none;
    width: 60%;
  }
  .ulist {
    max-height: 215px;
    overflow: auto;
  }

  .ulist.typeahead-object-list {
    position: absolute;
    padding: 0;
    margin-top: 0;
    width: 60%;
    z-index: 99;
    margin-right: auto;
    background-color: whitesmoke;
  }

</style>

<nav aria-label="breadcrumb">
  <ol class="breadcrumb">
    <li class="breadcrumb-item"><a href="/">Home</a></li>
    <li class="breadcrumb-item active" aria-current="page">Share Collection</li>
  </ol>
</nav>

<div class="container-fluid">
  <div class="row">
    <div class="col-5 mx-auto">
      <div>
        <img
          src="/files/img/CollectionImg.svg"
          class="img-thumbnail"
          alt="Profile" />
      </div>
      <br />
      <div class="row">
        {#if collection && !err}
          <h3 style="width: 90%;">Share collection "{collection.title}"</h3>
        {:else if err}
          <h5>Error: {err}</h5>
        {/if}

        {#if makePublic}
          <button
            class="tag tag-public mb-2 mx-right"
            id="tag-text"
            type="submit"
            on:click={handleSubmit}>
            PUBLIC
          </button>
        {:else}
          <button
            class="tag tag-private mb-2 mx-right"
            id="tag-text"
            type="submit"
            on:click={handleSubmit}>
            PRIVATE
          </button>
        {/if}
      </div>
      <br />

      <!-- <label>
        <input
          type="checkbox"
          name="Make Public"
          disabled={access}
          bind:checked={makePublic} />
        Make Public
      </label> -->
      <br />

      <h4>Select users to give access to:</h4>
      <div class="typeahead row">
        <input
          id="searchfield"
          type="text"
          name="searchfield"
          placeholder="Search"
          class="ip"
          value={searchInput.user_name}
          on:input={typeahead}
          on:focus={onFocus}
          on:blur={onBlur}
          style="width: 40%;" />
        <label>
          <input
            type="checkbox"
            name="full access"
            disabled={readOnly}
            bind:checked={access} />
          Full Access
        </label>
        <label>
          <input
            type="checkbox"
            name="read only"
            disabled={access}
            bind:checked={readOnly} />
          Read Only
        </label>
        <button
          type="button"
          class="btn btn-primary"
          disabled={!((access || readOnly) && searchInput != '')}
          on:click={handleSubmit}>
          Submit
        </button>
      </div>
      <ul class="ulist typeahead-object-list">
        {#if isFocused === true}
          {#if searchInput.length === 0}
            {#each possibleResults as oneResult}
              <Objects
                object={oneResult.user_name}
                on:mousedown={() => newSearchInput(oneResult)} />
            {/each}
          {:else}
            {#each results as oneResult}
              <Objects
                object={oneResult.user_name}
                on:mousedown={() => newSearchInput(oneResult)} />
            {/each}
          {/if}
        {/if}
      </ul>
      <br />
      <br />
      <div>
        {#if colUsers && !err}
          <h4>People with access</h4>
          <div class="row">
            {#each peopleWithAccess as pWA}
              <span style="width: 20%;">{pWA.user_name}</span>
              <span style="width: 50%;">{pWA.user_email}</span>
              {#if pWA.is_write}
                <span style="width: 20%;">READ/WRITE</span>
              {:else}<span style="width: 20%;">READ ONLY</span>{/if}
              <button
                type="submit"
                class="btn btn-primary"
                on:click={() => handleDeleteColUser(pWA.user_id)}>
                DEL
              </button>
            {/each}
          </div>
        {:else if err}
          <h5>Error: {err}</h5>
        {/if}
      </div>
    </div>
  </div>
</div>
