import fs from 'fs';
import { resolve } from 'path';
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue';
import { VitePWA } from 'vite-plugin-pwa';
import git from 'git-rev-sync';

import packageFile from './package.json';
import thirdPartyLicenseFile from './third-patry-dependencies.json';

const SRC_DIR = resolve(__dirname, './src');
const PUBLIC_DIR = resolve(__dirname, './public');
const BUILD_DIR = resolve(__dirname, './dist',);

export default defineConfig(async () => {
    const licenseContent = fs.readFileSync('./LICENSE', 'UTF-8');
    let buildUnixTime = '';

    for (let i = 0; i < process.argv.length; i++) {
        if (process.argv[i].indexOf('--') !== 0) {
            continue;
        }

        const pairs = process.argv[i].split('=');

        if (pairs[0] === '--buildUnixTime') {
            buildUnixTime = pairs[1];
        }
    }

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
                        'css/desktop-*.js',
                        'js/desktop-*.js'
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
                        if (/[\\/]node_modules[\\/]/i.test(id)) {
                            return 'vendor';
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
                '/api': {
                    target: 'http://127.0.0.1:8080/',
                    changeOrigin: true
                }
            }
        },
    };
})
