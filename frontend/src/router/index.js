import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import About from '@/views/About'
import store from '@/store/index'
import Roles from "@/components/Roles";

Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home,
        meta: {
            requiresAuth: true,
            roles: [],
        }
    },
    {
        path: '/about',
        name: 'About',
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        // component: () => import(/* webpackChunkName: "about" */ '@/views/About.vue'),
        component: About,
        meta: {
            requiresAuth: false,
        }
    },
    {
        path: '/studentCourses',
        name: 'StudentCourses',
        component: () => import('@/views/student/StudentCourses'),
        meta: {
            requiresAuth: true,
            roles: [Roles.student],
        }
    },
    {
        path: '/adminCourses',
        name: 'AdminCourses',
        component: () => import('@/views/teacherAndAssistant/AdminCourses'),
        meta: {
            requiresAuth: true,
            roles: [Roles.sysAdmin, Roles.teacher, Roles.assistant],
        }
    },
    {
        path: '/allCourses',
        name: 'AllCourses',
        component: () => import('@/views/sysAdmin/AllCourses'),
        meta: {
            requiresAuth: true,
            roles: [Roles.sysAdmin, Roles.teachAdmin],
        }
    },
    {
        path: '/course/:id',
        name: 'CourseInfo',
        component: () => import('@/views/student/CourseInfo'),
        meta: {
            requiresAuth: true,
            roles: [Roles.student],
        }
    },
    {
        path: '/adminCourse/:id',
        name: 'AdminCourseInfo',
        component: () => import('@/views/teacherAndAssistant/AdminCourseInfo'),
        meta: {
            requiresAuth: true,
            roles: [Roles.teachAdmin, Roles.assistant, Roles.teacher, Roles.sysAdmin],
        }
    },
    {
        path: '/experiment/:id',
        name: 'ExperimentInfo',
        component: () => import('@/views/student/ExperimentInfo'),
        meta: {
            requiresAuth: true,
            roles: [Roles.student],
        },
    },
    {
        path: '/adminExperiment/:id',
        name: 'AdminExperimentInfo',
        component: () => import('@/views/teacherAndAssistant/AdminExperimentInfo'),
        meta: {
            requiresAuth: true,
            roles: [Roles.teacher, Roles.assistant, Roles.sysAdmin],
        },
    },
    {
        path: '/students',
        name: 'Students',
        component: () => import('@/views/ManageStudents'),
        meta: {
            requiresAuth: true,
            roles: [Roles.teachAdmin, Roles.sysAdmin],
        }
    },
    {
        path: '/users',
        name: 'Users',
        component: () => import('@/views/sysAdmin/ManageUsers'),
        meta: {
            requiresAuth: true,
            roles: [Roles.sysAdmin],
        }
    },
    {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard'),
        meta: {
            requiresAuth: true,
            roles: [],
        }
    },
    {
        path: '/allResources',
        name: 'AllResources',
        component: () => import('@/views/sysAdmin/AllResources'),
        meta: {
            requiresAuth: true,
            roles: [Roles.sysAdmin],
        }
    },
    {
        path: '*', // 404 页面
        name: 'NotFound',
        component: () => import('@/views/NotFound'),
        meta: {
            requiresAuth: true,
            roles: [Roles.sysAdmin, Roles.teachAdmin, Roles.assistant, Roles.teacher, Roles.student],
        }
    },
]

const router = new VueRouter({
    routes
})

function auth(to) {
    if (to.meta.roles.length === 0) {
        return true
    }
    //求交集
    let set1 = new Set(to.meta.roles);
    let set2 = new Set(store.state.roles);
    let result = new Set([...set1].filter(x => set2.has(x)));
    console.log(set1, set2, result)
    return result.size !== 0;
}

router.beforeEach((to, from, next) => {
    if (to.matched.some(record => record.meta.requiresAuth)) {
        if (store.state.login) {
            if (auth(to)) {
                next();
            } else {
                next(false);
            }
        } else {
            axios.get("home").then(({data}) => {
                store.commit("login", data.data)
                next();
            }).catch(() => {
                next({name: "About"});
            })
        }
    } else {
        next();
    }
});

export default router
