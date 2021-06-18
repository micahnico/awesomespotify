export default class Client {
  // request wraps fetch in a non-throwing interface. It automatically serializes options['body'] into JSON and sets
  // the appropriate Content-Type header unless Content-Type is already set. It automatically includes the API key.
  async request(method, url, options = {}) {
    let fetchOptions = {
      method: method,
      headers: new Headers(options.headers),
      body: options.body,
    }
    if (!fetchOptions.headers.has('Content-Type') && fetchOptions.body) {
      fetchOptions.headers.set('Content-Type', 'application/json; charset=utf-8')
      fetchOptions.body = JSON.stringify(fetchOptions.body)
    }

    let httpResponse
    try {
      httpResponse = await fetch(url, fetchOptions)

      if (httpResponse.ok) {
        let body
        if (httpResponse.headers.get('Content-Type') === 'application/json; charset=utf-8') {
          body = await httpResponse.json()
        } else {
          body = await httpResponse.text()
        }
        return { ok: true, body: body, response: httpResponse }
      } else {
        let error
        if (httpResponse.headers.get('Content-Type') === 'application/json; charset=utf-8') {
          const body = await httpResponse.json()
          error = body.error
        } else {
          const body = await httpResponse.text()
          error = {code: 'unknown_error', message: body}
        }
        return { ok: false, error: error, response: httpResponse }
      }
    } catch (err) {
      return { ok: false, error: {code: 'client_exception', message: err}, response: httpResponse}
    }
  }

  // get is a convenience wrapper around request.
  get(url, options = {}) {
    return this.request('GET', url, options)
  }

  // post is a convenience wrapper around request.
  post(url, body, options= {}) {
    options.body = body
    return this.request('POST', url, options)
  }

  // put is a convenience wrapper around request.
  put(url, body, options = {}) {
    options.body = body
    return this.request('PUT', url, options)
  }

  // patch is a convenience wrapper around request.
  patch(url, body, options= {}) {
    options.body = body
    return this.request('PATCH', url, options)
  }

  // delete is a convenience wrapper around request.
  delete(url, options= {}) {
    return this.request('DELETE', url, options)
  }
}
