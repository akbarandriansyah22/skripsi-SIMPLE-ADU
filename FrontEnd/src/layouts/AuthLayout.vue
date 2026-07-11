<template>
  <div
    class="min-h-screen bg-slate-50 flex items-center justify-center px-4 py-8"
  >
    <div
      class="w-full max-w-5xl overflow-hidden rounded-2xl bg-white shadow-elevated ring-1 ring-slate-200/50"
    >
      <div class="grid grid-cols-1 lg:grid-cols-[1.1fr_0.9fr]">
        <div
          class="hidden lg:flex flex-col justify-center gap-8 bg-gradient-to-br from-slate-900 to-slate-800 px-10 py-12 text-white"
        >
          <div>
            <h1 class="text-4xl font-bold">SIMPEL-ADU</h1>
            <p class="mt-4 max-w-xs text-slate-300 leading-relaxed">
              Platform layanan pengaduan terpadu untuk mahasiswa Fakultas Teknik
              UMJ dengan akses cepat, sistematis, dan aman.
            </p>
          </div>
          <div
            class="grid gap-4 rounded-2xl border border-slate-700/50 bg-slate-800/60 p-6 shadow-lg"
          >
            <div class="space-y-2">
              <p
                class="text-xs uppercase tracking-widest text-blue-400 font-semibold"
              >
                Langkah Mudah
              </p>
              <p class="text-base font-semibold text-white leading-snug">
                Masuk dan sampaikan aspirasi Anda dalam beberapa langkah.
              </p>
            </div>
            <div class="grid gap-3 text-sm text-slate-400 leading-relaxed">
              <p>1. Gunakan email kampus Anda.</p>
              <p>2. Isi detail aduan dengan jelas.</p>
              <p>3. Pantau status langsung dari dashboard.</p>
            </div>
          </div>
        </div>
        <div class="px-6 py-8 sm:px-10 sm:py-12">
          <div class="max-w-md mx-auto">
            <slot />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup></script>
