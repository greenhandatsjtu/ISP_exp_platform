<template>
  <v-dialog
      v-model="dialog"
      max-width="700"
  >
    <template v-slot:activator="{ on, attrs }">
      <v-btn class="ma-auto" color="indigo" dark v-bind="attrs" v-on="on">
        <v-icon>
          mdi-account-cog
        </v-icon>
        管理学生
      </v-btn>
    </template>

    <v-card>
      <v-card-title>
        管理学生
      </v-card-title>
      <v-text-field
          class="mx-2 my-0"
          v-model="search"
          append-icon="mdi-magnify"
          label="搜索学生"
          single-line
          hide-details
      ></v-text-field>
      <v-data-table
          v-model="selected"
          :search="search"
          class="ma-2 elevation-10"
          :headers="headers"
          :items="allStudents"
          :items-per-page="5"
          show-group-by
          item-key="student_number"
          :sort-by="['student_number']"
          show-select
      >
      </v-data-table>
      <v-card-actions>
        <v-btn color="primary" @click="updateAssistants">
          确定
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: "ManageStudentDialog",
  inject: ['message','reload'],
  data: function () {
    return {
      dialog: false,
      search: '',
      headers: [
        {text: '学号', value: 'student_number'},
        {text: '班级', value: 'class_number'},
        {text: '姓名', value: 'user.name'},
        {text: '邮箱', value: 'user.email'},
      ],
      allStudents: [],
      selected: []
    }
  },
  created() {
    this.$axios.get('admin/student')
        .then(({data}) => {
          this.allStudents = data.data
        })
        .catch(err => {
          console.error(err)
          this.message('warning', '获取学生失败')
        })
  },
  mounted() {
    this.$axios.get(`admin/course/${this.$route.params.id}/student`)
        .then(({data}) => {
          this.selected = data.data
        })
        .catch(err => {
          console.error(err)
          this.message('warning', '获取已分配学生失败')
        })
  },
  methods: {
    updateAssistants() {
      this.$axios.post(`admin/course/${this.$route.params.id}/students`, {
        students: this.selected.map(x => x.ID)
      })
          .then(() => {
            this.message('success', '分配成功')
            this.dialog = false
            this.reload()
          })
          .catch(() => {
            this.message('warning', '分配失败')
          })
    }
  }
}
</script>

<style scoped>

</style>
