<template>
  <v-col :cols="12" :sm="4" class="ma-auto">
    <v-hover v-slot:default="{hover}">
      <v-card
          rounded
          :elevation="hover?20:5"
          class="mb-1"
      >
        <v-img :src="require('@/assets/course.png')" contain max-height="200"/>
        <v-card-title v-text="course.course_name" class="font-weight-bold"></v-card-title>
        <v-card-actions>
          <v-btn
              text
              color="info"
              :to="{name: admin?'AdminCourseInfo':'CourseInfo', params: {id: course.ID}}"
          >
            进入课程
          </v-btn>
          <v-spacer></v-spacer>
          <v-dialog
              light
              v-if="remove"
              v-model="dialog"
              max-width="300"
          >
            <template v-slot:activator="{ on, attrs }">
              <v-btn icon v-bind="attrs" v-on="on">
                <v-icon color="error">mdi-delete</v-icon>
              </v-btn>
            </template>

            <v-card>
              <v-card-title class="headline grey lighten-2">
                确定删除课程？
              </v-card-title>
              <v-divider/>
              <v-card-actions>
                <v-btn color="primary" @click="dialog=false">
                  取消
                </v-btn>
                <v-spacer/>
                <v-btn color="warning" @click="deleteCourse">
                  确定
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
        </v-card-actions>
      </v-card>
    </v-hover>
  </v-col>
</template>

<script>
export default {
  name: "CourseCard",
  props: ['course', 'admin', 'remove'],
  inject: ['message'],
  data:()=>({
    dialog: false,
  }),
  methods: {
    deleteCourse() {
      this.$axios.get(`admin/course/${this.course.ID}/delete`)
          .then(() => {
            this.message('success', '成功删除课程')
            this.dialog = false
          })
          .catch(err => {
            console.log(err)
            this.message('warning', '删除课程失败')
          })
    }
  }
}
</script>

<style scoped>

</style>
