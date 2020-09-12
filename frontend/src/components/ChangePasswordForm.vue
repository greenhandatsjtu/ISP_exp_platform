<template>
  <v-bottom-sheet v-model="sheet" inset>
    <template v-slot:activator="{ on, attrs }">
      <v-btn icon color="info" v-bind="attrs" v-on="on">
        <v-tooltip top>
          <template v-slot:activator="{on, attrs}">
            <v-icon large v-bind="attrs" v-on="on">
              mdi-form-textbox-password
            </v-icon>
          </template>
          <span>修改密码</span>
        </v-tooltip>
      </v-btn>
    </template>
    <v-sheet class="text-center pa-4">
      <v-form v-model="valid" ref="form">
        <v-text-field v-model="password"
                      clearable
                      label="密码"
                      :type="showPassword ? 'text' : 'password'"
                      :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                      @click:append="showPassword = !showPassword"
                      :rules="rules.passwordRules"></v-text-field>
        <v-text-field
            required
            v-model="confirmPassword"
            label="确认密码"
            type="Password"
            :rules="rules.confirmPasswordRules"></v-text-field>
        <v-row class="mx-auto">
          <v-btn
              color="error"
              @click="closeForm"
          >取消
          </v-btn>
          <v-spacer></v-spacer>
          <!--                    表单填完之前，提交按钮禁用-->
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
export default {
  name: "ChangePasswordForm",
  inject: ['message', 'reload'],
  data: () => ({
    valid: true,
    password: '',
    sheet: false,
    confirmPassword: '',
    showPassword: false,
    rules: {
      passwordRules: [
        v => !!v || "密码不能为空！",
        v => (v && v.length >= 6) || "密码不能低于6位！",
      ],
      confirmPasswordRules: [],
    },
  }),
  methods: {
    submit() {
      let data = new FormData()
      data.append('pwd', this.password)
      this.$axios.post('password', data)
          .then(() => {
            this.message('success', '成功修改密码')
            this.$refs.form.reset()
            this.showPassword = false
            this.sheet = false
            this.reload()
          })
          .catch(err => {
            console.log(err)
            this.message('warning', '修改失败')
          })
    },
    closeForm() {
      this.$refs.form.reset()
      this.showPassword = false
      this.sheet = false
    }
  },
  mounted: function () {
    this.rules.confirmPasswordRules = [
      v => !!v || "密码不能为空！",
      v => v === this.password || "密码不一致！"]
  }
}
</script>

<style scoped>

</style>
