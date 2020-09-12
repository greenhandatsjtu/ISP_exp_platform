import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
    state: {
        login: false, // 登录状态
        user: null, // 用户信息
        roles: [], //角色
    },
    mutations: {
        login(state, user) {
            state.user = user
            user.roles.forEach((value) => {
                state.roles.push(value.description)
            })
            state.login = true
        },
        logout(state) {
            state.login = false
            state.user = null
            state.roles = []
        },
        message(state, data) {
            state.snackbar.message = data.message
            state.snackbar.status = data.status
            state.snackbar.open = true
        }
    },
    actions: {},
    modules: {}
})
