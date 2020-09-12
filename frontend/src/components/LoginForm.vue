<template>
  <v-bottom-sheet v-model="open" inset persistent>
    <template v-slot:activator="{ on, attrs }">
      <v-btn block class="info" v-on="on" v-bind="attrs">{{ $store.state.login ? "注销" : "登录" }}</v-btn>
    </template>
    <v-sheet class="text-center pa-4">
      <v-form v-model="valid" ref="form">
        <v-text-field v-model="user.email"
                      required
                      clearable
                      label="邮箱"
                      type="Email"
                      v-on:keydown.enter.exact="submit"
                      :rules="rules.emailRules"></v-text-field>
        <v-text-field v-model="user.password"
                      clearable
                      label="密码"
                      :type="showPassword ? 'text' : 'password'"
                      :append-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                      @click:append="showPassword = !showPassword"
                      v-on:keydown.enter.exact="submit"
                      :rules="rules.passwordRules"></v-text-field>
        <v-row class="mx-auto">
          <v-btn
              color="error"
              @click="closeForm"
          >取消
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn
              :disabled="!valid"
              color="info"
              @click="submit"
          >确认
          </v-btn>
        </v-row>
      </v-form>
    </v-sheet>
  </v-bottom-sheet>
</template>

<script>
import MySnackBar from "@/components/MySnackBar";

export default {
  name: "LoginForm",
  components: {MySnackBar},
  props: ['loginForm'],
  inject:['message'],
  data: () => ({
    user: {
      email: "",
      password: "",
    },
    open: false,
    valid: true,
    showPassword: false,
    // 表单规则
    rules: {
      emailRules: [
        v => !!v || "邮件不能为空！",
        v => /.+@.+\..+/.test(v) || '非法邮件地址！'],
      passwordRules: [
        v => !!v || "密码不能为空！",
        v => (v && v.length >= 4) || "密码不能低于4位！",
      ],
    },
  }),
  methods: {
    //提交登录表单
    submit: function () {
      if(!this.valid) return
      this.$axios.post("login", this.user)
          .then(({data}) => {
            this.$store.commit("login", data.data)
            this.$refs.form.reset()
            this.open = false
            this.showPassword = false
            this.$router.push({name: 'Dashboard'})
            this.message('info', `欢迎您，${data.data.name}`)
          })
          .catch(error => {
            this.message('error', '登录失败')
            console.log(error);
          })
    },
    closeForm: function () {
      this.open = false
      this.showPassword = false
    }
  },
}
</script>
