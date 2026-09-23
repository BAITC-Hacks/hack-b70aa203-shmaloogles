<template>
  <div class="hero-art" aria-hidden="true">
    <div class="art-orbit orbit-one" /><div class="art-orbit orbit-two" />
    <svg class="bridge-illustration" viewBox="0 0 480 340" fill="none">
      <defs>
        <linearGradient id="bridge-face" x1="135" y1="90" x2="320" y2="290" gradientUnits="userSpaceOnUse"><stop stop-color="#74ACFF" /><stop offset="1" stop-color="#3777F4" /></linearGradient>
        <linearGradient id="bridge-side" x1="275" y1="138" x2="360" y2="277" gradientUnits="userSpaceOnUse"><stop stop-color="#3A7AEE" /><stop offset="1" stop-color="#2860CF" /></linearGradient>
        <linearGradient id="bridge-base" x1="80" y1="214" x2="387" y2="302" gradientUnits="userSpaceOnUse"><stop stop-color="#D0E3FF" /><stop offset="1" stop-color="#EAF3FF" /></linearGradient>
      </defs>
      <ellipse cx="249" cy="292" rx="156" ry="25" fill="#6594D5" opacity=".09" />
      <path d="m72 235 168-73 179 89-169 76L72 235Z" fill="url(#bridge-base)" />
      <path d="m72 235 178 81 169-76v13l-169 76-178-81v-13Z" fill="#C0D8FB" />
      <path d="m110 232 25-12 143 71-25 12-143-71Z" fill="#ACCDF9" />
      <path d="m155 192 155-71 39 20-155 71-39-20Z" fill="#BAD8FF" />
      <path d="M194 212v-28c0-70 98-100 98-26v24l43 21v-58c0-97-141-85-141 4v63Z" fill="url(#bridge-side)" />
      <path d="m155 252 39-18v-85c0-89 141-101 141-4v58l-39 18v-58c0-65-102-64-102 3v68l-39 18Z" fill="#255FCE" />
      <path d="M155 252v-89c0-99 141-101 141-3v61l-43-21v-58c0-48-55-48-55 0v89l-43 21Z" fill="url(#bridge-face)" />
      <path d="M155 164c0-99 141-101 141-3" stroke="#A7CEFF" stroke-width="2" />
      <path d="m156 252 42-21 1-88" stroke="#C0DCFF" stroke-opacity=".45" stroke-width="2" />
      <path d="m94 216 8-4m5-2 8-4m5-2 8-4 M333 274l8-4m5-2 8-4m5-2 8-4" stroke="white" stroke-width="3" stroke-linecap="round" />
      <path d="M86 150c-39 34-43 87 2 101 M374 112c60 9 77 69 28 106" stroke="#80ACED" stroke-width="1.5" stroke-dasharray="5 6" />
      <circle cx="86" cy="150" r="5" fill="#6CA5F7" /><circle cx="402" cy="218" r="5" fill="#6CA5F7" />
      <path d="m365 83 4 12 12 4-12 4-4 12-4-12-12-4 12-4 4-12Z" fill="#78A9F8" />
      <path d="m119 96 3 8 8 3-8 3-3 8-3-8-8-3 8-3 3-8Z" fill="#A8CBFF" />
    </svg>
    <div class="art-note business-note"><span class="art-note-icon"><AppIcon name="briefcase" :size="21" /></span><div><small>ЕСТЬ ЗАДАЧА</small><strong>Идея бизнеса</strong></div><span class="art-note-dot" /></div>
    <div class="art-note team-note"><span class="art-note-icon"><AppIcon name="users" :size="21" /></span><div><small>ЕСТЬ РЕШЕНИЕ</small><strong>Энергия команды</strong></div></div>
    <div class="art-match"><span><AppIcon name="check" :size="13" /></span>Вместе — больше</div>
    <span class="art-coordinate coordinate-left">55° IDEAS</span><span class="art-coordinate coordinate-right">71° OPPORTUNITIES</span>
  </div>
</template>
