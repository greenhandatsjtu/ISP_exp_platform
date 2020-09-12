<template>
  <v-container fluid class="fill-height">
    <v-row>
      <course-card v-for="(course, index) in courses" :key="index" :course="course" :admin="true" :remove="true"/>
    </v-row>
    <v-bottom-sheet v-model="sheet" inset max-width="500">
      <template v-slot:activator="{ on, attrs }">
        <v-btn fab fixed bottom right color="primary" v-bind="attrs" v-on="on">
          <v-icon x-large>
            mdi-plus
          </v-icon>
        </v-btn>
      </template>
      <v-sheet class="text-center pa-4">
        <v-form v-model="valid" ref="form">
          <v-text-field v-model="name"
                        required
                        counter
                        label="课程名"
                        :rules="nameRules"></v-text-field>
          <v-divider/>
          <v-row class="mx-auto">
            <v-btn
                color="error"
                @click="sheet=false"
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
import CourseCard from "@/components/CourseCard";

export default {
  name: "AllCourses",
  components: {CourseCard},
  inject: ['message','reload'],
  data: () => ({
    courses: [],
    sheet: false,
    valid: true,
    name: '',
    nameRules: [
      v => !!v || "课程名不能为空！",
      v => v.length <= 25 || '课程名过长！',],
  }),
  methods: {
    submit() {
      this.$axios.post('admin/course', {
        course_name: this.name
      })
          .then(() => {
            this.message('success', '新增成功')
            this.sheet = false
            this.reload()
          })
          .catch(err => {
            console.log(err)
            this.message('warning', '新增失败')
          })
    }
  },
  created() {
    this.$axios.get('course')
        .then(({data}) => {
          this.courses = data.data
        })
        .catch(error => {
          this.message('error', '获取课程失败')
          console.log(error);
        });
  }
}
</script>

<style scoped>

</style>
