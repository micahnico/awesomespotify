<script>
  import Header from '../components/header.svelte'
  import {link} from 'svelte-spa-router'
  import { getContext, onDestroy, onMount } from 'svelte'

  const client = getContext('client')

  let artists
  let song
  let urlSafeSong
  let lyrics
  let imageUrl
  let bgHex
  let txtHex
  let findLyricsError

  let user
  let loading = true
  let lookupInterval

  onMount(async () => {
    const userResponse = await $client.get(`/api/user/get`)
    if (userResponse.ok) {
      user = userResponse.body
    }

    if (user) {
      const lyricResponse = await $client.get(`/api/find`)
      if (lyricResponse.ok) {
        artists = lyricResponse.body.Artists
        song = lyricResponse.body.Song
        urlSafeSong = lyricResponse.body.URLSafeSong || ""
        lyrics = lyricResponse.body.Lyrics
        imageUrl = lyricResponse.body.ImageURL
        bgHex = lyricResponse.body.BgHex
        txtHex = lyricResponse.body.TxtHex
        findLyricsError = lyricResponse.body.Error
      }
    }

    lookupInterval = setInterval(autoDetectNewSong, 5000)
    loading = false
  })

  onDestroy(() => {
    clearInterval(lookupInterval)
  })

  const logIn = async () => {
    const response = await $client.get(`/api/login`)
    if (response.ok) {
      location.href = response.body.url
    }
  }

  const autoDetectNewSong = async () => {
    if (!loading && user && !document.hidden) {
      const lyricResponse = await $client.get(`/api/find?currentSong=${urlSafeSong.replace(" ", "%20")}`)
      if (lyricResponse.ok && lyricResponse.body.Error != "Already fetched lyrics" && lyricResponse.body.Error != "No currently playing song") {
        artists = lyricResponse.body.Artists
        song = lyricResponse.body.Song
        urlSafeSong = lyricResponse.body.URLSafeSong || ""
        lyrics = lyricResponse.body.Lyrics
        imageUrl = lyricResponse.body.ImageURL
        bgHex = lyricResponse.body.BgHex
        txtHex = lyricResponse.body.TxtHex
        findLyricsError = lyricResponse.body.Error
      }
    }
  }
</script>

<style>
  @keyframes rotation {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(359deg);
    }
  }
  .rotate {
    animation: rotation 5s infinite linear;
  }

  .custom-img-size {
    height: 200px;
    width: 200px;
  }

  .currently-playing-info {
    display: flex;
    flex-direction: column;
  }

  @media (min-width: 640px) {
    .custom-img-size {
      height: 250px;
      width: 250px;
    }
  }

  @media (min-width: 1280px) {
    .currently-playing-info {
      flex-direction: row;
    }
  }

  @media (min-width: 1536px) {
    .custom-img-size {
      height: 300px;
      width: 300px;
    }
  }
</style>

<Header/>

<div class="min-h-screen" style="background: #111111;">
  <div class="w-full flex justify-center p-5">
    <div class='w-full lg:w-3/4 xl:w-3/5 bg-white shadow-lg'>
      {#if loading}
        <div class="currently-playing-info items-center p-10 bg-spotify-green" style="color: #111111;">
          <img src="vinyl.png" alt="" class="custom-img-size rotate">
          <div class="mt-5 xl:mt-0 xl:ml-5">
            <p class='text-2xl sm:text-3xl md:text-4xl font-bold mb-1'>Finding Lyrics...</p>
          </div>
        </div>
      {:else}
        {#if user}
          {#if findLyricsError != "No currently playing song"}
            <div class="currently-playing-info items-center p-10" style="background: {bgHex}; color: {txtHex};">
              <img src={imageUrl} alt="" class="border custom-img-size">
              <div class="mt-5 xl:mt-0 xl:ml-5">
                <p class='text-2xl sm:text-3xl md:text-4xl font-bold mb-1'>{song}</p>
                <p class='text-xl md:text-2xl'>
                  {#each artists as artist, i}
                    {#if i}, {/if}{artist}
                  {/each}
                </p>
              </div>
            </div>
          {/if}
          {#if lyrics}
            <div class="pt-7 pb-10 px-10 text-white" style="background: #212020;">
              {@html lyrics}
              <p class="text-gray-500 mt-7">These lyrics were taken from <a href="https://genius.com" class="text-blue-500 hover:text-blue-700">genius.com</a></p>
            </div>
          {:else}
            <p class='py-7 px-10 text-xl font-bold text-white' style="background: #212020;">{findLyricsError} 🙁</p>
          {/if}
        {:else}
          <div class="p-10 text-white" style="background: #212020;">
            <p class='text-3xl font-bold mb-1'>Find Lyrics</p>
            <p><span class="cursor-pointer link" on:click={logIn}>Log in with Spotify</span> to access</p>
          </div>
        {/if}
      {/if}
    </div>
  </div>
</div>
