<script lang="ts">
  import {onMount} from "svelte";
  import {themeChange} from "theme-change";
  import {browser} from "$app/environment";
  import {navigating} from "$app/stores";

  function update() {
    const checkbox: any = document.getElementById('theme-selector-ckb')

    checkbox.checked = localStorage.getItem('theme') === 'light';
  }

  onMount(() => {
    themeChange(false)

    // if the html tag does not have the data-theme attribute, set it to the preferred theme (set by browser)
    if (!document.documentElement.hasAttribute('data-theme')) {
      if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        document.documentElement.setAttribute('data-theme', 'dark')
        localStorage.setItem('theme', 'dark')
      } else {
        document.documentElement.setAttribute('data-theme', 'light')
        localStorage.setItem('theme', 'light')
      }
    }

    setTimeout(() => {
      update()
    }, 25)
  })
</script>

<label class="swap swap-flip text-3xl">
    <input on:change={update} id="theme-selector-ckb" type="checkbox"/>

    <div data-set-theme="dark" class="swap-on">🌞</div>
    <div data-set-theme="light" class="swap-off">🌚</div>
</label>
