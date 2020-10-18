import MainPage from './components/Main.vue'
import MainPageHomeTab from './components/main/Home.vue'

const router = [
    {
        path: '/',
        component: MainPage,
        tabs: [
            {
                path: '/',
                id: 'main-tab-home',
                component: MainPageHomeTab,
            }
        ]
    },
    {
        path: '(.*)',
        redirect: '/'
    }
]

export default router
