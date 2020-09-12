<!--课程详情页-->
<template>
  <v-container fluid class="fill-height pa-0">
    <h1 class="text-truncate">{{ course.course_name }}</h1>
    <v-row v-if="!!grade">
      <v-card-title class="headline"><span class="font-weight-bold">成绩：</span>{{ grade }}</v-card-title>
    </v-row>
    <v-row>
      <v-col>
        <course-teacher-card :course="course"/>
      </v-col>
      <v-col>
        <course-assistant-card :course="course"/>
      </v-col>
    </v-row>
    <v-row>
      <experiment-card v-for="(experiment, index) in course.experiments" :key="index"
                       :experiment="experiment"></experiment-card>
    </v-row>
  </v-container>
</template>

<script>
import ExperimentCard from "@/components/ExperimentCard";
import CourseTeacherCard from "@/components/CourseTeacherCard";
import CourseAssistantCard from "@/components/CourseAssistantCard";

export default {
  name: "CourseInfo",
  components: {CourseAssistantCard, CourseTeacherCard, ExperimentCard},
  data() {
    return {
      course: {
        course_name: "",
        experiments: [],
        assistant: [],
        teachers: [],
      },
      grade: null,
    }
  },
  methods: {
    getGrade() {
      this.$axios.get(`student/course/${this.$route.params.id}/grade`)
          .then(({data}) => {
            this.grade = data.data.grade
          })
    },
  },
  beforeCreate() {
    this.$axios.get(`course/${this.$route.params.id}`)
        .then(({data}) => {
          this.course = data.data
        })
  },
  created() {
    this.getGrade()
  }
}
</script>

<style scoped>

</style>
