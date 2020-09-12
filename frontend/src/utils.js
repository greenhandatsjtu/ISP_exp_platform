import store from "@/store";

export default {
    install(Vue, options) {
        Vue.prototype.auth = function (roles) {
            if (roles.length === 0) {
                return true
            }
            //求交集
            let set1 = new Set(roles);
            let set2 = new Set(store.state.roles);
            let result = new Set([...set1].filter(x => set2.has(x)));
            console.log(set1, set2, result)
            return result.size !== 0;
        }
    }
}
