<!--教师和助教的课程页面-->
<template>
  <v-row>
    <course-card v-for="(course, index) in courses" :key="index" :course="course" :admin="true"
                 :remove="auth(['teach_admin'])"/>
  </v-row>
</template>

<script>

import CourseCard from "@/components/CourseCard";

export default {
  name: "AdminCourses",
  components: {CourseCard},
  inject: ['message'],
  data: () => ({
    courses: [],
  }),
  created() {
    this.$axios.get('admin/courses')
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
