const fs = require('fs');
const path = require('path');

const FRONTEND_LOCALES_DIR = path.join(__dirname, '..', '..', 'src', 'locales');
const BACKEND_LOCALES_DIR = path.join(__dirname, '..', '..', 'pkg', 'locales');
const OUTPUT_DIR = process.argv[2] || path.join(__dirname, '..', '..', 'i18n-badge');

const DEFAULT_LANGUAGE_TAG = 'en';

const BACKEND_SKIP_STRUCTS = new Set([
    'GlobalTextItems',
    'DefaultTypes',
    'DataConverterTextItems',
]);

function discoverFrontendLanguages() {
    const indexPath = path.join(FRONTEND_LOCALES_DIR, 'index.ts');
    const content = fs.readFileSync(indexPath, 'utf-8');

    const importMap = {};
    const importRegex = /import\s+(\w+)\s+from\s+['"]\.\/([\w_]+\.json)['"]/g;
    let match;

    while ((match = importRegex.exec(content)) !== null) {
        importMap[match[1]] = match[2];
    }

    const result = {};
    const langRegex = /['"]([^'"]+)['"]\s*:\s*\{[^}]*content\s*:\s*(\w+)/g;

    while ((match = langRegex.exec(content)) !== null) {
        const tag = match[1];
        const varName = match[2];

        if (importMap[varName]) {
            result[tag] = importMap[varName];
        }
    }

    return result;
}

function discoverBackendLanguages() {
    const allLocalesPath = path.join(BACKEND_LOCALES_DIR, 'all_locales.go');
    const content = fs.readFileSync(allLocalesPath, 'utf-8');

    const result = {};
    const entryRegex = /"([^"]+)"\s*:\s*\{[^}]*Content\s*:\s*(\w+)/g;
    let match;

    while ((match = entryRegex.exec(content)) !== null) {
        const tag = match[1];
        const fileName = tag.toLowerCase().replace(/-/g, '_') + '.go';
        const filePath = path.join(BACKEND_LOCALES_DIR, fileName);

        if (fs.existsSync(filePath)) {
            result[tag] = fileName;
        }
    }

    return result;
}

function flattenJSON(obj, prefix) {
    const result = {};

    for (const key of Object.keys(obj)) {
        const fullKey = prefix ? prefix + '.' + key : key;

        if (typeof obj[key] === 'object' && obj[key] !== null) {
            Object.assign(result, flattenJSON(obj[key], fullKey));
        } else {
            result[fullKey] = obj[key];
        }
    }

    return result;
}

function shouldSkipFrontendKey(key) {
    if (key.startsWith('global.')) {
        return true;
    } else if (key.startsWith('default.')) {
        return true;
    } else if (key.startsWith('currency.')) {
        if (key.startsWith('currency.unit.')) {
            return true;
        } else {
            return false;
        }
    } else if (key.startsWith('mapprovider.')) {
        return true;
    } else if (key.startsWith('encoding.')) {
        return true;
    } else if (key.startsWith('document.')) {
        if (key.startsWith('document.anchor.')) {
            return true;
        } else {
            return false;
        }
    } else {
        return false;
    }
}

function isFrontendAlwaysTranslatedKey(key) {
    if (key.startsWith('language.')) {
        return true;
    } else if (key.startsWith('format.')) {
        if (key.startsWith('format.misc.')) {
            if (key === 'format.misc.multiTextJoinSeparator') {
                return true;
            } else if (key === 'format.misc.eachMonthDayInMonthDays') {
                return true;
            } else {
                return false;
            }
        } else {
            return true;
        }
    } else if (key.startsWith('datetime.')) {
        return true;
    } else if (key.startsWith('timezone.')) {
        return true;
    } else if (key.startsWith('currency.')) {
        if (key === 'currency.name.EUR') {
            return true;
        } else {
            return false;
        }
    } else if (key.startsWith('parameter.')) {
        if (key === 'parameter.id') {
            return true;
        } else {
            return false;
        }
    } else {
        if (key === 'OK') {
            return true;
        } else {
            return false;
        }
    }
}

function extractGoStringFields(content) {
    const fields = [];
    const structBlockRegex = /(\w+):\s*&\w+\{([^}]*)\}/gs;
    let blockMatch;

    while ((blockMatch = structBlockRegex.exec(content)) !== null) {
        const structName = blockMatch[1];
        const blockBody = blockMatch[2];
        const fieldRegex = /(\w+):\s+"((?:[^"\\]|\\.)*)"/g;
        let fieldMatch;

        while ((fieldMatch = fieldRegex.exec(blockBody)) !== null) {
            fields.push({
                struct: structName,
                name: fieldMatch[1],
                value: fieldMatch[2],
            });
        }
    }

    return fields;
}

function getProgressColor(progress) {
    if (progress >= 95) {
        return 'brightgreen';
    } else if (progress >= 90) {
        return 'green';
    } else if (progress >= 70) {
        return 'yellowgreen';
    } else if (progress >= 50) {
        return 'yellow';
    } else if (progress >= 20) {
        return 'orange';
    } else {
        return 'red';
    }
}

function main() {
    const frontendLangs = discoverFrontendLanguages();
    const backendLangs = discoverBackendLanguages();
    const allTags = new Set([...Object.keys(frontendLangs), ...Object.keys(backendLangs)]);

    console.log('Discovered ' + allTags.size + ' languages: ' + [...allTags].sort().join(', '));

    const defaultFrontendJSON = JSON.parse(fs.readFileSync(path.join(FRONTEND_LOCALES_DIR, `${DEFAULT_LANGUAGE_TAG}.json`), 'utf-8'));
    const defaultFrontendItemsMap = flattenJSON(defaultFrontendJSON, '');
    const defaultFrontendKeys = Object.keys(defaultFrontendItemsMap);
    const frontendTranslatableKeys = defaultFrontendKeys.filter(function (k) {
        return !shouldSkipFrontendKey(k);
    });
    const frontendSkippedCount = defaultFrontendKeys.length - frontendTranslatableKeys.length;
    const frontendTotal = frontendTranslatableKeys.length;

    const defaultBackendContent = fs.readFileSync(path.join(BACKEND_LOCALES_DIR, `${DEFAULT_LANGUAGE_TAG}.go`), 'utf-8');
    const defaultBackendItems = extractGoStringFields(defaultBackendContent);
    const defaultBackendTranslatableItems = defaultBackendItems.filter(function (f) {
        return !BACKEND_SKIP_STRUCTS.has(f.struct);
    });
    const backendSkippedCount = defaultBackendItems.length - defaultBackendTranslatableItems.length;
    const backendTotal = defaultBackendTranslatableItems.length;

    console.log('Frontend: ' + frontendTotal + ' translatable keys (' + frontendSkippedCount + ' excluded)');
    console.log('Backend: ' + backendTotal + ' translatable fields (' + backendSkippedCount + ' excluded)');

    const results = {};
    const untranslatedKeys = {};

    for (const tag of allTags) {
        results[tag] = {
            languageTag: tag,
            frontendTranslated: 0,
            frontendTotal: frontendTotal,
            backendTranslated: 0,
            backendTotal: backendTotal
        };
        untranslatedKeys[tag] = [];
    }

    for (const tag of Object.keys(frontendLangs)) {
        if (tag === DEFAULT_LANGUAGE_TAG) {
            results[tag].frontendTranslated = frontendTotal;
            continue;
        }

        const file = frontendLangs[tag];
        const filePath = path.join(FRONTEND_LOCALES_DIR, file);

        if (!fs.existsSync(filePath)) {
            continue;
        }

        const json = JSON.parse(fs.readFileSync(filePath, 'utf-8'));
        const kv = flattenJSON(json, '');
        let translated = 0;

        for (const key of frontendTranslatableKeys) {
            if (kv[key] !== undefined && kv[key] !== '' && (kv[key] !== defaultFrontendItemsMap[key] || isFrontendAlwaysTranslatedKey(key))) {
                translated++;
            } else {
                untranslatedKeys[tag].push({ source: path.join('src', 'locales', file), key: key, defaultValue: defaultFrontendItemsMap[key], value: kv[key] });
            }
        }

        results[tag].frontendTranslated = translated;
    }

    for (const tag of Object.keys(backendLangs)) {
        if (tag === DEFAULT_LANGUAGE_TAG) {
            results[tag].backendTranslated = backendTotal;
            continue;
        }

        const file = backendLangs[tag];
        const filePath = path.join(BACKEND_LOCALES_DIR, file);

        if (!fs.existsSync(filePath)) {
            continue;
        }

        const content = fs.readFileSync(filePath, 'utf-8');
        const fields = extractGoStringFields(content).filter(function (f) {
            return !BACKEND_SKIP_STRUCTS.has(f.struct);
        });
        let translated = 0;

        for (let i = 0; i < defaultBackendTranslatableItems.length; i++) {
            if (i < fields.length && fields[i].value !== defaultBackendTranslatableItems[i].value) {
                translated++;
            } else {
                untranslatedKeys[tag].push({ source: path.join('pkg', 'locales', file), key: defaultBackendTranslatableItems[i].struct + '.' + defaultBackendTranslatableItems[i].name, defaultValue: defaultBackendTranslatableItems[i].value, value: (i < fields.length) ? fields[i].value : null });
            }
        }

        results[tag].backendTranslated = translated;
    }

    for (const tag of Object.keys(results)) {
        const r = results[tag];
        const totalTranslated = r.frontendTranslated + r.backendTranslated;
        const totalItems = r.frontendTotal + r.backendTotal;
        r.totalProgress = Math.round((totalTranslated / totalItems) * 10000) / 100;
    }

    const sortedResults = {};
    var sortedTags = Object.keys(results).sort();

    for (const tag of sortedTags) {
        sortedResults[tag] = results[tag];
    }

    if (!fs.existsSync(OUTPUT_DIR)) {
        fs.mkdirSync(OUTPUT_DIR, { recursive: true });
    }

    var badgesDir = path.join(OUTPUT_DIR, 'badges');

    if (!fs.existsSync(badgesDir)) {
        fs.mkdirSync(badgesDir, { recursive: true });
    }

    fs.writeFileSync(
        path.join(OUTPUT_DIR, 'i18n-progress.json'),
        JSON.stringify(sortedResults, null, 4) + '\n'
    );

    for (const tag of sortedTags) {
        const data = sortedResults[tag];
        const badge = {
            schemaVersion: 1,
            label: 'translation',
            message: data.totalProgress + '%',
            color: getProgressColor(data.totalProgress)
        };

        fs.writeFileSync(
            path.join(badgesDir, tag + '.json'),
            JSON.stringify(badge, null, 4) + '\n'
        );
    }

    var untranslatedDir = path.join(OUTPUT_DIR, 'untranslated');

    if (!fs.existsSync(untranslatedDir)) {
        fs.mkdirSync(untranslatedDir, { recursive: true });
    }

    for (const tag of sortedTags) {
        const items = untranslatedKeys[tag] || [];

        fs.writeFileSync(
            path.join(untranslatedDir, tag + '.json'),
            JSON.stringify(items, null, 4) + '\n'
        );
    }

    for (const tag of sortedTags) {
        const data = sortedResults[tag];
        const missingCount = (untranslatedKeys[tag] || []).length;
        console.log(tag + ': ' + data.totalProgress + '% (frontend: ' + data.frontendTranslated + '/' + data.frontendTotal + ', backend: ' + data.backendTranslated + '/' + data.backendTotal + ', untranslated: ' + missingCount + ')');
    }

    console.log('\nResults written to ' + OUTPUT_DIR);
}

main();
