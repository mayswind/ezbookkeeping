import fs from 'fs';
import { resolve } from 'path';
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue';
import vuetify from 'vite-plugin-vuetify';
import { VitePWA } from 'vite-plugin-pwa';
import git from 'git-rev-sync';

import packageFile from './package.json';
import thirdPartyLicenseFile from './third-patry-dependencies.json';

const SRC_DIR = resolve(__dirname, './src');
const PUBLIC_DIR = resolve(__dirname, './public');
const BUILD_DIR = resolve(__dirname, './dist',);

export default defineConfig(async () => {
    const licenseContent = fs.readFileSync('./LICENSE', 'UTF-8');
    let buildUnixTime = process.env.buildUnixTime || '';

    return {
        root: SRC_DIR,
        publicDir: PUBLIC_DIR,
        base: './',
        define: {
            __EZBOOKKEEPING_VERSION__: JSON.stringify(packageFile.version),
            __EZBOOKKEEPING_BUILD_UNIX_TIME__: JSON.stringify(buildUnixTime),
            __EZBOOKKEEPING_BUILD_COMMIT_HASH__: JSON.stringify(git.short()),
            __EZBOOKKEEPING_LICENSE__: JSON.stringify(licenseContent),
            __EZBOOKKEEPING_THIRD_PARTY_LICENSES__: JSON.stringify(thirdPartyLicenseFile)
        },
        plugins: [
            vue({
                template: {
                    compilerOptions: {
                        isCustomElement: (tag) => tag.includes('swiper-')
                    }
                }
            }),
            vuetify({
                styles: {
                    configFile: 'styles/desktop/configured-variables/_vuetify.scss'
                }
            }),
            VitePWA({
                filename: 'sw.js',
                manifestFilename: 'manifest.json',
                strategies: 'generateSW',
                injectRegister: 'null',
                manifest: {
                    name: 'ezBookkeeping',
                    short_name: 'ezBookkeeping',
                    description: 'A lightweight personal bookkeeping app hosted by yourself.',
                    theme_color: '#C67E48',
                    background_color: '#F6F7F8',
                    start_url: './',
                    scope: './',
                    display: 'standalone',
                    related_applications: [],
                    prefer_related_applications: false,
                    icons: [
                        {
                            src: 'img/ezbookkeeping-192.png',
                            sizes: '192x192',
                            type: 'image/png'
                        },
                        {
                            src: 'img/ezbookkeeping-512.png',
                            sizes: '512x512',
                            type: 'image/png'
                        }
                    ]
                },
                workbox: {
                    globDirectory: 'dist/',
                    globPatterns: ['**/*.{js,css,html,ico,png,jpg,jpeg,gif,tiff,bmp,ttf,woff,woff2,svg,eot}'],
                    globIgnores: [
                        'index.html',
                        'mobile.html',
                        'desktop.html',
                        'robots.txt',
                        'img/splash_screens/*',
                        'img/desktop/*',
                        'fonts/*.eot',
                        'fonts/*.ttf',
                        'fonts/*.svg',
                        'fonts/*.woff',
                        'css/vendor-desktop-*.css',
                        'css/desktop-*.css',
                        'js/vendor-desktop-*.js',
                        'js/desktop-*.js'
                    ],
                    runtimeCaching: [
                        {
                            urlPattern: /.*\/(mobile\/|desktop\/|index\.html|mobile\.html|desktop\.html)/,
                            handler: 'NetworkFirst'
                        },
                        {
                            urlPattern: /.*\/img\/(splash_screens|desktop)\/.*\.(png|jpg|jpeg|gif|tiff|bmp|svg)/,
                            handler: 'StaleWhileRevalidate'
                        },
                        {
                            urlPattern: /.*\/fonts\/.*\.(eot|ttf|svg|woff)/,
                            handler: 'CacheFirst'
                        },
                        {
                            urlPattern: /.*\/css\/(vendor-desktop-\.*|desktop-\.*)\.css/,
                            handler: 'CacheFirst'
                        },
                        {
                            urlPattern: /.*\/js\/(vendor-desktop-\.*|desktop-\.*)\.js/,
                            handler: 'CacheFirst'
                        }
                    ],
                    navigateFallback: '',
                    skipWaiting: true,
                    clientsClaim: true
                }
            })
        ],
        build: {
            outDir: BUILD_DIR,
            sourcemap: false,
            assetsInlineLimit: 0,
            emptyOutDir: true,
            rollupOptions: {
                input: {
                    index: resolve(SRC_DIR, 'index.html'),
                    desktop: resolve(SRC_DIR, 'desktop.html'),
                    mobile: resolve(SRC_DIR, 'mobile.html')
                },
                output: {
                    assetFileNames: (assetInfo) => {
                        const fileExt = assetInfo.name.split('.').at(1);
                        let assetType = fileExt;

                        if (/png|jpe?g|gif|tiff|bmp|ico/i.test(fileExt)) {
                            assetType = 'img';
                        } else if (/ttf|woff|woff2|svg|eot/i.test(fileExt)) {
                            assetType = 'fonts';
                        }

                        return `${assetType}/[name]-[hash][extname]`;
                    },
                    chunkFileNames: 'js/[name]-[hash].js',
                    entryFileNames: 'js/[name]-[hash].js',
                    manualChunks: function (id) {
                        if (/[\\/]node_modules[\\/]leaflet[\\/]/i.test(id)) {
                            return 'leaflet';
                        } else if (/[\\/]node_modules[\\/](moment|moment-timezone)[\\/]/i.test(id)) {
                            return 'moment';
                        } else if (/[\\/]node_modules[\\/](dom7|framework7.*|skeleton-elements|swiper)[\\/]/i.test(id)) {
                            return 'vendor-mobile';
                        } else if (/[\\/]node_modules[\\/](vuetify|vue-router|vue3-perfect-scrollbar|perfect-scrollbar|vuedraggable|sortablejs|@mdi.*)[\\/]/i.test(id)) {
                            return 'vendor-desktop';
                        } else if (/[\\/]node_modules[\\/](echarts|zrender|tslib|resize-detector|vue-echarts)[\\/]/i.test(id)) {
                            return 'vendor-desktop';
                        } else if (/plugin-vuetify:/i.test(id)) {
                            return 'vendor-desktop';
                        } else if (/[\\/]node_modules[\\/]/i.test(id)) {
                            return 'vendor-common';
                        } else if (/[\\/]src[\\/]locales[\\/]/i.test(id)) {
                            return 'locales';
                        } else if (/[\\/]src[\\/](consts|stores)[\\/]/i.test(id)) {
                            return 'common';
                        } else if (/[\\/]src[\\/]lib[\\/](map[\\/]|[a-zA-Z0-9-_]+\.js)/i.test(id)) {
                            return 'common';
                        } else if (/[\\/]src[\\/]components[\\/]common[\\/]/i.test(id)) {
                            return 'common';
                        }
                    }
                },
                treeshake: false
            },
        },
        resolve: {
            alias: {
                '@': SRC_DIR,
            },
        },
        server: {
            host: '0.0.0.0',
            port: 8081,
            strictPort: true,
            proxy: {
                '/dev': {
                    target: 'http://127.0.0.1:8080/',
                    changeOrigin: true
                },
                '/api': {
                    target: 'http://127.0.0.1:8080/',
                    changeOrigin: true
                },
                '/qrcode': {
                    target: 'http://127.0.0.1:8080/',
                    changeOrigin: true
                },
                '/proxy': {
                    target: 'http://127.0.0.1:8080/',
                    changeOrigin: true
                },
                '/_AMapService': {
                    target: 'http://127.0.0.1:8080/',
                    changeOrigin: true
                }
            }
        },
    };
})
