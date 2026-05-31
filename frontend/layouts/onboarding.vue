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
      class="mm-onboarding-layout-main relative z-10 flex min-h-[100dvh] w-full items-start justify-center px-4 py-5 sm:px-6 sm:py-6 md:px-8 lg:py-8"
    >
      <NuxtPage />
    </main>
  </div>
</template>
