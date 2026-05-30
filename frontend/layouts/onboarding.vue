<template>
  <div
    class="font-sans antialiased relative min-h-screen overflow-x-hidden bg-[color:var(--mm-bg)] text-[color:var(--mm-text)]"
  >
    <SharedBackgroundFlow />
    <div
      class="pointer-events-none absolute inset-0 z-[1] bg-[color:var(--mm-bg)]/92"
      aria-hidden="true"
    />

    <main
      class="relative z-10 flex min-h-[100dvh] items-center justify-center px-4 py-10 sm:py-14"
    >
      <NuxtPage />
    </main>
  </div>
</template>
