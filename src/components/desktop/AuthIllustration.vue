<template>
    <div class="auth-illustration" :class="`auth-illustration--${variant}`" aria-hidden="true">
        <svg class="auth-illustration__art" :viewBox="variant === 'signup' ? '0 0 560 760' : '0 0 720 720'"
             fill="none" xmlns="http://www.w3.org/2000/svg">
            <defs>
                <clipPath id="auth-signup-avatar-clip">
                    <circle cx="280" cy="178" r="45" />
                </clipPath>
                <pattern id="auth-grid-pattern" width="42" height="42" patternUnits="userSpaceOnUse">
                    <path class="auth-illustration__grid-line" d="M42 0H0V42" />
                </pattern>
                <radialGradient id="auth-grid-fade" cx="0" cy="0" r="1"
                                gradientTransform="translate(0.5 0.5) rotate(90) scale(0.52 0.58)">
                    <stop offset="0.2" stop-color="white" />
                    <stop offset="0.72" stop-color="white" stop-opacity="0.62" />
                    <stop offset="1" stop-color="white" stop-opacity="0" />
                </radialGradient>
                <mask id="auth-grid-mask" maskContentUnits="objectBoundingBox">
                    <rect width="1" height="1" fill="url(#auth-grid-fade)" />
                </mask>
            </defs>

            <g v-if="variant !== 'signup'">
                <circle class="auth-illustration__halo" cx="360" cy="351" r="238" />
                <circle class="auth-illustration__orbit" cx="360" cy="351" r="286" />
                <path class="auth-illustration__orbit auth-illustration__orbit--inner" d="M116 386C153 194 314 88 491 151C601 191 654 291 628 408" />
                <circle class="auth-illustration__dot auth-illustration__dot--one" cx="109" cy="301" r="7" />
                <circle class="auth-illustration__dot auth-illustration__dot--two" cx="618" cy="225" r="11" />
                <rect class="auth-illustration__grid" x="82" y="112" width="556" height="496" rx="48"
                      fill="url(#auth-grid-pattern)" mask="url(#auth-grid-mask)" />
                <circle class="auth-illustration__grid-node" cx="150" cy="574" r="5" />
                <circle class="auth-illustration__grid-node" cx="612" cy="448" r="4" />
            </g>

            <g v-else>
                <ellipse class="auth-illustration__halo" cx="280" cy="373" rx="206" ry="270" />
                <ellipse class="auth-illustration__orbit" cx="280" cy="373" rx="240" ry="310" />
                <path class="auth-illustration__orbit auth-illustration__orbit--inner" d="M73 425C92 220 218 92 369 129C461 151 508 243 496 364" />
                <circle class="auth-illustration__dot auth-illustration__dot--one" cx="72" cy="278" r="7" />
                <circle class="auth-illustration__dot auth-illustration__dot--two" cx="486" cy="203" r="10" />
                <rect class="auth-illustration__grid" x="42" y="84" width="476" height="596" rx="48"
                      fill="url(#auth-grid-pattern)" mask="url(#auth-grid-mask)" />
                <circle class="auth-illustration__grid-node" cx="92" cy="628" r="5" />
                <circle class="auth-illustration__grid-node" cx="490" cy="454" r="4" />
            </g>

            <g v-if="variant === 'login'" class="auth-illustration__scene">
                <rect class="auth-illustration__shadow" x="151" y="175" width="428" height="344" rx="38" />
                <rect class="auth-illustration__surface" x="140" y="160" width="428" height="344" rx="38" />
                <path class="auth-illustration__surface-muted" d="M140 198C140 177 157 160 178 160H530C551 160 568 177 568 198V224H140V198Z" />
                <circle class="auth-illustration__accent" cx="180" cy="192" r="8" />
                <circle class="auth-illustration__soft" cx="206" cy="192" r="8" />
                <circle class="auth-illustration__line" cx="232" cy="192" r="8" />
                <rect class="auth-illustration__soft" x="173" y="258" width="148" height="94" rx="22" />
                <path class="auth-illustration__accent-stroke" d="M196 322L222 294L246 309L277 273L299 289" />
                <rect class="auth-illustration__soft" x="341" y="258" width="194" height="94" rx="22" />
                <rect class="auth-illustration__accent" x="369" y="305" width="22" height="24" rx="6" />
                <rect class="auth-illustration__accent-faded" x="404" y="286" width="22" height="43" rx="6" />
                <rect class="auth-illustration__accent" x="439" y="271" width="22" height="58" rx="6" />
                <rect class="auth-illustration__accent-faded" x="474" y="298" width="22" height="31" rx="6" />
                <rect class="auth-illustration__line" x="173" y="383" width="225" height="13" rx="6.5" />
                <rect class="auth-illustration__line auth-illustration__line--short" x="173" y="414" width="172" height="13" rx="6.5" />
                <rect class="auth-illustration__accent" x="439" y="384" width="96" height="54" rx="17" />
                <path class="auth-illustration__on-accent" d="M465 411H509M495 397L509 411L495 425" />
                <circle class="auth-illustration__coin" cx="565" cy="489" r="58" />
                <path class="auth-illustration__coin-mark" d="M565 458V520M583 470C578 465 572 462 564 462C552 462 545 468 545 477C545 487 553 491 565 494C578 497 585 502 585 511C585 521 577 527 565 527C555 527 547 523 541 516" />
            </g>

            <g v-else-if="variant === 'signup'" class="auth-illustration__scene">
                <rect class="auth-illustration__shadow" x="126" y="101" width="330" height="532" rx="42" />
                <rect class="auth-illustration__surface" x="114" y="86" width="330" height="532" rx="42" />

                <circle class="auth-illustration__soft" cx="280" cy="178" r="48" />
                <g clip-path="url(#auth-signup-avatar-clip)">
                    <circle class="auth-illustration__accent" cx="280" cy="165" r="15" />
                    <path class="auth-illustration__accent-faded" d="M244 215C248 186 260 176 280 176C300 176 312 186 316 215H244Z" />
                </g>
                <rect class="auth-illustration__line" x="190" y="247" width="180" height="12" rx="6" />
                <rect class="auth-illustration__line auth-illustration__line--short" x="220" y="276" width="120" height="12" rx="6" />

                <path class="auth-illustration__step-line" d="M165 337V530" />
                <circle class="auth-illustration__accent" cx="165" cy="337" r="13" />
                <circle class="auth-illustration__soft auth-illustration__step-dot" cx="165" cy="432" r="13" />
                <circle class="auth-illustration__soft auth-illustration__step-dot" cx="165" cy="527" r="13" />

                <rect class="auth-illustration__soft" x="196" y="303" width="190" height="68" rx="18" />
                <circle class="auth-illustration__accent-faded" cx="225" cy="337" r="12" />
                <rect class="auth-illustration__line" x="253" y="326" width="95" height="10" rx="5" />
                <rect class="auth-illustration__line auth-illustration__line--short" x="253" y="347" width="68" height="8" rx="4" />

                <rect class="auth-illustration__soft" x="196" y="398" width="190" height="68" rx="18" />
                <rect class="auth-illustration__accent-faded" x="215" y="418" width="16" height="16" rx="5" />
                <rect class="auth-illustration__accent" x="239" y="418" width="16" height="16" rx="5" />
                <rect class="auth-illustration__accent-faded" x="215" y="442" width="16" height="8" rx="4" />
                <rect class="auth-illustration__line" x="276" y="426" width="74" height="10" rx="5" />

                <rect class="auth-illustration__soft" x="196" y="493" width="190" height="68" rx="18" />
                <circle class="auth-illustration__accent" cx="226" cy="527" r="15" />
                <path class="auth-illustration__on-accent auth-illustration__on-accent--small" d="M218 527L224 533L234 520" />
                <rect class="auth-illustration__line" x="258" y="522" width="91" height="10" rx="5" />

                <circle class="auth-illustration__coin" cx="406" cy="582" r="49" />
                <path class="auth-illustration__coin-mark" d="M384 582L401 599L430 565" />
                <path class="auth-illustration__accent-stroke" d="M73 210H119M96 187V233" />
            </g>

            <g v-else-if="variant === 'unlock'" class="auth-illustration__scene">
                <rect class="auth-illustration__shadow" x="160" y="209" width="400" height="302" rx="42" />
                <rect class="auth-illustration__surface" x="148" y="194" width="400" height="302" rx="42" />
                <path class="auth-illustration__line" d="M200 255H325" />
                <path class="auth-illustration__line" d="M200 291H288" />
                <rect class="auth-illustration__soft" x="197" y="348" width="58" height="68" rx="18" />
                <rect class="auth-illustration__soft" x="273" y="348" width="58" height="68" rx="18" />
                <rect class="auth-illustration__soft" x="349" y="348" width="58" height="68" rx="18" />
                <rect class="auth-illustration__soft" x="425" y="348" width="58" height="68" rx="18" />
                <circle class="auth-illustration__accent" cx="226" cy="382" r="7" />
                <circle class="auth-illustration__accent" cx="302" cy="382" r="7" />
                <circle class="auth-illustration__accent" cx="378" cy="382" r="7" />
                <circle class="auth-illustration__accent" cx="454" cy="382" r="7" />
                <path class="auth-illustration__shield" d="M493 159L592 197V277C592 346 548 387 493 408C438 387 394 346 394 277V197L493 159Z" />
                <path class="auth-illustration__on-accent" d="M493 228C477 228 464 241 464 257C464 267 469 276 478 281L472 320H514L508 281C517 276 522 267 522 257C522 241 509 228 493 228Z" />
            </g>

            <g v-else-if="variant === 'recovery'" class="auth-illustration__scene">
                <rect class="auth-illustration__shadow" x="129" y="219" width="462" height="292" rx="40" />
                <rect class="auth-illustration__surface" x="117" y="204" width="462" height="292" rx="40" />
                <path class="auth-illustration__soft" d="M147 252L348 402L549 252V444C549 456 539 466 527 466H169C157 466 147 456 147 444V252Z" />
                <path class="auth-illustration__accent-faded" d="M153 240L328 371C340 380 356 380 368 371L543 240" />
                <circle class="auth-illustration__coin" cx="526" cy="211" r="69" />
                <path class="auth-illustration__coin-mark" d="M550 196C550 181 539 171 524 171C509 171 498 181 498 196C498 208 505 216 516 220V249H535V220C544 216 550 208 550 196Z" />
                <path class="auth-illustration__accent-stroke" d="M172 166H232M202 136V196" />
                <circle class="auth-illustration__dot" cx="184" cy="532" r="13" />
            </g>

            <g v-else class="auth-illustration__scene">
                <rect class="auth-illustration__shadow" x="153" y="181" width="414" height="350" rx="42" />
                <rect class="auth-illustration__surface" x="141" y="166" width="414" height="350" rx="42" />
                <circle class="auth-illustration__soft" cx="348" cy="268" r="70" />
                <path class="auth-illustration__accent-stroke auth-illustration__accent-stroke--wide" d="M321 281L300 302C282 320 253 320 235 302C217 284 217 255 235 237L276 196C294 178 323 178 341 196" />
                <path class="auth-illustration__accent-stroke auth-illustration__accent-stroke--wide" d="M375 255L396 234C414 216 443 216 461 234C479 252 479 281 461 299L420 340C402 358 373 358 355 340" />
                <path class="auth-illustration__line" d="M207 402H489" />
                <path class="auth-illustration__line" d="M248 443H448" />
                <circle class="auth-illustration__coin" cx="522" cy="476" r="58" />
                <path class="auth-illustration__coin-mark" d="M495 476L515 496L550 458" />
            </g>
        </svg>
    </div>
</template>

<script setup lang="ts">
type AuthIllustrationVariant = 'login' | 'signup' | 'unlock' | 'recovery' | 'connect';

defineProps<{
    variant: AuthIllustrationVariant;
}>();
</script>
