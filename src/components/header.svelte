<script>
  import { getContext, onMount } from 'svelte'

  const client: any = getContext('client')

  let user: any
  let loading: boolean = true

  onMount(async () => {
    const response = await $client.get(`/api/user/get`)
    if (response.ok) {
      user = response.body
      loading = false
    }
  })

  const logIn = async () => {
    const response = await $client.get(`/api/login`)
    if (response.ok) {
      location.reload() // TODO: temporary, should be fixed after resolving auth opening in other tab
    }
  }
</script>

<nav class="bg-spotify-green">
  <div class="max-w-7xl mx-auto px-5 lg:px-12">
    <div class="relative flex items-center justify-between h-16">
      <div class="flex-1 flex items-center sm:items-stretch justify-start">
        <div class="flex-shrink-0 flex items-center">
          <p class="text-white text-2xl font-bold">Awesome Spotify</p>
        </div>
      </div>
      <div class="absolute inset-y-0 right-0 flex items-center sm:static sm:inset-auto sm:ml-6 sm:pr-0">
        <div class="ml-3 relative">
          <div>
            {#if !loading && user == null}
              <button class="bg-spotify-dark-green text-white px-3 py-2 rounded-md text-sm font-medium" on:click={logIn}>Log In</button>
            {/if}
          </div>
        </div>
      </div>
    </div>
  </div>
</nav>
