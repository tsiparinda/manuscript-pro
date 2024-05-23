function s3_bucket_url_js(key) {
  const s3 = new AWS.S3()
  var params = { Bucket: 'brucheion', Key: key, Expires: 60 * 300 }
  var s3_presigned = s3.getSignedUrl('getObject', params)
  return s3_presigned
}

async function listDZIFiles(bucketName, prefix) {
  const s3 = new AWS.S3()
  const params = {
    Bucket: bucketName,
    Prefix: prefix,
    MaxKeys: 100000,
  }

  try {
    const result = await s3.listObjectsV2(params).promise()

    if (!result) {
      console.error('Error: The result object is undefined.')
      return []
    }

    if (!result.Contents) {
      console.error('Error: The result.Contents property is undefined.', result)
      return []
    }

    console.log('Result:', result)
    const dziFiles = result.Contents.filter((file) =>
      file.Key.endsWith('.dzi')
    ).map((file) => file.Key)
    return dziFiles
  } catch (error) {
    console.error('Error listing DZI files:', error)
    throw error
  }
}

async function listDZIcollection(bucketName, prefix) {
  const s3 = new AWS.S3()

  let result
  let isTruncated = true
  let continuationToken
  let allContents = []
  const delimiter = '/'

  while (isTruncated) {
    const params = {
      Bucket: bucketName,
      Prefix: prefix,
      //Delimiter: delimiter,
      ContinuationToken: continuationToken,
    }

    try {
      result = await s3.listObjectsV2(params).promise()
      isTruncated = result.IsTruncated
      continuationToken = result.NextContinuationToken

      if (!result || !result.Contents) {
        console.error(
          'Error: The result object or its Contents property is undefined.',
          result
        )
        return []
      }

      allContents = allContents.concat(result.Contents)
    } catch (error) {
      console.error('Error listing DZI files:', error)
      throw error
    }
  }

  const dziFiles = allContents
    .filter((file) => file.Key.endsWith('.dzi'))
    .map((file) => file.Key)
  // return  array
  // nbh/J1img/positive/J1_55v.dzi
  return dziFiles
}

export { listDZIFiles, s3_bucket_url_js, listDZIcollection }
