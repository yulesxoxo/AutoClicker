<script lang="ts">
  import { Greet } from "@wails/greetservice";
  let name = $state("");
  let greetResult: string | null = $state(null);

  async function doGreet() {
    if (name) {
      greetResult = await Greet(name);
    }
  }
</script>

<h1 class="pb-2 text-center h1 text-3xl select-none">
  Wails + TS + SvelteKit + Skeleton
</h1>
<div
  class="flex flex-col space-y-4 overflow-hidden card bg-surface-100-900/50 p-4 text-center select-none"
>
  {#if greetResult}
    {greetResult}
  {:else}
    Please enter your name below 👇
  {/if}

  <input bind:value={name} class="input" type="text" autocomplete="off" />
  <button
    class="btn preset-filled-primary-100-900 hover:preset-filled-primary-500"
    onclick={doGreet}
  >
    Greet
  </button>
</div>
