<template>
  <v-dialog
      v-model="dialog"
      max-width="400"
  >
    <template v-slot:activator="{ on, attrs }">
      <v-btn class="ma-auto primary" v-bind="attrs" v-on="on">
        添加学生
      </v-btn>
    </template>

    <v-card class="px-4">
      <v-card-title>
        新增学生
      </v-card-title>
      <v-form
          ref="form"
          v-model="valid">
        <v-text-field
            v-model="student.student_number"
            label="学号"
            required
            :rules="IDRules"
            counter
        ></v-text-field>
        <v-text-field
            v-model="student.name"
            label="姓名"
            required
            :rules="nameRules"
            counter
        ></v-text-field>
        <v-text-field
            v-model="student.class_number"
            label="班号"
            required
            :rules="classRules"
            counter
        ></v-text-field>
        <v-text-field
            v-model="student.email"
            label="邮箱"
            required
            type="email"
            :rules="emailRules"
            counter
        ></v-text-field>
      </v-form>
      <v-card-actions>
        <v-btn color="success" @click="dialog=false">
          取消
        </v-btn>
        <v-spacer/>
        <v-btn color="primary" :disabled="!valid" @click="addStudent">
          确定
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: "AddStudentDialog",
  inject: ['reload', 'message'],
  data: function () {
    return {
      dialog: false,
      valid: true,
      student: {
        name: '',
        email: '',
        class_number: '',
        student_number: '',
      },
      nameRules: [
        v => !!v || '学生名未填！',
        v => v.length <= 50 || '学生名过长！',
      ],
      emailRules: [
        v => !!v || "邮件不能为空！",
        v => v.length <= 100 || '邮箱过长！',
        v => /.+@.+\..+/.test(v) || '非法邮件地址！'
      ],
      classRules: [
        v => !!v || "班级未填！",
        v => v.length <= 20 || '班号过长！',
      ],
      IDRules: [
        v => !!v || "学号未填！",
        v => v.length <= 20 || '学号过长！',
      ],
    }
  },
  methods: {
    addStudent() {
      this.$axios.post('admin/students', [this.student]).then(() => {
        this.message('success', '添加成功')
        this.reload()
      }).catch(err => {
        console.log(err)
        this.message('error', '添加失败')
      })
    }
  }
}
</script>

<style scoped>

</style>
