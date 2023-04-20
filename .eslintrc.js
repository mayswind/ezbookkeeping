module.exports = {
    'root': true,
    'env': {
        'node': true
    },
    'extends': [
        'eslint:recommended',
        'plugin:vue/vue3-essential'
    ],
    'rules': {
        "vue/multi-word-component-names": "off",
        "vue/no-use-v-if-with-v-for": "off"
    }
}
