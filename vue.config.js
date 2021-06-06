const fs = require('fs');

const GitRevisionPlugin = require('git-revision-webpack-plugin');
const MomentLocalesPlugin = require('moment-locales-webpack-plugin');

const pkgFile = require('./package.json');
const thirdPartyLicenseFile = require('./third-patry-dependencies.json');

const licenseFile = fs.readFileSync('./LICENSE', 'UTF-8');

module.exports = {
    pages: {
        index: {
            entry: 'src/index-main.js',
            template: 'src/public/index.html',
            filename: 'index.html',
            chunks: ['vendors-common-bundle', 'vendors-index-bundle', 'common-bundle', 'index']
        },
        desktop: {
            entry: 'src/desktop-main.js',
            template: 'src/public/desktop.html',
            filename: 'desktop.html',
            chunks: ['vendors-common-bundle', 'vendors-desktop-bundle', 'common-bundle', 'desktop']
        },
        mobile: {
            entry: 'src/mobile-main.js',
            template: 'src/public/mobile.html',
            filename: 'mobile.html',
            chunks: ['vendors-common-bundle', 'vendors-mobile-bundle', 'common-bundle', 'mobile']
        }
    },
    publicPath: '',
    productionSourceMap: false,
    configureWebpack: {
        plugins: [
            new MomentLocalesPlugin()
        ]
    },
    chainWebpack: config => {
        config.optimization.splitChunks({
            cacheGroups: {
                'vendors-common-bundle': {
                    name: 'vendors-common-bundle',
                    test: /[\\/]node_modules[\\/]/,
                    chunks: 'initial',
                    priority: 10,
                    minChunks: 2
                },
                'vendors-bundle': {
                    name: (module, chunks) => {
                        const allChunksNames = chunks.map((item) => item.name).join('-');
                        return `vendors-${allChunksNames}-bundle`;
                    },
                    test: /[\\/]node_modules[\\/]/,
                    chunks: 'initial',
                    priority: 5,
                    minChunks: 1
                },
                'common-bundle': {
                    name: 'common-bundle',
                    chunks: 'initial',
                    priority: 1,
                    minChunks: 2
                }
            }
        });

        config.plugin('define').tap(definitions => {
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

            const gitRevisionPlugin = new GitRevisionPlugin();
            definitions[0]['process.env']['VERSION'] = JSON.stringify(pkgFile.version);
            definitions[0]['process.env']['COMMIT_HASH'] = JSON.stringify(gitRevisionPlugin.commithash());
            definitions[0]['process.env']['BUILD_UNIXTIME'] = buildUnixTime;
            definitions[0]['process.env']['LICENSE'] = JSON.stringify(licenseFile.trim());
            definitions[0]['process.env']['THIRD_PARTY_LICENSES'] = JSON.stringify(thirdPartyLicenseFile);

            return definitions;
        });
    },
    pwa: {
        name: 'ezBookkeeping',
        themeColor: '#C67E48',
        appleMobileWebAppCapable: 'yes',
        appleMobileWebAppStatusBarStyle: 'default',
        workboxPluginMode: 'GenerateSW',
        manifestPath: 'manifest.json',
        manifestOptions: {
            short_name: 'ezBookkeeping',
            icons: [
                {
                    src: "img/ezbookkeeping-192.png",
                    sizes: "192x192",
                    type: "image/png"
                },
                {
                    src: "img/ezbookkeeping-512.png",
                    sizes: "512x512",
                    type: "image/png"
                }
            ],
            start_url: './',
            scope: "./",
            display: 'standalone',
            background_color: "#F6F7F8",
            related_applications: [],
            prefer_related_applications: false
        },
        iconPaths: {
            favicon32: 'favicon.png',
            favicon16: 'favicon.ico',
            appleTouchIcon: 'touchicon.png',
        },
        workboxOptions: {
            importWorkboxFrom: 'local',
            importsDirectory: 'sw',
            precacheManifestFilename: 'precache-manifest.[manifestHash].js',
            directoryIndex: '/',
            globDirectory: '.',
            templatedURLs: {
                '/': [
                    'src/public/mobile.html'
                ]
            },
            exclude: [
                /^index\.html$/,
                /^mobile\.html$/,
                /^desktop\.html$/,
                /^robots\.txt$/,
                /^(css|js)\/desktop\.[a-z0-9]+\.js$/,
                /^(css|js)\/vendors-desktop-bundle\.[a-z0-9]+\.js$/,
            ],
            skipWaiting: true,
            clientsClaim: true,
            swDest: 'sw.js'
        }
    },
    devServer: {
        host: '0.0.0.0',
        port: 8081,
        disableHostCheck: true,
        proxy: {
            '/api': {
                target: 'http://127.0.0.1:8080/',
                changeOrigin: true
            }
        }
    }
}
