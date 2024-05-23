const URN = require('urn-lib')

const citeNid = 'cite2'
const citeUrn = URN.create('urn', {
  components: ['nid', 'namespace', 'work', 'passage'],
})

function validateUrn(urn, opts = {}) {
  const noPassage = opts.noPassage || false
  const components = citeUrn.parse(urn)
  const nid = opts.nid || citeNid

  return (
    !!components &&
    components.nid === nid &&
    !!components.namespace &&
    !!components.work &&
    ((noPassage && !components.passage) || (!noPassage && !!components.passage))
  )
}
// replace last part urn to another schema value
// urn:cts:sktlit:skt0001.nyaya002.M3D:3.1.1 3.1.2 -> urn:cts:sktlit:skt0001.nyaya002.M3D:3.1.2
function replaceLastPart(str, replacement) {
  var lastIndex = str.lastIndexOf(':')
  return str.slice(0, lastIndex + 1) + replacement
}

function validateNoColon(str) {
  // Checks if the string contains a colon
  if (str.includes(':')) {
    // If it does, return false
    return false
  } else {
    // If it doesn't, return true
    return true
  }
}

module.exports = { validateUrn, replaceLastPart, validateNoColon }
