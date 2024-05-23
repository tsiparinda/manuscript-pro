<script>
  export let lineIndex
  export let transcriptionEditNav
  export let currentTranscriptionLines = []
  export let isLineEditor
  export let transcriptionFontSize

  $: transcriptionLineNumbers = updateLineNumber(currentTranscriptionLines)

  $: if (isLineEditor) {
    transcriptionEditNav = true
  } else {
    transcriptionEditNav = false
  }

  // $: console.log(
  //   'Selected line index:',
  //   lineIndex,
  //   currentTranscriptionLines.length,
  //   currentTranscriptionLines,
  //   isLineEditor
  // )

  function handleEditorAddLineBeforeClick() {
    currentTranscriptionLines.splice(lineIndex, 0, '#')
    currentTranscriptionLines = [...currentTranscriptionLines]
    console.log(
      '!handleEditorAddLineBeforeClick',
      lineIndex,
      currentTranscriptionLines
    )
  }

  function handleEditorAddLineAfterClick() {
    currentTranscriptionLines.splice(parseInt(lineIndex) + 1, 0, '#')
    currentTranscriptionLines = [...currentTranscriptionLines]
    lineIndex = parseInt(lineIndex) + 1
    console.log(
      '!handleEditorAddLineAfterClick',
      lineIndex,
      currentTranscriptionLines
    )
  }

  function handleEditorDelLineClick() {
    currentTranscriptionLines.splice(lineIndex, 1)
    currentTranscriptionLines = [...currentTranscriptionLines]
    //  passage.transcriptionlines = currentTranscriptionLines.map((obj) => obj.line)
  }

  function updateLineNumber(strarray) {
    let count = 0
    let lineNumbers = []
    strarray.forEach((line, index) => {
      let val = getNumeralBeforeClosingBrace(line)
      if (val > 0) {
        count = val
      } else if (val === null) {
        count = count + 1
      } else {
        count = 1
      }
      lineNumbers.push(count)
    })
    // console.log('!!!!', lineNumbers)
    return lineNumbers
  }

  function getNumeralBeforeClosingBrace(code) {
    const closingBraceIndex = code.indexOf('}')
    if (closingBraceIndex === -1 || closingBraceIndex === 0) {
      return null
    } else {
      let numeral = ''
      let i = closingBraceIndex - 1
      while (i >= 0 && /[0-9]/.test(code.charAt(i))) {
        numeral = code.charAt(i) + numeral
        i--
      }
      return numeral === '' ? null : parseInt(numeral, 10)
    }
  }

</script>

<style>
  .toolbar_editor {
    align-items: center;
    height: 40px;
    margin: 0;
    margin-left: 20px;
    list-style: none;
    padding: 0;
  }

  .btn_toolbar {
    background: #147cb4;
    border-radius: 20px;
    border: 1px solid #ffffff;
    font-family: 'Montserrat';
    font-style: normal;
    font-weight: 600;
    font-size: 12px;
    line-height: 13px;
    text-align: center;
    color: #ffffff;
  }

  .btn_toolbar:hover {
    background: #1ca6f1;
    color: #ffffff;
  }

  .btn_toolbar:active {
    background: #0e5981;
    color: #ffffff;
  }

  a.tip {
    border-bottom: 1px dashed;
    text-decoration: none;
  }
  a.tip:hover {
    cursor: context-menu;
    position: relative;
    transition: ease 1s;
  }
  a.tip span {
    display: none;
  }
  a.tip:hover span {
    border: #fbb7b7 1px;
    opacity: 0.9;
    border-radius: 20px;
    padding: 5px 20px 5px 5px;
    display: inline-block;
    z-index: 100;
    background: url(../images/status-info.png) #fff1a0 no-repeat 100% 5%;
    left: 0px;
    margin: 10px;
    width: 150px;
    position: absolute;
    top: 20px;
    text-decoration: none;
    color: #020202;
    font-style: italic;
    word-wrap: break-word;
    white-space: break-spaces;
    line-height: 1.5;
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

<div class="transcription" style="font-size: {transcriptionFontSize};">
  {#each currentTranscriptionLines as line, i}
    {#if transcriptionEditNav && lineIndex == i}
      <div class="row toolbar_editor">
        <a
          class="btn btn_toolbar tip my-2 my-sm-0"
          data-title="This button will insert one line above"
          href="#top"
          on:click={handleEditorAddLineBeforeClick}>
          Add
          <i class="fa fa-level-up fa-lg" />
          <!-- <span>This button will insert one line above</span> -->
        </a>
        <a
          class="btn btn_toolbar tip my-2 my-sm-0"
          data-title="This button will insert one line below"
          href="#top"
          on:click={handleEditorAddLineAfterClick}>
          Add
          <i class="fa fa-level-down fa-lg" />
          <!-- <span>This button will insert one line below</span> -->
        </a>
        <a
          class="btn btn_toolbar tip my-2 my-sm-0"
          data-title="This button will delete the current line"
          href="#top"
          on:click={handleEditorDelLineClick}>
          Del
          <i class="fa fa-remove fa-lg" />
          <!-- <span>This button will delete the current line</span> -->
        </a>
      </div>
    {/if}
    <div>
      <span class="font-weight-bold" style="display: table-cell; width: 40px;">
        {transcriptionLineNumbers[i]}:
      </span>
      <span
        style="margin-left: 25px; display: table-cell; webkit-line-break:
        normal;"
        contenteditable={`${isLineEditor ? true : false}`}
        data-index={i}>
        {line}
      </span>
    </div>
  {/each}
</div>
