<script>
  import Header from '../components/header.svelte'
  import {link} from 'svelte-spa-router'
  import { getContext, onMount } from 'svelte'

  const client: any = getContext('client')

  let artists: string
  let song: string
  let lyrics: string
  let imageUrl: string
  let bgHex: string
  let findLyricsError: string

  let user: any
  let loading: boolean = true

  onMount(async () => {
    const userResponse = await $client.get(`/api/user/get`)
    if (userResponse.ok) {
      user = userResponse.body
    }

    if (user) {
      const lyricResponse = await $client.get(`/api/lyrics/find`)
      if (lyricResponse.ok) {
        artists = lyricResponse.body.Artists
        song = lyricResponse.body.Song
        lyrics = lyricResponse.body.Lyrics
        imageUrl = lyricResponse.body.ImageURL
        bgHex = lyricResponse.body.BgHex
        findLyricsError = lyricResponse.body.Error
      }
    }

    loading = false
  })

  const logIn = async () => {
    const response = await $client.get(`/api/login`)
    if (response.ok) {
      location.reload() // TODO: temporary, should be fixed after resolving auth opening in other tab
    }
  }
</script>

<style>
  .custom-img-size {
    height: 150px;
    width: 150px;
  }

  @media (min-width: 1280px) {
  .custom-img-size {
    height: 200px;
    width: 200px;
  }
}

  @media (min-width: 1536px) {
    .custom-img-size {
      height: 250px;
      width: 250px;
    }
  }
</style>

<Header/>

<div class="w-full flex justify-center p-5 bg-gray-100">
  <div class='w-full lg:w-3/4 xl:w-3/5 bg-white rounded-md shadow-lg'>
    {#if loading}
      <p class='p-10 text-xl font-bold mb-1'>Finding Lyrics...</p>
    {:else}
      {#if user}
        {#if !findLyricsError}
          <div class="flex items-center p-10 rounded-t-md mb-7" style="background-color: {bgHex};">
            <img src={imageUrl} alt="" class="border custom-img-size">
            <div class="ml-5">
              <p class='text-2xl sm:text-3xl md:text-4xl font-bold mb-1'>{song}</p>
              <p class='text-xl md:text-2xl text-gray-500'>
                {#each artists as artist, i}
                  {#if i}, {/if}{artist}
                {/each}
              </p>
            </div>
          </div>
          {#if lyrics}
            <div class="pb-10 px-10">
              {@html lyrics}
              <hr class="my-5">
              <p class="text-gray-500">These lyrics were taken from <a href="https://genius.com" class="text-blue-500 hover:text-blue-700">genius.com</a></p>
            </div>
          {/if}
        {:else}
          <p class='p-10 text-xl font-bold mb-1'>{findLyricsError}</p>
        {/if}
      {:else}
        <div class="p-10">
          <p class='text-3xl font-bold mb-1'>Lyrics</p>
          <p><span class="cursor-pointer link" on:click={logIn}>Log in with Spotify</span> to access</p>
        </div>
      {/if}
    {/if}
  </div>
</div>
