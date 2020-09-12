<template>
  <v-container fluid class="fill-height">
    <users-table/>
    <v-bottom-sheet light v-model="sheet" inset max-width="500">
      <v-card-title class="grey lighten-1 headline black--text">
        新建用户
      </v-card-title>
      <template v-slot:activator="{ on, attrs }">
        <v-btn fab fixed bottom right color="primary" v-bind="attrs" v-on="on">
          <v-icon x-large>
            mdi-plus
          </v-icon>
        </v-btn>
      </template>
      <v-sheet class="text-center pa-4">
        <v-form v-model="valid" ref="form">
          <v-text-field v-model="user.name"
                        required
                        clearable
                        label="姓名"
                        :rules="rules.nameRules"></v-text-field>
          <v-text-field v-model="user.email"
                        required
                        clearable
                        label="邮箱"
                        type="Email"
                        :rules="rules.emailRules"></v-text-field>
          <v-text-field v-model="user.password"
                        clearable
                        label="密码"
                        :rules="rules.passwordRules"></v-text-field>
          <v-row class="mx-auto">
            <v-btn
                color="error"
                @click="sheet = false"
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
  </v-container>
</template>

<script>
import UsersTable from "@/components/UsersTable";

export default {
  name: "ManageUsers",
  components: {UsersTable},
  inject: ['message', 'reload'],
  data() {
    return {
      valid: false,
      rules: {
        nameRules: [
          v => !!v || "姓名不能为空！",
          v => (v && v.length <= 20) || "姓名过长！",
        ],
        emailRules: [
          v => !!v || "邮件不能为空！",
        ],
        passwordRules: [
          v => !!v || "密码不能为空！",
          v => (v && v.length >= 6) || "密码不能低于6位！",
        ],
      },
      sheet: false,
      user: {
        email: '',
        password: '',
        name: '',
      }
    }
  },
  methods: {
    submit() {
      this.$axios.post('admin/user', this.user)
          .then(() => {
            this.message('success', `成功添加用户：${this.user.name}`)
            this.reload()
          })
          .catch(() => {
            this.message('error', '添加用户失败')
          })
    }
  }
}
</script>

<style scoped>

</style>
