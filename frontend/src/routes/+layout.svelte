<script>
  import "../app.css";

  let { children } = $props();
</script>

<div class="flex h-screen flex-col overflow-hidden">
  <main class="w-full pt-2 pr-4 pb-4 pl-4">
    {@render children()}
  </main>
</div>
