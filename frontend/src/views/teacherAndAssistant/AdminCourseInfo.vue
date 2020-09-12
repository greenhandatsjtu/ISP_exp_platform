<!--课程详情页-->
<template>
  <v-container fluid class="fill-height pa-0">
    <h1 class="text-truncate">{{ course.course_name }}</h1>
    <v-row>
        <manage-student-dialog/>
        <add-experiment-dialog/>
    </v-row>
    <v-row>
      <v-col>
        <course-teacher-card :course="course"/>
      </v-col>
      <v-col>
        <course-assistant-card :course="course"/>
      </v-col>
    </v-row>
    <v-expansion-panels focusable accordion popout>
      <v-expansion-panel>
        <v-expansion-panel-header>课程实验情况</v-expansion-panel-header>
        <v-expansion-panel-content>
          <resources-table scope="course"/>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
    <v-row>
      <experiment-card v-for="(experiment, index) in course.experiments" :key="index"
                       :experiment="experiment" :admin="auth(['sys_admin','teacher','assistant'])"/>
    </v-row>
  </v-container>
</template>

<script>
import ExperimentCard from "@/components/ExperimentCard";
import CourseTeacherCard from "@/components/CourseTeacherCard";
import CourseAssistantCard from "@/components/CourseAssistantCard";
import ManageStudentDialog from "@/components/ManageStudentDialog";
import AddExperimentDialog from "@/components/AddExperimentDialog";
import ResourcesTable from "@/components/ResourcesTable";

export default {
  name: "AdminCourseInfo",
  components: {
    ResourcesTable,
    AddExperimentDialog, ManageStudentDialog, CourseAssistantCard, CourseTeacherCard, ExperimentCard},
  data() {
    return {
      course: {
        course_name: "",
        experiments: [],
        assistant: [],
        teachers: [],
      },
    }
  },
  methods: {},
  beforeCreate() {
    this.$axios.get(`course/${this.$route.params.id}`)
        .then(({data}) => {
          this.course = data.data
        })
  },
}
</script>

<style scoped>

</style>
