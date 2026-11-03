<template>
    <div class="auth-illustration" :class="`auth-illustration--${variant}`" aria-hidden="true">
        <svg class="auth-illustration__art" :viewBox="variant === 'signup' ? '0 0 560 760' : '0 0 720 720'"
             fill="none" xmlns="http://www.w3.org/2000/svg" focusable="false">
            <g v-if="variant === 'login'" class="auth-illustration__scene">
                <rect class="auth-illustration__backdrop" x="110" y="159" width="500" height="384" rx="138" transform="rotate(-8 360 352)" />

                <g class="auth-illustration__shadow" transform="rotate(-4 360 353)">
                    <g transform="rotate(7 506 224)">
                        <path class="auth-illustration__paper" d="M479 170H531C536 170 540 174 540 179V271L532 266L524 271L516 266L508 271L500 266L492 271L484 266L476 271V179C476 174 477 170 479 170Z" />
                        <rect class="auth-illustration__accent" x="490" y="186" width="24" height="4" rx="2" />
                        <rect class="auth-illustration__rule" x="490" y="199" width="36" height="3" rx="1.5" />
                    </g>

                    <!-- Lower leaves are exact translations of the top-page silhouette. -->
                    <path class="auth-illustration__tint" transform="translate(0 7)" d="M176 216H350C354 216 357 218 360 221C363 218 366 216 370 216H544C550 216 554 220 554 226V474C554 480 550 484 544 484H370C366 484 363 486 360 489C357 486 354 484 350 484H176C170 484 166 480 166 474V226C166 220 170 216 176 216Z" />
                    <path class="auth-illustration__paper" transform="translate(0 3)" d="M176 216H350C354 216 357 218 360 221C363 218 366 216 370 216H544C550 216 554 220 554 226V474C554 480 550 484 544 484H370C366 484 363 486 360 489C357 486 354 484 350 484H176C170 484 166 480 166 474V226C166 220 170 216 176 216Z" />
                    <path class="auth-illustration__paper" d="M176 216H350C354 216 357 218 360 221C363 218 366 216 370 216H544C550 216 554 220 554 226V474C554 480 550 484 544 484H370C366 484 363 486 360 489C357 486 354 484 350 484H176C170 484 166 480 166 474V226C166 220 170 216 176 216Z" />

                    <!-- Curvature is confined to the unprinted 20px gutter. -->
                    <path class="auth-illustration__soft" d="M350 216C354 216 357 218 360 221C363 218 366 216 370 216V484C366 484 363 486 360 489C357 486 354 484 350 484V216Z" />
                    <path class="auth-illustration__fine" d="M360 222V488" stroke-width="1" />

                    <rect class="auth-illustration__soft" x="191" y="243" width="30" height="33" rx="8" />
                    <path class="auth-illustration__accent-outline" d="M200 251H212V269L209 266L206 269L203 266L200 269V251ZM204 256H208M204 261H208" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
                    <rect class="auth-illustration__ink" x="234" y="251" width="64" height="6" rx="3" />
                    <rect class="auth-illustration__rule" x="234" y="266" width="42" height="4" rx="2" />

                    <circle class="auth-illustration__soft" cx="204" cy="314" r="12" />
                    <path class="auth-illustration__accent-outline" d="M199 311H208V316C208 320 199 320 199 316V311ZM208 312H210C213 312 213 316 209 316M198 322H210" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
                    <rect class="auth-illustration__muted" x="230" y="306" width="49" height="5" rx="2.5" />
                    <rect class="auth-illustration__rule" x="230" y="319" width="33" height="4" rx="2" />
                    <rect class="auth-illustration__muted" x="305" y="309" width="21" height="5" rx="2.5" />
                    <path class="auth-illustration__fine" d="M192 338H326" stroke-width="1" />

                    <circle class="auth-illustration__soft" cx="204" cy="365" r="12" />
                    <path class="auth-illustration__accent-outline" d="M198 363H210L208 370H200L198 363ZM200 363L204 359L208 363" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
                    <rect class="auth-illustration__muted" x="230" y="357" width="43" height="5" rx="2.5" />
                    <rect class="auth-illustration__rule" x="230" y="370" width="29" height="4" rx="2" />
                    <rect class="auth-illustration__muted" x="302" y="360" width="24" height="5" rx="2.5" />
                    <path class="auth-illustration__fine" d="M192 389H326" stroke-width="1" />

                    <circle class="auth-illustration__soft" cx="204" cy="416" r="12" />
                    <path class="auth-illustration__accent-outline" d="M199 420V414L204 410L209 414V420H205V416H203V420H199Z" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
                    <rect class="auth-illustration__muted" x="230" y="408" width="55" height="5" rx="2.5" />
                    <rect class="auth-illustration__rule" x="230" y="421" width="37" height="4" rx="2" />
                    <rect class="auth-illustration__muted" x="307" y="411" width="19" height="5" rx="2.5" />

                    <!-- Right-page baselines are parallel to the page's top and bottom edges. -->
                    <path class="auth-illustration__accent" d="M504 216H526V279L515 271L504 279V216Z" />
                    <rect class="auth-illustration__ink" x="391" y="250" width="65" height="6" rx="3" />
                    <rect class="auth-illustration__rule" x="391" y="265" width="43" height="4" rx="2" />
                    <path class="auth-illustration__soft" d="M391 375L415 354L438 361L462 332L486 342L522 306V391H391V375Z" />
                    <path class="auth-illustration__accent-outline" d="M391 375L415 354L438 361L462 332L486 342L522 306" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" />
                    <path class="auth-illustration__fine" d="M391 391H523" stroke-width="1" />
                    <circle class="auth-illustration__accent" cx="522" cy="306" r="4" />
                    <rect class="auth-illustration__rule" x="391" y="419" width="43" height="4" rx="2" />
                    <rect class="auth-illustration__accent" x="470" y="415" width="53" height="10" rx="5" />
                    <rect class="auth-illustration__rule" x="391" y="443" width="83" height="4" rx="2" />
                </g>
            </g>

            <g v-else-if="variant === 'signup'" class="auth-illustration__scene">
                <rect class="auth-illustration__backdrop" x="84" y="137" width="392" height="491" rx="132" transform="rotate(-6 280 382)" />
                <g class="auth-illustration__shadow" transform="rotate(-4 280 375)">
                    <rect class="auth-illustration__tint" x="145" y="161" width="282" height="442" rx="24" />
                    <rect class="auth-illustration__paper" x="142" y="155" width="282" height="442" rx="24" />
                    <rect class="auth-illustration__paper" x="138" y="148" width="282" height="442" rx="24" />
                    <path class="auth-illustration__soft" d="M162 148H177V590H162C148.7 590 138 579.3 138 566V172C138 158.7 148.7 148 162 148Z" />
                    <path class="auth-illustration__fine" d="M177 166V572" stroke-width="1" />
                    <path class="auth-illustration__accent" d="M352 148H377V226L364.5 217L352 226V148Z" />

                    <!-- Center the cover artwork within the white area, excluding the spine. -->
                    <g transform="translate(18.5 0)">
                        <rect class="auth-illustration__soft" x="238" y="258" width="84" height="92" rx="23" />
                        <path class="auth-illustration__accent-outline" d="M264 280H296V322L290 317L283 322L277 317L270 322L264 317V280ZM272 292H288M272 301H285M272 310H281" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                        <rect class="auth-illustration__ink" x="222" y="378" width="116" height="8" rx="4" />
                        <rect class="auth-illustration__rule" x="245" y="400" width="70" height="5" rx="2.5" />
                        <path class="auth-illustration__fine" d="M247 519H313" stroke-width="1.5" stroke-linecap="round" />
                    </g>
                </g>
            </g>

            <g v-else-if="variant === 'verify'" class="auth-illustration__scene">
                <rect class="auth-illustration__backdrop" x="126" y="180" width="468" height="358" rx="130" transform="rotate(-7 360 359)" />
                <g class="auth-illustration__shadow" transform="translate(360 354) scale(1.1) rotate(-5) translate(-360 -354)">
                    <path class="auth-illustration__paper" d="M195 241H525C536 241 545 250 545 261V324C534 324 526 333 526 344C526 355 534 364 545 364V449C545 460 536 469 525 469H195C184 469 175 460 175 449V364C186 364 194 355 194 344C194 333 186 324 175 324V261C175 250 184 241 195 241Z" />
                    <rect class="auth-illustration__soft" x="209" y="275" width="44" height="42" rx="12" />
                    <path class="auth-illustration__accent-outline" d="M219 288H243V304H219V288ZM219 289L231 297L243 289" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
                    <rect class="auth-illustration__ink" x="269" y="283" width="114" height="6" rx="3" />
                    <rect class="auth-illustration__rule" x="269" y="301" width="76" height="4" rx="2" />
                    <path class="auth-illustration__fine" d="M209 337H395" stroke-width="1" />
                    <g class="auth-illustration__accent-outline" transform="translate(211 353)" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="12" cy="12" r="4" />
                        <path d="M16 8V13A3 3 0 0 0 22 13V12A10 10 0 1 0 18 20" />
                    </g>
                    <rect class="auth-illustration__muted" x="248" y="359" width="133" height="5" rx="2.5" />
                    <rect class="auth-illustration__rule" x="248" y="375" width="91" height="4" rx="2" />
                    <path class="auth-illustration__fine" d="M209 416H363" stroke-width="1" />
                    <rect class="auth-illustration__rule" x="209" y="428" width="67" height="3" rx="1.5" />
                    <path class="auth-illustration__fine" d="M427 259V451" stroke-width="1.2" stroke-dasharray="2 6" stroke-linecap="round" />
                    <circle class="auth-illustration__soft" cx="478" cy="313" r="31" />
                    <circle class="auth-illustration__accent-outline" cx="478" cy="313" r="24" stroke-width="1" opacity="0.4" />
                    <g v-if="verified" class="auth-illustration__verification-mark auth-illustration__accent-outline" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M468 313L475 320L488 306" />
                    </g>
                    <g v-else class="auth-illustration__accent-outline" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
                        <rect x="465" y="304" width="26" height="18" rx="2" />
                        <path d="M466 305L478 314L490 305" />
                    </g>
                    <rect class="auth-illustration__rule" x="460" y="372" width="36" height="4" rx="2" />
                    <path class="auth-illustration__accent-outline" d="M459 414H497" stroke-width="1.3" opacity="0.4" stroke-linecap="round" />
                </g>
            </g>

            <g v-else-if="variant === 'unlock'" class="auth-illustration__scene">
                <rect class="auth-illustration__backdrop" x="142" y="162" width="436" height="398" rx="130" transform="rotate(-7 360 361)" />
                <g class="auth-illustration__shadow" transform="rotate(-4 360 360)">
                    <rect class="auth-illustration__tint" x="228" y="209" width="244" height="316" rx="23" />
                    <rect class="auth-illustration__paper" x="224" y="202" width="244" height="316" rx="23" />
                    <rect class="auth-illustration__paper" x="220" y="195" width="244" height="316" rx="23" />
                    <path class="auth-illustration__soft" d="M243 195H256V511H243C230.3 511 220 500.7 220 488V218C220 205.3 230.3 195 243 195Z" />
                    <path class="auth-illustration__fine" d="M256 211V495" stroke-width="1" />
                    <path class="auth-illustration__accent" d="M405 195H426V254L415.5 246L405 254V195Z" />
                    <rect class="auth-illustration__soft" x="330" y="275" width="60" height="65" rx="17" />
                    <path class="auth-illustration__accent-outline" d="M349 291H372V322L366 318L360 322L355 318L349 322V291ZM355 300H366M355 307H364" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
                    <rect class="auth-illustration__ink" x="320" y="366" width="80" height="6" rx="3" />
                    <rect class="auth-illustration__rule" x="333" y="383" width="54" height="4" rx="2" />
                    <!-- The shackle is visibly detached from its left socket, while its right leg stays seated. -->
                    <g transform="translate(14 0)">
                        <path class="auth-illustration__paper" d="M423 357V355C423 336 437 321 455 321C473 321 487 336 487 355V379H474V355C474 343 466 334 455 334C444 334 436 343 436 355V357H423Z" />
                        <path class="auth-illustration__accent-outline" d="M429.5 357V355C429.5 339 440 327.5 455 327.5C470 327.5 480.5 339 480.5 355V386" stroke-width="2" stroke-linecap="round" />
                        <rect class="auth-illustration__soft" x="398" y="395" width="93" height="10" rx="5" />
                        <rect class="auth-illustration__paper" x="402" y="377" width="98" height="91" rx="19" />
                        <rect class="auth-illustration__soft" x="413" y="388" width="76" height="69" rx="12" />
                        <path class="auth-illustration__accent-outline" d="M451 408C445 408 440 412 440 418C440 423 443 426 447 428L445 440H457L455 428C459 426 462 423 462 418C462 412 457 408 451 408Z" stroke-width="1.7" stroke-linejoin="round" />
                    </g>
                </g>
            </g>

            <g v-else-if="variant === 'recovery'" class="auth-illustration__scene">
                <rect class="auth-illustration__backdrop" x="119" y="183" width="482" height="354" rx="130" transform="rotate(-7 360 360)" />
                <g class="auth-illustration__shadow" transform="rotate(-3 360 360)">
                    <rect class="auth-illustration__paper" x="160" y="252" width="400" height="240" rx="23" />
                    <path class="auth-illustration__soft" d="M174 268L348 392C355 398 365 398 372 392L546 268V468C546 474 541 479 535 479H185C179 479 174 474 174 468V268Z" />
                    <path class="auth-illustration__fine" d="M174 475L288 372M546 475L432 372" stroke-width="1.1" stroke-linecap="round" />
                    <path class="auth-illustration__accent-outline" d="M174 268L348 392C355 398 365 398 372 392L546 268" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" opacity="0.45" />
                    <circle class="auth-illustration__paper" cx="535" cy="251" r="43" />
                    <path class="auth-illustration__accent-outline" d="M535 230C525 230 517 238 517 248C517 256 521 262 528 265V282H542V265C549 262 553 256 553 248C553 238 545 230 535 230Z" transform="translate(0 -5)" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" />
                </g>
            </g>

            <g v-else class="auth-illustration__scene">
                <rect class="auth-illustration__backdrop" x="126" y="183" width="468" height="354" rx="128" transform="rotate(-7 360 360)" />
                <g class="auth-illustration__shadow" transform="translate(360 366) scale(1.1) rotate(-5) translate(-360 -352)">
                    <path class="auth-illustration__paper" d="M195 241H521C532 241 541 250 541 261V289C529 289 521 298 521 310C521 322 529 331 541 331V432C541 443 532 452 521 452H195C184 452 175 443 175 432V331C187 331 195 322 195 310C195 298 187 289 175 289V261C175 250 184 241 195 241Z" />
                    <rect class="auth-illustration__soft" x="210" y="273" width="43" height="46" rx="12" />
                    <path class="auth-illustration__accent-outline" d="M223 286H240V307L236 304L232 307L227 304L223 307V286ZM227 292H236M227 297H234" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
                    <rect class="auth-illustration__ink" x="269" y="283" width="113" height="7" rx="3.5" />
                    <rect class="auth-illustration__rule" x="269" y="302" width="76" height="4" rx="2" />
                    <rect class="auth-illustration__muted" x="210" y="352" width="132" height="5" rx="2.5" />
                    <rect class="auth-illustration__rule" x="210" y="369" width="95" height="4" rx="2" />
                    <g class="auth-illustration__muted" opacity="0.6">
                        <rect x="211" y="405" width="3" height="16" rx="1" />
                        <rect x="219" y="405" width="2" height="16" rx="1" />
                        <rect x="226" y="405" width="4" height="16" rx="1" />
                        <rect x="235" y="405" width="2" height="16" rx="1" />
                        <rect x="242" y="405" width="3" height="16" rx="1" />
                        <rect x="250" y="405" width="2" height="16" rx="1" />
                        <rect x="257" y="405" width="4" height="16" rx="1" />
                        <rect x="266" y="405" width="2" height="16" rx="1" />
                    </g>
                    <g transform="translate(-16 0)">
                        <path class="auth-illustration__fine" d="M447 259V434" stroke-width="1.2" stroke-dasharray="2 6" stroke-linecap="round" />
                        <rect class="auth-illustration__soft" x="466" y="281" width="52" height="52" rx="15" />
                        <g class="auth-illustration__accent-outline" transform="translate(478 293) scale(1.166667)" stroke-width="1.65" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M10 13A5 5 0 0 0 17.54 13.54L20.54 10.54A5 5 0 0 0 13.47 3.47L11.75 5.18" />
                            <path d="M14 11A5 5 0 0 0 6.46 10.46L3.46 13.46A5 5 0 0 0 10.53 20.53L12.24 18.82" />
                        </g>
                        <rect class="auth-illustration__rule" x="474" y="367" width="35" height="4" rx="2" />
                        <rect class="auth-illustration__accent" x="482" y="404" width="19" height="5" rx="2.5" opacity="0.55" />
                    </g>
                </g>
            </g>
        </svg>
    </div>
</template>

<script setup lang="ts">
type AuthIllustrationVariant = 'login' | 'signup' | 'verify' | 'unlock' | 'recovery' | 'connect';

defineProps<{
    variant: AuthIllustrationVariant;
    verified?: boolean;
}>();
</script>
